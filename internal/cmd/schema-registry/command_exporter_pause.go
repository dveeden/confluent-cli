package schemaregistry

import (
	"context"

	srsdk "github.com/confluentinc/schema-registry-sdk-go"
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/utils"
)

func (c *exporterCommand) newPauseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause <name>",
		Short: "Pause schema exporter.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.pause,
	}

	pcmd.AddApiKeyFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddApiSecretFlag(cmd)
	pcmd.AddContextFlag(cmd, c.CLICommand)
	pcmd.AddEnvironmentFlag(cmd, c.AuthenticatedCLICommand)
	pcmd.AddOutputFlag(cmd)

	return cmd
}

func (c *exporterCommand) pause(cmd *cobra.Command, args []string) error {
	srClient, ctx, err := getApiClient(cmd, c.srClient, c.Config, c.Version)
	if err != nil {
		return err
	}

	return pauseExporter(cmd, args[0], srClient, ctx)
}

func pauseExporter(cmd *cobra.Command, name string, srClient *srsdk.APIClient, ctx context.Context) error {
	if _, _, err := srClient.DefaultApi.PauseExporter(ctx, name); err != nil {
		return err
	}

	utils.Printf(cmd, errors.ExporterActionMsg, "Paused", name)
	return nil
}
