package bridge

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc/metadata"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	batchv1 "k8s.io/api/batch/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"pixielabs.ai/pixielabs/src/cloud/vzconn/vzconnpb"
	"pixielabs.ai/pixielabs/src/shared/cvmsgspb"
	"pixielabs.ai/pixielabs/src/utils"
	"pixielabs.ai/pixielabs/src/vizier/utils/messagebus"
)

// UpdaterJobYAML is the YAML that should be applied for the updater job.
const UpdaterJobYAML string = `---
apiVersion: batch/v1
kind: Job
metadata:
  name: vizier-upgrade-job
spec:
  template:
    metadata:
      name: vizier-upgrade-job
    spec:
      serviceAccountName: updater-service-account
      containers:
      - name: updater
        image: gcr.io/pl-dev-infra/vizier/vizier_updater_image:__VIZIER_UPDATER_IMAGE_TAG__
        envFrom:
        - configMapRef:
            name: pl-cloud-config
        env:
        - name: PL_CLOUD_TOKEN
          valueFrom:
            secretKeyRef:
              name: pl-update-job-secrets
              key: cloud-token
        - name: PL_VIZIER_VERSION
          value: __PL_VIZIER_VERSION__
        - name: PL_REDEPLOY_ETCD
          value: __PL_REDEPLOY_ETCD__
        - name: PL_CLIENT_TLS_CERT
          value: /certs/client.crt
        - name: PL_CLIENT_TLS_KEY
          value: /certs/client.key
        - name: PL_SERVER_TLS_CERT
          value: /certs/server.crt
        - name: PL_SERVER_TLS_KEY
          value: /certs/server.key
        - name: PL_TLS_CA_CERT
          value: /certs/ca.crt
        volumeMounts:
        - name: certs
          mountPath: /certs
      imagePullSecrets:
      - name: pl-image-secret
      volumes:
      - name: certs
        secret:
          secretName: service-tls-certs
      restartPolicy: "Never"
  backoffLimit: 1
  parallelism: 1
  completions: 1`

const (
	heartbeatIntervalS = 5 * time.Second
	// HeartbeatTopic is the topic that heartbeats are written to.
	HeartbeatTopic                = "heartbeat"
	registrationTimeout           = 30 * time.Second
	passthroughReplySubjectPrefix = "v2c.reply-"
	vizStatusCheckFailInterval    = 10 * time.Second
)

// ErrRegistrationTimeout is the registration timeout error.
var ErrRegistrationTimeout = errors.New("Registration timeout")

const upgradeJobName = "vizier-upgrade-job"

// VizierInfo fetches information about Vizier.
type VizierInfo interface {
	GetAddress() (string, int32, error)
	GetVizierClusterInfo() (*cvmsgspb.VizierClusterInfo, error)
	GetK8sState() (map[string]*cvmsgspb.PodStatus, int32, time.Time)
	ParseJobYAML(yamlStr string, imageTag map[string]string, envSubtitutions map[string]string) (*batchv1.Job, error)
	LaunchJob(j *batchv1.Job) (*batchv1.Job, error)
	CreateSecret(string, map[string]string) error
	WaitForJobCompletion(string) (bool, error)
	DeleteJob(string) error
	GetJob(string) (*batchv1.Job, error)
	GetClusterUID() (string, error)
	UpdateClusterID(string) error
}

// VizierHealthChecker is the interface that gets information on health of a Vizier.
type VizierHealthChecker interface {
	GetStatus() (time.Time, error)
}

// Bridge is the NATS<->GRPC bridge.
type Bridge struct {
	vizierID      uuid.UUID
	jwtSigningKey string
	sessionID     int64
	deployKey     string

	vzConnClient vzconnpb.VZConnServiceClient
	vzInfo       VizierInfo
	vizChecker   VizierHealthChecker

	hbSeqNum int64

	nc         *nats.Conn
	natsCh     chan *nats.Msg
	registered bool
	// There are a two sets of streams that we manage for the GRPC side. The incoming
	// data and the outgoing data. GRPC does not natively provide a channel based interface
	// so we wrap the Send/Recv calls with goroutines that are responsible for
	// performing the read/write operations.
	//
	// If the GRPC connection gets disrupted, we close all the readers and writer routines, but we leave the
	// channels in place so that data does not get lost. The data will simply be resent
	// once the connection comes back alive. If data is lost due to a crash, the rest of the system
	// is resilient to this loss, but reducing it is optimal to prevent a lot of replay traffic.

	grpcOutCh chan *vzconnpb.V2CBridgeMessage
	grpcInCh  chan *vzconnpb.C2VBridgeMessage
	// Explicitly prioritize passthrough traffic to prevent script failure under load.
	ptOutCh chan *vzconnpb.V2CBridgeMessage
	// This tracks the message we are trying to send, but has not been sent yet.
	pendingGRPCOutMsg *vzconnpb.V2CBridgeMessage

	quitCh chan bool      // Channel is used to signal that things should shutdown.
	wg     sync.WaitGroup // Tracks all the active goroutines.
	wdWg   sync.WaitGroup // Tracks all the active goroutines.

	updateRunning bool // True if an update is running.
	updateFailed  bool // True if an update has failed (sticky).
}

// New creates a cloud connector to cloud bridge.
func New(vizierID uuid.UUID, jwtSigningKey string, deployKey string, sessionID int64, vzClient vzconnpb.VZConnServiceClient, vzInfo VizierInfo, nc *nats.Conn, checker VizierHealthChecker) *Bridge {
	return &Bridge{
		vizierID:      vizierID,
		jwtSigningKey: jwtSigningKey,
		deployKey:     deployKey,
		sessionID:     sessionID,
		vzConnClient:  vzClient,
		vizChecker:    checker,
		vzInfo:        vzInfo,
		hbSeqNum:      0,
		nc:            nc,
		// Buffer NATS channels to make sure we don't back-pressure NATS
		natsCh:            make(chan *nats.Msg, 5000),
		registered:        false,
		ptOutCh:           make(chan *vzconnpb.V2CBridgeMessage, 5000),
		grpcOutCh:         make(chan *vzconnpb.V2CBridgeMessage, 5000),
		grpcInCh:          make(chan *vzconnpb.C2VBridgeMessage, 5000),
		pendingGRPCOutMsg: nil,
		quitCh:            make(chan bool),
		wg:                sync.WaitGroup{},
		wdWg:              sync.WaitGroup{},
	}
}

// WatchDog watches and make sure the bridge is functioning. If not commits suicide to try to self-heal.
func (s *Bridge) WatchDog() {
	defer s.wdWg.Done()
	t := time.NewTicker(30 * time.Second)

	for {
		lastHbSeq := atomic.LoadInt64(&s.hbSeqNum)
		select {
		case <-s.quitCh:
			log.Trace("Quitting watchdog")
			return
		case <-t.C:
			currentHbSeqNum := atomic.LoadInt64(&s.hbSeqNum)
			if currentHbSeqNum == lastHbSeq {
				log.Fatal("Heartbeat messages failed, assuming stream is dead. Killing self to restart...")
			}
		}
	}
}

// WaitForUpdater waits for the update job to complete, if any.
func (s *Bridge) WaitForUpdater() {
	defer func() {
		s.updateRunning = false
	}()
	ok, err := s.vzInfo.WaitForJobCompletion(upgradeJobName)
	if err != nil {
		log.WithError(err).Error("Error while trying to watch vizier-upgrade-job")
		return
	}
	s.updateFailed = !ok
	err = s.vzInfo.DeleteJob(upgradeJobName)
	if err != nil {
		log.WithError(err).Error("Error deleting upgrade job")
	}
}

// RegisterDeployment registers the vizier cluster using the deployment key.
func (s *Bridge) RegisterDeployment() error {
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "X-API-KEY", s.deployKey)
	clusterInfo, err := s.vzInfo.GetVizierClusterInfo()
	if err != nil {
		return err
	}
	resp, err := s.vzConnClient.RegisterVizierDeployment(ctx, &vzconnpb.RegisterVizierDeploymentRequest{
		K8sClusterUID:     clusterInfo.ClusterUID,
		K8sClusterName:    clusterInfo.ClusterName,
		K8sClusterVersion: clusterInfo.ClusterVersion,
	})
	if err != nil {
		return err
	}

	// Get cluster ID and assign to secrets.
	s.vizierID = utils.UUIDFromProtoOrNil(resp.VizierID)

	return s.vzInfo.UpdateClusterID(s.vizierID.String())
}

// RunStream manages starting and restarting the stream to VZConn.
func (s *Bridge) RunStream() {
	if s.vzConnClient == nil {
		var vzClient vzconnpb.VZConnServiceClient
		var err error

		connect := func() error {
			log.Info("Connecting to VZConn...")
			vzClient, err = NewVZConnClient()
			if err != nil {
				log.WithError(err).Error("Failed to connect to VZConn")
			}
			return err
		}

		backOffOpts := backoff.NewExponentialBackOff()
		backOffOpts.InitialInterval = 30 * time.Second
		backOffOpts.Multiplier = 2
		backOffOpts.MaxElapsedTime = 30 * time.Minute
		err = backoff.Retry(connect, backOffOpts)
		if err != nil {
			log.WithError(err).Fatal("Could not connect to VZConn")
		}
		log.Info("Successfully connected to VZConn")
		s.vzConnClient = vzClient
	}

	natsTopic := messagebus.V2CTopic("*")
	log.WithField("topic", natsTopic).Trace("Subscribing to NATS")
	natsSub, err := s.nc.ChanSubscribe(natsTopic, s.natsCh)
	if err != nil {
		log.WithError(err).Fatal("Failed to subscribe to NATS.")
	}
	defer natsSub.Unsubscribe()
	// Set large limits on message size and count.
	natsSub.SetPendingLimits(1e7, 1e7)

	// Check if there is an existing update job. If so, then set the status to "UPDATING".
	_, err = s.vzInfo.GetJob(upgradeJobName)
	if err != nil && !k8sErrors.IsNotFound(err) {
		log.WithError(err).Fatal("Could not check for upgrade job")
	}
	if err == nil { // There is an upgrade job running.
		s.updateRunning = true
		go s.WaitForUpdater()
	}

	// Get the cluster ID, if not already specified.
	if s.vizierID == uuid.Nil {
		err = s.RegisterDeployment()
		if err != nil {
			log.WithError(err).Fatal("Failed to register vizier deployment")
		}
	}

	s.wdWg.Add(1)
	go s.WatchDog()

	for {
		s.registered = false
		select {
		case <-s.quitCh:
			return
		default:
			log.Trace("Starting stream")
			errCh := make(chan error)
			err := s.StartStream(errCh)
			if err == nil {
				log.Trace("Stream ending")
			} else {
				log.WithError(err).Error("Stream errored. Restarting stream")
			}
			close(errCh)
		}
	}
}

func (s *Bridge) handleUpdateMessage(msg *types.Any) error {
	pb := &cvmsgspb.UpdateOrInstallVizierRequest{}
	err := types.UnmarshalAny(msg, pb)
	if err != nil {
		log.WithError(err).Error("Could not unmarshal update req message")
		return err
	}

	// TODO(michelle): Fill in the YAML contents.
	job, err := s.vzInfo.ParseJobYAML(UpdaterJobYAML, map[string]string{"updater": pb.Version}, map[string]string{
		"PL_VIZIER_VERSION": pb.Version,
		"PL_REDEPLOY_ETCD":  fmt.Sprintf("%v", pb.RedeployEtcd),
	})
	if err != nil {
		log.WithError(err).Error("Could not parse job")
		return err
	}

	err = s.vzInfo.CreateSecret("pl-update-job-secrets", map[string]string{
		"cloud-token": pb.Token,
	})
	if err != nil {
		log.WithError(err).Error("Failed to create job secrets")
		return err
	}

	_, err = s.vzInfo.LaunchJob(job)
	if err != nil {
		log.WithError(err).Error("Could not launch job")
		return err
	}

	// Send response message to indicate update job has started.
	m := cvmsgspb.UpdateOrInstallVizierResponse{
		UpdateStarted: true,
	}
	reqAnyMsg, err := types.MarshalAny(&m)
	if err != nil {
		return err
	}

	v2cMsg := cvmsgspb.V2CMessage{
		Msg: reqAnyMsg,
	}
	b, err := v2cMsg.Marshal()
	if err != nil {
		return err
	}
	err = s.nc.Publish(messagebus.V2CTopic("VizierUpdateResponse"), b)
	if err != nil {
		log.WithError(err).Error("Failed to publish VizierUpdateResponse")
		return err
	}

	return nil
}

func (s *Bridge) doRegistrationHandshake(stream vzconnpb.VZConnService_NATSBridgeClient) error {
	addr, _, err := s.vzInfo.GetAddress()
	if err != nil {
		log.WithError(err).Error("Unable to get vizier proxy address")
	}

	clusterInfo, err := s.vzInfo.GetVizierClusterInfo()
	if err != nil {
		log.WithError(err).Error("Unable to get k8s cluster info")
	}
	// Send over a registration request and wait for ACK.
	regReq := &cvmsgspb.RegisterVizierRequest{
		VizierID:    utils.ProtoFromUUID(&s.vizierID),
		JwtKey:      s.jwtSigningKey,
		Address:     addr,
		ClusterInfo: clusterInfo,
	}

	err = s.publishBridgeSync(stream, "register", regReq)
	if err != nil {
		return err
	}

	for {
		select {
		case <-time.After(registrationTimeout):
			log.Error("Timeout with registration terminating stream")
			return ErrRegistrationTimeout
		case resp := <-s.grpcInCh:
			// Try to receive the registerAck.
			if resp.Topic != "registerAck" {
				log.Error("Unexpected message type while waiting for ACK")
			}
			registerAck := &cvmsgspb.RegisterVizierAck{}
			err = types.UnmarshalAny(resp.Msg, registerAck)
			if err != nil {
				return err
			}
			switch registerAck.Status {
			case cvmsgspb.ST_FAILED_NOT_FOUND:
				return errors.New("registration not found, cluster unknown in pixie-cloud")
			case cvmsgspb.ST_OK:
				s.registered = true
				return nil
			default:
				return errors.New("registration unsuccessful: " + err.Error())
			}
		}
	}
}

// StartStream starts the stream between the cloud connector and Vizier connector.
func (s *Bridge) StartStream(errCh chan error) error {
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := s.vzConnClient.NATSBridge(ctx)
	if err != nil {
		log.WithError(err).Error("Error starting stream")
		return err
	}
	// Wait for  all goroutines to terminate.
	defer func() {
		s.wg.Wait()
	}()

	// Setup the stream reader go routine.
	done := make(chan bool)
	defer close(done)
	// Cancel the stream to make sure everything get shutdown properly.
	defer func() {
		cancel()
	}()

	s.wg.Add(1)
	go s.startStreamGRPCReader(stream, done, errCh)
	s.wg.Add(1)
	go s.startStreamGRPCWriter(stream, done, errCh)

	if !s.registered {
		// Need to do registration handshake before we allow any cvmsgs.
		err := s.doRegistrationHandshake(stream)
		if err != nil {
			return err
		}
	}

	log.Trace("Registration Complete.")
	s.wg.Add(1)
	err = s.HandleNATSBridging(stream, done, errCh)
	return err
}

func (s *Bridge) startStreamGRPCReader(stream vzconnpb.VZConnService_NATSBridgeClient, done chan bool, errCh chan<- error) {
	defer s.wg.Done()
	log.Trace("Starting GRPC reader stream")
	defer log.Trace("Closing GRPC read stream")
	for {
		select {
		case <-s.quitCh:
			return
		case <-stream.Context().Done():
			return
		case <-done:
			log.Info("Closing GRPC reader because of <-done")
			return
		default:
			log.Trace("Waiting for next message")
			msg, err := stream.Recv()
			if err != nil && err == io.EOF {
				log.Trace("Stream has closed(Read)")
				// stream closed.
				return
			}
			if err != nil && errors.Is(err, context.Canceled) {
				log.Trace("Stream has been cancelled")
				return
			}
			if err != nil {
				log.WithError(err).Error("Got a stream read error")
				return
			}
			s.grpcInCh <- msg
		}
	}
}

func (s *Bridge) startStreamGRPCWriter(stream vzconnpb.VZConnService_NATSBridgeClient, done chan bool, errCh chan<- error) {
	defer s.wg.Done()
	log.Trace("Starting GRPC writer stream")
	defer log.Trace("Closing GRPC writer stream")

	sendMsg := func(m *vzconnpb.V2CBridgeMessage) {
		s.pendingGRPCOutMsg = m
		// Write message to GRPC if it exists.
		err := stream.Send(s.pendingGRPCOutMsg)
		if err != nil {
			// Need to resend this message.
			return
		}
		s.pendingGRPCOutMsg = nil
		return
	}

	for {
		// Pending message try to send it first.
		if s.pendingGRPCOutMsg != nil {
			err := stream.Send(s.pendingGRPCOutMsg)
			if err != nil {
				// Error sending message. The stream might terminate in the middle so the select
				// guards against exited goroutines to prevent a hang.
				select {
				case errCh <- err:
				case <-done:
				case <-s.quitCh:
				}

				return
			}
			s.pendingGRPCOutMsg = nil
		}
		// Try to send PT traffic first.
		select {
		case <-s.quitCh:
			return
		case <-stream.Context().Done():
			log.Trace("Write stream has closed")
			return
		case <-done:
			log.Trace("Closing GRPC writer because of <-done")
			stream.CloseSend()
			// Quit called.
			return
		case m := <-s.ptOutCh:
			sendMsg(m)
			break
		default:
		}

		select {
		case <-stream.Context().Done():
			log.Trace("Write stream has closed")
			return
		case <-done:
			log.Trace("Closing GRPC writer because of <-done")
			stream.CloseSend()
			// Quit called.
			return
		case m := <-s.ptOutCh:
			sendMsg(m)
		case m := <-s.grpcOutCh:
			sendMsg(m)
		}
	}
}

func (s *Bridge) parseV2CNatsMsg(data *nats.Msg) (*cvmsgspb.V2CMessage, string, error) {
	v2cPrefix := messagebus.V2CTopic("")
	topic := strings.TrimPrefix(data.Subject, v2cPrefix)

	// Message over nats should be wrapped in a V2CMessage.
	v2cMsg := &cvmsgspb.V2CMessage{}
	err := proto.Unmarshal(data.Data, v2cMsg)
	if err != nil {
		return nil, "", err
	}
	return v2cMsg, topic, nil
}

// HandleNATSBridging routes message to and from cloud NATS.
func (s *Bridge) HandleNATSBridging(stream vzconnpb.VZConnService_NATSBridgeClient, done chan bool, errCh chan error) error {
	defer s.wg.Done()
	defer log.Info("Closing NATS Bridge")
	// Vizier -> Cloud side:
	// 1. Listen to NATS on v2c.<topic>.
	// 2. Extract Topic from the stream name above.
	// 3. Wrap the message and throw it over the wire.

	// Cloud -> Vizier side:
	// 1. Read the stream.
	// 2. For cvmsgs of type: C2VBridgeMessage, read the topic
	//    and throw it onto nats under c2v.topic

	log.Info("Starting NATS bridge.")
	hbChan := s.generateHeartbeats(done)

	for {
		select {
		case <-s.quitCh:
			return nil
		case <-done:
			return nil
		case e := <-errCh:
			log.WithError(e).Error("GRPC error, terminating stream")
			return e
		case data := <-s.natsCh:
			v2cPrefix := messagebus.V2CTopic("")
			if !strings.HasPrefix(data.Subject, v2cPrefix) {
				return errors.New("invalid subject: " + data.Subject)
			}

			v2cMsg, topic, err := s.parseV2CNatsMsg(data)
			if err != nil {
				log.WithError(err).Error("Failed to parse message")
				return err
			}

			if strings.HasPrefix(data.Subject, passthroughReplySubjectPrefix) {
				// Passthrough message.
				err = s.publishPTBridgeCh(topic, v2cMsg.Msg)
				if err != nil {
					return err
				}
			} else {
				err = s.publishBridgeCh(topic, v2cMsg.Msg)
				if err != nil {
					return err
				}
			}
		case bridgeMsg := <-s.grpcInCh:
			log.
				WithField("type", bridgeMsg.Msg.TypeUrl).
				Info("Got Message on GRPC channel")
			if bridgeMsg == nil {
				return nil
			}

			log.
				WithField("msg", bridgeMsg.String()).
				Trace("Got Message on GRPC channel")

			if bridgeMsg.Topic == "VizierUpdate" {
				err := s.handleUpdateMessage(bridgeMsg.Msg)
				if err != nil {
					log.WithError(err).Error("Failed to launch vizier update job")
				}
				continue
			}

			topic := messagebus.C2VTopic(bridgeMsg.Topic)

			natsMsg := &cvmsgspb.C2VMessage{
				VizierID: s.vizierID.String(),
				Msg:      bridgeMsg.Msg,
			}
			b, err := natsMsg.Marshal()
			if err != nil {
				log.WithError(err).Error("Failed to marshal")
				return err
			}

			log.WithField("topic", topic).
				WithField("msg", natsMsg.String()).
				Trace("Publishing to NATS")
			err = s.nc.Publish(topic, b)
			if err != nil {
				log.WithError(err).Error("Failed to publish")
				return err
			}
		case hbMsg := <-hbChan:
			log.WithField("heartbeat", hbMsg.GoString()).Trace("Sending heartbeat")
			err := s.publishProtoToBridgeCh(HeartbeatTopic, hbMsg)
			if err != nil {
				return err
			}
		case <-stream.Context().Done():
			log.Info("Stream has been closed, shutting down grpc readers")
			return nil
		}
	}
	return nil
}

// Stop terminates the server. Don't reuse this server object after stop has been called.
func (s *Bridge) Stop() {
	close(s.quitCh)
	// Wait fo all goroutines to stop.
	s.wg.Wait()
	s.wdWg.Wait()
}

func (s *Bridge) publishBridgeCh(topic string, msg *types.Any) error {
	wrappedReq := &vzconnpb.V2CBridgeMessage{
		Topic:     topic,
		SessionId: s.sessionID,
		Msg:       msg,
	}

	// Don't stall the queue for regular message.
	select {
	case s.grpcOutCh <- wrappedReq:
	default:
		log.WithField("Topic", wrappedReq.Topic).Error("Dropping message because of queue backoff")
	}
	return nil
}

func (s *Bridge) publishPTBridgeCh(topic string, msg *types.Any) error {
	wrappedReq := &vzconnpb.V2CBridgeMessage{
		Topic:     topic,
		SessionId: s.sessionID,
		Msg:       msg,
	}
	s.ptOutCh <- wrappedReq
	return nil
}

func (s *Bridge) publishProtoToBridgeCh(topic string, msg proto.Message) error {
	anyMsg, err := types.MarshalAny(msg)
	if err != nil {
		return err
	}

	return s.publishBridgeCh(topic, anyMsg)
}

func (s *Bridge) publishBridgeSync(stream vzconnpb.VZConnService_NATSBridgeClient, topic string, msg proto.Message) error {
	anyMsg, err := types.MarshalAny(msg)
	if err != nil {
		return err
	}

	wrappedReq := &vzconnpb.V2CBridgeMessage{
		Topic:     topic,
		SessionId: s.sessionID,
		Msg:       anyMsg,
	}

	if err := stream.Send(wrappedReq); err != nil {
		return err
	}
	return nil
}

func (s *Bridge) generateHeartbeats(done <-chan bool) (hbCh chan *cvmsgspb.VizierHeartbeat) {
	hbCh = make(chan *cvmsgspb.VizierHeartbeat)

	sendHeartbeat := func() {
		addr, port, err := s.vzInfo.GetAddress()
		if err != nil {
			log.WithError(err).Error("Failed to get vizier address")
		}
		podStatuses, numNodes, updatedTime := s.vzInfo.GetK8sState()
		hbMsg := &cvmsgspb.VizierHeartbeat{
			VizierID:               utils.ProtoFromUUID(&s.vizierID),
			Time:                   time.Now().UnixNano(),
			SequenceNumber:         atomic.LoadInt64(&s.hbSeqNum),
			Address:                addr,
			Port:                   port,
			NumNodes:               numNodes,
			PodStatuses:            podStatuses,
			PodStatusesLastUpdated: updatedTime.UnixNano(),
			Status:                 s.currentStatus(),
			BootstrapMode:          viper.GetBool("bootstrap_mode"),
			BootstrapVersion:       viper.GetString("bootstrap_version"),
		}
		select {
		case <-s.quitCh:
			return
		case <-done:
			return
		case hbCh <- hbMsg:
			atomic.AddInt64(&s.hbSeqNum, 1)
		}
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(heartbeatIntervalS)
		defer ticker.Stop()

		// Send first heartbeat.
		sendHeartbeat()

		for {
			select {
			case <-s.quitCh:
				log.Info("Stopping heartbeat routine")
				return
			case <-done:
				log.Info("Stopping heartbeat routine")
				return
			case <-ticker.C:
				sendHeartbeat()
			}
		}
	}()
	return
}

func (s *Bridge) currentStatus() cvmsgspb.VizierStatus {
	if s.updateRunning && !s.updateFailed {
		return cvmsgspb.VZ_ST_UPDATING
	} else if s.updateFailed {
		return cvmsgspb.VZ_ST_UPDATE_FAILED
	}

	t, status := s.vizChecker.GetStatus()
	if time.Now().Sub(t) > vizStatusCheckFailInterval {
		return cvmsgspb.VZ_ST_UNKNOWN
	}
	if status != nil {
		return cvmsgspb.VZ_ST_UNHEALTHY
	}
	return cvmsgspb.VZ_ST_HEALTHY
}
