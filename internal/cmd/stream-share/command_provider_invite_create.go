package streamshare

import (
	"strings"

	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/examples"
	"github.com/confluentinc/cli/internal/pkg/output"
)

func (c *command) newCreateEmailInviteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Invite a consumer with email.",
		Args:  cobra.NoArgs,
		RunE:  c.createEmailInvite,
		Example: examples.BuildExampleString(
			examples.Example{
				Text: `Invite a user with email "user@example.com":`,
				Code: "confluent stream-share provider invite create --email user@example.com --topic topic-12345 --environment env-12345 --cluster lkc-12345",
			},
		),
	}

	cmd.Flags().String("email", "", "Email of the user with whom to share the topic.")
	cmd.Flags().String("topic", "", "Topic to be shared.")
	cmd.Flags().String("schema-registry-subjects", "", "A comma separated list of Schema Registry subjects")
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddClusterFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddOutputFlag(cmd)

	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("topic")
	_ = cmd.MarkFlagRequired("environment")
	_ = cmd.MarkFlagRequired("cluster")

	return cmd
}

func (c *command) createEmailInvite(cmd *cobra.Command, _ []string) error {
	environment, err := cmd.Flags().GetString("environment")
	if err != nil {
		return err
	}

	kafkaCluster, err := cmd.Flags().GetString("cluster")
	if err != nil {
		return err
	}

	topic, err := cmd.Flags().GetString("topic")
	if err != nil {
		return err
	}

	email, err := cmd.Flags().GetString("email")
	if err != nil {
		return err
	}

	srSubjects, err := cmd.Flags().GetString("schema-registry-subjects")
	if err != nil {
		return err
	}
	subjectsList := strings.Split(srSubjects, ",")

	srCluster, err := c.Context.FetchSchemaRegistryByAccountId(cmd.Context(), c.EnvironmentId())
	if err != nil {
		return err
	}

	invite, httpResp, err := c.V2Client.CreateInvite(environment, kafkaCluster, topic, email, srCluster.Id,
		c.Config.GetLastUsedOrgId(), subjectsList)
	if err != nil {
		return errors.CatchCCloudV2Error(err, httpResp)
	}

	return output.DescribeObject(cmd, c.buildProviderShare(invite), providerShareListFields, providerHumanLabelMap, providerStructuredLabelMap)
}
