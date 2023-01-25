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
Kafka Quotas Management API

API to manage various Kafka Quotas.

API version: 0.0.1-alpha1
Contact: kafka-cloud-fundament-aaaacmo35d4fp7t7p45tvuw6uq@confluent.slack.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package v1

import (
	"encoding/json"
)

import (
	"reflect"
)

// KafkaQuotasV1ClientQuotaSpec The desired state of the Client Quota
type KafkaQuotasV1ClientQuotaSpec struct {
	// The name of the client quota.
	DisplayName *string `json:"display_name,omitempty"`
	// A human readable description for the client quota.
	Description *string `json:"description,omitempty"`
	// Throughput for the client quota.
	Throughput *KafkaQuotasV1Throughput `json:"throughput,omitempty"`
	// The ID of the Dedicated Kafka cluster where the client quota is applied.
	Cluster *EnvScopedObjectReference `json:"cluster,omitempty"`
	// A list of principals to apply a client quota to. Use \"<default>\" to apply a client quota to all service accounts (see [Control application usage with Client Quotas](https://docs.confluent.io/cloud/current/clusters/client-quotas.html#control-application-usage-with-client-quotas) for more details).
	Principals *[]GlobalObjectReference `json:"principals,omitempty"`
	// The environment to which this belongs.
	Environment *GlobalObjectReference `json:"environment,omitempty"`
}

// NewKafkaQuotasV1ClientQuotaSpec instantiates a new KafkaQuotasV1ClientQuotaSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKafkaQuotasV1ClientQuotaSpec() *KafkaQuotasV1ClientQuotaSpec {
	this := KafkaQuotasV1ClientQuotaSpec{}
	return &this
}

// NewKafkaQuotasV1ClientQuotaSpecWithDefaults instantiates a new KafkaQuotasV1ClientQuotaSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKafkaQuotasV1ClientQuotaSpecWithDefaults() *KafkaQuotasV1ClientQuotaSpec {
	this := KafkaQuotasV1ClientQuotaSpec{}
	return &this
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaSpec) GetDisplayName() string {
	if o == nil || o.DisplayName == nil {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) GetDisplayNameOk() (*string, bool) {
	if o == nil || o.DisplayName == nil {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) HasDisplayName() bool {
	if o != nil && o.DisplayName != nil {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *KafkaQuotasV1ClientQuotaSpec) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaSpec) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *KafkaQuotasV1ClientQuotaSpec) SetDescription(v string) {
	o.Description = &v
}

// GetThroughput returns the Throughput field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaSpec) GetThroughput() KafkaQuotasV1Throughput {
	if o == nil || o.Throughput == nil {
		var ret KafkaQuotasV1Throughput
		return ret
	}
	return *o.Throughput
}

// GetThroughputOk returns a tuple with the Throughput field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) GetThroughputOk() (*KafkaQuotasV1Throughput, bool) {
	if o == nil || o.Throughput == nil {
		return nil, false
	}
	return o.Throughput, true
}

// HasThroughput returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) HasThroughput() bool {
	if o != nil && o.Throughput != nil {
		return true
	}

	return false
}

// SetThroughput gets a reference to the given KafkaQuotasV1Throughput and assigns it to the Throughput field.
func (o *KafkaQuotasV1ClientQuotaSpec) SetThroughput(v KafkaQuotasV1Throughput) {
	o.Throughput = &v
}

// GetCluster returns the Cluster field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaSpec) GetCluster() EnvScopedObjectReference {
	if o == nil || o.Cluster == nil {
		var ret EnvScopedObjectReference
		return ret
	}
	return *o.Cluster
}

// GetClusterOk returns a tuple with the Cluster field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) GetClusterOk() (*EnvScopedObjectReference, bool) {
	if o == nil || o.Cluster == nil {
		return nil, false
	}
	return o.Cluster, true
}

// HasCluster returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) HasCluster() bool {
	if o != nil && o.Cluster != nil {
		return true
	}

	return false
}

// SetCluster gets a reference to the given EnvScopedObjectReference and assigns it to the Cluster field.
func (o *KafkaQuotasV1ClientQuotaSpec) SetCluster(v EnvScopedObjectReference) {
	o.Cluster = &v
}

// GetPrincipals returns the Principals field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaSpec) GetPrincipals() []GlobalObjectReference {
	if o == nil || o.Principals == nil {
		var ret []GlobalObjectReference
		return ret
	}
	return *o.Principals
}

// GetPrincipalsOk returns a tuple with the Principals field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) GetPrincipalsOk() (*[]GlobalObjectReference, bool) {
	if o == nil || o.Principals == nil {
		return nil, false
	}
	return o.Principals, true
}

// HasPrincipals returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) HasPrincipals() bool {
	if o != nil && o.Principals != nil {
		return true
	}

	return false
}

// SetPrincipals gets a reference to the given []GlobalObjectReference and assigns it to the Principals field.
func (o *KafkaQuotasV1ClientQuotaSpec) SetPrincipals(v []GlobalObjectReference) {
	o.Principals = &v
}

// GetEnvironment returns the Environment field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaSpec) GetEnvironment() GlobalObjectReference {
	if o == nil || o.Environment == nil {
		var ret GlobalObjectReference
		return ret
	}
	return *o.Environment
}

// GetEnvironmentOk returns a tuple with the Environment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) GetEnvironmentOk() (*GlobalObjectReference, bool) {
	if o == nil || o.Environment == nil {
		return nil, false
	}
	return o.Environment, true
}

// HasEnvironment returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaSpec) HasEnvironment() bool {
	if o != nil && o.Environment != nil {
		return true
	}

	return false
}

// SetEnvironment gets a reference to the given GlobalObjectReference and assigns it to the Environment field.
func (o *KafkaQuotasV1ClientQuotaSpec) SetEnvironment(v GlobalObjectReference) {
	o.Environment = &v
}

// Redact resets all sensitive fields to their zero value.
func (o *KafkaQuotasV1ClientQuotaSpec) Redact() {
	o.recurseRedact(o.DisplayName)
	o.recurseRedact(o.Description)
	o.recurseRedact(o.Throughput)
	o.recurseRedact(o.Cluster)
	o.recurseRedact(o.Principals)
	o.recurseRedact(o.Environment)
}

func (o *KafkaQuotasV1ClientQuotaSpec) recurseRedact(v interface{}) {
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

func (o KafkaQuotasV1ClientQuotaSpec) zeroField(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (o KafkaQuotasV1ClientQuotaSpec) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.DisplayName != nil {
		toSerialize["display_name"] = o.DisplayName
	}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	if o.Throughput != nil {
		toSerialize["throughput"] = o.Throughput
	}
	if o.Cluster != nil {
		toSerialize["cluster"] = o.Cluster
	}
	if o.Principals != nil {
		toSerialize["principals"] = o.Principals
	}
	if o.Environment != nil {
		toSerialize["environment"] = o.Environment
	}
	return json.Marshal(toSerialize)
}

type NullableKafkaQuotasV1ClientQuotaSpec struct {
	value *KafkaQuotasV1ClientQuotaSpec
	isSet bool
}

func (v NullableKafkaQuotasV1ClientQuotaSpec) Get() *KafkaQuotasV1ClientQuotaSpec {
	return v.value
}

func (v *NullableKafkaQuotasV1ClientQuotaSpec) Set(val *KafkaQuotasV1ClientQuotaSpec) {
	v.value = val
	v.isSet = true
}

func (v NullableKafkaQuotasV1ClientQuotaSpec) IsSet() bool {
	return v.isSet
}

func (v *NullableKafkaQuotasV1ClientQuotaSpec) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKafkaQuotasV1ClientQuotaSpec(val *KafkaQuotasV1ClientQuotaSpec) *NullableKafkaQuotasV1ClientQuotaSpec {
	return &NullableKafkaQuotasV1ClientQuotaSpec{value: val, isSet: true}
}

func (v NullableKafkaQuotasV1ClientQuotaSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKafkaQuotasV1ClientQuotaSpec) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
