package login

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/confluentinc/ccloud-sdk-go-v1"
	"github.com/confluentinc/cli/internal/cmd/admin"
	"github.com/spf13/cobra"

	"github.com/confluentinc/cli/internal/pkg/ccloudv2"

	orgv1 "github.com/confluentinc/cc-structs/kafka/org/v1"
	pauth "github.com/confluentinc/cli/internal/pkg/auth"
	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	v1 "github.com/confluentinc/cli/internal/pkg/config/v1"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/log"
	"github.com/confluentinc/cli/internal/pkg/netrc"
	"github.com/confluentinc/cli/internal/pkg/utils"
)

type command struct {
	*pcmd.CLICommand
	cfg                      *v1.Config
	ccloudClientFactory      pauth.CCloudClientFactory
	mdsClientManager         pauth.MDSClientManager
	netrcHandler             netrc.NetrcHandler
	loginCredentialsManager  pauth.LoginCredentialsManager
	loginOrganizationManager pauth.LoginOrganizationManager
	authTokenHandler         pauth.AuthTokenHandler
	isTest                   bool
}

func New(cfg *v1.Config, prerunner pcmd.PreRunner, ccloudClientFactory pauth.CCloudClientFactory, mdsClientManager pauth.MDSClientManager, netrcHandler netrc.NetrcHandler, loginCredentialsManager pauth.LoginCredentialsManager, loginOrganizationManager pauth.LoginOrganizationManager, authTokenHandler pauth.AuthTokenHandler, isTest bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to Confluent Cloud or Confluent Platform.",
		Long: fmt.Sprintf("Log in to Confluent Cloud using your email and password, or non-interactively using the `%s` and `%s` environment variables.\n\n", pauth.ConfluentCloudEmail, pauth.ConfluentCloudPassword) +
			fmt.Sprintf("Log in to a specific Confluent Cloud organization using the `--organization-id` flag, or by setting the environment variable `%s`.\n\n", pauth.ConfluentCloudOrganizationId) +
			fmt.Sprintf("Log in to Confluent Platform with your username and password, or non-interactively using `%s`, `%s`, `%s`, and `%s`. ", pauth.ConfluentPlatformUsername, pauth.ConfluentPlatformPassword, pauth.ConfluentPlatformMDSURL, pauth.ConfluentPlatformCACertPath) +
			fmt.Sprintf("In a non-interactive login, `%s` replaces the `--url` flag, and `%s` replaces the `--ca-cert-path` flag.\n\n", pauth.ConfluentPlatformMDSURL, pauth.ConfluentPlatformCACertPath) +
			"Even with the environment variables set, you can force an interactive login using the `--prompt` flag.",
		Args: cobra.NoArgs,
	}

	cmd.Flags().String("url", "", "Metadata Service (MDS) URL, for on-prem deployments.")
	cmd.Flags().String("ca-cert-path", "", "Self-signed certificate chain in PEM format, for on-prem deployments.")
	cmd.Flags().Bool("no-browser", false, "Do not open a browser window when authenticating via Single Sign-On (SSO).")
	cmd.Flags().String("organization-id", "", "The Confluent Cloud organization to log in to. If empty, log in to the default organization.")
	cmd.Flags().Bool("prompt", false, "Bypass non-interactive login and prompt for login credentials.")
	cmd.Flags().Bool("save", false, "Save username/password (non-SSO) credentials to the .netrc file in your $HOME directory. Use the flag to get automatically logged back in when your token expires, after one hour for Confluent Cloud or after six hours for Confluent Platform.")

	c := &command{
		CLICommand:               pcmd.NewAnonymousCLICommand(cmd, prerunner),
		cfg:                      cfg,
		mdsClientManager:         mdsClientManager,
		ccloudClientFactory:      ccloudClientFactory,
		netrcHandler:             netrcHandler,
		loginCredentialsManager:  loginCredentialsManager,
		loginOrganizationManager: loginOrganizationManager,
		authTokenHandler:         authTokenHandler,
		isTest:                   isTest,
	}
	cmd.RunE = c.login

	return cmd
}

func (c *command) login(cmd *cobra.Command, _ []string) error {
	url, err := c.getURL(cmd)
	if err != nil {
		return err
	}

	isCCloud := ccloudv2.IsCCloudURL(url, c.isTest)

	url, warningMsg, err := validateURL(url, isCCloud)
	if err != nil {
		return err
	}
	if warningMsg != "" {
		utils.ErrPrintf(cmd, errors.UsingLoginURLDefaults, warningMsg)
	}

	if isCCloud {
		return c.loginCCloud(cmd, url)
	} else {
		return c.loginMDS(cmd, url)
	}
}

func (c *command) loginCCloud(cmd *cobra.Command, url string) error {
	orgResourceId, err := c.getOrgResourceId(cmd)
	if err != nil {
		return err
	}

	credentials, err := c.getCCloudCredentials(cmd, url, orgResourceId)
	if err != nil {
		return err
	}

	noBrowser, err := cmd.Flags().GetBool("no-browser")
	if err != nil {
		return err
	}

	token, refreshToken, err := c.authTokenHandler.GetCCloudTokens(c.ccloudClientFactory, url, credentials, noBrowser, orgResourceId)

	endOfFreeTrialErr, isEndOfFreeTrialErr := err.(*errors.EndOfFreeTrialError)

	// for orgs that are suspended because it reached end of free trial due to paywall removal, they should still be
	// able to log in.
	if err != nil && !isEndOfFreeTrialErr {
		// for orgs that are suspended due to other reason, they shouldn't be able to log in and should return error immediately.
		if err, ok := err.(*ccloud.SuspendedOrganizationError); ok {
			return errors.NewErrorWithSuggestions(err.Error(), errors.SuspendedOrganizationSuggestions)
		}
		return err
	}

	client := c.ccloudClientFactory.JwtHTTPClientFactory(context.Background(), token, url)
	credentials.AuthToken = token
	credentials.AuthRefreshToken = refreshToken

	currentEnv, currentOrg, err := pauth.PersistCCloudCredentialsToConfig(c.Config.Config, client, url, credentials)
	if err != nil {
		return err
	}

	utils.Printf(cmd, errors.LoggedInAsMsgWithOrg, credentials.Username, currentOrg.ResourceId, currentOrg.Name)
	log.CliLogger.Debugf(errors.LoggedInUsingEnvMsg, currentEnv.Id, currentEnv.Name)

	// if org is at the end of free trial, print instruction about how to add payment method to unsuspend the org.
	// otherwise, print remaining free credit upon each login.
	if isEndOfFreeTrialErr {
		// only print error and do not return it, since end-of-free-trial users should still be able to log in.
		utils.ErrPrintln(cmd, fmt.Sprintf("Error: %s", endOfFreeTrialErr.Error()))
		errors.DisplaySuggestionsMessage(endOfFreeTrialErr.UserFacingError(), os.Stderr)
	} else {
		if err := c.printRemainingFreeCredit(cmd, client); err != nil {
			return err
		}
	}

	return c.saveLoginToNetrc(cmd, true, credentials)
}

func (c *command) printRemainingFreeCredit(cmd *cobra.Command, client *ccloud.Client) error {
	org := &orgv1.Organization{Id: c.Config.Context().State.Auth.Account.OrganizationId}
	promoCodes, err := client.Billing.GetClaimedPromoCodes(context.Background(), org, true)
	if err != nil {
		return err
	}

	// only print remaining free credit if there is any unexpired promo code
	if len(promoCodes) != 0 {
		var remainingFreeCredit int64
		for _, promoCode := range promoCodes {
			remainingFreeCredit += promoCode.Balance
		}
		utils.Println(cmd, fmt.Sprintf(errors.RemainingFreeCreditMsg, admin.ConvertToUSD(remainingFreeCredit)))
	}

	return nil
}

// Order of precedence: env vars > config file > netrc file > prompt
// i.e. if login credentials found in env vars then acquire token using env vars and skip checking for credentials else where
func (c *command) getCCloudCredentials(cmd *cobra.Command, url, orgResourceId string) (*pauth.Credentials, error) {
	client := c.ccloudClientFactory.AnonHTTPClientFactory(url)
	c.loginCredentialsManager.SetCloudClient(client)

	promptOnly, err := cmd.Flags().GetBool("prompt")
	if err != nil {
		return nil, err
	}
	if promptOnly {
		return pauth.GetLoginCredentials(c.loginCredentialsManager.GetCloudCredentialsFromPrompt(cmd, orgResourceId))
	}

	netrcFilterParams := netrc.NetrcMachineParams{
		IsCloud: true,
		URL:     url,
	}
	return pauth.GetLoginCredentials(
		c.loginCredentialsManager.GetCloudCredentialsFromEnvVar(orgResourceId),
		c.loginCredentialsManager.GetCredentialsFromConfig(c.cfg),
		c.loginCredentialsManager.GetCredentialsFromNetrc(cmd, netrcFilterParams),
		c.loginCredentialsManager.GetCloudCredentialsFromPrompt(cmd, orgResourceId),
	)
}

func (c *command) loginMDS(cmd *cobra.Command, url string) error {
	credentials, err := c.getConfluentCredentials(cmd, url)
	if err != nil {
		return err
	}

	// Current functionality:
	// empty ca-cert-path is equivalent to not using ca-cert-path flag
	// if users want to login with ca-cert-path they must explicilty use the flag every time they login
	//
	// For legacy users:
	// if ca-cert-path flag is not used, then return caCertPath value stored in config for the login context
	// if user passes empty string for ca-cert-path flag then reset the ca-cert-path value in config for the context
	// (only for legacy contexts is it still possible for the context name without ca-cert-path to have ca-cert-path)
	var isLegacyContext bool
	caCertPath, err := getCACertPath(cmd)
	if err != nil {
		return err
	}
	if caCertPath == "" {
		contextName := pauth.GenerateContextName(credentials.Username, url, "")
		caCertPath, err = c.checkLegacyContextCACertPath(cmd, contextName)
		if err != nil {
			return err
		}
		isLegacyContext = caCertPath != ""
	}

	if caCertPath != "" {
		caCertPath, err = filepath.Abs(caCertPath)
		if err != nil {
			return err
		}
	}

	client, err := c.mdsClientManager.GetMDSClient(url, caCertPath)
	if err != nil {
		return err
	}

	token, err := c.authTokenHandler.GetConfluentToken(client, credentials)
	if err != nil {
		return err
	}

	err = pauth.PersistConfluentLoginToConfig(c.Config.Config, credentials.Username, url, token, caCertPath, isLegacyContext)
	if err != nil {
		return err
	}

	if err := c.saveLoginToNetrc(cmd, false, credentials); err != nil {
		return err
	}

	log.CliLogger.Debugf(errors.LoggedInAsMsg, credentials.Username)
	return nil
}

func getCACertPath(cmd *cobra.Command) (string, error) {
	if path, err := cmd.Flags().GetString("ca-cert-path"); path != "" || err != nil {
		return path, err
	}

	return pauth.GetEnvWithFallback(pauth.ConfluentPlatformCACertPath, pauth.DeprecatedConfluentPlatformCACertPath), nil
}

// Order of precedence: env vars > netrc > prompt
// i.e. if login credentials found in env vars then acquire token using env vars and skip checking for credentials else where
func (c *command) getConfluentCredentials(cmd *cobra.Command, url string) (*pauth.Credentials, error) {
	promptOnly, err := cmd.Flags().GetBool("prompt")
	if err != nil {
		return nil, err
	}
	if promptOnly {
		return pauth.GetLoginCredentials(c.loginCredentialsManager.GetOnPremCredentialsFromPrompt(cmd))
	}

	netrcFilterParams := netrc.NetrcMachineParams{
		IsCloud: false,
		URL:     url,
	}

	return pauth.GetLoginCredentials(
		c.loginCredentialsManager.GetOnPremCredentialsFromEnvVar(),
		c.loginCredentialsManager.GetCredentialsFromNetrc(cmd, netrcFilterParams),
		c.loginCredentialsManager.GetOnPremCredentialsFromPrompt(cmd),
	)
}

func (c *command) checkLegacyContextCACertPath(cmd *cobra.Command, contextName string) (string, error) {
	changed := cmd.Flags().Changed("ca-cert-path")
	// if flag used but empty string is passed then user intends to reset the ca-cert-path
	if changed {
		return "", nil
	}
	ctx, ok := c.Config.Contexts[contextName]
	if !ok {
		return "", nil
	}
	return ctx.Platform.CaCertPath, nil
}

func (c *command) getURL(cmd *cobra.Command) (string, error) {
	if url, err := cmd.Flags().GetString("url"); url != "" || err != nil {
		return url, err
	}

	if url := pauth.GetEnvWithFallback(pauth.ConfluentPlatformMDSURL, pauth.DeprecatedConfluentPlatformMDSURL); url != "" {
		return url, nil
	}

	return pauth.CCloudURL, nil
}

func (c *command) saveLoginToNetrc(cmd *cobra.Command, isCloud bool, credentials *pauth.Credentials) error {
	save, err := cmd.Flags().GetBool("save")
	if err != nil {
		return err
	}

	if save {
		if credentials.IsSSO {
			utils.ErrPrintln(cmd, "The `--save` flag was ignored since SSO credentials are not stored locally.")
			return nil
		}

		if err := c.netrcHandler.WriteNetrcCredentials(isCloud, credentials.IsSSO, c.Config.Config.Context().NetrcMachineName, credentials.Username, credentials.Password); err != nil {
			return err
		}

		utils.ErrPrintf(cmd, errors.WroteCredentialsToNetrcMsg, c.netrcHandler.GetFileName())
	}

	return nil
}

func validateURL(url string, isCCloud bool) (string, string, error) {
	if isCCloud {
		for _, hostname := range ccloudv2.Hostnames {
			if strings.Contains(url, hostname) {
				if !strings.HasSuffix(strings.TrimSuffix(url, "/"), hostname) {
					return url, "", errors.NewErrorWithSuggestions(errors.UnneccessaryUrlFlagForCloudLoginErrorMsg, errors.UnneccessaryUrlFlagForCloudLoginSuggestions)
				} else {
					break
				}
			}
		}
	}
	protocolRgx, _ := regexp.Compile(`(\w+)://`)
	portRgx, _ := regexp.Compile(`:(\d+\/?)`)

	protocolMatch := protocolRgx.MatchString(url)
	portMatch := portRgx.MatchString(url)

	var msg []string
	if !protocolMatch {
		if isCCloud {
			url = "https://" + url
			msg = append(msg, "https protocol")
		} else {
			url = "http://" + url
			msg = append(msg, "http protocol")
		}
	}
	if !portMatch && !isCCloud {
		url = url + ":8090"
		msg = append(msg, "default MDS port 8090")
	}

	var pattern string
	if isCCloud {
		pattern = `^\w+://[^/ ]+`
	} else {
		pattern = `^\w+://[^/ ]+:\d+(?:\/|$)`
	}
	matched, _ := regexp.MatchString(pattern, url)
	if !matched {
		return url, "", errors.New(errors.InvalidLoginURLMsg)
	}

	return url, strings.Join(msg, " and "), nil
}

func (c *command) getOrgResourceId(cmd *cobra.Command) (string, error) {
	return pauth.GetLoginOrganization(
		c.loginOrganizationManager.GetLoginOrganizationFromArgs(cmd),
		c.loginOrganizationManager.GetLoginOrganizationFromEnvVar(cmd),
		c.loginOrganizationManager.GetDefaultLoginOrganization(),
	)
}
