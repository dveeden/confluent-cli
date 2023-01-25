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
Stream Share APIs

# Introduction

API version: 0.1.0-alpha0
Contact: cdx@confluent.io
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package v1

import (
	"encoding/json"
)

import (
	"reflect"
)

// CdxV1RedeemTokenRequest Redeem share with token request parameters
type CdxV1RedeemTokenRequest struct {
	// APIVersion defines the schema version of this representation of a resource.
	ApiVersion *string `json:"api_version,omitempty"`
	// Kind defines the object this REST resource represents.
	Kind *string `json:"kind,omitempty"`
	// ID is the \"natural identifier\" for an object within its scope/namespace; it is normally unique across time but not space. That is, you can assume that the ID will not be reclaimed and reused after an object is deleted (\"time\"); however, it may collide with IDs for other object `kinds` or objects of the same `kind` within a different scope/namespace (\"space\").
	Id       *string     `json:"id,omitempty"`
	Metadata *ObjectMeta `json:"metadata,omitempty"`
	// The encrypted token
	Token *string `json:"token,omitempty"`
	// Consumer's AWS account ID for PrivateLink access.
	AwsAccount *string `json:"aws_account,omitempty"`
	// Consumer's Azure subscription ID for PrivateLink access.
	AzureSubscription *string `json:"azure_subscription,omitempty"`
	// Consumer's GCP project ID for Private Service Connect access.
	GcpProject *string `json:"gcp_project,omitempty"`
}

// NewCdxV1RedeemTokenRequest instantiates a new CdxV1RedeemTokenRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCdxV1RedeemTokenRequest() *CdxV1RedeemTokenRequest {
	this := CdxV1RedeemTokenRequest{}
	return &this
}

// NewCdxV1RedeemTokenRequestWithDefaults instantiates a new CdxV1RedeemTokenRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCdxV1RedeemTokenRequestWithDefaults() *CdxV1RedeemTokenRequest {
	this := CdxV1RedeemTokenRequest{}
	return &this
}

// GetApiVersion returns the ApiVersion field value if set, zero value otherwise.
func (o *CdxV1RedeemTokenRequest) GetApiVersion() string {
	if o == nil || o.ApiVersion == nil {
		var ret string
		return ret
	}
	return *o.ApiVersion
}

// GetApiVersionOk returns a tuple with the ApiVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CdxV1RedeemTokenRequest) GetApiVersionOk() (*string, bool) {
	if o == nil || o.ApiVersion == nil {
		return nil, false
	}
	return o.ApiVersion, true
}

// HasApiVersion returns a boolean if a field has been set.
func (o *CdxV1RedeemTokenRequest) HasApiVersion() bool {
	if o != nil && o.ApiVersion != nil {
		return true
	}

	return false
}

// SetApiVersion gets a reference to the given string and assigns it to the ApiVersion field.
func (o *CdxV1RedeemTokenRequest) SetApiVersion(v string) {
	o.ApiVersion = &v
}

// GetKind returns the Kind field value if set, zero value otherwise.
func (o *CdxV1RedeemTokenRequest) GetKind() string {
	if o == nil || o.Kind == nil {
		var ret string
		return ret
	}
	return *o.Kind
}

// GetKindOk returns a tuple with the Kind field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CdxV1RedeemTokenRequest) GetKindOk() (*string, bool) {
	if o == nil || o.Kind == nil {
		return nil, false
	}
	return o.Kind, true
}

// HasKind returns a boolean if a field has been set.
func (o *CdxV1RedeemTokenRequest) HasKind() bool {
	if o != nil && o.Kind != nil {
		return true
	}

	return false
}

// SetKind gets a reference to the given string and assigns it to the Kind field.
func (o *CdxV1RedeemTokenRequest) SetKind(v string) {
	o.Kind = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *CdxV1RedeemTokenRequest) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CdxV1RedeemTokenRequest) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CdxV1RedeemTokenRequest) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *CdxV1RedeemTokenRequest) SetId(v string) {
	o.Id = &v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *CdxV1RedeemTokenRequest) GetMetadata() ObjectMeta {
	if o == nil || o.Metadata == nil {
		var ret ObjectMeta
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CdxV1RedeemTokenRequest) GetMetadataOk() (*ObjectMeta, bool) {
	if o == nil || o.Metadata == nil {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *CdxV1RedeemTokenRequest) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given ObjectMeta and assigns it to the Metadata field.
func (o *CdxV1RedeemTokenRequest) SetMetadata(v ObjectMeta) {
	o.Metadata = &v
}

// GetToken returns the Token field value if set, zero value otherwise.
func (o *CdxV1RedeemTokenRequest) GetToken() string {
	if o == nil || o.Token == nil {
		var ret string
		return ret
	}
	return *o.Token
}

// GetTokenOk returns a tuple with the Token field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CdxV1RedeemTokenRequest) GetTokenOk() (*string, bool) {
	if o == nil || o.Token == nil {
		return nil, false
	}
	return o.Token, true
}

// HasToken returns a boolean if a field has been set.
func (o *CdxV1RedeemTokenRequest) HasToken() bool {
	if o != nil && o.Token != nil {
		return true
	}

	return false
}

// SetToken gets a reference to the given string and assigns it to the Token field.
func (o *CdxV1RedeemTokenRequest) SetToken(v string) {
	o.Token = &v
}

// GetAwsAccount returns the AwsAccount field value if set, zero value otherwise.
func (o *CdxV1RedeemTokenRequest) GetAwsAccount() string {
	if o == nil || o.AwsAccount == nil {
		var ret string
		return ret
	}
	return *o.AwsAccount
}

// GetAwsAccountOk returns a tuple with the AwsAccount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CdxV1RedeemTokenRequest) GetAwsAccountOk() (*string, bool) {
	if o == nil || o.AwsAccount == nil {
		return nil, false
	}
	return o.AwsAccount, true
}

// HasAwsAccount returns a boolean if a field has been set.
func (o *CdxV1RedeemTokenRequest) HasAwsAccount() bool {
	if o != nil && o.AwsAccount != nil {
		return true
	}

	return false
}

// SetAwsAccount gets a reference to the given string and assigns it to the AwsAccount field.
func (o *CdxV1RedeemTokenRequest) SetAwsAccount(v string) {
	o.AwsAccount = &v
}

// GetAzureSubscription returns the AzureSubscription field value if set, zero value otherwise.
func (o *CdxV1RedeemTokenRequest) GetAzureSubscription() string {
	if o == nil || o.AzureSubscription == nil {
		var ret string
		return ret
	}
	return *o.AzureSubscription
}

// GetAzureSubscriptionOk returns a tuple with the AzureSubscription field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CdxV1RedeemTokenRequest) GetAzureSubscriptionOk() (*string, bool) {
	if o == nil || o.AzureSubscription == nil {
		return nil, false
	}
	return o.AzureSubscription, true
}

// HasAzureSubscription returns a boolean if a field has been set.
func (o *CdxV1RedeemTokenRequest) HasAzureSubscription() bool {
	if o != nil && o.AzureSubscription != nil {
		return true
	}

	return false
}

// SetAzureSubscription gets a reference to the given string and assigns it to the AzureSubscription field.
func (o *CdxV1RedeemTokenRequest) SetAzureSubscription(v string) {
	o.AzureSubscription = &v
}

// GetGcpProject returns the GcpProject field value if set, zero value otherwise.
func (o *CdxV1RedeemTokenRequest) GetGcpProject() string {
	if o == nil || o.GcpProject == nil {
		var ret string
		return ret
	}
	return *o.GcpProject
}

// GetGcpProjectOk returns a tuple with the GcpProject field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CdxV1RedeemTokenRequest) GetGcpProjectOk() (*string, bool) {
	if o == nil || o.GcpProject == nil {
		return nil, false
	}
	return o.GcpProject, true
}

// HasGcpProject returns a boolean if a field has been set.
func (o *CdxV1RedeemTokenRequest) HasGcpProject() bool {
	if o != nil && o.GcpProject != nil {
		return true
	}

	return false
}

// SetGcpProject gets a reference to the given string and assigns it to the GcpProject field.
func (o *CdxV1RedeemTokenRequest) SetGcpProject(v string) {
	o.GcpProject = &v
}

// Redact resets all sensitive fields to their zero value.
func (o *CdxV1RedeemTokenRequest) Redact() {
	o.recurseRedact(o.ApiVersion)
	o.recurseRedact(o.Kind)
	o.recurseRedact(o.Id)
	o.recurseRedact(o.Metadata)
	o.recurseRedact(o.Token)
	o.recurseRedact(o.AwsAccount)
	o.recurseRedact(o.AzureSubscription)
	o.recurseRedact(o.GcpProject)
}

func (o *CdxV1RedeemTokenRequest) recurseRedact(v interface{}) {
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

func (o CdxV1RedeemTokenRequest) zeroField(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (o CdxV1RedeemTokenRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.ApiVersion != nil {
		toSerialize["api_version"] = o.ApiVersion
	}
	if o.Kind != nil {
		toSerialize["kind"] = o.Kind
	}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}
	if o.Token != nil {
		toSerialize["token"] = o.Token
	}
	if o.AwsAccount != nil {
		toSerialize["aws_account"] = o.AwsAccount
	}
	if o.AzureSubscription != nil {
		toSerialize["azure_subscription"] = o.AzureSubscription
	}
	if o.GcpProject != nil {
		toSerialize["gcp_project"] = o.GcpProject
	}
	return json.Marshal(toSerialize)
}

type NullableCdxV1RedeemTokenRequest struct {
	value *CdxV1RedeemTokenRequest
	isSet bool
}

func (v NullableCdxV1RedeemTokenRequest) Get() *CdxV1RedeemTokenRequest {
	return v.value
}

func (v *NullableCdxV1RedeemTokenRequest) Set(val *CdxV1RedeemTokenRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableCdxV1RedeemTokenRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableCdxV1RedeemTokenRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCdxV1RedeemTokenRequest(val *CdxV1RedeemTokenRequest) *NullableCdxV1RedeemTokenRequest {
	return &NullableCdxV1RedeemTokenRequest{value: val, isSet: true}
}

func (v NullableCdxV1RedeemTokenRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCdxV1RedeemTokenRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
