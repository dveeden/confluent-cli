package iam

import (
	"net/http"
	"os"
	"sort"
	"strings"

	mdsv2 "github.com/confluentinc/ccloud-sdk-go-v2/mds/v2"
	mds "github.com/confluentinc/mds-sdk-go/mdsv1"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/examples"
	"github.com/confluentinc/cli/internal/pkg/output"
)

type roleBindingListOut struct {
	Principal   string `human:"Principal" serialized:"principal"`
	Email       string `human:"Email" serialized:"email"`
	ServiceName string `human:"Service Name" serialized:"service_name"`
	PoolName    string `human:"Pool Name" serialized:"pool_name"`
}

func (c *roleBindingCommand) newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List role bindings.",
		Long:  "List the role bindings for a particular principal and/or role, and a particular scope.",
		Args:  cobra.NoArgs,
		RunE:  c.list,
	}

	if c.cfg.IsCloudLogin() {
		cmd.Example = examples.BuildExampleString(
			examples.Example{
				Text: "List the role bindings for current user:",
				Code: "confluent iam rbac role-binding list --current-user",
			},
			examples.Example{
				Text: `List the role bindings for user "u-123456":`,
				Code: "confluent iam rbac role-binding list --principal User:u-123456",
			},
			examples.Example{
				Text: `List the role bindings for principals with role "CloudClusterAdmin":`,
				Code: "confluent iam rbac role-binding list --role CloudClusterAdmin --current-environment --cloud-cluster lkc-123456",
			},
			examples.Example{
				Text: `List the role bindings for user "u-123456" with role "CloudClusterAdmin":`,
				Code: "confluent iam rbac role-binding list --principal User:u-123456 --role CloudClusterAdmin --environment env-12345 --cloud-cluster lkc-123456",
			},
		)
	} else {
		cmd.Example = examples.BuildExampleString(
			examples.Example{
				Text: "Only use the `--resource` flag when specifying a `--role` with no `--principal` specified. If specifying a `--principal`, then the `--resource` flag is ignored. To list role bindings for a specific role on an identified resource:",
				Code: "confluent iam rbac role-binding list --kafka-cluster $KAFKA_CLUSTER_ID --role DeveloperRead --resource Topic",
			},
			examples.Example{
				Text: "List the role bindings for a specific principal:",
				Code: "confluent iam rbac role-binding list --kafka-cluster $KAFKA_CLUSTER_ID --principal User:my-user",
			},
			examples.Example{
				Text: "List the role bindings for a specific principal, filtered to a specific role:",
				Code: "confluent iam rbac role-binding list --kafka-cluster $KAFKA_CLUSTER_ID --principal User:my-user --role DeveloperRead",
			},
			examples.Example{
				Text: "List the principals bound to a specific role:",
				Code: "confluent iam rbac role-binding list --kafka-cluster $KAFKA_CLUSTER_ID --role DeveloperWrite",
			},
			examples.Example{
				Text: "List the principals bound to a specific resource with a specific role:",
				Code: "confluent iam rbac role-binding list --kafka-cluster $KAFKA_CLUSTER_ID --role DeveloperWrite --resource Topic:my-topic",
			},
		)
	}

	cmd.Flags().String("principal", "", "Principal whose role bindings should be listed.")
	cmd.Flags().Bool("current-user", false, "Show role bindings belonging to current user.")
	cmd.Flags().String("role", "", "List role bindings under a specific role given to a principal. Or if no principal is specified, list principals with the role.")

	if c.cfg.IsCloudLogin() {
		cmd.Flags().String("environment", "", "Environment ID for scope of role binding listings.")
		cmd.Flags().Bool("current-environment", false, "Use current environment ID for scope.")
		cmd.Flags().String("cloud-cluster", "", "Cloud cluster ID for scope of role binding listings.")
		cmd.Flags().String("kafka-cluster", "", "Kafka cluster ID for scope of role binding listings.")
		if os.Getenv("XX_DATAPLANE_3_ENABLE") != "" {
			cmd.Flags().String("schema-registry-cluster", "", "Schema Registry cluster ID for the role binding listings.")
			cmd.Flags().String("ksql-cluster", "", "ksqlDB cluster ID for the role binding listings.")
		}
	} else {
		cmd.Flags().String("kafka-cluster", "", "Kafka cluster ID for scope of role binding listings.")
		cmd.Flags().String("schema-registry-cluster", "", "Schema Registry cluster ID for scope of role binding listings.")
		cmd.Flags().String("ksql-cluster", "", "ksqlDB cluster ID for scope of role binding listings.")
		cmd.Flags().String("connect-cluster", "", "Kafka Connect cluster ID for scope of role binding listings.")
		cmd.Flags().String("cluster-name", "", "Cluster name to uniquely identify the cluster for role binding listings.")
		pcmd.AddContextFlag(cmd, c.CLICommand)
	}

	cmd.Flags().String("resource", "", "If specified with a role and no principals, list principals with role bindings to the role for this qualified resource.")
	cmd.Flags().Bool("inclusive", true, "List all role bindings in a specific scope and its nested scopes.")
	pcmd.AddOutputFlag(cmd)

	return cmd
}

var (
	ksqlOrSchemaRegistryRoleBindingError = errors.New(errors.KsqlOrSchemaRegistryRoleBindingErrorMsg)
)

func (c *roleBindingCommand) list(cmd *cobra.Command, _ []string) error {
	if c.cfg.IsCloudLogin() {
		listRoleBinding, err := c.parseV2RoleBinding(cmd)
		if err != nil {
			return err
		}
		err = c.ccloudListV2(cmd, listRoleBinding)
		if err == ksqlOrSchemaRegistryRoleBindingError {
			options, err := c.parseCommon(cmd)
			if err != nil {
				return err
			}
			return c.ccloudList(cmd, options)
		}
		return err
	} else {
		options, err := c.parseCommon(cmd)
		if err != nil {
			return err
		}
		return c.confluentList(cmd, options)
	}
}

func (c *roleBindingCommand) ccloudList(cmd *cobra.Command, options *roleBindingOptions) error {
	if cmd.Flags().Changed("principal") || cmd.Flags().Changed("current-user") {
		return c.listMyRoleBindings(cmd, options)
	} else if cmd.Flags().Changed("role") {
		return c.ccloudListRolePrincipals(cmd, options)
	}
	return errors.New(errors.PrincipalOrRoleRequiredErrorMsg)
}

func (c *roleBindingCommand) listMyRoleBindings(cmd *cobra.Command, options *roleBindingOptions) error {
	scopeV2 := &options.scopeV2

	currentUser, err := cmd.Flags().GetBool("current-user")
	if err != nil {
		return err
	}

	principal := options.principal
	if currentUser {
		principal = "User:" + c.Context.GetUser().GetResourceId()
	}

	scopedRoleBindingMappings, _, err := c.MDSv2Client.RBACRoleBindingSummariesApi.MyRoleBindings(c.createContext(), principal, *scopeV2)
	if err != nil {
		return err
	}

	userToEmailMap, err := c.getUserIdToEmailMap()
	if err != nil {
		return err
	}

	role, err := cmd.Flags().GetString("role")
	if err != nil {
		return err
	}

	list := output.NewList(cmd)
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
					list.Add(&roleBindingOut{
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
					list.Add(&roleBindingOut{
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
	return list.Print()
}

func (c *roleBindingCommand) getPoolToNameMap() (map[string]string, error) {
	providers, err := c.V2Client.ListIdentityProviders()
	if err != nil {
		return map[string]string{}, err
	}
	poolToName := make(map[string]string)
	for _, provider := range providers {
		pools, err := c.V2Client.ListIdentityPools(provider.GetId())
		if err != nil {
			return map[string]string{}, err
		}
		for _, pool := range pools {
			poolToName["User:"+pool.GetId()] = pool.GetDisplayName()
		}
	}
	return poolToName, nil
}

func (c *roleBindingCommand) getUserIdToEmailMap() (map[string]string, error) {
	userToEmailMap := make(map[string]string)
	users, err := c.V2Client.ListIamUsers()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		userToEmailMap["User:"+u.GetId()] = u.GetEmail()
	}
	return userToEmailMap, nil
}

func (c *roleBindingCommand) getServiceAccountIdToNameMap() (map[string]string, error) {
	serviceAccounts, err := c.V2Client.ListIamServiceAccounts()
	if err != nil {
		return nil, err
	}

	serviceAccountToNameMap := make(map[string]string)
	for _, sa := range serviceAccounts {
		serviceAccountToNameMap["User:"+sa.GetId()] = sa.GetDisplayName()
	}
	return serviceAccountToNameMap, nil
}

func (c *roleBindingCommand) ccloudListRolePrincipals(cmd *cobra.Command, options *roleBindingOptions) error {
	scopeV2 := &options.scopeV2
	role := options.role

	var principals []string
	var err error
	if cmd.Flags().Changed("resource") {
		r, err := cmd.Flags().GetString("resource")
		if err != nil {
			return err
		}
		resource, err := parseAndValidateResourcePatternV2(r, false)
		if err != nil {
			return err
		}
		//  skip validation when role != "" before migrating to v2 role api
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

	userToEmailMap, err := c.getUserIdToEmailMap()
	if err != nil {
		return err
	}

	serviceAccountToNameMap, err := c.getServiceAccountIdToNameMap()
	if err != nil {
		return err
	}

	// TODO: Catch this error once Identity Providers goes GA
	poolToNameMap, _ := c.getPoolToNameMap()

	list := output.NewList(cmd)
	for _, principal := range principals {
		row := &roleBindingListOut{Principal: principal}
		if email, ok := userToEmailMap[principal]; ok {
			row.Email = email
			list.Add(row)
		}
		if name, ok := serviceAccountToNameMap[principal]; ok {
			row.ServiceName = name
			list.Add(row)
		}
		if name, ok := poolToNameMap[principal]; ok {
			row.PoolName = name
			list.Add(row)
		}
	}
	return list.Print()
}

func (c *roleBindingCommand) confluentList(cmd *cobra.Command, options *roleBindingOptions) error {
	if cmd.Flags().Changed("principal") {
		return c.listPrincipalResources(cmd, options)
	} else if cmd.Flags().Changed("role") {
		return c.confluentListRolePrincipals(cmd, options)
	}
	return errors.New(errors.PrincipalOrRoleRequiredErrorMsg)
}

func (c *roleBindingCommand) listPrincipalResources(cmd *cobra.Command, options *roleBindingOptions) error {
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

	principalsRolesResourcePatterns, response, err := c.MDSClient.RBACRoleBindingSummariesApi.LookupResourcesForPrincipal(c.createContext(), principal, *scope)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			return c.listPrincipalResourcesV1(cmd, scope, principal, role)
		}
		return err
	}

	list := output.NewList(cmd)
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
						list.Add(&roleBindingOut{
							Principal:    principalName,
							Role:         roleName,
							ResourceType: resourcePattern.ResourceType,
							Name:         resourcePattern.Name,
							PatternType:  resourcePattern.PatternType,
						})
					}
				}
				if len(resourcePatterns) == 0 && clusterScopedRoles[roleName] {
					list.Add(&roleBindingOut{
						Principal:    principalName,
						Role:         roleName,
						ResourceType: "Cluster",
					})
				}
			}
		}
	}
	list.Filter(resourcePatternListFields)
	return list.Print()
}

func (c *roleBindingCommand) listPrincipalResourcesV1(cmd *cobra.Command, mdsScope *mds.MdsScope, principal string, role string) error {
	var err error
	roleNames := []string{role}
	if role == "*" {
		roleNames, _, err = c.MDSClient.RBACRoleBindingSummariesApi.ScopedPrincipalRolenames(c.createContext(), principal, *mdsScope)
		if err != nil {
			return err
		}
	}

	list := output.NewList(cmd)
	for _, roleName := range roleNames {
		resourcePatterns, _, err := c.MDSClient.RBACRoleBindingCRUDApi.GetRoleResourcesForPrincipal(c.createContext(), principal, roleName, *mdsScope)
		if err != nil {
			return err
		}
		for _, pattern := range resourcePatterns {
			list.Add(&roleBindingOut{
				Role:         roleName,
				ResourceType: pattern.ResourceType,
				Name:         pattern.Name,
				PatternType:  pattern.PatternType,
			})
		}
		if len(resourcePatterns) == 0 && clusterScopedRoles[roleName] {
			list.Add(&roleBindingOut{
				Role:         roleName,
				ResourceType: "Cluster",
			})
		}
	}
	list.Filter([]string{"Role", "ResourceType", "Name", "PatternType"})
	return list.Print()
}

func (c *roleBindingCommand) confluentListRolePrincipals(cmd *cobra.Command, options *roleBindingOptions) error {
	scope := &options.mdsScope
	role := options.role

	var principals []string
	if cmd.Flags().Changed("resource") {
		r, err := cmd.Flags().GetString("resource")
		if err != nil {
			return err
		}

		resource, err := parseAndValidateResourcePattern(r, false)
		if err != nil {
			return err
		}

		if err := c.validateRoleAndResourceTypeV1(role, resource.ResourceType); err != nil {
			return err
		}

		principals, _, err = c.MDSClient.RBACRoleBindingSummariesApi.LookupPrincipalsWithRoleOnResource(c.createContext(), role, resource.ResourceType, resource.Name, *scope)
		if err != nil {
			return err
		}
	} else {
		var err error
		principals, _, err = c.MDSClient.RBACRoleBindingSummariesApi.LookupPrincipalsWithRole(c.createContext(), role, *scope)
		if err != nil {
			return err
		}
	}

	list := output.NewList(cmd)
	for _, principal := range principals {
		list.Add(&roleBindingOut{Principal: principal})
	}
	list.Filter([]string{"Principal"})
	return list.Print()
}

func (c *roleBindingCommand) ccloudListV2(cmd *cobra.Command, listRoleBinding *mdsv2.IamV2RoleBinding) error {
	if cmd.Flags().Changed("principal") || cmd.Flags().Changed("current-user") {
		return c.listMyRoleBindingsV2(cmd, listRoleBinding)
	} else if cmd.Flags().Changed("role") {
		return c.ccloudListRolePrincipalsV2(cmd, listRoleBinding)
	}
	return errors.New(errors.PrincipalOrRoleRequiredErrorMsg)
}

func (c *roleBindingCommand) listMyRoleBindingsV2(cmd *cobra.Command, listRoleBinding *mdsv2.IamV2RoleBinding) error {
	list := output.NewList(cmd)

	currentUser, err := cmd.Flags().GetBool("current-user")
	if err != nil {
		return err
	}

	if currentUser {
		listRoleBinding.Principal = mdsv2.PtrString("User:" + c.State.Auth.User.GetResourceId())
	}

	inclusive, err := cmd.Flags().GetBool("inclusive")
	if err != nil {
		return err
	}

	if inclusive {
		listRoleBinding.CrnPattern = mdsv2.PtrString(listRoleBinding.GetCrnPattern() + "/*")
	} else {
		listRoleBinding.CrnPattern = mdsv2.PtrString(listRoleBinding.GetCrnPattern())
	}

	resp, err := c.V2Client.ListIamRoleBindings(listRoleBinding.GetPrincipal(), listRoleBinding.GetRoleName(), listRoleBinding.GetCrnPattern())
	if err != nil {
		return err
	}
	roleBindings := resp.Data

	userToEmailMap, err := c.getUserIdToEmailMap()
	if err != nil {
		return err
	}

	role, err := cmd.Flags().GetString("role")
	if err != nil {
		return err
	}

	for _, rolebinding := range roleBindings {
		roleName := rolebinding.GetRoleName()
		if role != "" && role != roleName {
			continue
		}

		principalName := rolebinding.GetPrincipal()
		principalEmail := userToEmailMap[principalName]

		crnPattern := rolebinding.GetCrnPattern()
		if strings.Contains(crnPattern, "ksql") || strings.Contains(crnPattern, "schema") {
			return ksqlOrSchemaRegistryRoleBindingError
		}

		var envName, cloudClusterName, clusterType, logicalCluster, resourceType, resourceName, patternType string
		for _, elem := range strings.Split(crnPattern, "/") {
			elemParts := strings.Split(elem, "=")
			if len(elemParts) < 2 {
				continue
			}

			prefix := elemParts[0]
			content := elemParts[1]

			switch prefix {
			case "organization":
				continue
			case "environment":
				envName = content
			case "cloud-cluster":
				cloudClusterName = content
			case "ksql":
				clusterType = "ksqlDB"
				logicalCluster = content
			case "schema-registry":
				clusterType = "Schema Registry"
				logicalCluster = content
			case "kafka":
				clusterType = "Kafka"
				logicalCluster = content
				resourceType = "Cluster"
				resourceName = "kafka-cluster"
				patternType = literalPatternType
			default:
				resourceType = cases.Title(language.Und).String(prefix)
				resourceName = strings.TrimSuffix(content, "*")
				patternType = literalPatternType
			}
		}

		if strings.Contains(crnPattern, "*") {
			patternType = prefixedPatternType
		}
		list.Add(&roleBindingOut{
			Principal:      principalName,
			Email:          principalEmail,
			Role:           roleName,
			Environment:    envName,
			CloudCluster:   cloudClusterName,
			ClusterType:    clusterType,
			LogicalCluster: logicalCluster,
			ResourceType:   resourceType,
			Name:           resourceName,
			PatternType:    patternType,
		})
	}

	list.Sort(true)
	return list.Print()
}

func (c *roleBindingCommand) ccloudListRolePrincipalsV2(cmd *cobra.Command, listRoleBinding *mdsv2.IamV2RoleBinding) error {
	list := output.NewList(cmd)

	inclusive, err := cmd.Flags().GetBool("inclusive")
	if err != nil {
		return err
	}

	if inclusive {
		listRoleBinding.CrnPattern = mdsv2.PtrString(listRoleBinding.GetCrnPattern() + "/*")
	} else {
		listRoleBinding.CrnPattern = mdsv2.PtrString(listRoleBinding.GetCrnPattern())
	}

	resp, err := c.V2Client.ListIamRoleBindings(listRoleBinding.GetPrincipal(), listRoleBinding.GetRoleName(), listRoleBinding.GetCrnPattern())
	if err != nil {
		return err
	}
	roleBindings := resp.Data

	principals := make(map[string]bool)
	principalStrings := []string{}

	for _, roleBinding := range roleBindings {
		if strings.Contains(roleBinding.GetCrnPattern(), "ksql") || strings.Contains(roleBinding.GetCrnPattern(), "schema") {
			return ksqlOrSchemaRegistryRoleBindingError
		}
		if !principals[roleBinding.GetPrincipal()] {
			principals[roleBinding.GetPrincipal()] = true
			principalStrings = append(principalStrings, roleBinding.GetPrincipal())
		}
	}

	userToEmailMap, err := c.getUserIdToEmailMap()
	if err != nil {
		return err
	}

	serviceAccountToNameMap, err := c.getServiceAccountIdToNameMap()
	if err != nil {
		return err
	}

	// TODO: Catch this error once Identity Providers goes GA
	poolToNameMap, _ := c.getPoolToNameMap()

	sort.Strings(principalStrings)
	for _, principal := range principalStrings {
		row := &roleBindingListOut{Principal: principal}
		if email, ok := userToEmailMap[principal]; ok {
			row.Email = email
			list.Add(row)
		}
		if name, ok := serviceAccountToNameMap[principal]; ok {
			row.ServiceName = name
			list.Add(row)
		}
		if name, ok := poolToNameMap[principal]; ok {
			row.PoolName = name
			list.Add(row)
		}
	}

	return list.Print()
}
