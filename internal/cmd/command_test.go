package cmd

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	orgv1 "github.com/confluentinc/cc-structs/kafka/org/v1"
	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/config"
	v1 "github.com/confluentinc/cli/internal/pkg/config/v1"
	"github.com/confluentinc/cli/internal/pkg/utils"
	pversion "github.com/confluentinc/cli/internal/pkg/version"
	testserver "github.com/confluentinc/cli/test/test-server"
)

var (
	mockBaseConfig = &config.BaseConfig{}
	mockVersion    = new(pversion.Version)

	regularOrgContextState = &v1.ContextState{
		Auth: &v1.AuthConfig{
			Organization: testserver.RegularOrg,
		},
	}
	suspendedOrgContextState = func(eventType orgv1.SuspensionEventType) *v1.ContextState {
		return &v1.ContextState{
			Auth: &v1.AuthConfig{
				Organization: testserver.SuspendedOrg(eventType),
			},
		}
	}
)

func TestHelp_NoContext(t *testing.T) {
	cfg := &v1.Config{BaseConfig: mockBaseConfig}

	out, err := runWithConfig(cfg)
	require.NoError(t, err)

	commands := []string{
		"cloud-signup", "completion", "context", "help", "kafka", "local", "login", "logout", "update", "version",
	}
	if runtime.GOOS == "windows" {
		commands = utils.Remove(commands, "local")
	}

	for _, command := range commands {
		require.Contains(t, out, command)
	}
}

func TestHelp_CloudSuspendedOrg(t *testing.T) {
	cfg := &v1.Config{
		BaseConfig: mockBaseConfig,
		Contexts: map[string]*v1.Context{"cloud": {
			PlatformName: "confluent.cloud",
			State:        suspendedOrgContextState(orgv1.SuspensionEventType_SUSPENSION_EVENT_CUSTOMER_INITIATED_ORG_DEACTIVATION),
		}},
		CurrentContext: "cloud",
	}

	out, err := runWithConfig(cfg)
	require.NoError(t, err)

	commands := []string{
		"cloud-signup", "completion", "context", "help", "kafka", "local", "login", "logout", "prompt", "shell", "update", "version",
	}
	if runtime.GOOS == "windows" {
		commands = utils.Remove(commands, "local")
	}

	for _, command := range commands {
		require.Contains(t, out, command)
	}
}

func TestHelp_CloudEndOfFreeTrialSuspendedOrg(t *testing.T) {
	cfg := &v1.Config{
		BaseConfig: mockBaseConfig,
		Contexts: map[string]*v1.Context{"cloud": {
			PlatformName: "confluent.cloud",
			State:        suspendedOrgContextState(orgv1.SuspensionEventType_SUSPENSION_EVENT_END_OF_FREE_TRIAL),
		}},
		CurrentContext: "cloud",
	}

	cli := NewConfluentCommand(cfg, true, mockVersion)
	out, err := pcmd.ExecuteCommand(cli.Command, "help")
	require.NoError(t, err)

	// note users can still run "confluent admin payment update" or "confluent admin promo add" if the org is suspended
	// but only due to end of free trial
	commands := []string{
		"admin", "cloud-signup", "completion", "context", "help", "kafka", "local", "login", "logout", "prompt", "shell", "update", "version",
	}
	if runtime.GOOS == "windows" {
		commands = utils.Remove(commands, "local")
	}

	for _, command := range commands {
		require.Contains(t, out, command)
	}

	// check that some top level cloud commands are not included
	cloudCommands := []string{"api-key", "iam", "schema-registry", "price"}
	for _, command := range cloudCommands {
		require.NotContains(t, out, command)
	}

	// check "admin payment" command
	out, err = pcmd.ExecuteCommand(cli.Command, "admin", "payment", "--help")
	require.NoError(t, err)
	require.Contains(t, out, "update")
	require.NotContains(t, out, "describe")

	// check "admin promo" command
	out, err = pcmd.ExecuteCommand(cli.Command, "admin", "promo", "--help")
	require.NoError(t, err)
	require.Contains(t, out, "add")
	require.NotContains(t, out, "list")
}

func TestHelp_Cloud(t *testing.T) {
	cfg := &v1.Config{
		BaseConfig: mockBaseConfig,
		Contexts: map[string]*v1.Context{"cloud": {
			PlatformName: "confluent.cloud",
			State:        regularOrgContextState,
		}},
		CurrentContext: "cloud",
	}

	out, err := runWithConfig(cfg)
	require.NoError(t, err)

	commands := []string{
		"admin", "api-key", "audit-log", "cloud-signup", "completion", "context", "connect", "environment", "help",
		"iam", "kafka", "ksql", "login", "logout", "price", "prompt", "schema-registry", "shell", "update", "version",
	}

	for _, command := range commands {
		require.Contains(t, out, command)
	}
}

func TestHelp_CloudWithAPIKey(t *testing.T) {
	cfg := &v1.Config{
		BaseConfig: mockBaseConfig,
		Contexts: map[string]*v1.Context{
			"cloud-with-api-key": {
				PlatformName: "confluent.cloud",
				Credential:   &v1.Credential{CredentialType: v1.APIKey},
				State:        regularOrgContextState,
			},
		},
		CurrentContext: "cloud-with-api-key",
	}

	out, err := runWithConfig(cfg)
	require.NoError(t, err)

	commands := []string{
		"admin", "audit-log", "cloud-signup", "completion", "context", "help", "kafka", "login", "logout", "update",
		"version",
	}

	for _, command := range commands {
		require.Contains(t, out, command)
	}
}

func TestHelp_OnPrem(t *testing.T) {
	cfg := &v1.Config{
		BaseConfig:     mockBaseConfig,
		Contexts:       map[string]*v1.Context{"on-prem": {PlatformName: "https://example.com"}},
		CurrentContext: "on-prem",
	}

	out, err := runWithConfig(cfg)
	require.NoError(t, err)

	commands := []string{
		"audit-log", "cloud-signup", "cluster", "completion", "context", "connect", "help", "iam", "kafka", "ksql",
		"local", "login", "logout", "schema-registry", "secret", "update", "version",
	}
	if runtime.GOOS == "windows" {
		commands = utils.Remove(commands, "local")
	}

	for _, command := range commands {
		require.Contains(t, out, command)
	}
}

func runWithConfig(cfg *v1.Config) (string, error) {
	cli := NewConfluentCommand(cfg, true, mockVersion)
	return pcmd.ExecuteCommand(cli.Command, "help")
}
