package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/bincover"
	linkv1 "github.com/confluentinc/cc-structs/kafka/clusterlink/v1"
	corev1 "github.com/confluentinc/cc-structs/kafka/core/v1"
	orgv1 "github.com/confluentinc/cc-structs/kafka/org/v1"
	productv1 "github.com/confluentinc/cc-structs/kafka/product/core/v1"
	schedv1 "github.com/confluentinc/cc-structs/kafka/scheduler/v1"
	utilv1 "github.com/confluentinc/cc-structs/kafka/util/v1"
	opv1 "github.com/confluentinc/cc-structs/operator/v1"
	"github.com/confluentinc/ccloud-sdk-go"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/confluentinc/cli/internal/pkg/config"
	v3 "github.com/confluentinc/cli/internal/pkg/config/v3"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/utils"
)

var (
	noRebuild           = flag.Bool("no-rebuild", false, "skip rebuilding CLI if it already exists")
	update              = flag.Bool("update", false, "update golden files")
	debug               = flag.Bool("debug", true, "enable verbose output")
	skipSsoBrowserTests = flag.Bool("skip-sso-browser-tests", false, "If flag is preset, run the tests that require a web browser.")
	ssoTestEmail        = *flag.String("sso-test-user-email", "ziru+paas-integ-sso@confluent.io", "The email of an sso enabled test user.")
	ssoTestPassword     = *flag.String("sso-test-user-password", "aWLw9eG+F", "The password for the sso enabled test user.")
	// this connection is preconfigured in Auth0 to hit a test Okta account
	ssoTestConnectionName = *flag.String("sso-test-connection-name", "confluent-dev", "The Auth0 SSO connection name.")
	// browser tests by default against devel
	ssoTestLoginUrl  = *flag.String("sso-test-login-url", "https://devel.cpdev.cloud", "The login url to use for the sso browser test.")
	cover            = false
	ccloudTestBin    = ccloudTestBinNormal
	confluentTestBin = confluentTestBinNormal
	covCollector     *bincover.CoverageCollector
	environments     = []*orgv1.Account{{Id: "a-595", Name: "default"}, {Id: "not-595", Name: "other"}}
	serviceAccountID = int32(12345)
)

const (
	confluentTestBinNormal = "confluent_test"
	ccloudTestBinNormal    = "ccloud_test"
	ccloudTestBinRace      = "ccloud_test_race"
	confluentTestBinRace   = "confluent_test_race"
	mergedCoverageFilename = "integ_coverage.txt"
)

// CLITest represents a test configuration
type CLITest struct {
	// Name to show in go test output; defaults to args if not set
	name string
	// The CLI command being tested; this is a string of args and flags passed to the binary
	args string
	// The set of environment variables to be set when the CLI is run
	env []string
	// "default" if you need to login, or "" otherwise
	login string
	// The kafka cluster ID to "use"
	useKafka string
	// The API Key to set as Kafka credentials
	authKafka string
	// Name of a golden output fixture containing expected output
	fixture string
	// True iff fixture represents a regex
	regex bool
	// Fixed string to check if output contains
	contains string
	// Fixed string to check that output does not contain
	notContains string
	// Expected exit code (e.g., 0 for success or 1 for failure)
	wantErrCode int
	// If true, don't reset the config/state between tests to enable testing CLI workflows
	workflow bool
	// An optional function that allows you to specify other calls
	wantFunc func(t *testing.T)
}

// CLITestSuite is the CLI integration tests.
type CLITestSuite struct {
	suite.Suite
}

// TestCLI runs the CLI integration test suite.
func TestCLI(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}

func init() {
	collectCoverage := os.Getenv("INTEG_COVER")
	cover = collectCoverage == "on"
	ciEnv := os.Getenv("CI")
	if ciEnv == "on" {
		ccloudTestBin = ccloudTestBinRace
		confluentTestBin = confluentTestBinRace
	}
	if runtime.GOOS == "windows" {
		ccloudTestBin = ccloudTestBin + ".exe"
		confluentTestBin = confluentTestBin + ".exe"
	}
}

// SetupSuite builds the CLI binary to test
func (s *CLITestSuite) SetupSuite() {
	covCollector = bincover.NewCoverageCollector(mergedCoverageFilename, cover)
	covCollector.Setup()
	req := require.New(s.T())

	// dumb but effective
	err := os.Chdir("..")
	req.NoError(err)
	err = os.Setenv("XX_CCLOUD_RBAC", "yes")
	req.NoError(err)
	for _, binary := range []string{ccloudTestBin, confluentTestBin} {
		if _, err = os.Stat(binaryPath(s.T(), binary)); os.IsNotExist(err) || !*noRebuild {
			var makeArgs string
			if ccloudTestBin == ccloudTestBinRace {
				makeArgs = "build-integ-race"
			} else {
				makeArgs = "build-integ-nonrace"
			}
			makeCmd := exec.Command("make", makeArgs)
			output, err := makeCmd.CombinedOutput()
			if err != nil {
				s.T().Log(string(output))
				req.NoError(err)
			}
		}
	}
}

func (s *CLITestSuite) TearDownSuite() {
	// Merge coverage profiles.
	_ = os.Unsetenv("XX_CCLOUD_RBAC")
	covCollector.TearDown()
}

func (s *CLITestSuite) TestConfluentHelp() {
	var tests []CLITest
	if runtime.GOOS == "windows" {
		tests = []CLITest{
			{name: "no args", fixture: "confluent-help-flag-windows.golden", wantErrCode: 1},
			{args: "help", fixture: "confluent-help-windows.golden"},
			{args: "--help", fixture: "confluent-help-flag-windows.golden"},
			{args: "version", fixture: "confluent-version.golden", regex: true},
		}
	} else {
		tests = []CLITest{
			{name: "no args", fixture: "confluent-help-flag.golden", wantErrCode: 1},
			{args: "help", fixture: "confluent-help.golden"},
			{args: "--help", fixture: "confluent-help-flag.golden"},
			{args: "version", fixture: "confluent-version.golden", regex: true},
		}
	}

	loginURL := serveMds(s.T()).URL

	for _, tt := range tests {
		s.runConfluentTest(tt, loginURL)
	}
}

func (s *CLITestSuite) TestCcloudHelp() {
	tests := []CLITest{
		{name: "no args", fixture: "help-flag-fail.golden", wantErrCode: 1},
		{args: "help", fixture: "help.golden"},
		{args: "--help", fixture: "help-flag.golden"},
		{args: "version", fixture: "version.golden", regex: true},
	}

	kafkaURL := serveKafkaAPI(s.T()).URL
	loginURL := serve(s.T(), kafkaURL).URL

	for _, tt := range tests {
		s.runCcloudTest(tt, loginURL)
	}
}

func assertUserAgent(t *testing.T, expected string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		require.Regexp(t, expected, r.Header.Get("User-Agent"))
	}
}

func (s *CLITestSuite) TestUserAgent() {
	checkUserAgent := func(t *testing.T, expected string) string {
		kafkaApiRouter := http.NewServeMux()
		kafkaApiRouter.HandleFunc("/", assertUserAgent(t, expected))
		kafkaApiServer := httptest.NewServer(kafkaApiRouter)
		cloudRouter := http.NewServeMux()
		cloudRouter.HandleFunc("/api/sessions", compose(assertUserAgent(t, expected), handleLogin(t)))
		cloudRouter.HandleFunc("/api/me", compose(assertUserAgent(t, expected), handleMe(t)))
		cloudRouter.HandleFunc("/api/check_email/", compose(assertUserAgent(t, expected), handleCheckEmail(t)))
		cloudRouter.HandleFunc("/api/clusters/", compose(assertUserAgent(t, expected), handleKafkaClusterGetListDeleteDescribe(t, kafkaApiServer.URL)))
		return httptest.NewServer(cloudRouter).URL
	}

	serverURL := checkUserAgent(s.T(), fmt.Sprintf("Confluent-Cloud-CLI/v(?:[0-9]\\.?){3}([^ ]*) \\(https://confluent.cloud; support@confluent.io\\) "+
		"ccloud-sdk-go/%s \\(%s/%s; go[^ ]*\\)", ccloud.SDKVersion, runtime.GOOS, runtime.GOARCH))
	env := []string{"XX_CCLOUD_EMAIL=valid@user.com", "XX_CCLOUD_PASSWORD=pass1"}

	s.T().Run("ccloud login", func(tt *testing.T) {
		_ = runCommand(tt, ccloudTestBin, env, "login --url "+serverURL, 0)
	})
	s.T().Run("ccloud cluster list", func(tt *testing.T) {
		_ = runCommand(tt, ccloudTestBin, env, "kafka cluster list", 0)
	})
	s.T().Run("ccloud topic list", func(tt *testing.T) {
		_ = runCommand(tt, ccloudTestBin, env, "kafka topic list --cluster lkc-abc123", 0)
	})
}

func (s *CLITestSuite) TestCcloudErrors() {
	type errorer interface {
		GetError() *corev1.Error
	}
	serveErrors := func(t *testing.T) string {
		req := require.New(t)
		write := func(w http.ResponseWriter, resp proto.Message) {
			if r, ok := resp.(errorer); ok {
				w.WriteHeader(int(r.GetError().Code))
			}
			b, err := utilv1.MarshalJSONToBytes(resp)
			req.NoError(err)
			_, err = io.WriteString(w, string(b))
			req.NoError(err)
		}
		router := http.NewServeMux()
		router.HandleFunc("/api/sessions", handleLogin(t))
		router.HandleFunc("/api/me", handleMe(t))
		router.HandleFunc("/api/check_email/", handleCheckEmail(t))
		router.HandleFunc("/api/clusters", func(w http.ResponseWriter, r *http.Request) {
			switch r.Header.Get("Authorization") {
			// TODO: these assume the upstream doesn't change its error responses. Fragile, fragile, fragile. :(
			// https://github.com/confluentinc/cc-auth-service/blob/06db0bebb13fb64c9bc3c6e2cf0b67709b966632/jwt/token.go#L23
			case "Bearer expired":
				write(w, &schedv1.GetKafkaClustersReply{Error: &corev1.Error{Message: "token is expired", Code: http.StatusUnauthorized}})
			case "Bearer malformed":
				write(w, &schedv1.GetKafkaClustersReply{Error: &corev1.Error{Message: "malformed token", Code: http.StatusBadRequest}})
			case "Bearer invalid":
				// TODO: The response for an invalid token should be 4xx, not 500 (e.g., if you take a working token from devel and try in stag)
				write(w, &schedv1.GetKafkaClustersReply{Error: &corev1.Error{Message: "Token parsing error: crypto/rsa: verification error", Code: http.StatusInternalServerError}})
			default:
				req.Fail("reached the unreachable", "auth=%s", r.Header.Get("Authorization"))
			}
		})
		server := httptest.NewServer(router)
		return server.URL
	}

	s.T().Run("invalid user or pass", func(tt *testing.T) {
		loginURL := serveErrors(tt)
		env := []string{"XX_CCLOUD_EMAIL=incorrect@user.com", "XX_CCLOUD_PASSWORD=pass1"}
		output := runCommand(tt, ccloudTestBin, env, "login --url "+loginURL, 1)
		require.Contains(tt, output, errors.InvalidLoginErrorMsg)
		require.Contains(tt, output, errors.ComposeSuggestionsMessage(errors.CCloudInvalidLoginSuggestions))
	})

	s.T().Run("expired token", func(tt *testing.T) {
		loginURL := serveErrors(tt)
		env := []string{"XX_CCLOUD_EMAIL=expired@user.com", "XX_CCLOUD_PASSWORD=pass1"}
		output := runCommand(tt, ccloudTestBin, env, "login --url "+loginURL, 0)
		require.Contains(tt, output, fmt.Sprintf(errors.LoggedInAsMsg, "expired@user.com"))
		require.Contains(tt, output, fmt.Sprintf(errors.LoggedInUsingEnvMsg, "a-595", "default"))
		output = runCommand(tt, ccloudTestBin, []string{}, "kafka cluster list", 1)
		require.Contains(tt, output, errors.TokenExpiredMsg)
		require.Contains(tt, output, errors.NotLoggedInErrorMsg)
	})

	s.T().Run("malformed token", func(tt *testing.T) {
		loginURL := serveErrors(tt)
		env := []string{"XX_CCLOUD_EMAIL=malformed@user.com", "XX_CCLOUD_PASSWORD=pass1"}
		output := runCommand(tt, ccloudTestBin, env, "login --url "+loginURL, 0)
		require.Contains(tt, output, fmt.Sprintf(errors.LoggedInAsMsg, "malformed@user.com"))
		require.Contains(tt, output, fmt.Sprintf(errors.LoggedInUsingEnvMsg, "a-595", "default"))

		output = runCommand(s.T(), ccloudTestBin, []string{}, "kafka cluster list", 1)
		require.Contains(tt, output, errors.CorruptedTokenErrorMsg)
		require.Contains(tt, output, errors.ComposeSuggestionsMessage(errors.CorruptedTokenSuggestions))
	})

	s.T().Run("invalid jwt", func(tt *testing.T) {
		loginURL := serveErrors(tt)
		env := []string{"XX_CCLOUD_EMAIL=invalid@user.com", "XX_CCLOUD_PASSWORD=pass1"}
		output := runCommand(tt, ccloudTestBin, env, "login --url "+loginURL, 0)
		require.Contains(tt, output, fmt.Sprintf(errors.LoggedInAsMsg, "invalid@user.com"))
		require.Contains(tt, output, fmt.Sprintf(errors.LoggedInUsingEnvMsg, "a-595", "default"))

		output = runCommand(s.T(), ccloudTestBin, []string{}, "kafka cluster list", 1)
		require.Contains(tt, output, errors.CorruptedTokenErrorMsg)
		require.Contains(tt, output, errors.ComposeSuggestionsMessage(errors.CorruptedTokenSuggestions))
	})
}

func (s *CLITestSuite) runCcloudTest(tt CLITest, loginURL string) {
	if tt.name == "" {
		tt.name = tt.args
	}
	if strings.HasPrefix(tt.name, "error") {
		tt.wantErrCode = 1
	}

	s.T().Run(tt.name, func(t *testing.T) {
		if !tt.workflow {
			resetConfiguration(t, "ccloud")
		}

		if tt.login == "default" {
			env := []string{"XX_CCLOUD_EMAIL=fake@user.com", "XX_CCLOUD_PASSWORD=pass1"}
			output := runCommand(t, ccloudTestBin, env, "login --url "+loginURL, 0)
			if *debug {
				fmt.Println(output)
			}
		}

		if tt.useKafka != "" {
			output := runCommand(t, ccloudTestBin, []string{}, "kafka cluster use "+tt.useKafka, 0)
			if *debug {
				fmt.Println(output)
			}
		}

		if tt.authKafka != "" {
			output := runCommand(t, ccloudTestBin, []string{}, "api-key create --resource "+tt.useKafka, 0)
			if *debug {
				fmt.Println(output)
			}
			// HACK: we don't have scriptable output yet so we parse it from the table
			key := strings.TrimSpace(strings.Split(strings.Split(output, "\n")[3], "|")[2])
			output = runCommand(t, ccloudTestBin, []string{}, fmt.Sprintf("api-key use %s --resource %s", key, tt.useKafka), 0)
			if *debug {
				fmt.Println(output)
			}
		}
		output := runCommand(t, ccloudTestBin, tt.env, tt.args, tt.wantErrCode)
		if *debug {
			fmt.Println(output)
		}

		if strings.HasPrefix(tt.args, "kafka cluster create") ||
			strings.HasPrefix(tt.args, "config context current") {
			re := regexp.MustCompile("https?://127.0.0.1:[0-9]+")
			output = re.ReplaceAllString(output, "http://127.0.0.1:12345")
		}

		if strings.HasPrefix(tt.args, "api-key list") {

		}

		s.validateTestOutput(tt, t, output)
	})
}

func (s *CLITestSuite) runConfluentTest(tt CLITest, loginURL string) {
	if tt.name == "" {
		tt.name = tt.args
	}
	if strings.HasPrefix(tt.name, "error") {
		tt.wantErrCode = 1
	}
	s.T().Run(tt.name, func(t *testing.T) {
		if !tt.workflow {
			resetConfiguration(t, "confluent")
		}

		if tt.login == "default" {
			env := []string{"XX_CONFLUENT_USERNAME=fake@user.com", "XX_CONFLUENT_PASSWORD=pass1"}
			output := runCommand(t, confluentTestBin, env, "login --url "+loginURL, 0)
			if *debug {
				fmt.Println(output)
			}
		}

		output := runCommand(t, confluentTestBin, []string{}, tt.args, tt.wantErrCode)

		if strings.HasPrefix(tt.args, "config context list") ||
			strings.HasPrefix(tt.args, "config context current") {
			re := regexp.MustCompile("https?://127.0.0.1:[0-9]+")
			output = re.ReplaceAllString(output, "http://127.0.0.1:12345")
		}

		s.validateTestOutput(tt, t, output)
	})
}

func (s *CLITestSuite) validateTestOutput(tt CLITest, t *testing.T, output string) {
	if *update && !tt.regex && tt.fixture != "" {
		writeFixture(t, tt.fixture, output)
	}
	actual := utils.NormalizeNewLines(output)
	if tt.contains != "" {
		require.Contains(t, actual, tt.contains)
	} else if tt.notContains != "" {
		require.NotContains(t, actual, tt.notContains)
	} else if tt.fixture != "" {
		expected := utils.NormalizeNewLines(LoadFixture(t, tt.fixture))
		if tt.regex {
			require.Regexp(t, expected, actual)
		} else if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\n   actual:\n%s\nexpected:\n%s", actual, expected)
		}
	}
	if tt.wantFunc != nil {
		tt.wantFunc(t)
	}
}

func runCommand(t *testing.T, binaryName string, env []string, args string, wantErrCode int) string {
	output, exitCode, err := covCollector.RunBinary(binaryPath(t, binaryName), "TestRunMain", env, strings.Split(args, " "))
	if err != nil && wantErrCode == 0 {
		require.Failf(t, "unexpected error",
			"exit %d: %s\n%s", exitCode, args, output)
	}
	require.Equal(t, wantErrCode, exitCode, output)
	return output
}

func resetConfiguration(t *testing.T, cliName string) {
	// HACK: delete your current config to isolate tests cases for non-workflow tests...
	// probably don't really want to do this or devs will get mad
	cfg := v3.New(&config.Params{
		CLIName: cliName,
	})
	err := cfg.Save()
	require.NoError(t, err)
}

func writeFixture(t *testing.T, fixture string, content string) {
	err := ioutil.WriteFile(FixturePath(t, fixture), []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func binaryPath(t *testing.T, binaryName string) string {
	dir, err := os.Getwd()
	require.NoError(t, err)
	return path.Join(dir, binaryName)
}

var (
	keyStore        = map[int32]*schedv1.ApiKey{}
	keyIndex        = int32(1)
	keyTimestamp, _ = types.TimestampProto(time.Date(1999, time.February, 24, 0, 0, 0, 0, time.UTC))
)

type ApiKeyList []*schedv1.ApiKey

// Len is part of sort.Interface.
func (d ApiKeyList) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d ApiKeyList) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use Key as the value to sort by
func (d ApiKeyList) Less(i, j int) bool {
	return d[i].Key < d[j].Key
}

func init() {
	keyStore[keyIndex] = &schedv1.ApiKey{
		Id:     keyIndex,
		Key:    "MYKEY1",
		Secret: "MYSECRET1",
		LogicalClusters: []*schedv1.ApiKey_Cluster{
			{Id: "lkc-bob", Type: "kafka"},
		},
		UserId: 12,
	}
	keyIndex += 1
	keyStore[keyIndex] = &schedv1.ApiKey{
		Id:     keyIndex,
		Key:    "MYKEY2",
		Secret: "MYSECRET2",
		LogicalClusters: []*schedv1.ApiKey_Cluster{
			{Id: "lkc-abc", Type: "kafka"},
		},
		UserId: 18,
	}
	keyIndex += 1
	keyStore[100] = &schedv1.ApiKey{
		Id:     keyIndex,
		Key:    "UIAPIKEY100",
		Secret: "UIAPISECRET100",
		LogicalClusters: []*schedv1.ApiKey_Cluster{
			{Id: "lkc-cool1", Type: "kafka"},
		},
		UserId: 25,
	}
	keyStore[101] = &schedv1.ApiKey{
		Id:     keyIndex,
		Key:    "UIAPIKEY101",
		Secret: "UIAPISECRET101",
		LogicalClusters: []*schedv1.ApiKey_Cluster{
			{Id: "lkc-other1", Type: "kafka"},
		},
		UserId: 25,
	}
	keyStore[102] = &schedv1.ApiKey{
		Id:     keyIndex,
		Key:    "UIAPIKEY102",
		Secret: "UIAPISECRET102",
		LogicalClusters: []*schedv1.ApiKey_Cluster{
			{Id: "lksqlc-ksql1", Type: "ksql"},
		},
		UserId: 25,
	}
	keyStore[103] = &schedv1.ApiKey{
		Id:     keyIndex,
		Key:    "UIAPIKEY103",
		Secret: "UIAPISECRET103",
		LogicalClusters: []*schedv1.ApiKey_Cluster{
			{Id: "lkc-cool1", Type: "kafka"},
		},
		UserId: 25,
	}
	keyStore[200] = &schedv1.ApiKey{
		Id:     keyIndex,
		Key:    "SERVICEACCOUNTKEY1",
		Secret: "SERVICEACCOUNTSECRET1",
		LogicalClusters: []*schedv1.ApiKey_Cluster{
			{Id: "lkc-bob", Type: "kafka"},
		},
		UserId: serviceAccountID,
	}
	for _, k := range keyStore {
		k.Created = keyTimestamp
	}
}

func serve(t *testing.T, kafkaAPIURL string) *httptest.Server {
	router := http.NewServeMux()
	router.HandleFunc("/api/sessions", handleLogin(t))
	router.HandleFunc("/api/check_email/", handleCheckEmail(t))
	router.HandleFunc("/api/me", handleMe(t))
	router.HandleFunc("/api/api_keys", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			req := &schedv1.CreateApiKeyRequest{}
			err := utilv1.UnmarshalJSON(r.Body, req)
			require.NoError(t, err)
			require.NotEmpty(t, req.ApiKey.AccountId)
			apiKey := req.ApiKey
			apiKey.Id = keyIndex
			apiKey.Key = fmt.Sprintf("MYKEY%d", keyIndex)
			apiKey.Secret = fmt.Sprintf("MYSECRET%d", keyIndex)
			apiKey.Created = keyTimestamp
			if req.ApiKey.UserId == 0 {
				apiKey.UserId = 23
			} else {
				apiKey.UserId = req.ApiKey.UserId
			}
			keyIndex++
			keyStore[apiKey.Id] = apiKey
			b, err := utilv1.MarshalJSONToBytes(&schedv1.CreateApiKeyReply{ApiKey: apiKey})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(b))
			require.NoError(t, err)
		} else if r.Method == "GET" {
			require.NotEmpty(t, r.URL.Query().Get("account_id"))
			apiKeys := apiKeysFilter(r.URL)
			// Return sorted data or the test output will not be stable
			sort.Sort(ApiKeyList(apiKeys))
			b, err := utilv1.MarshalJSONToBytes(&schedv1.GetApiKeysReply{ApiKeys: apiKeys})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(b))
			require.NoError(t, err)
		}
	})
	router.HandleFunc("/api/api_keys/", handleAPIKeyUpdateAndDelete(t))
	router.HandleFunc("/api/accounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			b, err := utilv1.MarshalJSONToBytes(&orgv1.ListAccountsReply{Accounts: environments})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(b))
			require.NoError(t, err)
		} else if r.Method == "POST" {
			req := &orgv1.CreateAccountRequest{}
			err := utilv1.UnmarshalJSON(r.Body, req)
			require.NoError(t, err)
			account := &orgv1.Account{
				Id:             "a-5555",
				Name:           req.Account.Name,
				OrganizationId: 0,
			}
			b, err := utilv1.MarshalJSONToBytes(&orgv1.CreateAccountReply{
				Account: account,
			})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(b))
			require.NoError(t, err)
		}
	})
	router.HandleFunc("/api/accounts/a-595", handleEnvironmentRequests(t, "a-595"))
	router.HandleFunc("/api/accounts/not-595", handleEnvironmentRequests(t, "not-595"))
	router.HandleFunc("/api/clusters/lkc-describe", handleKafkaClusterDescribeTest(t))
	router.HandleFunc("/api/clusters/lkc-describe-dedicated", handleKafkaClusterDescribeTest(t))
	router.HandleFunc("/api/clusters/lkc-describe-dedicated-pending", handleKafkaClusterDescribeTest(t))
	router.HandleFunc("/api/clusters/lkc-describe-dedicated-with-encryption", handleKafkaClusterDescribeTest(t))
	router.HandleFunc("/api/clusters/lkc-update", handleKafkaClusterUpdateTest(t))
	router.HandleFunc("/api/clusters/lkc-update-dedicated", handleKafkaDedicatedClusterUpdateTest(t))
	router.HandleFunc("/api/clusters/", handleKafkaClusterGetListDeleteDescribe(t, kafkaAPIURL))
	router.HandleFunc("/api/clusters", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			handleKafkaClusterCreate(t, kafkaAPIURL)(w, r)
		} else if r.Method == "GET" {
			cluster := schedv1.KafkaCluster{
				Id:              "lkc-123",
				Name:            "abc",
				Deployment:      &schedv1.Deployment{Sku: productv1.Sku_BASIC},
				Durability:      0,
				Status:          0,
				Region:          "us-central1",
				ServiceProvider: "gcp",
			}
			clusterMultizone := schedv1.KafkaCluster{
				Id:              "lkc-456",
				Name:            "def",
				Deployment:      &schedv1.Deployment{Sku: productv1.Sku_BASIC},
				Durability:      1,
				Status:          0,
				Region:          "us-central1",
				ServiceProvider: "gcp",
			}
			b, err := utilv1.MarshalJSONToBytes(&schedv1.GetKafkaClustersReply{
				Clusters: []*schedv1.KafkaCluster{&cluster, &clusterMultizone},
			})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(b))
			require.NoError(t, err)
		}
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"error": {"message": "unexpected call to `+r.URL.Path+`"}}`)
		require.NoError(t, err)
	})
	router.HandleFunc("/api/schema_registries/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		id := q.Get("id")
		if id == "" {
			id = "lsrc-1234"
		}
		accountId := q.Get("account_id")
		srCluster := &schedv1.SchemaRegistryCluster{
			Id:        id,
			AccountId: accountId,
			Name:      "account schema-registry",
			Endpoint:  "SASL_SSL://sr-endpoint",
		}
		fmt.Println(srCluster)
		b, err := utilv1.MarshalJSONToBytes(&schedv1.GetSchemaRegistryClusterReply{
			Cluster: srCluster,
		})
		require.NoError(t, err)
		_, err = io.WriteString(w, string(b))
		require.NoError(t, err)
	})
	router.HandleFunc("/api/service_accounts", handleServiceAccountRequests(t))
	router.HandleFunc("/api/accounts/a-595/clusters/lkc-123/connectors/az-connector/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		return
	})
	router.HandleFunc("/api/accounts/a-595/clusters/lkc-123/connectors", handleConnect(t))
	router.HandleFunc("/api/accounts/a-595/clusters/lkc-123/connector-plugins/GcsSink/config/validate", handleConnectorCatalogDescribe(t))
	router.HandleFunc("/api/accounts/a-595/clusters/lkc-123/connector-plugins", handleConnectPlugins(t))
	router.HandleFunc("/api/ksqls", handleKSQLCreateList(t))
	router.HandleFunc("/api/ksqls/lksqlc-ksql1/", func(w http.ResponseWriter, r *http.Request) {
		ksqlCluster := &schedv1.KSQLCluster{
			Id:                "lksqlc-ksql1",
			AccountId:         "25",
			KafkaClusterId:    "lkc-12345",
			OutputTopicPrefix: "pksqlc-abcde",
			Name:              "account ksql",
			Storage:           101,
			Endpoint:          "SASL_SSL://ksql-endpoint",
		}
		reply, err := utilv1.MarshalJSONToBytes(&schedv1.GetKSQLClusterReply{
			Cluster: ksqlCluster,
		})
		require.NoError(t, err)
		_, err = io.WriteString(w, string(reply))
		require.NoError(t, err)
	})
	router.HandleFunc("/api/ksqls/lksqlc-12345", func(w http.ResponseWriter, r *http.Request) {
		ksqlCluster := &schedv1.KSQLCluster{
			Id:                "lksqlc-12345",
			AccountId:         "25",
			KafkaClusterId:    "lkc-abcde",
			OutputTopicPrefix: "pksqlc-zxcvb",
			Name:              "account ksql",
			Storage:           130,
			Endpoint:          "SASL_SSL://ksql-endpoint",
		}
		reply, err := utilv1.MarshalJSONToBytes(&schedv1.GetKSQLClusterReply{
			Cluster: ksqlCluster,
		})
		require.NoError(t, err)
		_, err = io.WriteString(w, string(reply))
		require.NoError(t, err)
	})
	router.HandleFunc("/api/env_metadata", func(w http.ResponseWriter, r *http.Request) {
		clouds := []*schedv1.CloudMetadata{
			{
				Id:   "gcp",
				Name: "Google Cloud Platform",
				Regions: []*schedv1.Region{
					{
						Id:            "asia-southeast1",
						Name:          "asia-southeast1 (Singapore)",
						IsSchedulable: true,
					},
					{
						Id:            "asia-east2",
						Name:          "asia-east2 (Hong Kong)",
						IsSchedulable: true,
					},
				},
			},
			{
				Id:   "aws",
				Name: "Amazon Web Services",
				Regions: []*schedv1.Region{
					{
						Id:            "ap-northeast-1",
						Name:          "ap-northeast-1 (Tokyo)",
						IsSchedulable: false,
					},
					{
						Id:            "us-east-1",
						Name:          "us-east-1 (N. Virginia)",
						IsSchedulable: true,
					},
				},
			},
			{
				Id:   "azure",
				Name: "Azure",
				Regions: []*schedv1.Region{
					{
						Id:            "southeastasia",
						Name:          "southeastasia (Singapore)",
						IsSchedulable: false,
					},
				},
			},
		}
		reply, err := utilv1.MarshalJSONToBytes(&schedv1.GetEnvironmentMetadataReply{
			Clouds: clouds,
		})
		require.NoError(t, err)
		_, err = io.WriteString(w, string(reply))
		require.NoError(t, err)
	})
	router.HandleFunc("/api/organizations/0/price_table", handlePriceTable(t))
	router.HandleFunc("/api/organizations/0/payment_info", handlePaymentInfo(t))
	router.HandleFunc("/api/users", handleUsers(t))
	addMdsv2alpha1(t, router)
	return httptest.NewServer(router)
}

func apiKeysFilter(url *url.URL) []*schedv1.ApiKey {
	var apiKeys []*schedv1.ApiKey
	q := url.Query()
	uid := q.Get("user_id")
	clusterIds := q["cluster_id"]

	for _, a := range keyStore {
		uidFilter := (uid == "0") || (uid == strconv.Itoa(int(a.UserId)))
		clusterFilter := (len(clusterIds) == 0) || func(clusterIds []string) bool {
			for _, c := range a.LogicalClusters {
				for _, clusterId := range clusterIds {
					if c.Id == clusterId {
						return true
					}
				}
			}
			return false
		}(clusterIds)

		if uidFilter && clusterFilter {
			apiKeys = append(apiKeys, a)
		}
	}
	return apiKeys
}

func serveKafkaAPI(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/2.0/kafka/lkc-acls/acls:search", handleKafkaACLsList(t))
	mux.HandleFunc("/2.0/kafka/lkc-acls/acls", handleKafkaACLsCreate(t))
	mux.HandleFunc("/2.0/kafka/lkc-acls/acls/delete", handleKafkaACLsDelete(t))

	mux.HandleFunc("/2.0/kafka/lkc-links/links/", handleKafkaLinks(t))

	mux.HandleFunc("/2.0/kafka/lkc-topics/topics/test-topic/mirror:stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("/2.0/kafka/lkc-topics/topics/not-found/mirror:stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// TODO: no idea how this "topic already exists" API request or response actually looks
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, err := io.WriteString(w, `{}`)
		require.NoError(t, err)
	})
	return httptest.NewServer(mux)
}

func handleKafkaLinks(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		lastElem := parts[len(parts)-1]

		if lastElem == "" {
			// No specific link here, we want a list of ALL links

			listResponsePayload := []*linkv1.ListLinksResponseItem{
				&linkv1.ListLinksResponseItem{LinkName: "link-1", LinkId: "1234", ClusterId: "Blah"},
				&linkv1.ListLinksResponseItem{LinkName: "link-2", LinkId: "4567", ClusterId: "blah"},
			}

			listReply, err := json.Marshal(listResponsePayload)
			require.NoError(t, err)
			_, err = io.WriteString(w, string(listReply))
			require.NoError(t, err)
		} else {
			// Return properties for the selected link.
			describeResponsePayload := linkv1.DescribeLinkResponse{
				Entries: []*linkv1.DescribeLinkResponseEntry{
					{
						Name:  "replica.fetch.max.bytes",
						Value: "1048576",
					},
				},
			}
			describeReply, err := json.Marshal(describeResponsePayload)
			require.NoError(t, err)
			_, err = io.WriteString(w, string(describeReply))
			require.NoError(t, err)
		}
	}
}

func handleLogin(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := require.New(t)
		b, err := ioutil.ReadAll(r.Body)
		req.NoError(err)
		auth := &struct {
			Email    string
			Password string
		}{}
		err = json.Unmarshal(b, auth)
		req.NoError(err)
		switch auth.Email {
		case "incorrect@user.com":
			w.WriteHeader(http.StatusForbidden)
		case "expired@user.com":
			http.SetCookie(w, &http.Cookie{Name: "auth_token", Value: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MzAxMjQ4NTcsImV4cCI6MTUzMDAzODQ1NywiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSJ9.Y2ui08GPxxuV9edXUBq-JKr1VPpMSnhjSFySczCby7Y"})
		case "malformed@user.com":
			http.SetCookie(w, &http.Cookie{Name: "auth_token", Value: "malformed"})
		case "invalid@user.com":
			http.SetCookie(w, &http.Cookie{Name: "auth_token", Value: "invalid"})
		default:
			http.SetCookie(w, &http.Cookie{Name: "auth_token", Value: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1NjE2NjA4NTcsImV4cCI6MjUzMzg2MDM4NDU3LCJhdWQiOiJ3d3cuZXhhbXBsZS5jb20iLCJzdWIiOiJqcm9ja2V0QGV4YW1wbGUuY29tIn0.G6IgrFm5i0mN7Lz9tkZQ2tZvuZ2U7HKnvxMuZAooPmE"})
		}
	}
}

func handleMe(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := utilv1.MarshalJSONToBytes(&orgv1.GetUserReply{
			User: &orgv1.User{
				Id:         23,
				Email:      "cody@confluent.io",
				FirstName:  "Cody",
				ResourceId: "u-11aaa",
			},
			Accounts: environments,
		})
		require.NoError(t, err)
		_, err = io.WriteString(w, string(b))
		require.NoError(t, err)
	}
}

func handleCheckEmail(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := require.New(t)
		email := strings.Replace(r.URL.String(), "/api/check_email/", "", 1)
		reply := &orgv1.GetUserReply{}
		switch email {
		case "cody@confluent.io":
			reply.User = &orgv1.User{
				Email: "cody@confluent.io",
			}
		}
		b, err := utilv1.MarshalJSONToBytes(reply)
		req.NoError(err)
		_, err = io.WriteString(w, string(b))
		req.NoError(err)
	}
}

func handleKafkaClusterGetListDeleteDescribe(t *testing.T, kafkaAPIURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		id := parts[len(parts)-1]
		if id == "lkc-unknown" {
			_, err := io.WriteString(w, `{"error":{"code":404,"message":"resource not found","nested_errors":{},"details":[],"stack":null},"cluster":null}`)
			require.NoError(t, err)
			return
		}
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			// this is in the body of delete requests
			require.NotEmpty(t, r.URL.Query().Get("account_id"))
		}
		// Now return the KafkaCluster with updated ApiEndpoint
		b, err := utilv1.MarshalJSONToBytes(&schedv1.GetKafkaClusterReply{
			Cluster: &schedv1.KafkaCluster{
				Id:              id,
				Name:            "kafka-cluster",
				Deployment:      &schedv1.Deployment{Sku: productv1.Sku_BASIC},
				NetworkIngress:  100,
				NetworkEgress:   100,
				Storage:         500,
				ServiceProvider: "aws",
				Region:          "us-west-2",
				Endpoint:        "SASL_SSL://kafka-endpoint",
				ApiEndpoint:     kafkaAPIURL,
			},
		})
		require.NoError(t, err)
		_, err = io.WriteString(w, string(b))
		require.NoError(t, err)
	}
}

func handleKafkaClusterDescribeTest(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		cluster := &schedv1.KafkaCluster{
			Id:              id,
			Name:            "kafka-cluster",
			Deployment:      &schedv1.Deployment{Sku: productv1.Sku_BASIC},
			NetworkIngress:  100,
			NetworkEgress:   100,
			Storage:         500,
			ServiceProvider: "aws",
			Region:          "us-west-2",
			Endpoint:        "SASL_SSL://kafka-endpoint",
			ApiEndpoint:     "http://kafka-api-url",
		}
		switch id {
		case "lkc-describe-dedicated":
			cluster.Cku = 1
			cluster.Deployment = &schedv1.Deployment{Sku: productv1.Sku_DEDICATED}
		case "lkc-describe-dedicated-pending":
			cluster.Cku = 1
			cluster.PendingCku = 2
			cluster.Deployment = &schedv1.Deployment{Sku: productv1.Sku_DEDICATED}
		case "lkc-describe-dedicated-with-encryption":
			cluster.Cku = 1
			cluster.EncryptionKeyId = "abc123"
			cluster.Deployment = &schedv1.Deployment{Sku: productv1.Sku_DEDICATED}
		}
		b, err := utilv1.MarshalJSONToBytes(&schedv1.GetKafkaClusterReply{
			Cluster: cluster,
		})
		require.NoError(t, err)
		_, err = io.WriteString(w, string(b))
		require.NoError(t, err)
	}
}

func handleKafkaClusterUpdateTest(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Describe client call
		var out []byte
		if r.Method == "GET" {
			id := r.URL.Query().Get("id")
			var err error
			out, err = utilv1.MarshalJSONToBytes(&schedv1.GetKafkaClusterReply{
				Cluster: &schedv1.KafkaCluster{
					Id:              id,
					Name:            "lkc-update",
					Deployment:      &schedv1.Deployment{Sku: productv1.Sku_BASIC},
					NetworkIngress:  100,
					NetworkEgress:   100,
					Storage:         500,
					Status:          schedv1.ClusterStatus_UP,
					ServiceProvider: "aws",
					Region:          "us-west-2",
					Endpoint:        "SASL_SSL://kafka-endpoint",
					ApiEndpoint:     "http://kafka-api-url",
				},
			})
			require.NoError(t, err)
		}
		// Update client call
		if r.Method == "PUT" {
			req := &schedv1.UpdateKafkaClusterRequest{}
			err := utilv1.UnmarshalJSON(r.Body, req)
			require.NoError(t, err)
			if req.Cluster.Cku > 0 {
				out, err = utilv1.MarshalJSONToBytes(&schedv1.GetKafkaClusterReply{
					Cluster: nil,
					Error: &corev1.Error{
						Message: "cluster expansion is supported for dedicated clusters only",
					},
				})
			} else {
				out, err = utilv1.MarshalJSONToBytes(&schedv1.GetKafkaClusterReply{
					Cluster: &schedv1.KafkaCluster{
						Id:              req.Cluster.Id,
						Name:            req.Cluster.Name,
						Deployment:      &schedv1.Deployment{Sku: productv1.Sku_BASIC},
						NetworkIngress:  100,
						NetworkEgress:   100,
						Storage:         500,
						Status:          schedv1.ClusterStatus_UP,
						ServiceProvider: "aws",
						Region:          "us-west-2",
						Endpoint:        "SASL_SSL://kafka-endpoint",
						ApiEndpoint:     "http://kafka-api-url",
					},
				})
			}
			require.NoError(t, err)
		}
		_, err := io.WriteString(w, string(out))
		require.NoError(t, err)
	}
}

func handleKafkaDedicatedClusterUpdateTest(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var out []byte
		if r.Method == "GET" {
			id := r.URL.Query().Get("id")
			var err error
			out, err = utilv1.MarshalJSONToBytes(&schedv1.GetKafkaClusterReply{
				Cluster: &schedv1.KafkaCluster{
					Id:              id,
					Name:            "lkc-update-dedicated",
					Cku:             1,
					Deployment:      &schedv1.Deployment{Sku: productv1.Sku_DEDICATED},
					NetworkIngress:  50,
					NetworkEgress:   150,
					Storage:         30000,
					Status:          schedv1.ClusterStatus_EXPANDING,
					ServiceProvider: "aws",
					Region:          "us-west-2",
					Endpoint:        "SASL_SSL://kafka-endpoint",
					ApiEndpoint:     "http://kafka-api-url",
				},
			})
			require.NoError(t, err)
		}
		// Update client call
		if r.Method == "PUT" {
			req := &schedv1.UpdateKafkaClusterRequest{}
			err := utilv1.UnmarshalJSON(r.Body, req)
			require.NoError(t, err)
			out, err = utilv1.MarshalJSONToBytes(&schedv1.GetKafkaClusterReply{
				Cluster: &schedv1.KafkaCluster{
					Id:              req.Cluster.Id,
					Name:            req.Cluster.Name,
					Cku:             1,
					PendingCku:      req.Cluster.Cku,
					Deployment:      &schedv1.Deployment{Sku: productv1.Sku_DEDICATED},
					NetworkIngress:  50 * req.Cluster.Cku,
					NetworkEgress:   150 * req.Cluster.Cku,
					Storage:         30000 * req.Cluster.Cku,
					Status:          schedv1.ClusterStatus_EXPANDING,
					ServiceProvider: "aws",
					Region:          "us-west-2",
					Endpoint:        "SASL_SSL://kafka-endpoint",
					ApiEndpoint:     "http://kafka-api-url",
				},
			})
			require.NoError(t, err)
		}
		_, err := io.WriteString(w, string(out))
		require.NoError(t, err)
	}
}

func handleKafkaClusterCreate(t *testing.T, kafkaAPIURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &schedv1.CreateKafkaClusterRequest{}
		err := utilv1.UnmarshalJSON(r.Body, req)
		require.NoError(t, err)
		var b []byte
		if req.Config.Deployment.Sku == productv1.Sku_DEDICATED {
			b, err = utilv1.MarshalJSONToBytes(&schedv1.GetKafkaClusterReply{
				Cluster: &schedv1.KafkaCluster{
					Id:              "lkc-def963",
					AccountId:       req.Config.AccountId,
					Name:            req.Config.Name,
					Cku:             req.Config.Cku,
					Deployment:      &schedv1.Deployment{Sku: productv1.Sku_DEDICATED},
					NetworkIngress:  50 * req.Config.Cku,
					NetworkEgress:   150 * req.Config.Cku,
					Storage:         30000 * req.Config.Cku,
					ServiceProvider: req.Config.ServiceProvider,
					Region:          req.Config.Region,
					Endpoint:        "SASL_SSL://kafka-endpoint",
					ApiEndpoint:     kafkaAPIURL,
				},
			})
		} else {
			b, err = utilv1.MarshalJSONToBytes(&schedv1.GetKafkaClusterReply{
				Cluster: &schedv1.KafkaCluster{
					Id:              "lkc-def963",
					AccountId:       req.Config.AccountId,
					Name:            req.Config.Name,
					Deployment:      &schedv1.Deployment{Sku: productv1.Sku_BASIC},
					NetworkIngress:  100,
					NetworkEgress:   100,
					Storage:         5000,
					ServiceProvider: req.Config.ServiceProvider,
					Region:          req.Config.Region,
					Endpoint:        "SASL_SSL://kafka-endpoint",
					ApiEndpoint:     kafkaAPIURL,
				},
			})
		}
		require.NoError(t, err)
		_, err = io.WriteString(w, string(b))
		require.NoError(t, err)
	}
}

func handleKafkaACLsList(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		results := []*schedv1.ACLBinding{
			{
				Pattern: &schedv1.ResourcePatternConfig{
					ResourceType: schedv1.ResourceTypes_TOPIC,
					Name:         "test-topic",
					PatternType:  schedv1.PatternTypes_LITERAL,
				},
				Entry: &schedv1.AccessControlEntryConfig{
					Operation:      schedv1.ACLOperations_READ,
					PermissionType: schedv1.ACLPermissionTypes_ALLOW,
				},
			},
		}
		reply, err := json.Marshal(results)
		require.NoError(t, err)
		_, err = io.WriteString(w, string(reply))
		require.NoError(t, err)
	}
}

func handleKafkaACLsCreate(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var bindings []*schedv1.ACLBinding
			err := json.NewDecoder(r.Body).Decode(&bindings)
			require.NoError(t, err)
			require.NotEmpty(t, bindings)
			for _, binding := range bindings {
				require.NotEmpty(t, binding.GetPattern())
				require.NotEmpty(t, binding.GetEntry())
			}
		}
	}
}

func handleKafkaACLsDelete(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var filters []*schedv1.ACLFilter
		err := json.NewDecoder(r.Body).Decode(&filters)
		require.NoError(t, err)
		require.NotEmpty(t, filters)
		for _, filter := range filters {
			require.NotEmpty(t, filter.GetEntryFilter())
			require.NotEmpty(t, filter.GetPatternFilter())
		}
	}
}

func handleKSQLCreateList(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ksqlCluster1 := &schedv1.KSQLCluster{
			Id:                "lksqlc-ksql5",
			AccountId:         "25",
			KafkaClusterId:    "lkc-qwert",
			OutputTopicPrefix: "pksqlc-abcde",
			Name:              "account ksql",
			Storage:           101,
			Endpoint:          "SASL_SSL://ksql-endpoint",
		}
		ksqlCluster2 := &schedv1.KSQLCluster{
			Id:                "lksqlc-woooo",
			AccountId:         "25",
			KafkaClusterId:    "lkc-zxcvb",
			OutputTopicPrefix: "pksqlc-ghjkl",
			Name:              "kay cee queue elle",
			Storage:           123,
			Endpoint:          "SASL_SSL://ksql-endpoint",
		}
		if r.Method == "POST" {
			reply, err := utilv1.MarshalJSONToBytes(&schedv1.GetKSQLClusterReply{
				Cluster: ksqlCluster1,
			})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(reply))
			require.NoError(t, err)
		} else if r.Method == "GET" {
			listReply, err := utilv1.MarshalJSONToBytes(&schedv1.GetKSQLClustersReply{
				Clusters: []*schedv1.KSQLCluster{ksqlCluster1, ksqlCluster2},
			})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(listReply))
			require.NoError(t, err)
		}
	}
}

func handleConnect(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			connectorExpansion := &opv1.ConnectorExpansion{
				Id: &opv1.ConnectorId{Id: "lcc-123"},
				Info: &opv1.ConnectorInfo{
					Name:   "az-connector",
					Type:   "Sink",
					Config: map[string]string{},
				},
				Status: &opv1.ConnectorStateInfo{Name: "az-connector", Connector: &opv1.ConnectorState{State: "Running"},
					Tasks: []*opv1.TaskState{{Id: 1, State: "Running"}},
				}}
			listReply, err := json.Marshal(map[string]*opv1.ConnectorExpansion{"lcc-123": connectorExpansion})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(listReply))
			require.NoError(t, err)
		} else if r.Method == "POST" {
			var request opv1.ConnectorInfo
			err := utilv1.UnmarshalJSON(r.Body, &request)
			require.NoError(t, err)
			connector1 := &schedv1.Connector{
				Name:           request.Name,
				KafkaClusterId: "lkc-123",
				AccountId:      "a-595",
				UserConfigs:    request.Config,
				Plugin:         request.Config["connector.class"],
			}
			reply, err := utilv1.MarshalJSONToBytes(connector1)
			require.NoError(t, err)
			_, err = io.WriteString(w, string(reply))
			require.NoError(t, err)
		}
	}
}

func handleConnectorCatalogDescribe(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		configInfos := &opv1.ConfigInfos{
			Name:       "",
			Groups:     nil,
			ErrorCount: 1,
			Configs: []*opv1.Configs{
				{
					Value: &opv1.ConfigValue{
						Name:   "kafka.api.key",
						Errors: []string{"\"kafka.api.key\" is required"},
					},
				},
				{
					Value: &opv1.ConfigValue{
						Name:   "kafka.api.secret",
						Errors: []string{"\"kafka.api.secret\" is required"},
					},
				},
				{
					Value: &opv1.ConfigValue{
						Name:   "topics",
						Errors: []string{"\"topics\" is required"},
					},
				},
				{
					Value: &opv1.ConfigValue{
						Name:   "data.format",
						Errors: []string{"\"data.format\" is required", "Value \"null\" doesn't belong to the property's \"data.format\" enum"},
					},
				},
				{
					Value: &opv1.ConfigValue{
						Name:   "gcs.credentials.config",
						Errors: []string{"\"gcs.credentials.config\" is required"},
					},
				},
				{
					Value: &opv1.ConfigValue{
						Name:   "gcs.bucket.name",
						Errors: []string{"\"gcs.bucket.name\" is required"},
					},
				},
				{
					Value: &opv1.ConfigValue{
						Name:   "time.interval",
						Errors: []string{"\"data.format\" is required", "Value \"null\" doesn't belong to the property's \"time.interval\" enum"},
					},
				},
				{
					Value: &opv1.ConfigValue{
						Name:   "tasks.max",
						Errors: []string{"\"tasks.max\" is required"},
					},
				},
			},
		}
		reply, err := json.Marshal(configInfos)
		require.NoError(t, err)
		_, err = io.WriteString(w, string(reply))
		require.NoError(t, err)
	}
}

func handleConnectPlugins(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			connectorPlugin1 := &opv1.ConnectorPluginInfo{
				Class: "AzureBlobSink",
				Type:  "Sink",
			}
			connectorPlugin2 := &opv1.ConnectorPluginInfo{
				Class: "GcsSink",
				Type:  "Sink",
			}
			listReply, err := json.Marshal([]*opv1.ConnectorPluginInfo{connectorPlugin1, connectorPlugin2})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(listReply))
			require.NoError(t, err)
		}
	}
}

func compose(funcs ...func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, f := range funcs {
			f(w, r)
		}
	}
}

func handleEnvironmentRequests(t *testing.T, id string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, env := range environments {
			if env.Id == id {
				// env found
				if r.Method == "GET" {
					b, err := utilv1.MarshalJSONToBytes(&orgv1.GetAccountReply{Account: env})
					require.NoError(t, err)
					_, err = io.WriteString(w, string(b))
					require.NoError(t, err)
				} else if r.Method == "PUT" {
					req := &orgv1.UpdateAccountRequest{}
					err := utilv1.UnmarshalJSON(r.Body, req)
					require.NoError(t, err)
					env.Name = req.Account.Name
					b, err := utilv1.MarshalJSONToBytes(&orgv1.UpdateAccountReply{Account: env})
					require.NoError(t, err)
					_, err = io.WriteString(w, string(b))
					require.NoError(t, err)
				} else if r.Method == "DELETE" {
					b, err := utilv1.MarshalJSONToBytes(&orgv1.DeleteAccountReply{})
					require.NoError(t, err)
					_, err = io.WriteString(w, string(b))
					require.NoError(t, err)
				}
				return
			}
		}
		// env not found
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleAPIKeyUpdateAndDelete(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		urlSplit := strings.Split(r.URL.Path, "/")
		keyId, err := strconv.Atoi(urlSplit[len(urlSplit)-1])
		require.NoError(t, err)
		index := int32(keyId)
		apiKey := keyStore[index]
		if r.Method == "PUT" {
			req := &schedv1.UpdateApiKeyRequest{}
			err = utilv1.UnmarshalJSON(r.Body, req)
			require.NoError(t, err)
			apiKey.Description = req.ApiKey.Description
			result := &schedv1.UpdateApiKeyReply{
				ApiKey: apiKey,
				Error:  nil,
			}
			b, err := utilv1.MarshalJSONToBytes(result)
			require.NoError(t, err)
			_, err = io.WriteString(w, string(b))
			require.NoError(t, err)
		} else if r.Method == "DELETE" {
			req := &schedv1.DeleteApiKeyRequest{}
			err = utilv1.UnmarshalJSON(r.Body, req)
			require.NoError(t, err)
			delete(keyStore, index)
			result := &schedv1.DeleteApiKeyReply{
				ApiKey: apiKey,
				Error:  nil,
			}
			b, err := utilv1.MarshalJSONToBytes(result)
			require.NoError(t, err)
			_, err = io.WriteString(w, string(b))
			require.NoError(t, err)
		}

	}
}

func handleServiceAccountRequests(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			serviceAccount := &orgv1.User{
				Id:                 serviceAccountID,
				ServiceName:        "service_account",
				ServiceDescription: "at your service.",
			}
			listReply, err := utilv1.MarshalJSONToBytes(&orgv1.GetServiceAccountsReply{
				Users: []*orgv1.User{serviceAccount},
			})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(listReply))
			require.NoError(t, err)
		case "POST":
			req := &orgv1.CreateServiceAccountRequest{}
			err := utilv1.UnmarshalJSON(r.Body, req)
			require.NoError(t, err)
			serviceAccount := &orgv1.User{
				Id:                 55555,
				ServiceName:        req.User.ServiceName,
				ServiceDescription: req.User.ServiceDescription,
			}
			createReply, err := utilv1.MarshalJSONToBytes(&orgv1.CreateServiceAccountReply{
				Error: nil,
				User:  serviceAccount,
			})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(createReply))
			require.NoError(t, err)
		case "PUT":
			req := &orgv1.UpdateServiceAccountRequest{}
			err := utilv1.UnmarshalJSON(r.Body, req)
			require.NoError(t, err)
			updateReply, err := utilv1.MarshalJSONToBytes(&orgv1.UpdateServiceAccountReply{
				Error: nil,
				User:  req.User,
			})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(updateReply))
			require.NoError(t, err)
		case "DELETE":
			req := &orgv1.DeleteServiceAccountRequest{}
			err := utilv1.UnmarshalJSON(r.Body, req)
			require.NoError(t, err)
			updateReply, err := utilv1.MarshalJSONToBytes(&orgv1.DeleteServiceAccountReply{
				Error: nil,
			})
			require.NoError(t, err)
			_, err = io.WriteString(w, string(updateReply))
			require.NoError(t, err)
		}
	}
}
