package kafka

import (
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	v1 "github.com/confluentinc/cli/internal/pkg/config/v1"
	dynamicconfig "github.com/confluentinc/cli/internal/pkg/dynamic-config"
	"github.com/confluentinc/cli/internal/pkg/errors"
	launchdarkly "github.com/confluentinc/cli/internal/pkg/featureflags"
	"github.com/confluentinc/cli/internal/pkg/output"

	kafkaquotas "github.com/confluentinc/ccloud-sdk-go-v2-internal/kafka-quotas/v1"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
)

type quotaCommand struct {
	*pcmd.AuthenticatedStateFlagCommand
}

func newQuotaCommand(config *v1.Config, prerunner pcmd.PreRunner) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "quota",
		Short:       "Manage Kafka client quotas.",
		Annotations: map[string]string{pcmd.RunRequirement: pcmd.RequireNonAPIKeyCloudLogin},
	}

	c := &quotaCommand{pcmd.NewAuthenticatedStateFlagCommand(cmd, prerunner)}

	dc := dynamicconfig.New(config, nil, nil)
	_ = dc.ParseFlagsIntoConfig(cmd)

	clientQuotasEnable := launchdarkly.Manager.BoolVariation("cli.client_quotas.enable", dc.Context(), v1.CliLaunchDarklyClient, true, false)
	c.Hidden = !clientQuotasEnable

	c.AddCommand(c.newCreateCommand())
	c.AddCommand(c.newDeleteCommand())
	c.AddCommand(c.newDescribeCommand())
	c.AddCommand(c.newListCommand())
	c.AddCommand(c.newUpdateCommand())

	return c.Command
}

func quotaErr(err error) error {
	if openAPIError, ok := err.(kafkaquotas.GenericOpenAPIError); ok {
		test, ok := openAPIError.Model().(kafkaquotas.Failure)
		if !ok {
			return err
		}
		var formattedErr error
		for _, e := range test.Errors {
			formattedErr = multierror.Append(formattedErr, errors.Errorf("%s: %s", err.Error(), e.GetDetail()))
		}
		return formattedErr
	}
	return err
}

func quotaToPrintable(quota kafkaquotas.KafkaQuotasV1ClientQuota, format string) interface{} {
	s := struct {
		Id          string
		DisplayName string
		Description string
		Ingress     string
		Egress      string
		Principals  string
		Cluster     string
		Environment string
	}{
		Id:          *quota.Id,
		DisplayName: *quota.DisplayName,
		Description: *quota.Description,
		Ingress:     *quota.Throughput.IngressByteRate,
		Egress:      *quota.Throughput.EgressByteRate,
		Principals:  principalsToString(*quota.Principals),
		Cluster:     quota.Cluster.Id,
		Environment: quota.Environment.Id,
	}
	if format == output.Human.String() {
		s.Ingress = s.Ingress + " B/s"
		s.Egress = s.Egress + " B/s"
	}
	return &s
}

func principalsToString(principals []kafkaquotas.ObjectReference) string {
	principalStr := ""
	for i, principal := range principals {
		principalStr += principal.Id
		if i < len(principals)-1 {
			principalStr += ","
		}
	}
	return principalStr
}
