package kafka

import (
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"strings"

	schedv1 "github.com/confluentinc/cc-structs/kafka/scheduler/v1"
	"github.com/confluentinc/kafka-rest-sdk-go/kafkarestv3"

	cloudkafkarestv3 "github.com/confluentinc/ccloud-sdk-go-v2/kafkarest/v3"
	"github.com/confluentinc/cli/internal/pkg/errors"
)

const KafkaRestBadRequestErrorCode = 40002
const KafkaRestUnknownTopicOrPartitionErrorCode = 40403
const SelfSignedCertError = "x509: certificate is not authorized to sign other certificates"
const UnauthorizedCertError = "x509: certificate signed by unknown authority"

type kafkaRestV3Error struct {
	Code    int    `json:"error_code"`
	Message string `json:"message"`
}

func kafkaRestHttpError(httpResp *http.Response) error {
	return errors.NewErrorWithSuggestions(
		fmt.Sprintf(errors.KafkaRestErrorMsg, httpResp.Request.Method, httpResp.Request.URL, httpResp.Status),
		errors.InternalServerErrorSuggestions)
}

func parseOpenAPIError(err error) (*kafkaRestV3Error, error) {
	if openAPIError, ok := err.(kafkarestv3.GenericOpenAPIError); ok {
		var decodedError kafkaRestV3Error
		err = json.Unmarshal(openAPIError.Body(), &decodedError)
		if err != nil {
			return nil, err
		}
		return &decodedError, nil
	}
	return nil, fmt.Errorf("unexpected type")
}

func kafkaRestError(url string, err error, httpResp *http.Response) error {
	switch e := err.(type) {
	case *neturl.Error:
		if strings.Contains(e.Error(), SelfSignedCertError) || strings.Contains(e.Error(), UnauthorizedCertError) {
			return errors.NewErrorWithSuggestions(fmt.Sprintf(errors.KafkaRestConnectionMsg, url, e.Err), errors.KafkaRestCertErrorSuggestions)
		}
		return errors.Errorf(errors.KafkaRestConnectionMsg, url, e.Err)
	case kafkarestv3.GenericOpenAPIError:
		openAPIError, parseErr := parseOpenAPIError(err)
		if parseErr == nil {
			if strings.Contains(openAPIError.Message, "invalid_token") {
				return errors.NewErrorWithSuggestions(errors.InvalidMDSToken, errors.InvalidMDSTokenSuggestions)
			}
			return fmt.Errorf("REST request failed: %v (%v)", openAPIError.Message, openAPIError.Code)
		}
		if httpResp != nil && httpResp.StatusCode >= 400 {
			return kafkaRestHttpError(httpResp)
		}
		return errors.NewErrorWithSuggestions(errors.UnknownErrorMsg, errors.InternalServerErrorSuggestions)
	}
	return err
}

// Converts ACLBinding to Kafka REST CreateAclRequestData
func aclBindingToClustersClusterIdAclsRequestData(acl *schedv1.ACLBinding) cloudkafkarestv3.CreateAclRequestData {
	var aclRequestData cloudkafkarestv3.CreateAclRequestData

	if acl.Pattern.ResourceType != schedv1.ResourceTypes_UNKNOWN {
		aclRequestData.ResourceType = cloudkafkarestv3.AclResourceType(acl.Pattern.ResourceType.String())
	}

	if acl.Pattern.PatternType != schedv1.PatternTypes_UNKNOWN {
		aclRequestData.PatternType = acl.Pattern.PatternType.String()
	}

	aclRequestData.ResourceName = acl.Pattern.Name
	aclRequestData.Principal = acl.Entry.Principal
	aclRequestData.Host = acl.Entry.Host

	if acl.Entry.Operation != schedv1.ACLOperations_UNKNOWN {
		aclRequestData.Operation = acl.Entry.Operation.String()
	}

	if acl.Entry.PermissionType != schedv1.ACLPermissionTypes_UNKNOWN {
		aclRequestData.Permission = acl.Entry.PermissionType.String()
	}

	return aclRequestData
}

// Converts ACLFilter to Kafka REST ClustersClusterIdAclsDeleteOpts
func aclFilterToClustersClusterIdAclsDeleteRequestData(acl *schedv1.ACLFilter) cloudkafkarestv3.CreateAclRequestData {
	var aclRequestData cloudkafkarestv3.CreateAclRequestData

	if acl.PatternFilter.ResourceType != schedv1.ResourceTypes_UNKNOWN {
		aclRequestData.ResourceType = cloudkafkarestv3.AclResourceType(acl.PatternFilter.ResourceType.String())
	}

	aclRequestData.ResourceName = acl.PatternFilter.Name

	if acl.PatternFilter.PatternType != schedv1.PatternTypes_UNKNOWN {
		aclRequestData.PatternType = acl.PatternFilter.PatternType.String()
	}

	aclRequestData.Principal = acl.EntryFilter.Principal
	aclRequestData.Host = acl.EntryFilter.Host

	if acl.EntryFilter.Operation != schedv1.ACLOperations_UNKNOWN {
		aclRequestData.Operation = acl.EntryFilter.Operation.String()
	}

	if acl.EntryFilter.PermissionType != schedv1.ACLPermissionTypes_UNKNOWN {
		aclRequestData.Permission = acl.EntryFilter.PermissionType.String()
	}

	return aclRequestData
}

func getDeleteAclRequestWithData(req cloudkafkarestv3.ApiDeleteKafkaAclsRequest, requestData cloudkafkarestv3.CreateAclRequestData) cloudkafkarestv3.ApiDeleteKafkaAclsRequest {
	return req.ResourceType(requestData.ResourceType).PatternType(requestData.PatternType).Principal(requestData.Principal).Host(requestData.Host).Operation(requestData.Operation).Permission(requestData.Permission)
}

func getGetAclsRequestWithData(req cloudkafkarestv3.ApiGetKafkaAclsRequest, requestData cloudkafkarestv3.CreateAclRequestData) cloudkafkarestv3.ApiGetKafkaAclsRequest {
	return req.ResourceType(requestData.ResourceType).PatternType(requestData.PatternType).Principal(requestData.Principal).Host(requestData.Host).Operation(requestData.Operation).Permission(requestData.Permission)
}
