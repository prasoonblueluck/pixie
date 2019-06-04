// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: src/vizier/services/metadata/datapb/data.proto

package datapb

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	proto1 "pixielabs.ai/pixielabs/src/common/uuid/proto"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type AgentData struct {
	AgentID         *proto1.UUID `protobuf:"bytes,1,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	HostInfo        *HostInfo    `protobuf:"bytes,2,opt,name=host_info,json=hostInfo,proto3" json:"host_info,omitempty"`
	CreateTimeNS    int64        `protobuf:"varint,3,opt,name=create_time_ns,json=createTimeNs,proto3" json:"create_time_ns,omitempty"`
	LastHeartbeatNS int64        `protobuf:"varint,4,opt,name=last_heartbeat_ns,json=lastHeartbeatNs,proto3" json:"last_heartbeat_ns,omitempty"`
}

func (m *AgentData) Reset()      { *m = AgentData{} }
func (*AgentData) ProtoMessage() {}
func (*AgentData) Descriptor() ([]byte, []int) {
	return fileDescriptor_98f8b3e566c670dc, []int{0}
}
func (m *AgentData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AgentData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AgentData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AgentData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AgentData.Merge(m, src)
}
func (m *AgentData) XXX_Size() int {
	return m.Size()
}
func (m *AgentData) XXX_DiscardUnknown() {
	xxx_messageInfo_AgentData.DiscardUnknown(m)
}

var xxx_messageInfo_AgentData proto.InternalMessageInfo

func (m *AgentData) GetAgentID() *proto1.UUID {
	if m != nil {
		return m.AgentID
	}
	return nil
}

func (m *AgentData) GetHostInfo() *HostInfo {
	if m != nil {
		return m.HostInfo
	}
	return nil
}

func (m *AgentData) GetCreateTimeNS() int64 {
	if m != nil {
		return m.CreateTimeNS
	}
	return 0
}

func (m *AgentData) GetLastHeartbeatNS() int64 {
	if m != nil {
		return m.LastHeartbeatNS
	}
	return 0
}

type HostInfo struct {
	Hostname string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
}

func (m *HostInfo) Reset()      { *m = HostInfo{} }
func (*HostInfo) ProtoMessage() {}
func (*HostInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_98f8b3e566c670dc, []int{1}
}
func (m *HostInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HostInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HostInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HostInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HostInfo.Merge(m, src)
}
func (m *HostInfo) XXX_Size() int {
	return m.Size()
}
func (m *HostInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_HostInfo.DiscardUnknown(m)
}

var xxx_messageInfo_HostInfo proto.InternalMessageInfo

func (m *HostInfo) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func init() {
	proto.RegisterType((*AgentData)(nil), "pl.vizier.services.metadata.AgentData")
	proto.RegisterType((*HostInfo)(nil), "pl.vizier.services.metadata.HostInfo")
}

func init() {
	proto.RegisterFile("src/vizier/services/metadata/datapb/data.proto", fileDescriptor_98f8b3e566c670dc)
}

var fileDescriptor_98f8b3e566c670dc = []byte{
	// 407 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xbf, 0xae, 0xd3, 0x30,
	0x14, 0xc6, 0xe3, 0x7b, 0x11, 0xb7, 0xf5, 0xad, 0x28, 0x84, 0xa5, 0x2a, 0x92, 0x5b, 0x55, 0x02,
	0x75, 0xc1, 0x96, 0x40, 0x82, 0x81, 0x01, 0x35, 0x74, 0x68, 0x24, 0xd4, 0x21, 0xa5, 0x0b, 0x4b,
	0xe4, 0x24, 0x6e, 0x62, 0x29, 0x89, 0xa3, 0xd8, 0x29, 0x88, 0x09, 0xde, 0x80, 0xc7, 0xe0, 0x51,
	0x18, 0x3b, 0x76, 0xaa, 0xa8, 0xbb, 0x30, 0xf6, 0x11, 0x50, 0x1c, 0x52, 0xc4, 0x82, 0x58, 0x7c,
	0xce, 0xe7, 0xf3, 0xfb, 0xec, 0xe3, 0x3f, 0x10, 0xcb, 0x32, 0x24, 0x5b, 0xfe, 0x89, 0xb3, 0x92,
	0x48, 0x56, 0x6e, 0x79, 0xc8, 0x24, 0xc9, 0x98, 0xa2, 0x11, 0x55, 0x94, 0xd4, 0x43, 0x11, 0x98,
	0x80, 0x8b, 0x52, 0x28, 0x61, 0x3f, 0x2a, 0x52, 0xdc, 0xe0, 0xb8, 0xc5, 0x71, 0x8b, 0x0f, 0x9f,
	0xc6, 0x5c, 0x25, 0x55, 0x80, 0x43, 0x91, 0x91, 0x58, 0xc4, 0x82, 0x18, 0x4f, 0x50, 0x6d, 0x8c,
	0x32, 0xc2, 0x64, 0xcd, 0x5a, 0xc3, 0x71, 0xbd, 0x77, 0x28, 0xb2, 0x4c, 0xe4, 0xa4, 0xaa, 0x78,
	0xd4, 0xe0, 0x26, 0x6d, 0x88, 0xc9, 0x97, 0x2b, 0xd8, 0x9d, 0xc5, 0x2c, 0x57, 0x73, 0xaa, 0xa8,
	0xfd, 0x12, 0x76, 0x68, 0x2d, 0x7c, 0x1e, 0x0d, 0xc0, 0x18, 0x4c, 0x6f, 0x9f, 0xf5, 0x71, 0x91,
	0xe2, 0x9a, 0x2f, 0x02, 0xbc, 0x5e, 0xbb, 0x73, 0xe7, 0x56, 0x1f, 0x46, 0x37, 0xc6, 0xe1, 0xce,
	0xbd, 0x1b, 0x43, 0xbb, 0x91, 0xed, 0xc0, 0x6e, 0x22, 0xa4, 0xf2, 0x79, 0xbe, 0x11, 0x83, 0x2b,
	0xe3, 0x7c, 0x8c, 0xff, 0x71, 0x10, 0xbc, 0x10, 0x52, 0xb9, 0xf9, 0x46, 0x78, 0x9d, 0xe4, 0x77,
	0x66, 0xbf, 0x80, 0xf7, 0xc2, 0x92, 0x51, 0xc5, 0x7c, 0xc5, 0x33, 0xe6, 0xe7, 0x72, 0x70, 0x3d,
	0x06, 0xd3, 0x6b, 0xe7, 0xbe, 0x3e, 0x8c, 0x7a, 0x6f, 0x4c, 0xe5, 0x1d, 0xcf, 0xd8, 0x72, 0xe5,
	0xf5, 0xc2, 0x3f, 0x4a, 0xda, 0xaf, 0xe1, 0x83, 0x94, 0x4a, 0xe5, 0x27, 0x8c, 0x96, 0x2a, 0x60,
	0x54, 0xd5, 0xd6, 0x3b, 0xc6, 0xfa, 0x50, 0x1f, 0x46, 0xfd, 0xb7, 0x54, 0xaa, 0x45, 0x5b, 0x5b,
	0xae, 0xbc, 0x7e, 0xfa, 0xd7, 0x84, 0x9c, 0x3c, 0x81, 0x9d, 0xb6, 0x1d, 0x7b, 0x08, 0x4d, 0x43,
	0x39, 0xcd, 0x98, 0xb9, 0x81, 0xae, 0x77, 0xd1, 0xce, 0x87, 0xdd, 0x11, 0x59, 0xfb, 0x23, 0xb2,
	0xce, 0x47, 0x04, 0x3e, 0x6b, 0x04, 0xbe, 0x69, 0x04, 0xbe, 0x6b, 0x04, 0x76, 0x1a, 0x81, 0x1f,
	0x1a, 0x81, 0x9f, 0x1a, 0x59, 0x67, 0x8d, 0xc0, 0xd7, 0x13, 0xb2, 0x76, 0x27, 0x64, 0xed, 0x4f,
	0xc8, 0x7a, 0x3f, 0x2b, 0xf8, 0x47, 0xce, 0x52, 0x1a, 0x48, 0x4c, 0x39, 0xb9, 0x08, 0xf2, 0x1f,
	0xdf, 0xe2, 0x55, 0x13, 0x82, 0xbb, 0xe6, 0xad, 0x9e, 0xff, 0x0a, 0x00, 0x00, 0xff, 0xff, 0xb2,
	0x2d, 0x45, 0x39, 0x4b, 0x02, 0x00, 0x00,
}

func (this *AgentData) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AgentData)
	if !ok {
		that2, ok := that.(AgentData)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.AgentID.Equal(that1.AgentID) {
		return false
	}
	if !this.HostInfo.Equal(that1.HostInfo) {
		return false
	}
	if this.CreateTimeNS != that1.CreateTimeNS {
		return false
	}
	if this.LastHeartbeatNS != that1.LastHeartbeatNS {
		return false
	}
	return true
}
func (this *HostInfo) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HostInfo)
	if !ok {
		that2, ok := that.(HostInfo)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Hostname != that1.Hostname {
		return false
	}
	return true
}
func (this *AgentData) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&datapb.AgentData{")
	if this.AgentID != nil {
		s = append(s, "AgentID: "+fmt.Sprintf("%#v", this.AgentID)+",\n")
	}
	if this.HostInfo != nil {
		s = append(s, "HostInfo: "+fmt.Sprintf("%#v", this.HostInfo)+",\n")
	}
	s = append(s, "CreateTimeNS: "+fmt.Sprintf("%#v", this.CreateTimeNS)+",\n")
	s = append(s, "LastHeartbeatNS: "+fmt.Sprintf("%#v", this.LastHeartbeatNS)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *HostInfo) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&datapb.HostInfo{")
	s = append(s, "Hostname: "+fmt.Sprintf("%#v", this.Hostname)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringData(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *AgentData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AgentData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.AgentID != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintData(dAtA, i, uint64(m.AgentID.Size()))
		n1, err := m.AgentID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.HostInfo != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintData(dAtA, i, uint64(m.HostInfo.Size()))
		n2, err := m.HostInfo.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.CreateTimeNS != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintData(dAtA, i, uint64(m.CreateTimeNS))
	}
	if m.LastHeartbeatNS != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintData(dAtA, i, uint64(m.LastHeartbeatNS))
	}
	return i, nil
}

func (m *HostInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HostInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Hostname) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintData(dAtA, i, uint64(len(m.Hostname)))
		i += copy(dAtA[i:], m.Hostname)
	}
	return i, nil
}

func encodeVarintData(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *AgentData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AgentID != nil {
		l = m.AgentID.Size()
		n += 1 + l + sovData(uint64(l))
	}
	if m.HostInfo != nil {
		l = m.HostInfo.Size()
		n += 1 + l + sovData(uint64(l))
	}
	if m.CreateTimeNS != 0 {
		n += 1 + sovData(uint64(m.CreateTimeNS))
	}
	if m.LastHeartbeatNS != 0 {
		n += 1 + sovData(uint64(m.LastHeartbeatNS))
	}
	return n
}

func (m *HostInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hostname)
	if l > 0 {
		n += 1 + l + sovData(uint64(l))
	}
	return n
}

func sovData(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozData(x uint64) (n int) {
	return sovData(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *AgentData) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AgentData{`,
		`AgentID:` + strings.Replace(fmt.Sprintf("%v", this.AgentID), "UUID", "proto1.UUID", 1) + `,`,
		`HostInfo:` + strings.Replace(fmt.Sprintf("%v", this.HostInfo), "HostInfo", "HostInfo", 1) + `,`,
		`CreateTimeNS:` + fmt.Sprintf("%v", this.CreateTimeNS) + `,`,
		`LastHeartbeatNS:` + fmt.Sprintf("%v", this.LastHeartbeatNS) + `,`,
		`}`,
	}, "")
	return s
}
func (this *HostInfo) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&HostInfo{`,
		`Hostname:` + fmt.Sprintf("%v", this.Hostname) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringData(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *AgentData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowData
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AgentData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AgentData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AgentID", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.AgentID == nil {
				m.AgentID = &proto1.UUID{}
			}
			if err := m.AgentID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HostInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.HostInfo == nil {
				m.HostInfo = &HostInfo{}
			}
			if err := m.HostInfo.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreateTimeNS", wireType)
			}
			m.CreateTimeNS = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreateTimeNS |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastHeartbeatNS", wireType)
			}
			m.LastHeartbeatNS = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastHeartbeatNS |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthData
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *HostInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowData
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HostInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HostInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hostname", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hostname = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthData
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipData(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowData
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowData
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowData
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthData
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthData
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowData
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipData(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthData
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthData = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowData   = fmt.Errorf("proto: integer overflow")
)
