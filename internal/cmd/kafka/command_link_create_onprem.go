package kafka

import (
	"fmt"

	"github.com/antihax/optional"
	"github.com/confluentinc/kafka-rest-sdk-go/kafkarestv3"
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/examples"
	"github.com/confluentinc/cli/internal/pkg/resource"
	"github.com/confluentinc/cli/internal/pkg/utils"
)

func (c *linkCommand) newCreateCommandOnPrem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <link>",
		Short: "Create a new cluster link.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.createOnPrem,
		Example: examples.BuildExampleString(
			examples.Example{
				Text: "Create a cluster link, using a configuration file.",
				Code: "confluent kafka link create my-link --destination-cluster 123456789 --config-file config.txt",
			},
			examples.Example{
				Text: "Create a cluster link using command line flags.",
				Code: "confluent kafka link create my-link --destination-cluster 123456789 --destination-bootstrap-server my-host:1234 --source-api-key my-key --source-api-secret my-secret",
			},
		),
	}

	cmd.Flags().String(destinationClusterIdFlagName, "", "Destination cluster ID.")
	cmd.Flags().String(destinationBootstrapServerFlagName, "", "Bootstrap server address of the destination cluster. Can alternatively be set in the config file using key bootstrap.servers.")
	cmd.Flags().String(sourceApiKeyFlagName, "", "An API key for the source cluster. For links at destination cluster this is used for remote cluster authentication. For links at source cluster this is used for local cluster authentication. "+authHelperMsg)
	cmd.Flags().String(sourceApiSecretFlagName, "", "An API secret for the source cluster. For links at destination cluster this is used for remote cluster authentication. For links at source cluster this is used for local cluster authentication. "+authHelperMsg)
	cmd.Flags().String(destinationApiKeyFlagName, "", "An API key for the destination cluster. This is used for remote cluster authentication links at the source cluster. "+authHelperMsg)
	cmd.Flags().String(destinationApiSecretFlagName, "", "An API secret for the destination cluster. This is used for remote cluster authentication for links at the source cluster. "+authHelperMsg)
	cmd.Flags().String(configFileFlagName, "", "Name of the file containing link configuration. Each property key-value pair should have the format of key=value. Properties are separated by new-line characters.")
	cmd.Flags().Bool(dryrunFlagName, false, "Validate a link, but do not create it.")
	cmd.Flags().Bool(noValidateFlagName, false, "Create a link even if the source cluster cannot be reached.")
	cmd.Flags().AddFlagSet(pcmd.OnPremKafkaRestSet())
	pcmd.AddContextFlag(cmd, c.CLICommand)

	_ = cmd.MarkFlagRequired(destinationClusterIdFlagName)

	return cmd
}

func (c *linkCommand) createOnPrem(cmd *cobra.Command, args []string) error {
	linkName := args[0]

	configFile, err := cmd.Flags().GetString(configFileFlagName)
	if err != nil {
		return err
	}

	dryRun, err := cmd.Flags().GetBool(dryrunFlagName)
	if err != nil {
		return err
	}

	noValidate, err := cmd.Flags().GetBool(noValidateFlagName)
	if err != nil {
		return err
	}

	configMap, linkMode, err := c.getConfigMapAndLinkMode(configFile)
	if err != nil {
		return err
	}

	// For on prem deployments, we only support source initiated links.
	if linkMode != Source {
		return errors.New("Confluent Platform only supports source initiated links.")
	}

	if err := c.addSecurityConfigToMap(cmd, linkMode, configMap); err != nil {
		return err
	}

	remoteClusterId, bootstrapServer, err := c.getRemoteClusterMetadata(cmd, linkMode)
	if err != nil {
		return err
	}

	if bootstrapServer != "" {
		configMap[bootstrapServersPropertyName] = bootstrapServer
	}

	data := kafkarestv3.CreateLinkRequestData{Configs: toCreateTopicConfigsOnPrem(configMap)}
	if linkMode == Destination {
		if remoteClusterId != "" {
			data.SourceClusterId = remoteClusterId
		}
	} else {
		if remoteClusterId != "" {
			data.DestinationClusterId = remoteClusterId
		}
	}

	client, ctx, err := initKafkaRest(c.AuthenticatedCLICommand, cmd)
	if err != nil {
		return err
	}

	clusterId, err := getClusterIdForRestRequests(client, ctx)
	if err != nil {
		return err
	}

	opts := &kafkarestv3.CreateKafkaLinkOpts{
		ValidateOnly:          optional.NewBool(dryRun),
		ValidateLink:          optional.NewBool(!noValidate),
		CreateLinkRequestData: optional.NewInterface(data),
	}

	if httpResp, err := client.ClusterLinkingV3Api.CreateKafkaLink(ctx, clusterId, linkName, opts); err != nil {
		return handleOpenApiError(httpResp, err, client)
	}

	msg := fmt.Sprintf(errors.CreatedResourceMsg, resource.ClusterLink, linkName)
	if dryRun {
		msg = "[DRY RUN]: " + msg
	}
	utils.Print(cmd, msg)

	return nil
}

func getListFieldsOnPrem(includeTopics bool) []string {
	x := []string{"Name"}

	if includeTopics {
		x = append(x, "TopicName")
	}

	return append(x, "DestinationClusterId")
}
