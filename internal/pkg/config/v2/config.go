package v2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/hashicorp/go-version"

	"github.com/confluentinc/cli/internal/pkg/config"
	v1 "github.com/confluentinc/cli/internal/pkg/config/v1"
	"github.com/confluentinc/cli/internal/pkg/errors"
)

const (
	defaultConfigFileFmt = "%s/.%s/config.json"
)

var Version, _ = version.NewVersion("2.0.0")

// Config represents the CLI configuration.
type Config struct {
	*config.BaseConfig
	DisableUpdateCheck bool                     `json:"disable_update_check"`
	DisableUpdates     bool                     `json:"disable_updates"`
	NoBrowser          bool                     `json:"no_browser" hcl:"no_browser"`
	Platforms          map[string]*Platform     `json:"platforms,omitempty"`
	Credentials        map[string]*Credential   `json:"credentials,omitempty"`
	Contexts           map[string]*Context      `json:"contexts,omitempty"`
	ContextStates      map[string]*ContextState `json:"context_states,omitempty"`
	CurrentContext     string                   `json:"current_context"`
	AnonymousId        string                   `json:"anonymous_id,omitempty"`
}

// NewBaseConfig initializes a new Config object
func New(params *config.Params) *Config {
	c := &Config{}
	baseCfg := config.NewBaseConfig(params, Version)
	c.BaseConfig = baseCfg
	if c.CLIName == "" {
		// HACK: this is a workaround while we're building multiple binaries off one codebase
		c.CLIName = "confluent"
	}
	c.Platforms = map[string]*Platform{}
	c.Credentials = map[string]*Credential{}
	c.Contexts = map[string]*Context{}
	c.ContextStates = map[string]*ContextState{}
	c.AnonymousId = uuid.New().String()
	return c
}

// Load reads the CLI config from disk.
// Save a default version if none exists yet.
func (c *Config) Load() error {
	filename, err := c.getFilename()
	if err != nil {
		return err
	}
	c.Filename = filename
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// Save a default version if none exists yet.
			if err := c.Save(); err != nil {
				return errors.Wrapf(err, errors.UnableToCreateConfigErrorMsg)
			}
			return nil
		}
		return errors.Wrapf(err, errors.UnableToReadConfigErrorMsg, filename)
	}
	err = json.Unmarshal(input, c)
	if err != nil {
		return errors.Wrapf(err, errors.ParseConfigErrorMsg, filename)
	}
	for _, context := range c.Contexts {
		// Some "pre-validation"
		if context.Name == "" {
			return errors.NewCorruptedConfigError(errors.NoNameContextErrorMsg, "", c.CLIName, c.Filename, c.Logger)
		}
		if context.CredentialName == "" {
			return errors.NewCorruptedConfigError(errors.UnspecifiedCredentialErrorMsg, context.Name, c.CLIName, c.Filename, c.Logger)
		}
		if context.PlatformName == "" {
			return errors.NewCorruptedConfigError(errors.UnspecifiedPlatformErrorMsg, context.Name, c.CLIName, c.Filename, c.Logger)
		}
		context.State = c.ContextStates[context.Name]
		context.Credential = c.Credentials[context.CredentialName]
		context.Platform = c.Platforms[context.PlatformName]
		context.Logger = c.Logger
		context.Config = c
	}
	err = c.Validate()
	if err != nil {
		return err
	}
	return nil
}

// Save writes the CLI config to disk.
func (c *Config) Save() error {
	err := c.Validate()
	if err != nil {
		return err
	}
	cfg, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return errors.Wrapf(err, errors.MarshalConfigErrorMsg)
	}
	filename, err := c.getFilename()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(filename), 0700)
	if err != nil {
		return errors.Wrapf(err, errors.CreateConfigDirectoryErrorMsg, filename)
	}
	err = ioutil.WriteFile(filename, cfg, 0600)
	if err != nil {
		return errors.Wrapf(err, errors.CreateConfigFileErrorMsg, filename)
	}
	return nil
}

func (c *Config) Validate() error {
	// Validate that current context exists.
	if c.CurrentContext != "" {
		if _, ok := c.Contexts[c.CurrentContext]; !ok {
			c.Logger.Trace("current context does not exist")
			return errors.NewCorruptedConfigError(errors.CurrentContextNotExistErrorMsg, c.CurrentContext, c.CLIName, c.Filename, c.Logger)
		}
	}
	// Validate that every context:
	// 1. Has no hanging references between the context and the config.
	// 2. Is mapped by name correctly in the config.
	for _, context := range c.Contexts {
		err := context.validate()
		if err != nil {
			c.Logger.Trace("context validation error")
			return err
		}
		if _, ok := c.Credentials[context.CredentialName]; !ok {
			c.Logger.Trace("unspecified credential error")
			return errors.NewCorruptedConfigError(errors.UnspecifiedCredentialErrorMsg, context.Name, c.CLIName, c.Filename, c.Logger)
		}
		if _, ok := c.Platforms[context.PlatformName]; !ok {
			c.Logger.Trace("unspecified platform error")
			return errors.NewCorruptedConfigError(errors.UnspecifiedPlatformErrorMsg, context.Name, c.CLIName, c.Filename, c.Logger)
		}
		if _, ok := c.ContextStates[context.Name]; !ok {
			c.ContextStates[context.Name] = new(ContextState)
		}
		if *c.ContextStates[context.Name] != *context.State {
			return errors.NewCorruptedConfigError(errors.ContextStateMismatchErrorMsg, context.Name, c.CLIName, c.Filename, c.Logger)
		}
	}
	// Validate that all context states are mapped to an existing context.
	for contextName := range c.ContextStates {
		if _, ok := c.Contexts[contextName]; !ok {
			c.Logger.Trace("context state mapped to nonexistent context")
			return errors.NewCorruptedConfigError(errors.ContextStateNotMappedErrorMsg, contextName, c.CLIName, c.Filename, c.Logger)
		}
	}

	return nil
}

// DeleteContext deletes the specified context, and returns an error if it's not found.
func (c *Config) DeleteContext(name string) error {
	_, err := c.FindContext(name)
	if err != nil {
		return err
	}
	delete(c.Contexts, name)
	if c.CurrentContext == name {
		c.CurrentContext = ""
	}
	delete(c.ContextStates, name)
	return nil
}

// FindContext finds a context by name, and returns nil if not found.
func (c *Config) FindContext(name string) (*Context, error) {
	context, ok := c.Contexts[name]
	if !ok {
		return nil, fmt.Errorf("context \"%s\" does not exist", name)
	}
	return context, nil
}

func (c *Config) AddContext(name string, platformName string, credentialName string,
	kafkaClusters map[string]*v1.KafkaClusterConfig, kafka string,
	schemaRegistryClusters map[string]*SchemaRegistryCluster, state *ContextState) error {
	if _, ok := c.Contexts[name]; ok {
		return fmt.Errorf("context \"%s\" already exists", name)
	}
	credential, ok := c.Credentials[credentialName]
	if !ok {
		return fmt.Errorf("credential \"%s\" not found", credentialName)
	}
	platform, ok := c.Platforms[platformName]
	if !ok {
		return fmt.Errorf("platform \"%s\" not found", platformName)
	}
	context, err := newContext(name, platform, credential, kafkaClusters, kafka,
		schemaRegistryClusters, state, c)
	if err != nil {
		return err
	}
	c.Contexts[name] = context
	c.ContextStates[name] = context.State
	err = c.Validate()
	if err != nil {
		return err
	}
	if c.CurrentContext == "" {
		c.CurrentContext = context.Name
	}
	return c.Save()
}

func (c *Config) SetContext(name string) error {
	_, err := c.FindContext(name)
	if err != nil {
		return err
	}
	c.CurrentContext = name
	return c.Save()
}

// Name returns the display name for the CLI
func (c *Config) Name() string {
	name := "Confluent CLI"
	if c.CLIName == "ccloud" {
		name = "Confluent Cloud CLI"
	}
	return name
}

func (c *Config) SaveCredential(credential *Credential) error {
	if credential.Name == "" {
		return fmt.Errorf("credential must have a name")
	}
	c.Credentials[credential.Name] = credential
	return c.Save()
}

func (c *Config) SavePlatform(platform *Platform) error {
	if platform.Name == "" {
		return fmt.Errorf("platform must have a name")
	}
	c.Platforms[platform.Name] = platform
	return c.Save()
}

func (c *Config) Support() string {
	support := "https://confluent.io; support@confluent.io"
	if c.CLIName == "ccloud" {
		support = "https://confluent.cloud; support@confluent.io"
	}
	return support
}

// APIName returns the display name of the remote API
// (e.g., Confluent Platform or Confluent Cloud)
func (c *Config) APIName() string {
	name := "Confluent Platform"
	if c.CLIName == "ccloud" {
		name = "Confluent Cloud"
	}
	return name
}

// Context returns the user specified context if it exists,
// the current Context, or nil if there's no context set.
func (c *Config) Context() *Context {
	return c.Contexts[c.CurrentContext]
}

func (c *Config) HasLogin() bool {
	ctx := c.Context()
	if ctx == nil {
		return false
	}
	return ctx.hasLogin()
}

func (c *Config) ResetAnonymousId() error {
	c.AnonymousId = uuid.New().String()
	return c.Save()
}

func (c *Config) getFilename() (string, error) {
	if c.Filename == "" {
		homedir, _ := os.UserHomeDir()
		c.Filename = filepath.FromSlash(fmt.Sprintf(defaultConfigFileFmt, homedir, c.CLIName))
	}
	return c.Filename, nil
}
