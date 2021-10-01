package iam

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"github.com/confluentinc/go-printer"
	"net/http"
	"os"
	"sort"
	"strings"

	orgv1 "github.com/confluentinc/cc-structs/kafka/org/v1"

	mds "github.com/confluentinc/mds-sdk-go/mdsv1"
	"github.com/confluentinc/mds-sdk-go/mdsv2alpha1"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/examples"
	"github.com/confluentinc/cli/internal/pkg/output"
)

var (
	resourcePatternListFields           = []string{"Principal", "Role", "ResourceType", "Name", "PatternType"}
	resourcePatternHumanListLabels      = []string{"Principal", "Role", "ResourceType", "Name", "PatternType"}
	resourcePatternStructuredListLabels = []string{"principal", "role", "resource_type", "name", "pattern_type"}

	// ccloud has Email as additional field
	ccloudResourcePatternListFields           = []string{"Principal", "Email", "Role", "Environment", "CloudCluster", "ClusterType", "LogicalCluster", "ResourceType", "Name", "PatternType"}
	ccloudResourcePatternHumanListLabels      = []string{"Principal", "Email", "Role", "Environment", "Cloud Cluster", "Cluster Type", "Logical Cluster", "Resource Type", "Name", "Pattern Type"}
	ccloudResourcePatternStructuredListLabels = []string{"principal", "email", "role", "environment", "cloud_cluster", "cluster_type", "logical_cluster", "resource_type", "resource_name", "pattern_type"}

	//TODO: please move this to a backend route (https://confluentinc.atlassian.net/browse/CIAM-890)
	clusterScopedRoles = map[string]bool{
		"SystemAdmin":   true,
		"ClusterAdmin":  true,
		"SecurityAdmin": true,
		"UserAdmin":     true,
		"Operator":      true,
	}

	clusterScopedRolesV2 = map[string]bool{
		"CloudClusterAdmin": true,
	}

	environmentScopedRoles = map[string]bool{
		"EnvironmentAdmin": true,
	}

	dataplaneNamespace =  optional.NewString("dataplane")
)

type rolebindingOptions struct {
	role               string
	resource           string
	prefix             bool
	principal          string
	scopeV2            mdsv2alpha1.Scope
	mdsScope           mds.MdsScope
	resourcesRequest   mds.ResourcesRequest
	resourcesRequestV2 mdsv2alpha1.ResourcesRequest
}

type rolebindingCommand struct {
	*cmd.AuthenticatedStateFlagCommand
	cliName string
	ccloudRbacDataplaneEnabled bool
}

type listDisplay struct {
	Principal      string `json:"principal"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Environment    string `json:"environment"`
	CloudCluster   string `json:"cloud_cluster"`
	ClusterType    string `json:"cluster_type"`
	LogicalCluster string `json:"logical_cluster"`
	ResourceType   string `json:"resource_type"`
	Name           string `json:"resource_name"`
	PatternType    string `json:"pattern_type"`
}

// NewRolebindingCommand returns the sub-command object for interacting with RBAC rolebindings.
func NewRolebindingCommand(cliName string, prerunner cmd.PreRunner) *cobra.Command {
	cobraRolebindingCmd := &cobra.Command{
		Use:   "rolebinding",
		Short: "Manage RBAC and IAM role bindings.",
		Long:  "Manage Role-Based Access Control (RBAC) and Identity and Access Management (IAM) role bindings.",
	}
	var cliCmd *cmd.AuthenticatedStateFlagCommand
	if cliName == "confluent" {
		cliCmd = cmd.NewAuthenticatedWithMDSStateFlagCommand(cobraRolebindingCmd, prerunner, RolebindingSubcommandFlags)
	} else {
		cliCmd = cmd.NewAuthenticatedStateFlagCommand(cobraRolebindingCmd, prerunner, nil)
	}
	ccloudRbacDataplaneEnabled := false
	if os.Getenv("XX_CCLOUD_RBAC_DATAPLANE") != "" {
		ccloudRbacDataplaneEnabled = true
	}
	roleBindingCmd := &rolebindingCommand{
		AuthenticatedStateFlagCommand: cliCmd,
		cliName:                       cliName,
		ccloudRbacDataplaneEnabled:    ccloudRbacDataplaneEnabled,
	}
	roleBindingCmd.init()
	return roleBindingCmd.Command
}

func (c *rolebindingCommand) init() {
	var example string
	if c.cliName == "ccloud" {
		example = examples.BuildExampleString(
			examples.Example{
				Text: "To list the role bindings for current user:",
				Code: "ccloud iam rolebinding list --current-user",
			},
			examples.Example{
				Text: "To list the role bindings for a specific principal:",
				Code: "ccloud iam rolebinding list --principal User:frodo",
			},
			examples.Example{
				Text: "To list the role bindings for a specific principal, filtered to a specific role:",
				Code: "ccloud iam rolebinding list --principal User:frodo --role CloudClusterAdmin --environment env-123 --cloud-cluster lkc-1111aaa",
			},
			examples.Example{
				Text: "To list the principals bound to a specific role:",
				Code: "ccloud iam rolebinding list --role CloudClusterAdmin --current-env --cloud-cluster lkc-1111aaa",
			},
		)
	} else {
		example = examples.BuildExampleString(
			examples.Example{
				Text: "Only use the `--resource` flag when specifying a `--role` with no `--principal` specified. If specifying a `--principal`, then the `--resource` flag is ignored. To list role bindings for a specific role on an identified resource:",
				Code: "confluent iam rolebinding list --kafka-cluster-id CID  --role DeveloperRead --resource Topic",
			},
			examples.Example{
				Text: "To list the role bindings for a specific principal:",
				Code: "confluent iam rolebinding list --kafka-cluster-id $CID --principal User:frodo",
			},
			examples.Example{
				Text: "To list the role bindings for a specific principal, filtered to a specific role:",
				Code: "confluent iam rolebinding list --kafka-cluster-id $CID --principal User:frodo --role DeveloperRead",
			},
			examples.Example{
				Text: "To list the principals bound to a specific role:",
				Code: "confluent iam rolebinding list --kafka-cluster-id $CID --role DeveloperWrite",
			},
			examples.Example{
				Text: "To list the principals bound to a specific resource with a specific role:",
				Code: "confluent iam rolebinding list --kafka-cluster-id $CID --role DeveloperWrite --resource Topic:shire-parties",
			},
		)

	}

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List role bindings.",
		Long:    "List the role bindings for a particular principal and/or role, and a particular scope.",
		Args:    cobra.NoArgs,
		RunE:    cmd.NewCLIRunE(c.list),
		Example: example,
	}
	listCmd.Flags().String("principal", "", "Principal whose role bindings should be listed.")
	listCmd.Flags().Bool("current-user", false, "Show role bindings belonging to current user.")
	listCmd.Flags().String("role", "", "List role bindings under a specific role given to a principal. Or if no principal is specified, list principals with the role.")
	if c.cliName == "ccloud" {
		listCmd.Flags().String("cloud-cluster", "", "Cloud cluster ID for scope of role binding listings.")
		listCmd.Flags().String("environment", "", "Environment ID for scope of role binding listings.")
		listCmd.Flags().Bool("current-env", false, "Use current environment ID for scope.")
		if c.ccloudRbacDataplaneEnabled {
			listCmd.Flags().String("resource", "", "If specified with a role and no principals, list principals with role bindings to the role for this qualified resource.")
			listCmd.Flags().String("kafka-cluster-id", "", "Kafka cluster ID for scope of role binding listings.")
		}
	} else {
		listCmd.Flags().String("resource", "", "If specified with a role and no principals, list principals with role bindings to the role for this qualified resource.")
		listCmd.Flags().String("kafka-cluster-id", "", "Kafka cluster ID for scope of role binding listings.")
		listCmd.Flags().String("schema-registry-cluster-id", "", "Schema Registry cluster ID for scope of role binding listings.")
		listCmd.Flags().String("ksql-cluster-id", "", "ksqlDB cluster ID for scope of role binding listings.")
		listCmd.Flags().String("connect-cluster-id", "", "Kafka Connect cluster ID for scope of role binding listings.")
		listCmd.Flags().String("cluster-name", "", "Cluster name to uniquely identify the cluster for role binding listings.")
	}
	listCmd.Flags().StringP(output.FlagName, output.ShortHandFlag, output.DefaultValue, output.Usage)
	listCmd.Flags().SortFlags = false

	c.AddCommand(listCmd)

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a role binding.",
		Args:  cobra.NoArgs,
		RunE:  cmd.NewCLIRunE(c.create),
		Example: examples.BuildExampleString(
			examples.Example{
				Text: "Create a role binding for the client permitting it produce to the topic users.",
				Code: c.cliName + " iam rolebinding create --principal User:appSA --role DeveloperWrite --resource Topic:users --kafka-cluster-id $KAFKA_CLUSTER_ID",
			},
		),
	}
	createCmd.Flags().String("role", "", "Role name of the new role binding.")
	createCmd.Flags().String("principal", "", "Qualified principal name for the role binding.")
	if c.cliName == "ccloud" {
		createCmd.Flags().String("cloud-cluster", "", "Cloud cluster ID for the role binding.")
		createCmd.Flags().String("environment", "", "Environment ID for scope of rolebinding create.")
		createCmd.Flags().Bool("current-env", false, "Use current environment ID for scope.")
		if c.ccloudRbacDataplaneEnabled{
			createCmd.Flags().Bool("prefix", false, "Whether the provided resource name is treated as a prefix pattern.")
			createCmd.Flags().String("resource", "", "Qualified resource name for the role binding.")
			createCmd.Flags().String("kafka-cluster-id", "", "Kafka cluster ID for the role binding.")
		}
	} else {
		createCmd.Flags().Bool("prefix", false, "Whether the provided resource name is treated as a prefix pattern.")
		createCmd.Flags().String("resource", "", "Qualified resource name for the role binding.")
		createCmd.Flags().String("kafka-cluster-id", "", "Kafka cluster ID for the role binding.")
		createCmd.Flags().String("schema-registry-cluster-id", "", "Schema Registry cluster ID for the role binding.")
		createCmd.Flags().String("ksql-cluster-id", "", "ksqlDB cluster ID for the role binding.")
		createCmd.Flags().String("connect-cluster-id", "", "Kafka Connect cluster ID for the role binding.")
		createCmd.Flags().String("cluster-name", "", "Cluster name to uniquely identify the cluster for rolebinding listings.")
	}
	createCmd.Flags().StringP(output.FlagName, output.ShortHandFlag, output.DefaultValue, output.Usage)
	createCmd.Flags().SortFlags = false
	check(createCmd.MarkFlagRequired("role"))
	check(createCmd.MarkFlagRequired("principal"))
	c.AddCommand(createCmd)

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an existing role binding.",
		Args:  cobra.NoArgs,
		RunE:  cmd.NewCLIRunE(c.delete),
	}
	deleteCmd.Flags().String("role", "", "Role name of the existing role binding.")
	deleteCmd.Flags().String("principal", "", "Qualified principal name associated with the role binding.")
	if c.cliName == "ccloud" {
		deleteCmd.Flags().String("cloud-cluster", "", "Cloud cluster ID for the role binding.")
		deleteCmd.Flags().String("environment", "", "Environment ID for scope of rolebinding delete.")
		deleteCmd.Flags().Bool("current-env", false, "Use current environment ID for scope.")
		if c.ccloudRbacDataplaneEnabled {
			deleteCmd.Flags().Bool("prefix", false, "Whether the provided resource name is treated as a prefix pattern.")
			deleteCmd.Flags().String("resource", "", "Qualified resource name for the role binding.")
			deleteCmd.Flags().String("kafka-cluster-id", "", "Kafka cluster ID for the role binding.")
		}
	} else {
		deleteCmd.Flags().Bool("prefix", false, "Whether the provided resource name is treated as a prefix pattern.")
		deleteCmd.Flags().String("resource", "", "Qualified resource name for the role binding.")
		deleteCmd.Flags().String("kafka-cluster-id", "", "Kafka cluster ID for the role binding.")
		deleteCmd.Flags().String("schema-registry-cluster-id", "", "Schema Registry cluster ID for the role binding.")
		deleteCmd.Flags().String("ksql-cluster-id", "", "ksqlDB cluster ID for the role binding.")
		deleteCmd.Flags().String("connect-cluster-id", "", "Kafka Connect cluster ID for the role binding.")
		deleteCmd.Flags().String("cluster-name", "", "Cluster name to uniquely identify the cluster for rolebinding listings.")
	}
	deleteCmd.Flags().StringP(output.FlagName, output.ShortHandFlag, output.DefaultValue, output.Usage)
	deleteCmd.Flags().SortFlags = false
	check(createCmd.MarkFlagRequired("role"))
	check(deleteCmd.MarkFlagRequired("principal"))
	check(deleteCmd.MarkFlagRequired("role"))
	c.AddCommand(deleteCmd)
}

func (c *rolebindingCommand) validatePrincipalFormat(principal string) error {
	if len(strings.Split(principal, ":")) == 1 {
		return errors.NewErrorWithSuggestions(errors.PrincipalFormatErrorMsg, errors.PrincipalFormatSuggestions)
	}

	return nil
}

func (c *rolebindingCommand) parseAndValidateResourcePattern(typename string, prefix bool) (mds.ResourcePattern, error) {
	var result mds.ResourcePattern
	if prefix {
		result.PatternType = "PREFIXED"
	} else {
		result.PatternType = "LITERAL"
	}

	parts := strings.Split(typename, ":")
	if len(parts) != 2 {
		return result, errors.NewErrorWithSuggestions(errors.ResourceFormatErrorMsg, errors.ResourceFormatSuggestions)
	}
	result.ResourceType = parts[0]
	result.Name = parts[1]

	return result, nil
}

func (c *rolebindingCommand) parseAndValidateResourcePatternV2(typename string, prefix bool) (mdsv2alpha1.ResourcePattern, error) {
	var result mdsv2alpha1.ResourcePattern
	if prefix {
		result.PatternType = "PREFIXED"
	} else {
		result.PatternType = "LITERAL"
	}

	parts := strings.Split(typename, ":")
	if len(parts) != 2 {
		return result, errors.NewErrorWithSuggestions(errors.ResourceFormatErrorMsg, errors.ResourceFormatSuggestions)
	}
	result.ResourceType = parts[0]
	result.Name = parts[1]

	return result, nil
}

func (c *rolebindingCommand) validateRoleAndResourceTypeV1(roleName string, resourceType string) error {
	ctx := c.createContext()
	role, resp, err := c.MDSClient.RBACRoleDefinitionsApi.RoleDetail(ctx, roleName)
	if err != nil || resp.StatusCode == 204 {
		if err == nil {
			return errors.NewErrorWithSuggestions(fmt.Sprintf(errors.LookUpRoleErrorMsg, roleName), errors.LookUpRoleSuggestions)
		} else {
			return errors.NewWrapErrorWithSuggestions(err, fmt.Sprintf(errors.LookUpRoleErrorMsg, roleName), errors.LookUpRoleSuggestions)
		}
	}

	var allResourceTypes []string
	found := false
	for _, operation := range role.AccessPolicy.AllowedOperations {
		allResourceTypes = append(allResourceTypes, operation.ResourceType)
		if operation.ResourceType == resourceType {
			found = true
			break
		}
	}

	if !found {
		suggestionsMsg := fmt.Sprintf(errors.InvalidResourceTypeSuggestions, strings.Join(allResourceTypes, ", "))
		return errors.NewErrorWithSuggestions(fmt.Sprintf(errors.InvalidResourceTypeErrorMsg, resourceType), suggestionsMsg)
	}

	return nil
}

func (c *rolebindingCommand) validateResourceTypeV1(resourceType string) error {
	ctx := c.createContext()
	roles, _, err := c.MDSClient.RBACRoleDefinitionsApi.Roles(ctx)
	if err != nil {
		return err
	}

	var allResourceTypes = make(map[string]bool)
	found := false
	for _, role := range roles {
		for _, operation := range role.AccessPolicy.AllowedOperations {
			allResourceTypes[operation.ResourceType] = true
			if operation.ResourceType == resourceType {
				found = true
				break
			}
		}
	}

	if !found {
		uniqueResourceTypes := []string{}
		for rt := range allResourceTypes {
			uniqueResourceTypes = append(uniqueResourceTypes, rt)
		}
		suggestionsMsg := fmt.Sprintf(errors.InvalidResourceTypeSuggestions, strings.Join(uniqueResourceTypes, ", "))
		return errors.NewErrorWithSuggestions(fmt.Sprintf(errors.InvalidResourceTypeErrorMsg, resourceType), suggestionsMsg)
	}

	return nil
}

func (c *rolebindingCommand) validateRoleAndResourceTypeV2(roleName string, resourceType string) error {
	ctx := c.createContext()
	roleDetail := mdsv2alpha1.RoleDetailOpts{}
	if c.ccloudRbacDataplaneEnabled {
		roleDetail.Namespace = dataplaneNamespace
	}
    // Currently we don't allow multiple namespace in roleDetail so as a workaround we first check with dataplane
    // namespace and if we get an error try without any namespace.
	role, resp, err := c.MDSv2Client.RBACRoleDefinitionsApi.RoleDetail(ctx, roleName, &roleDetail)
	if err != nil || resp.StatusCode == http.StatusNoContent {
		role, resp, err = c.MDSv2Client.RBACRoleDefinitionsApi.RoleDetail(ctx, roleName, nil)
		if err != nil || resp.StatusCode == http.StatusNoContent {
			if err == nil {
				return errors.NewErrorWithSuggestions(fmt.Sprintf(errors.LookUpRoleErrorMsg, roleName), errors.LookUpRoleSuggestions)
			} else {
				return errors.NewWrapErrorWithSuggestions(err, fmt.Sprintf(errors.LookUpRoleErrorMsg, roleName), errors.LookUpRoleSuggestions)
			}
		}
	}

	var allResourceTypes []string
	for _, policies := range role.Policies {
		for _, operation := range policies.AllowedOperations {
			allResourceTypes = append(allResourceTypes, operation.ResourceType)
			if operation.ResourceType == resourceType {
				return nil
			}
		}
	}

	suggestionsMsg := fmt.Sprintf(errors.InvalidResourceTypeSuggestions, strings.Join(allResourceTypes, ", "))
	return errors.NewErrorWithSuggestions(fmt.Sprintf(errors.InvalidResourceTypeErrorMsg, resourceType), suggestionsMsg)
}

func (c *rolebindingCommand) validateResourceTypeV2(resourceType string) error {
	ctx := c.createContext()
	roles, _, err := c.MDSv2Client.RBACRoleDefinitionsApi.Roles(ctx, nil)
	if err != nil {
		return err
	}

	var allResourceTypes = make(map[string]bool)
	found := false
	for _, role := range roles {
		for _, policies := range role.Policies {
			for _, operation := range policies.AllowedOperations {
				allResourceTypes[operation.ResourceType] = true
				if operation.ResourceType == resourceType {
					found = true
					break
				}
			}
		}
	}

	if !found {
		uniqueResourceTypes := []string{}
		for rt := range allResourceTypes {
			uniqueResourceTypes = append(uniqueResourceTypes, rt)
		}
		suggestionsMsg := fmt.Sprintf(errors.InvalidResourceTypeSuggestions, strings.Join(uniqueResourceTypes, ", "))
		return errors.NewErrorWithSuggestions(fmt.Sprintf(errors.InvalidResourceTypeErrorMsg, resourceType), suggestionsMsg)
	}

	return nil
}

func (c *rolebindingCommand) parseAndValidateScope(cmd *cobra.Command) (*mds.MdsScope, error) {
	scope := &mds.MdsScopeClusters{}
	nonKafkaScopesSet := 0

	clusterName, err := cmd.Flags().GetString("cluster-name")
	if err != nil {
		return nil, err
	}

	cmd.Flags().Visit(func(flag *pflag.Flag) {
		switch flag.Name {
		case "kafka-cluster-id":
			scope.KafkaCluster = flag.Value.String()
		case "schema-registry-cluster-id":
			scope.SchemaRegistryCluster = flag.Value.String()
			nonKafkaScopesSet++
		case "ksql-cluster-id":
			scope.KsqlCluster = flag.Value.String()
			nonKafkaScopesSet++
		case "connect-cluster-id":
			scope.ConnectCluster = flag.Value.String()
			nonKafkaScopesSet++
		}
	})

	if clusterName != "" && (scope.KafkaCluster != "" || nonKafkaScopesSet > 0) {
		return nil, errors.New(errors.BothClusterNameAndScopeErrorMsg)
	}

	if clusterName == "" {
		if scope.KafkaCluster == "" && nonKafkaScopesSet > 0 {
			return nil, errors.New(errors.SpecifyKafkaIDErrorMsg)
		}

		if scope.KafkaCluster == "" && nonKafkaScopesSet == 0 {
			return nil, errors.New(errors.SpecifyClusterErrorMsg)
		}

		if nonKafkaScopesSet > 1 {
			return nil, errors.New(errors.MoreThanOneNonKafkaErrorMsg)
		}
		return &mds.MdsScope{Clusters: *scope}, nil
	}

	return &mds.MdsScope{ClusterName: clusterName}, nil
}

func (c *rolebindingCommand) parseAndValidateScopeV2(cmd *cobra.Command) (*mdsv2alpha1.Scope, error) {
	scopeV2 := &mdsv2alpha1.Scope{}
	orgResourceId := c.State.Auth.Organization.GetResourceId()
	scopeV2.Path = []string{"organization=" + orgResourceId}

	if cmd.Flags().Changed("current-env") {
		scopeV2.Path = append(scopeV2.Path, "environment="+c.EnvironmentId())
	} else if cmd.Flags().Changed("environment") {
		env, err := cmd.Flags().GetString("environment")
		if err != nil {
			return nil, err
		}
		scopeV2.Path = append(scopeV2.Path, "environment="+env)
	}

	if cmd.Flags().Changed("cloud-cluster") {
		cluster, err := cmd.Flags().GetString("cloud-cluster")
		if err != nil {
			return nil, err
		}
		scopeV2.Path = append(scopeV2.Path, "cloud-cluster="+cluster)
	}

	if c.ccloudRbacDataplaneEnabled && cmd.Flags().Changed("kafka-cluster-id") {
		kafkaCluster, err := cmd.Flags().GetString("kafka-cluster-id")
		if err != nil {
			return nil, err
		}
		scopeV2.Clusters.KafkaCluster = kafkaCluster
	}

	if cmd.Flags().Changed("role") {
		role, err := cmd.Flags().GetString("role")
		if err != nil {
			return nil, err
		}
		if clusterScopedRolesV2[role] && !cmd.Flags().Changed("cloud-cluster") {
			return nil, errors.New(errors.SpecifyCloudClusterErrorMsg)
		}
		if (environmentScopedRoles[role] || clusterScopedRolesV2[role]) && !cmd.Flags().Changed("current-env") && !cmd.Flags().Changed("environment") {
			return nil, errors.New(errors.SpecifyEnvironmentErrorMsg)
		}
	}

	if cmd.Flags().Changed("cloud-cluster") && !cmd.Flags().Changed("current-env") && !cmd.Flags().Changed("environment") {
		return nil, errors.New(errors.SpecifyCloudClusterErrorMsg)
	}
	return scopeV2, nil
}

func (c *rolebindingCommand) confluentList(cmd *cobra.Command, options *rolebindingOptions) error {
	if cmd.Flags().Changed("principal") {
		return c.listPrincipalResources(cmd, options)
	} else if cmd.Flags().Changed("role") {
		return c.confluentListRolePrincipals(cmd, options)
	}
	return errors.New(errors.PrincipalOrRoleRequiredErrorMsg)
}

func (c *rolebindingCommand) listMyRoleBindings(cmd *cobra.Command, options *rolebindingOptions) error {
	scopeV2 := &options.scopeV2
	var principal string
	currentUser, err := cmd.Flags().GetBool("current-user")
	if err != nil {
		return err
	}
	if currentUser {
		principal = "User:" + c.State.Auth.User.ResourceId
	} else {
		principal = options.principal
	}
	scopedRoleBindingMappings, _, err := c.MDSv2Client.RBACRoleBindingSummariesApi.MyRoleBindings(
		c.createContext(),
		principal,
		*scopeV2)
	if err != nil {
		return err
	}

	userToEmailMap, err := c.userIdToEmailMap()
	if err != nil {
		return err
	}

	outputWriter, err := output.NewListOutputWriter(cmd, ccloudResourcePatternListFields, ccloudResourcePatternHumanListLabels, ccloudResourcePatternStructuredListLabels)
	if err != nil {
		return err
	}

	role, err := cmd.Flags().GetString("role")
	if err != nil {
		return err
	}

	for _, scopedRoleBindingMapping := range scopedRoleBindingMappings {
		roleBindingScope := scopedRoleBindingMapping.Scope
		for principalName, roleBindings := range scopedRoleBindingMapping.Rolebindings {
			principalEmail := userToEmailMap[principalName]
			for roleName, resourcePatterns := range roleBindings {
				if role != "" && role != roleName {
					continue
				}

				envName := ""
				cloudClusterName := ""
				for _, elem := range roleBindingScope.Path {
					// we don't capture the organization name because it's always this organization
					if strings.HasPrefix(elem, "environment=") {
						envName = strings.TrimPrefix(elem, "environment=")
					}
					if strings.HasPrefix(elem, "cloud-cluster=") {
						cloudClusterName = strings.TrimPrefix(elem, "cloud-cluster=")
					}
				}
				clusterType := ""
				logicalCluster := ""
				if roleBindingScope.Clusters.ConnectCluster != "" {
					clusterType = "Connect"
					logicalCluster = roleBindingScope.Clusters.ConnectCluster
				} else if roleBindingScope.Clusters.KsqlCluster != "" {
					clusterType = "ksqlDB"
					logicalCluster = roleBindingScope.Clusters.KsqlCluster
				} else if roleBindingScope.Clusters.SchemaRegistryCluster != "" {
					clusterType = "Schema Registry"
					logicalCluster = roleBindingScope.Clusters.SchemaRegistryCluster
				} else if roleBindingScope.Clusters.KafkaCluster != "" {
					clusterType = "Kafka"
					logicalCluster = roleBindingScope.Clusters.KafkaCluster
				}

				for _, resourcePattern := range resourcePatterns {
					if cmd.Flags().Changed("resource") {
						resource, err := cmd.Flags().GetString("resource")
						if err != nil {
							return err
						}
						if resource != resourcePattern.ResourceType {
							continue
						}
					}
					outputWriter.AddElement(&listDisplay{
						Principal:      principalName,
						Email:          principalEmail,
						Role:           roleName,
						Environment:    envName,
						CloudCluster:   cloudClusterName,
						ClusterType:    clusterType,
						LogicalCluster: logicalCluster,
						ResourceType:   resourcePattern.ResourceType,
						Name:           resourcePattern.Name,
						PatternType:    resourcePattern.PatternType,
					})
				}

				if len(resourcePatterns) == 0 {
					outputWriter.AddElement(&listDisplay{
						Principal:    principalName,
						Email:        principalEmail,
						Role:         roleName,
						Environment:  envName,
						CloudCluster: cloudClusterName,
					})
				}
			}
		}
	}

	outputWriter.StableSort()

	return outputWriter.Out()
}

func (c *rolebindingCommand) ccloudList(cmd *cobra.Command, options *rolebindingOptions) error {
	if cmd.Flags().Changed("principal") || cmd.Flags().Changed("current-user") {
		return c.listMyRoleBindings(cmd, options)
	} else if cmd.Flags().Changed("role") {
		return c.ccloudListRolePrincipals(cmd, options)
	} else {
		return errors.New(errors.PrincipalOrRoleRequiredErrorMsg)
	}
}

func (c *rolebindingCommand) list(cmd *cobra.Command, _ []string) error {
	options, err := c.parseCommon(cmd)
	if err != nil {
		return err
	}
	if c.cliName == "ccloud" {
		return c.ccloudList(cmd, options)
	} else {
		return c.confluentList(cmd, options)
	}
}

func (c *rolebindingCommand) listPrincipalResources(cmd *cobra.Command, options *rolebindingOptions) error {
	scope := &options.mdsScope
	principal := options.principal

	role := "*"
	if cmd.Flags().Changed("role") {
		r, err := cmd.Flags().GetString("role")
		if err != nil {
			return err
		}
		role = r
	}
	principalsRolesResourcePatterns, response, err := c.MDSClient.RBACRoleBindingSummariesApi.LookupResourcesForPrincipal(
		c.createContext(),
		principal,
		*scope)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			return c.listPrincipalResourcesV1(cmd, scope, principal, role)
		}
		return err
	}

	outputWriter, err := output.NewListOutputWriter(cmd, resourcePatternListFields, resourcePatternHumanListLabels, resourcePatternStructuredListLabels)
	if err != nil {
		return err
	}

	for principalName, rolesResourcePatterns := range principalsRolesResourcePatterns {
		for roleName, resourcePatterns := range rolesResourcePatterns {
			if role == "*" || roleName == role {
				for _, resourcePattern := range resourcePatterns {
					add := true
					if options.resource != "" {
						add = false
						for _, rp := range options.resourcesRequest.ResourcePatterns {
							if rp == resourcePattern {
								add = true
							}
						}
					}
					if add {
						outputWriter.AddElement(&listDisplay{
							Principal:    principalName,
							Role:         roleName,
							ResourceType: resourcePattern.ResourceType,
							Name:         resourcePattern.Name,
							PatternType:  resourcePattern.PatternType,
						})
					}
				}
				if len(resourcePatterns) == 0 && clusterScopedRoles[roleName] {
					outputWriter.AddElement(&listDisplay{
						Principal:    principalName,
						Role:         roleName,
						ResourceType: "Cluster",
						Name:         "",
						PatternType:  "",
					})
				}
			}
		}
	}

	outputWriter.StableSort()

	return outputWriter.Out()
}

func (c *rolebindingCommand) listPrincipalResourcesV1(cmd *cobra.Command, mdsScope *mds.MdsScope, principal string, role string) error {
	var err error
	roleNames := []string{role}
	if role == "*" {
		roleNames, _, err = c.MDSClient.RBACRoleBindingSummariesApi.ScopedPrincipalRolenames(
			c.createContext(),
			principal,
			*mdsScope)
		if err != nil {
			return err
		}
	}

	var data [][]string
	for _, roleName := range roleNames {
		rps, _, err := c.MDSClient.RBACRoleBindingCRUDApi.GetRoleResourcesForPrincipal(
			c.createContext(),
			principal,
			roleName,
			*mdsScope)
		if err != nil {
			return err
		}
		for _, pattern := range rps {
			data = append(data, []string{roleName, pattern.ResourceType, pattern.Name, pattern.PatternType})
		}
		if len(rps) == 0 && clusterScopedRoles[roleName] {
			data = append(data, []string{roleName, "Cluster", "", ""})
		}
	}

	printer.RenderCollectionTable(data, []string{"Role", "ResourceType", "Name", "PatternType"})
	return nil
}

func (c *rolebindingCommand) confluentListRolePrincipals(cmd *cobra.Command, options *rolebindingOptions) error {
	scope := &options.mdsScope
	role := options.role

	var principals []string
	var err error
	if cmd.Flags().Changed("resource") {
		r, err := cmd.Flags().GetString("resource")
		if err != nil {
			return err
		}
		resource, err := c.parseAndValidateResourcePattern(r, false)
		if err != nil {
			return err
		}
		err = c.validateRoleAndResourceTypeV1(role, resource.ResourceType)
		if err != nil {
			return err
		}
		principals, _, err = c.MDSClient.RBACRoleBindingSummariesApi.LookupPrincipalsWithRoleOnResource(
			c.createContext(),
			role,
			resource.ResourceType,
			resource.Name,
			*scope)
		if err != nil {
			return err
		}
	} else {
		principals, _, err = c.MDSClient.RBACRoleBindingSummariesApi.LookupPrincipalsWithRole(
			c.createContext(),
			role,
			*scope)
		if err != nil {
			return err
		}
	}

	sort.Strings(principals)
	outputWriter, err := output.NewListOutputWriter(cmd, []string{"Principal"}, []string{"Principal"}, []string{"principal"})
	if err != nil {
		return err
	}
	for _, principal := range principals {
		displayStruct := &struct {
			Principal string
		}{
			Principal: principal,
		}
		outputWriter.AddElement(displayStruct)
	}
	return outputWriter.Out()
}

func (c *rolebindingCommand) ccloudListRolePrincipals(cmd *cobra.Command, options *rolebindingOptions) error {
	scopeV2 := &options.scopeV2
	role := options.role

	var principals []string
	var err error
	if c.ccloudRbacDataplaneEnabled && cmd.Flags().Changed("resource") {
		r, err := cmd.Flags().GetString("resource")
		if err != nil {
			return err
		}
		resource, err := c.parseAndValidateResourcePatternV2(r, false)
		if err != nil {
			return err
		}
		err = c.validateRoleAndResourceTypeV2(role, resource.ResourceType)
		if err != nil {
			return err
		}
		principals, _, err = c.MDSv2Client.RBACRoleBindingSummariesApi.LookupPrincipalsWithRoleOnResource(
			c.createContext(),
			role,
			resource.ResourceType,
			resource.Name,
			*scopeV2)
		if err != nil {
			return err
		}
	} else {
		principals, _, err = c.MDSv2Client.RBACRoleBindingSummariesApi.LookupPrincipalsWithRole(
			c.createContext(),
			role,
			*scopeV2)
		if err != nil {
			return err
		}
	}

	userToEmailMap, err := c.userIdToEmailMap()
	if err != nil {
		return err
	}

	sort.Strings(principals)
	outputWriter, err := output.NewListOutputWriter(cmd, []string{"Principal", "Email"}, []string{"Principal", "Email"}, []string{"principal", "email"})
	if err != nil {
		return err
	}
	for _, principal := range principals {
		displayStruct := &struct {
			Principal string
			Email     string
		}{
			Principal: principal,
			Email:     userToEmailMap[principal],
		}
		outputWriter.AddElement(displayStruct)
	}
	return outputWriter.Out()
}

func (c *rolebindingCommand) userIdToEmailMap() (map[string]string, error) {
	userToEmailMap := make(map[string]string)
	users, err := c.Client.User.List(context.Background())
	if err != nil {
		return userToEmailMap, err
	}
	for _, u := range users {
		userToEmailMap["User:"+u.ResourceId] = u.Email
	}
	return userToEmailMap, nil
}

func (c *rolebindingCommand) parseCommon(cmd *cobra.Command) (*rolebindingOptions, error) {
	role, err := cmd.Flags().GetString("role")
	if err != nil {
		return nil, err
	}

	resource := ""
	prefix := false
	if c.cliName == "confluent" || c.ccloudRbacDataplaneEnabled {
		resource, err = cmd.Flags().GetString("resource")
		if err != nil {
			return nil, err
		}
		prefix = cmd.Flags().Changed("prefix")
	}

	principal, err := cmd.Flags().GetString("principal")
	if err != nil {
		return nil, err
	}
	if c.cliName == "ccloud" {
		if strings.HasPrefix(principal, "User:") {
			principalValue := strings.TrimLeft(principal, "User:")
			if strings.Contains(principalValue, "@") {
				user, err := c.Client.User.Describe(context.Background(), &orgv1.User{Email: principalValue})
				if err != nil {
					return nil, err
				}
				principal = "User:" + user.ResourceId
			}
		}
	}

	if cmd.Flags().Changed("principal") {
		err = c.validatePrincipalFormat(principal)
		if err != nil {
			return nil, err
		}
	}

	scope := &mds.MdsScope{}
	scopeV2 := &mdsv2alpha1.Scope{}
	if c.cliName != "ccloud" {
		scope, err = c.parseAndValidateScope(cmd)
	} else {
		scopeV2, err = c.parseAndValidateScopeV2(cmd)
	}
	if err != nil {
		return nil, err
	}

	resourcesRequest := mds.ResourcesRequest{}
	resourcesRequestV2 := mdsv2alpha1.ResourcesRequest{}
	if resource != "" {
		if c.cliName == "ccloud"  && c.ccloudRbacDataplaneEnabled {
			parsedResourcePattern, err := c.parseAndValidateResourcePatternV2(resource, prefix)
			if err != nil {
				return nil, err
			}
			// Resource types are defined under roles' access policies, so if no role is specified,
			// we have to loop over the possible resource types for all roles (this is what
			// validateResourceTypeV2 does).
			if role != "" {
				err = c.validateRoleAndResourceTypeV2(role, parsedResourcePattern.ResourceType)
				if err != nil {
					return nil, err
				}
			} else {
				err = c.validateResourceTypeV2(parsedResourcePattern.ResourceType)
				if err != nil {
					return nil, err
				}
			}
			resourcePatterns := []mdsv2alpha1.ResourcePattern{
				parsedResourcePattern,
			}
			resourcesRequestV2 = mdsv2alpha1.ResourcesRequest{
				Scope:            *scopeV2,
				ResourcePatterns: resourcePatterns,
			}

		} else if c.cliName == "confluent" {
			parsedResourcePattern, err := c.parseAndValidateResourcePattern(resource, prefix)
			if err != nil {
				return nil, err
			}

			// Resource types are defined under roles' access policies, so if no role is specified,
			// we have to loop over the possible resource types for all roles (this is what
			// validateResourceTypeV1 does).
			if role != "" {
				err = c.validateRoleAndResourceTypeV1(role, parsedResourcePattern.ResourceType)
				if err != nil {
					return nil, err
				}
			} else {
				err = c.validateResourceTypeV1(parsedResourcePattern.ResourceType)
				if err != nil {
					return nil, err
				}
			}

			resourcePatterns := []mds.ResourcePattern{
				parsedResourcePattern,
			}
			resourcesRequest = mds.ResourcesRequest{
				Scope:            *scope,
				ResourcePatterns: resourcePatterns,
			}
		}
	}
	return &rolebindingOptions{
			role,
			resource,
			prefix,
			principal,
			*scopeV2,
			*scope,
			resourcesRequest,
			resourcesRequestV2,
		},
		nil
}

func (c *rolebindingCommand) confluentCreate(options *rolebindingOptions) (resp *http.Response, err error) {
	if options.resource != "" {
		resp, err = c.MDSClient.RBACRoleBindingCRUDApi.AddRoleResourcesForPrincipal(
			c.createContext(),
			options.principal,
			options.role,
			options.resourcesRequest)
	} else {
		resp, err = c.MDSClient.RBACRoleBindingCRUDApi.AddRoleForPrincipal(
			c.createContext(),
			options.principal,
			options.role,
			options.mdsScope)
	}
	return
}

func (c *rolebindingCommand) ccloudCreate(options *rolebindingOptions) (*http.Response, error) {
	if c.ccloudRbacDataplaneEnabled && options.resource != "" {
		return c.MDSv2Client.RBACRoleBindingCRUDApi.AddRoleResourcesForPrincipal(
			c.createContext(),
			options.principal,
			options.role,
			options.resourcesRequestV2)
	} else {
		return c.MDSv2Client.RBACRoleBindingCRUDApi.AddRoleForPrincipal(
			c.createContext(),
			options.principal,
			options.role,
			options.scopeV2)
	}
}

func (c *rolebindingCommand) create(cmd *cobra.Command, _ []string) error {
	options, err := c.parseCommon(cmd)
	if err != nil {
		return err
	}

	var resp *http.Response
	if c.cliName == "ccloud" {
		resp, err = c.ccloudCreate(options)
	} else {
		resp, err = c.confluentCreate(options)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.NewErrorWithSuggestions(fmt.Sprintf(errors.HTTPStatusCodeErrorMsg, resp.StatusCode), errors.HTTPStatusCodeSuggestions)
	}
	if c.cliName == "ccloud" {
		return c.displayCCloudCreateAndDeleteOutput(cmd, options)
	} else {
		return displayCreateAndDeleteOutput(cmd, options)
	}
}

func (c *rolebindingCommand) displayCCloudCreateAndDeleteOutput(cmd *cobra.Command, options *rolebindingOptions) error {
	var fieldsSelected []string
	structuredRename := map[string]string{"Principal": "principal", "Email": "email", "Role": "role", "ResourceType": "resource_type", "Name": "name", "PatternType": "pattern_type"}
	userResourceId := strings.TrimLeft(options.principal, "User:")
	user, err := c.Client.User.Describe(context.Background(), &orgv1.User{ResourceId: userResourceId})
	displayStruct := &listDisplay{
		Principal: options.principal,
		Role:      options.role,
	}

	if c.ccloudRbacDataplaneEnabled && options.resource != "" {
		if len(options.resourcesRequestV2.ResourcePatterns) != 1 {
			return errors.New("display error: number of resource pattern is not 1")
		}
		resourcePattern := options.resourcesRequestV2.ResourcePatterns[0]
		displayStruct.ResourceType = resourcePattern.ResourceType
		displayStruct.Name = resourcePattern.Name
		displayStruct.PatternType = resourcePattern.PatternType
	}

	if err != nil {
		if c.ccloudRbacDataplaneEnabled && options.resource != "" {
			fieldsSelected = resourcePatternListFields
		} else {
			fieldsSelected = []string{"Principal", "Role"}
		}
	} else {
		if c.ccloudRbacDataplaneEnabled && options.resource != "" {
			fieldsSelected = ccloudResourcePatternListFields
		} else {
			displayStruct.Email = user.Email
			fieldsSelected = []string{"Principal", "Email", "Role"}
		}
	}
	return output.DescribeObject(cmd, displayStruct, fieldsSelected, map[string]string{}, structuredRename)
}

func displayCreateAndDeleteOutput(cmd *cobra.Command, options *rolebindingOptions) error {
	var fieldsSelected []string
	structuredRename := map[string]string{"Principal": "principal", "Role": "role", "ResourceType": "resource_type", "Name": "name", "PatternType": "pattern_type"}
	displayStruct := &listDisplay{
		Principal: options.principal,
		Role:      options.role,
	}
	if options.resource != "" {
		fieldsSelected = resourcePatternListFields
		if len(options.resourcesRequest.ResourcePatterns) != 1 {
			return errors.New("display error: number of resource pattern is not 1")
		}
		resourcePattern := options.resourcesRequest.ResourcePatterns[0]
		displayStruct.ResourceType = resourcePattern.ResourceType
		displayStruct.Name = resourcePattern.Name
		displayStruct.PatternType = resourcePattern.PatternType
	} else {
		fieldsSelected = []string{"Principal", "Role", "ResourceType"}
		displayStruct.ResourceType = "Cluster"
	}
	return output.DescribeObject(cmd, displayStruct, fieldsSelected, map[string]string{}, structuredRename)
}

func (c *rolebindingCommand) confluentDelete(options *rolebindingOptions) (resp *http.Response, err error) {
	if options.resource != "" {
		resp, err = c.MDSClient.RBACRoleBindingCRUDApi.RemoveRoleResourcesForPrincipal(
			c.createContext(),
			options.principal,
			options.role,
			options.resourcesRequest)
	} else {
		resp, err = c.MDSClient.RBACRoleBindingCRUDApi.DeleteRoleForPrincipal(
			c.createContext(),
			options.principal,
			options.role,
			options.mdsScope)
	}
	return
}

func (c *rolebindingCommand) ccloudDelete(options *rolebindingOptions) (*http.Response, error) {
	if c.ccloudRbacDataplaneEnabled && options.resource != "" {
		return c.MDSv2Client.RBACRoleBindingCRUDApi.RemoveRoleResourcesForPrincipal (
			c.createContext(),
			options.principal,
			options.role,
			options.resourcesRequestV2)
	} else {
		return c.MDSv2Client.RBACRoleBindingCRUDApi.DeleteRoleForPrincipal(
			c.createContext(),
			options.principal,
			options.role,
			options.scopeV2)
	}
}

func (c *rolebindingCommand) delete(cmd *cobra.Command, _ []string) error {
	options, err := c.parseCommon(cmd)
	if err != nil {
		return err
	}

	var resp *http.Response
	if c.cliName == "ccloud" {
		resp, err = c.ccloudDelete(options)
	} else {
		resp, err = c.confluentDelete(options)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return errors.NewErrorWithSuggestions(fmt.Sprintf(errors.HTTPStatusCodeErrorMsg, resp.StatusCode), errors.HTTPStatusCodeSuggestions)
	}

	if c.cliName == "ccloud" {
		return c.displayCCloudCreateAndDeleteOutput(cmd, options)
	} else {
		return displayCreateAndDeleteOutput(cmd, options)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *rolebindingCommand) createContext() context.Context {
	if c.cliName == "ccloud" {
		return context.WithValue(context.Background(), mdsv2alpha1.ContextAccessToken, c.AuthToken())
	} else {
		return context.WithValue(context.Background(), mds.ContextAccessToken, c.AuthToken())
	}
}
