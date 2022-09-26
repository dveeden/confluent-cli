package pipeline

import (
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	v1 "github.com/confluentinc/cli/internal/pkg/config/v1"
	dynamicconfig "github.com/confluentinc/cli/internal/pkg/dynamic-config"
	launchdarkly "github.com/confluentinc/cli/internal/pkg/featureflags"
)

type Pipeline struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

type command struct {
	*pcmd.AuthenticatedCLICommand
	prerunner pcmd.PreRunner
}

func New(cfg *v1.Config, prerunner pcmd.PreRunner) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "pipeline",
		Short:       "Manage stream designer pipelines.",
		Annotations: map[string]string{pcmd.RunRequirement: pcmd.RequireNonAPIKeyCloudLogin},
	}

	c := &command{
		AuthenticatedCLICommand: pcmd.NewAuthenticatedCLICommand(cmd, prerunner),
		prerunner:               prerunner,
	}

	c.AddCommand(c.newActivateCommand(prerunner))
	c.AddCommand(c.newCreateCommand(prerunner))
	c.AddCommand(c.newDeactivateCommand(prerunner))
	c.AddCommand(c.newDeleteCommand(prerunner))
	c.AddCommand(c.newDescribeCommand(prerunner))
	c.AddCommand(c.newListCommand(prerunner))
	c.AddCommand(c.newUpdateCommand(prerunner))

	dc := dynamicconfig.New(cfg, nil, nil)
	_ = dc.ParseFlagsIntoConfig(cmd)

	c.Hidden = !cfg.IsTest && !launchdarkly.Manager.BoolVariation("cli.client_stream_designer.enable", dc.Context(), v1.CliLaunchDarklyClient, true, false)

	return c.Command
}
