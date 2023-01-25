// Copyright 2021 Confluent Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
REST Admin API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 3.0.0
Contact: kafka-clients-proxy-team@confluent.io
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package v3

import (
	"encoding/json"
)

import (
	"reflect"
)

// AclData struct for AclData
type AclData struct {
	Kind         string           `json:"kind"`
	Metadata     ResourceMetadata `json:"metadata"`
	ClusterId    string           `json:"cluster_id"`
	ResourceType AclResourceType  `json:"resource_type"`
	ResourceName string           `json:"resource_name"`
	PatternType  string           `json:"pattern_type"`
	Principal    string           `json:"principal"`
	Host         string           `json:"host"`
	Operation    string           `json:"operation"`
	Permission   string           `json:"permission"`
}

// NewAclData instantiates a new AclData object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAclData(kind string, metadata ResourceMetadata, clusterId string, resourceType AclResourceType, resourceName string, patternType string, principal string, host string, operation string, permission string) *AclData {
	this := AclData{}
	this.Kind = kind
	this.Metadata = metadata
	this.ClusterId = clusterId
	this.ResourceType = resourceType
	this.ResourceName = resourceName
	this.PatternType = patternType
	this.Principal = principal
	this.Host = host
	this.Operation = operation
	this.Permission = permission
	return &this
}

// NewAclDataWithDefaults instantiates a new AclData object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAclDataWithDefaults() *AclData {
	this := AclData{}
	return &this
}

// GetKind returns the Kind field value
func (o *AclData) GetKind() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Kind
}

// GetKindOk returns a tuple with the Kind field value
// and a boolean to check if the value has been set.
func (o *AclData) GetKindOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Kind, true
}

// SetKind sets field value
func (o *AclData) SetKind(v string) {
	o.Kind = v
}

// GetMetadata returns the Metadata field value
func (o *AclData) GetMetadata() ResourceMetadata {
	if o == nil {
		var ret ResourceMetadata
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *AclData) GetMetadataOk() (*ResourceMetadata, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *AclData) SetMetadata(v ResourceMetadata) {
	o.Metadata = v
}

// GetClusterId returns the ClusterId field value
func (o *AclData) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *AclData) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *AclData) SetClusterId(v string) {
	o.ClusterId = v
}

// GetResourceType returns the ResourceType field value
func (o *AclData) GetResourceType() AclResourceType {
	if o == nil {
		var ret AclResourceType
		return ret
	}

	return o.ResourceType
}

// GetResourceTypeOk returns a tuple with the ResourceType field value
// and a boolean to check if the value has been set.
func (o *AclData) GetResourceTypeOk() (*AclResourceType, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ResourceType, true
}

// SetResourceType sets field value
func (o *AclData) SetResourceType(v AclResourceType) {
	o.ResourceType = v
}

// GetResourceName returns the ResourceName field value
func (o *AclData) GetResourceName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ResourceName
}

// GetResourceNameOk returns a tuple with the ResourceName field value
// and a boolean to check if the value has been set.
func (o *AclData) GetResourceNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ResourceName, true
}

// SetResourceName sets field value
func (o *AclData) SetResourceName(v string) {
	o.ResourceName = v
}

// GetPatternType returns the PatternType field value
func (o *AclData) GetPatternType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PatternType
}

// GetPatternTypeOk returns a tuple with the PatternType field value
// and a boolean to check if the value has been set.
func (o *AclData) GetPatternTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PatternType, true
}

// SetPatternType sets field value
func (o *AclData) SetPatternType(v string) {
	o.PatternType = v
}

// GetPrincipal returns the Principal field value
func (o *AclData) GetPrincipal() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Principal
}

// GetPrincipalOk returns a tuple with the Principal field value
// and a boolean to check if the value has been set.
func (o *AclData) GetPrincipalOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Principal, true
}

// SetPrincipal sets field value
func (o *AclData) SetPrincipal(v string) {
	o.Principal = v
}

// GetHost returns the Host field value
func (o *AclData) GetHost() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Host
}

// GetHostOk returns a tuple with the Host field value
// and a boolean to check if the value has been set.
func (o *AclData) GetHostOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Host, true
}

// SetHost sets field value
func (o *AclData) SetHost(v string) {
	o.Host = v
}

// GetOperation returns the Operation field value
func (o *AclData) GetOperation() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Operation
}

// GetOperationOk returns a tuple with the Operation field value
// and a boolean to check if the value has been set.
func (o *AclData) GetOperationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Operation, true
}

// SetOperation sets field value
func (o *AclData) SetOperation(v string) {
	o.Operation = v
}

// GetPermission returns the Permission field value
func (o *AclData) GetPermission() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Permission
}

// GetPermissionOk returns a tuple with the Permission field value
// and a boolean to check if the value has been set.
func (o *AclData) GetPermissionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Permission, true
}

// SetPermission sets field value
func (o *AclData) SetPermission(v string) {
	o.Permission = v
}

// Redact resets all sensitive fields to their zero value.
func (o *AclData) Redact() {
	o.recurseRedact(&o.Kind)
	o.recurseRedact(&o.Metadata)
	o.recurseRedact(&o.ClusterId)
	o.recurseRedact(&o.ResourceType)
	o.recurseRedact(&o.ResourceName)
	o.recurseRedact(&o.PatternType)
	o.recurseRedact(&o.Principal)
	o.recurseRedact(&o.Host)
	o.recurseRedact(&o.Operation)
	o.recurseRedact(&o.Permission)
}

func (o *AclData) recurseRedact(v interface{}) {
	type redactor interface {
		Redact()
	}
	if r, ok := v.(redactor); ok {
		r.Redact()
	} else {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		switch val.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < val.Len(); i++ {
				// support data types declared without pointers
				o.recurseRedact(val.Index(i).Interface())
				// ... and data types that were declared without but need pointers (for Redact)
				if val.Index(i).CanAddr() {
					o.recurseRedact(val.Index(i).Addr().Interface())
				}
			}
		}
	}
}

func (o AclData) zeroField(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (o AclData) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["kind"] = o.Kind
	}
	if true {
		toSerialize["metadata"] = o.Metadata
	}
	if true {
		toSerialize["cluster_id"] = o.ClusterId
	}
	if true {
		toSerialize["resource_type"] = o.ResourceType
	}
	if true {
		toSerialize["resource_name"] = o.ResourceName
	}
	if true {
		toSerialize["pattern_type"] = o.PatternType
	}
	if true {
		toSerialize["principal"] = o.Principal
	}
	if true {
		toSerialize["host"] = o.Host
	}
	if true {
		toSerialize["operation"] = o.Operation
	}
	if true {
		toSerialize["permission"] = o.Permission
	}
	return json.Marshal(toSerialize)
}

type NullableAclData struct {
	value *AclData
	isSet bool
}

func (v NullableAclData) Get() *AclData {
	return v.value
}

func (v *NullableAclData) Set(val *AclData) {
	v.value = val
	v.isSet = true
}

func (v NullableAclData) IsSet() bool {
	return v.isSet
}

func (v *NullableAclData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAclData(val *AclData) *NullableAclData {
	return &NullableAclData{value: val, isSet: true}
}

func (v NullableAclData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAclData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
