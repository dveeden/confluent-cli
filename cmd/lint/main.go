// This is a set of "linters" defining the CLI specification
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/client9/gospell"
	"github.com/hashicorp/go-multierror"

	"github.com/confluentinc/cli/internal/cmd"
	linter "github.com/confluentinc/cli/internal/pkg/lint-cli"
	"github.com/confluentinc/cli/internal/pkg/version"
)

var (
	debug   = flag.Bool("debug", false, "print debug output")
	affFile = flag.String("aff-file", "", "hunspell .aff file")
	dicFile = flag.String("dic-file", "", "hunspell .dic file")

	vocab *gospell.GoSpell

	cliNames = []string{"confluent", "ccloud"}

	properNouns = []string{
		"ACL", "ACLs", "API", "Apache", "CCloud CLI", "CLI", "Confluent Cloud", "Confluent Platform", "Confluent",
		"Connect", "Control Center", "Enterprise", "IAM", "ksqlDB Server", "ksqlDB", "Kafka REST", "Kafka", "RBAC",
		"Schema Registry", "ZooKeeper", "ZooKeeper™", "cku",
	}
	vocabWords = []string{
		"ack", "acks", "acl", "acls", "apac", "api", "apikey", "apisecret", "auth", "avro", "aws", "backoff", "ccloud", "cku", "cli", "codec",
		"config", "configs", "connect", "connect-catalog", "consumer.config", "crn", "csu", "decrypt", "deserializer",
		"deserializers", "env", "eu", "formatter", "gcp", "geo", "gzip", "hostname", "html", "https", "iam", "init", "io",
		"json", "jsonschema", "kafka", "ksql", "lifecycle", "lz4", "mds", "multi-zone", "netrc", "pem", "plaintext",
		"producer.config", "protobuf", "rbac", "readwrite", "recv", "rolebinding", "rolebindings", "signup",
		"single-zone", "sr", "sso", "stdin", "systest", "tcp", "tmp", "transactional", "txt", "url", "us", "v2", "vpc",
		"whitelist", "yaml", "zstd",
	}
	utilityCommands = []string{
		"login", "logout", "version", "completion <shell>", "prompt", "update", "init <context-name>", "shell",
	}
	ccloudClusterScopedCommands = []linter.RuleFilter{
		linter.IncludeCommandContains("ccloud kafka acl", "ccloud kafka topic"),
		// only on children of kafka topic commands
		linter.ExcludeCommand("kafka topic"),
		//only on children of kafka acl commands
		linter.ExcludeCommand("kafka acl"),
	}
	confluentClusterScopedCommands = []linter.RuleFilter{
		linter.IncludeCommandContains("confluent kafka topic"),
		linter.ExcludeCommand("kafka topic"),
	}
	resourceScopedCommands = []linter.RuleFilter{
		linter.IncludeCommandContains("api-key use", "api-key create", "api-key store"),
	}
)

var rules = []linter.Rule{
	linter.Filter(
		linter.RequireNamedArgument(
			linter.NamedArgumentConfig{CreateCommandArg: "<name>", OtherCommandsArg: "<id>"},
			map[string]linter.NamedArgumentConfig{
				"environment":    {CreateCommandArg: "<name>", OtherCommandsArg: "<environment-id>"},
				"role":           {CreateCommandArg: "<name>", OtherCommandsArg: "<name>"},
				"topic":          {CreateCommandArg: "<topic>", OtherCommandsArg: "<topic>"},
				"api-key":        {CreateCommandArg: "N/A", OtherCommandsArg: "<api-key>"},
				"consumer-group": {CreateCommandArg: "<consumer-group>", OtherCommandsArg: "<consumer-group>"},
				"lag":            {CreateCommandArg: "N/A", OtherCommandsArg: "<consumer-group>"},
			},
		),
		linter.OnlyLeafCommands, linter.ExcludeCommand(utilityCommands...),
		// skip resource container commands
		linter.ExcludeUse("list", "auth"),
		// skip ACLs which don't have an identity (value objects rather than entities)
		linter.ExcludeCommandContains("kafka acl"),
		linter.ExcludeCommandContains("iam acl"),
		// skip api-key create since you don't get to choose a name for API keys
		linter.ExcludeCommandContains("api-key create"),
		// skip connector create since you don't get to choose id for connector
		linter.ExcludeCommandContains("connector create"),
		// skip local which delegates to external bash scripts
		linter.ExcludeCommandContains("local"),
		// skip for api-key store command since KEY is not last argument
		linter.ExcludeCommand("api-key store <api-key> <secret>"),
		// skip for rolebindings since they don't have names/IDs
		linter.ExcludeCommandContains("iam rolebinding"),
		// skip for register command since they don't have names/IDs
		linter.ExcludeCommandContains("cluster register"),
		// skip for unregister command since they don't have names/IDs
		linter.ExcludeCommandContains("cluster unregister"),
		// skip secret commands
		linter.ExcludeCommandContains("secret"),
		// skip schema-registry commands which do not use names/ID's
		linter.ExcludeCommandContains("schema-registry"),
		// skip ksql configure-acls command as it can take any number of topic arguments
		linter.ExcludeCommandContains("ksql app configure-acls"),
		// skip cluster describe as it takes a URL as a flag instead of a resource identity
		linter.ExcludeCommandContains("cluster describe"),
		// skip connector-catalog describe as it connector plugin name
		linter.ExcludeCommandContains("connector-catalog describe"),
		// skip connector event describe as it shows connector log events configuration for an org
		linter.ExcludeCommandContains("connector event describe"),
		// skip feedback command
		linter.ExcludeCommand("feedback"),
		// skip signup command
		linter.ExcludeCommandContains("signup"),
		// config context commands
		linter.ExcludeCommand("config context current"),
		linter.ExcludeCommandContains("config context get"),
		linter.ExcludeCommandContains("config context set"),
		linter.ExcludeCommandContains("audit-log"),
		// skip admin commands since they have two args
		linter.ExcludeCommandContains("admin"),
		// skip cluster linking commands
		linter.ExcludeCommandContains("kafka link"),
		linter.ExcludeCommandContains("kafka mirror"),
	),
	// TODO: ensuring --cluster is optional DOES NOT actually ensure that the cluster context is used
	linter.Filter(linter.RequireFlag("cluster", true), ccloudClusterScopedCommands...),
	linter.Filter(linter.RequireFlagType("cluster", "string"), ccloudClusterScopedCommands...),
	linter.Filter(linter.RequireFlagDescription("cluster", "Kafka cluster ID."), ccloudClusterScopedCommands...),
	// Require on-prem kafka topic commands to have required --url flag to specify rest API endpoint.
	linter.Filter(linter.RequireFlag("url", true), confluentClusterScopedCommands...),
	linter.Filter(linter.RequireFlagType("url", "string"), confluentClusterScopedCommands...),
	linter.Filter(linter.RequireFlagDescription("url", "Base URL of REST Proxy Endpoint of Kafka Cluster (include /kafka for embedded Rest Proxy). Must set flag or CONFLUENT_REST_URL."), confluentClusterScopedCommands...),
	linter.Filter(linter.RequireFlag("resource", false), resourceScopedCommands...),
	linter.Filter(linter.RequireFlag("resource", true), linter.IncludeCommandContains("api-key list")),
	linter.Filter(linter.RequireFlagType("resource", "string"), resourceScopedCommands...),
	linter.Filter(linter.RequireFlagType("resource", "string"), linter.IncludeCommandContains("api-key list")),
	linter.Filter(
		linter.RequireFlagSort(false),
		linter.OnlyLeafCommands,
		linter.ExcludeCommandContains("local"),
	),
	linter.RequireLowerCase("Use"),
	linter.Filter(
		linter.RequireSingular("Use"),
		linter.ExcludeCommandContains("local"),
	),
	linter.Filter(
		linter.RequireLengthBetween("Short", 13, 60),
		linter.ExcludeCommandContains("secret", "mirror"),
	),
	linter.RequireStartWithCapital("Short"),
	linter.RequireEndWithPunctuation("Short", false),
	linter.RequireCapitalizeProperNouns("Short", linter.SetDifferenceIgnoresCase(properNouns, cliNames)),
	linter.RequireStartWithCapital("Long"),
	linter.RequireEndWithPunctuation("Long", true),
	linter.RequireCapitalizeProperNouns("Long", linter.SetDifferenceIgnoresCase(properNouns, cliNames)),
	linter.Filter(
		linter.RequireNotTitleCase("Short", properNouns),
		linter.ExcludeCommandContains("secret", "mirror"),
	),
	linter.Filter(
		linter.RequireRealWords("Use", '-'),
		linter.ExcludeCommandContains("unregister"),
		linter.ExcludeCommandContains("audit-log", "failover"),
	),
}

var flagRules = []linter.FlagRule{
	linter.FlagFilter(
		linter.RequireFlagNameLength(2, 16),
		linter.ExcludeFlag(
			"compression-codec", "connect-cluster-id", "consumer-property", "enable-systest-events",
			"local-secrets-file", "max-partition-memory-bytes", "message-send-max-retries", "metadata-expiry-ms",
			"producer-property", "remote-secrets-file", "replication-factor", "request-required-acks", "request-timeout-ms", // TODO: change back if replication-factor is too long
			"schema-registry-cluster-id", "service-account", "skip-message-on-error", "socket-buffer-size",
			"value-deserializer", "bootstrap-servers", "source-bootstrap-server", "source-cluster-id", "source-api-secret",
		),
	),
	linter.FlagFilter(
		linter.RequireFlagUsageMessage,
		linter.ExcludeFlag("key-deserializer", "value-deserializer"),
	),
	linter.FlagFilter(
		linter.RequireFlagUsageStartWithCapital,
		linter.ExcludeFlag("ksql-cluster-id"),
	),
	linter.FlagFilter(
		linter.RequireFlagUsageEndWithPunctuation,
		linter.ExcludeFlag(
			"batch-size", "enable-systest-events", "formatter", "isolation-level", "line-reader", "max-block-ms",
			"max-memory-bytes", "max-partition-memory-bytes", "message-send-max-retries", "metadata-expiry-ms",
			"offset", "property", "request-required-acks", "request-timeout-ms", "retry-backoff-ms",
			"socket-buffer-size", "timeout",
		),
	),
	linter.RequireFlagKebabCase,
	linter.FlagFilter(
		linter.RequireFlagCharacters('-'),
		linter.ExcludeFlag("consumer.config", "producer.config"),
	),
	linter.FlagFilter(
		linter.RequireFlagDelimiter('-', 1),
		linter.ExcludeFlag(
			"ca-cert-path", "connect-cluster-id", "enable-systest-events", "if-not-exists", "kafka-cluster-id",
			"ksql-cluster-id", "local-secrets-file", "max-block-ms", "max-memory-bytes", "max-partition-memory-bytes",
			"message-send-max-retries", "metadata-expiry-ms", "remote-secrets-file", "request-required-acks",
			"request-timeout-ms", "retry-backoff-ms", "schema-registry-cluster-id", "service-account",
			"skip-message-on-error", "socket-buffer-size", "client-cert-path", "client-key-path",
			"source-bootstrap-server", "source-cluster-id", "source-api-secret", "source-api-key",
		),
	),
	linter.RequireFlagRealWords('-'),
	linter.FlagFilter(
		linter.RequireFlagUsageRealWords,
		linter.ExcludeFlag("source-api-key", "source-api-secret"),
	),
}

func main() {
	flag.Parse()

	var err error
	vocab, err = gospell.NewGoSpell(*affFile, *dicFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, w := range vocabWords {
		vocab.AddWordRaw(strings.ToLower(w))
		vocab.AddWordRaw(strings.ToUpper(w))
	}
	linter.SetVocab(vocab)

	l := linter.Linter{
		Rules:     rules,
		FlagRules: flagRules,
		Vocab:     vocab,
		Debug:     *debug,
	}

	var issues *multierror.Error
	for _, cliName := range cliNames {
		cli := cmd.NewConfluentCommand(cliName, true, &version.Version{Binary: cliName})
		if err := l.Lint(cli.Command); err != nil {
			issues = multierror.Append(issues, err)
		}
	}
	if issues.ErrorOrNil() != nil {
		fmt.Println(issues)
		os.Exit(1)
	}
}
