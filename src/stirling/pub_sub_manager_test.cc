#include <google/protobuf/text_format.h>
#include <google/protobuf/util/message_differencer.h>
#include <gtest/gtest.h>

#include <cstring>
#include <utility>
#include <vector>

#include "src/common/testing/testing.h"
#include "src/stirling/info_class_manager.h"
#include "src/stirling/pub_sub_manager.h"
#include "src/stirling/source_connector.h"
#include "src/stirling/stirling.h"

namespace pl {
namespace stirling {

using ::google::protobuf::TextFormat;
using ::google::protobuf::util::MessageDifferencer;
using stirlingpb::InfoClass;
using types::DataType;
using types::PatternType;
using types::SemanticType;

const char* kInfoClass0 = R"(
  name: "cpu"
  schema {
    name: "cpu"
    elements {
      name: "user_percentage"
      type: FLOAT64
      stype: ST_NONE
      ptype: METRIC_GAUGE
      desc: "User percentage"
    }
    elements {
      name: "system_percentage"
      type: FLOAT64
      stype: ST_NONE
      ptype: METRIC_GAUGE
      desc: "System percentage"
    }
    elements {
      name: "io_percentage"
      type: FLOAT64
      stype: ST_NONE
      ptype: METRIC_GAUGE
      desc: "IO percentage"
    }
    tabletized: false
    tabletization_key: 18446744073709551615
  }
  sampling_period_millis: 100
  push_period_millis: 1000
)";

const char* kInfoClass1 = R"(
  name: "my_table"
  schema {
    name: "my_table"
    elements {
      name: "a"
      type: FLOAT64
      stype: ST_NONE
      ptype: METRIC_GAUGE
      desc: ""
    }
    elements {
      name: "b"
      type: FLOAT64
      stype: ST_NONE
      ptype: METRIC_GAUGE
      desc: ""
    }
    elements {
      name: "c"
      type: FLOAT64
      stype: ST_NONE
      ptype: METRIC_GAUGE
      desc: ""
    }
    tabletized: false
    tabletization_key: 18446744073709551615
  }
  sampling_period_millis: 100
  push_period_millis: 1000
)";

// A test source connector to be used for testing.
class TestSourceConnector : public SourceConnector {
 public:
  static constexpr DataElement kElements[] = {
      {"user_percentage", "User percentage", DataType::FLOAT64, SemanticType::ST_NONE,
       PatternType::METRIC_GAUGE},
      {"system_percentage", "System percentage", DataType::FLOAT64, SemanticType::ST_NONE,
       PatternType::METRIC_GAUGE},
      {"io_percentage", "IO percentage", DataType::FLOAT64, SemanticType::ST_NONE,
       PatternType::METRIC_GAUGE}};

  static constexpr auto kTable = DataTableSchema("cpu", kElements, std::chrono::milliseconds{100},
                                                 std::chrono::milliseconds{1000});
  static constexpr auto kTables = MakeArray(kTable);

  static std::unique_ptr<SourceConnector> Create(std::string_view name) {
    return std::unique_ptr<SourceConnector>(new TestSourceConnector(name));
  }

  Status InitImpl() override { return Status::OK(); }
  Status StopImpl() override { return Status::OK(); }
  void TransferDataImpl(ConnectorContext* /* ctx */, uint32_t /* table_num */,
                        DataTable* /* data_table */) override{};

 protected:
  explicit TestSourceConnector(std::string_view name) : SourceConnector(name, kTables) {}
};

// A test source connector to be used for testing.
class TestSourceConnector2 : public SourceConnector {
 public:
  static constexpr DataElement kElements[] = {
      {"a", "", DataType::FLOAT64, SemanticType::ST_NONE, PatternType::METRIC_GAUGE},
      {"b", "", DataType::FLOAT64, SemanticType::ST_NONE, PatternType::METRIC_GAUGE},
      {"c", "", DataType::FLOAT64, SemanticType::ST_NONE, PatternType::METRIC_GAUGE}};

  static constexpr auto kTable = DataTableSchema(
      "my_table", kElements, std::chrono::milliseconds{100}, std::chrono::milliseconds{1000});
  static constexpr auto kTables = MakeArray(kTable);

  static std::unique_ptr<SourceConnector> Create(std::string_view name) {
    return std::unique_ptr<SourceConnector>(new TestSourceConnector2(name));
  }

  Status InitImpl() override { return Status::OK(); }
  Status StopImpl() override { return Status::OK(); }
  void TransferDataImpl(ConnectorContext* /* ctx */, uint32_t /* table_num */,
                        DataTable* /* data_table */) override{};

 protected:
  explicit TestSourceConnector2(std::string_view name) : SourceConnector(name, kTables) {}
};

class PubSubManagerTest : public ::testing::Test {
 protected:
  PubSubManagerTest() = default;
  void SetUp() override {
    {
      std::string name = "source0";
      std::unique_ptr<SourceConnector> source = TestSourceConnector::Create(name);
      auto info_class_mgr = std::make_unique<InfoClassManager>(TestSourceConnector::kTable);
      info_class_mgr->SetSourceConnector(source.get(), /* table_num */ 0);
      info_class_mgrs_.push_back(std::move(info_class_mgr));
      sources_.push_back(std::move(source));
    }

    {
      std::string name = "source1";
      std::unique_ptr<SourceConnector> source = TestSourceConnector2::Create(name);
      auto info_class_mgr = std::make_unique<InfoClassManager>(TestSourceConnector2::kTable);
      info_class_mgr->SetSourceConnector(source.get(), /* table_num */ 0);
      info_class_mgrs_.push_back(std::move(info_class_mgr));
      sources_.push_back(std::move(source));
    }

    pub_sub_manager_ = std::make_unique<PubSubManager>();
  }
  std::vector<std::unique_ptr<SourceConnector>> sources_;
  std::unique_ptr<PubSubManager> pub_sub_manager_;
  InfoClassManagerVec info_class_mgrs_;
};

// This test validates that the Publish proto generated by the PubSubManager
// matches the expected Publish proto message (based on kInfoClass proto
// and with some fields added in the test).
TEST_F(PubSubManagerTest, publish_test) {
  // Publish info classes using proto message.
  stirlingpb::Publish actual_publish_pb;
  pub_sub_manager_->PopulatePublishProto(&actual_publish_pb, info_class_mgrs_);

  // Set expectations for the publish message.
  stirlingpb::Publish expected_publish_pb;
  auto* info_class = expected_publish_pb.add_published_info_classes();
  ASSERT_TRUE(TextFormat::MergeFromString(kInfoClass0, info_class));
  info_class->set_id(0);

  info_class = expected_publish_pb.add_published_info_classes();
  ASSERT_TRUE(TextFormat::MergeFromString(kInfoClass1, info_class));
  info_class->set_id(1);

  EXPECT_TRUE(MessageDifferencer::Equals(actual_publish_pb, expected_publish_pb));
}

TEST_F(PubSubManagerTest, partial_publish_test) {
  // Publish info classes using proto message.
  stirlingpb::Publish actual_publish_pb;
  pub_sub_manager_->PopulatePublishProto(&actual_publish_pb, info_class_mgrs_, "cpu");

  // Set expectations for the publish message.
  stirlingpb::Publish expected_publish_pb;
  auto* info_class = expected_publish_pb.add_published_info_classes();
  ASSERT_TRUE(TextFormat::MergeFromString(kInfoClass0, info_class));

  // Copy ID from publication as the expectation.
  info_class->set_id(actual_publish_pb.published_info_classes(0).id());

  EXPECT_TRUE(MessageDifferencer::Equals(actual_publish_pb, expected_publish_pb));
}

// This test validates that the InfoClassManager objects have their subscriptions
// updated after the PubSubManager reads a subscribe message (from an agent). The
// subscribe message is created from the Publish proto message.
TEST_F(PubSubManagerTest, subscribe_test) {
  // Get publication.
  stirlingpb::Publish publish_pb;
  pub_sub_manager_->PopulatePublishProto(&publish_pb, info_class_mgrs_);

  // Send subscription.
  stirlingpb::Subscribe subscribe_pb = SubscribeToAllInfoClasses(publish_pb);
  ASSERT_OK(pub_sub_manager_->UpdateSchemaFromSubscribe(subscribe_pb, info_class_mgrs_));

  // Verify updated subscriptions.
  for (auto& info_class_mgr : info_class_mgrs_) {
    EXPECT_TRUE(info_class_mgr->subscribed());
  }
}

TEST_F(PubSubManagerTest, partial_subscribe_test) {
  // Get publication.
  stirlingpb::Publish publish_pb;
  pub_sub_manager_->PopulatePublishProto(&publish_pb, info_class_mgrs_);

  // Send subscription.
  stirlingpb::Subscribe subscribe_pb = SubscribeToInfoClass(publish_pb, "my_table");
  ASSERT_OK(pub_sub_manager_->UpdateSchemaFromSubscribe(subscribe_pb, info_class_mgrs_));

  // Verify updated subscriptions.
  ASSERT_EQ(info_class_mgrs_.size(), 2);
  EXPECT_FALSE(info_class_mgrs_[0]->subscribed());
  EXPECT_TRUE(info_class_mgrs_[1]->subscribed());
}

TEST_F(PubSubManagerTest, delta_subscribe_test) {
  // Get publication.
  stirlingpb::Publish publish_pb;
  pub_sub_manager_->PopulatePublishProto(&publish_pb, info_class_mgrs_);

  // Split the publication into subsciption pieces (one per info class).
  std::vector<stirlingpb::Subscribe> subs;
  for (const auto& p : publish_pb.published_info_classes()) {
    stirlingpb::Publish partial_pub;
    auto* info_class = partial_pub.add_published_info_classes();
    info_class->CopyFrom(p);

    subs.push_back(SubscribeToAllInfoClasses(partial_pub));
  }
  ASSERT_EQ(subs.size(), 2);

  // Perform first delta subscription.
  ASSERT_OK(pub_sub_manager_->UpdateSchemaFromSubscribe(subs[1], info_class_mgrs_));

  // Verify updated subscriptions.
  ASSERT_EQ(info_class_mgrs_.size(), 2);
  EXPECT_FALSE(info_class_mgrs_[0]->subscribed());
  EXPECT_TRUE(info_class_mgrs_[1]->subscribed());

  // Perform second delta subscription.
  ASSERT_OK(pub_sub_manager_->UpdateSchemaFromSubscribe(subs[0], info_class_mgrs_));

  // Verify updated subscriptions.
  ASSERT_EQ(info_class_mgrs_.size(), 2);
  EXPECT_TRUE(info_class_mgrs_[0]->subscribed());
  EXPECT_TRUE(info_class_mgrs_[1]->subscribed());
}

}  // namespace stirling
}  // namespace pl
