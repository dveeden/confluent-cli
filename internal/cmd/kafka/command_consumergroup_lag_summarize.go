package kafka

import (
	"github.com/spf13/cobra"

	cloudkafkarestv3 "github.com/confluentinc/ccloud-sdk-go-v2/kafkarest/v3"
	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/examples"
	"github.com/confluentinc/cli/internal/pkg/output"
)

var (
	lagSummaryFields       = []string{"ClusterId", "ConsumerGroupId", "TotalLag", "MaxLag", "MaxLagConsumerId", "MaxLagInstanceId", "MaxLagClientId", "MaxLagTopicName", "MaxLagPartitionId"}
	lagSummaryHumanRenames = map[string]string{
		"ClusterId":         "Cluster",
		"ConsumerGroupId":   "ConsumerGroup",
		"MaxLagConsumerId":  "MaxLagConsumer",
		"MaxLagInstanceId":  "MaxLagInstance",
		"MaxLagClientId":    "MaxLagClient",
		"MaxLagTopicName":   "MaxLagTopic",
		"MaxLagPartitionId": "MaxLagPartition",
	}
	lagSummaryStructuredRenames = map[string]string{
		"ClusterId":         "cluster",
		"ConsumerGroupId":   "consumer_group",
		"TotalLag":          "total_lag",
		"MaxLag":            "max_lag",
		"MaxLagConsumerId":  "max_lag_consumer",
		"MaxLagInstanceId":  "max_lag_instance",
		"MaxLagClientId":    "max_lag_client",
		"MaxLagTopicName":   "max_lag_topic",
		"MaxLagPartitionId": "max_lag_partition",
	}
)

func (c *lagCommand) newSummarizeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "summarize <consumer-group>",
		Short:             "Summarize consumer lag for a Kafka consumer group.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: pcmd.NewValidArgsFunction(c.validArgs),
		RunE:              c.summarize,
		Example: examples.BuildExampleString(
			examples.Example{
				Text: "Summarize the lag for the `my-consumer-group` consumer-group.",
				Code: "confluent kafka consumer-group lag summarize my-consumer-group",
			},
		),
		Hidden: true,
	}

	pcmd.AddClusterFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddContextFlag(cmd, c.CLICommand)
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddOutputFlag(cmd)

	return cmd
}

func (c *lagCommand) summarize(cmd *cobra.Command, args []string) error {
	consumerGroupId := args[0]

	kafkaREST, lkc, err := getKafkaRestProxyAndLkcId(c.AuthenticatedStateFlagCommand)
	if err != nil {
		return err
	}

	lagSummaryResp, httpResp, err := kafkaREST.Client.ConsumerGroupV3Api.GetKafkaConsumerGroupLagSummary(kafkaREST.Context, lkc, consumerGroupId).Execute()
	if err != nil {
		return kafkaRestError(kafkaREST.Client.GetConfig().Servers[0].URL, err, httpResp)
	}

	return output.DescribeObject(cmd, convertLagSummaryToStruct(lagSummaryResp), lagSummaryFields, lagSummaryHumanRenames, lagSummaryStructuredRenames)
}

func convertLagSummaryToStruct(lagSummaryData cloudkafkarestv3.ConsumerGroupLagSummaryData) *lagSummaryStruct {
	maxLagInstanceId := ""
	if lagSummaryData.MaxLagInstanceId.IsSet() {
		maxLagInstanceId = *lagSummaryData.MaxLagInstanceId.Get()
	}

	return &lagSummaryStruct{
		ClusterId:         lagSummaryData.ClusterId,
		ConsumerGroupId:   lagSummaryData.ConsumerGroupId,
		TotalLag:          lagSummaryData.TotalLag,
		MaxLag:            lagSummaryData.MaxLag,
		MaxLagConsumerId:  lagSummaryData.MaxLagConsumerId,
		MaxLagInstanceId:  maxLagInstanceId,
		MaxLagClientId:    lagSummaryData.MaxLagClientId,
		MaxLagTopicName:   lagSummaryData.MaxLagTopicName,
		MaxLagPartitionId: lagSummaryData.MaxLagPartitionId,
	}
}
