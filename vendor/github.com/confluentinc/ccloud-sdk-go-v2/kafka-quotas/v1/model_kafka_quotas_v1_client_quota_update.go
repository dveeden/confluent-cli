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

// KafkaQuotasV1ClientQuotaUpdate `ClientQuota` objects represent Client Quotas you can set at the service account level.  The API allows you to list, create, read, update, and delete your client quotas.   Related guide: [Client Quotas in Confluent Cloud](https://docs.confluent.io/cloud/current/clusters/client-quotas.html).  ## The Client Quotas Model <SchemaDefinition schemaRef=\"#/components/schemas/kafka-quotas.v1.ClientQuota\" />
type KafkaQuotasV1ClientQuotaUpdate struct {
	// APIVersion defines the schema version of this representation of a resource.
	ApiVersion *string `json:"api_version,omitempty"`
	// Kind defines the object this REST resource represents.
	Kind *string `json:"kind,omitempty"`
	// ID is the \"natural identifier\" for an object within its scope/namespace; it is normally unique across time but not space. That is, you can assume that the ID will not be reclaimed and reused after an object is deleted (\"time\"); however, it may collide with IDs for other object `kinds` or objects of the same `kind` within a different scope/namespace (\"space\").
	Id       *string                             `json:"id,omitempty"`
	Metadata *ObjectMeta                         `json:"metadata,omitempty"`
	Spec     *KafkaQuotasV1ClientQuotaSpecUpdate `json:"spec,omitempty"`
}

// NewKafkaQuotasV1ClientQuotaUpdate instantiates a new KafkaQuotasV1ClientQuotaUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKafkaQuotasV1ClientQuotaUpdate() *KafkaQuotasV1ClientQuotaUpdate {
	this := KafkaQuotasV1ClientQuotaUpdate{}
	return &this
}

// NewKafkaQuotasV1ClientQuotaUpdateWithDefaults instantiates a new KafkaQuotasV1ClientQuotaUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKafkaQuotasV1ClientQuotaUpdateWithDefaults() *KafkaQuotasV1ClientQuotaUpdate {
	this := KafkaQuotasV1ClientQuotaUpdate{}
	return &this
}

// GetApiVersion returns the ApiVersion field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaUpdate) GetApiVersion() string {
	if o == nil || o.ApiVersion == nil {
		var ret string
		return ret
	}
	return *o.ApiVersion
}

// GetApiVersionOk returns a tuple with the ApiVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaUpdate) GetApiVersionOk() (*string, bool) {
	if o == nil || o.ApiVersion == nil {
		return nil, false
	}
	return o.ApiVersion, true
}

// HasApiVersion returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaUpdate) HasApiVersion() bool {
	if o != nil && o.ApiVersion != nil {
		return true
	}

	return false
}

// SetApiVersion gets a reference to the given string and assigns it to the ApiVersion field.
func (o *KafkaQuotasV1ClientQuotaUpdate) SetApiVersion(v string) {
	o.ApiVersion = &v
}

// GetKind returns the Kind field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaUpdate) GetKind() string {
	if o == nil || o.Kind == nil {
		var ret string
		return ret
	}
	return *o.Kind
}

// GetKindOk returns a tuple with the Kind field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaUpdate) GetKindOk() (*string, bool) {
	if o == nil || o.Kind == nil {
		return nil, false
	}
	return o.Kind, true
}

// HasKind returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaUpdate) HasKind() bool {
	if o != nil && o.Kind != nil {
		return true
	}

	return false
}

// SetKind gets a reference to the given string and assigns it to the Kind field.
func (o *KafkaQuotasV1ClientQuotaUpdate) SetKind(v string) {
	o.Kind = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaUpdate) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaUpdate) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaUpdate) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *KafkaQuotasV1ClientQuotaUpdate) SetId(v string) {
	o.Id = &v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaUpdate) GetMetadata() ObjectMeta {
	if o == nil || o.Metadata == nil {
		var ret ObjectMeta
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaUpdate) GetMetadataOk() (*ObjectMeta, bool) {
	if o == nil || o.Metadata == nil {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaUpdate) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given ObjectMeta and assigns it to the Metadata field.
func (o *KafkaQuotasV1ClientQuotaUpdate) SetMetadata(v ObjectMeta) {
	o.Metadata = &v
}

// GetSpec returns the Spec field value if set, zero value otherwise.
func (o *KafkaQuotasV1ClientQuotaUpdate) GetSpec() KafkaQuotasV1ClientQuotaSpecUpdate {
	if o == nil || o.Spec == nil {
		var ret KafkaQuotasV1ClientQuotaSpecUpdate
		return ret
	}
	return *o.Spec
}

// GetSpecOk returns a tuple with the Spec field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KafkaQuotasV1ClientQuotaUpdate) GetSpecOk() (*KafkaQuotasV1ClientQuotaSpecUpdate, bool) {
	if o == nil || o.Spec == nil {
		return nil, false
	}
	return o.Spec, true
}

// HasSpec returns a boolean if a field has been set.
func (o *KafkaQuotasV1ClientQuotaUpdate) HasSpec() bool {
	if o != nil && o.Spec != nil {
		return true
	}

	return false
}

// SetSpec gets a reference to the given KafkaQuotasV1ClientQuotaSpecUpdate and assigns it to the Spec field.
func (o *KafkaQuotasV1ClientQuotaUpdate) SetSpec(v KafkaQuotasV1ClientQuotaSpecUpdate) {
	o.Spec = &v
}

// Redact resets all sensitive fields to their zero value.
func (o *KafkaQuotasV1ClientQuotaUpdate) Redact() {
	o.recurseRedact(o.ApiVersion)
	o.recurseRedact(o.Kind)
	o.recurseRedact(o.Id)
	o.recurseRedact(o.Metadata)
	o.recurseRedact(o.Spec)
}

func (o *KafkaQuotasV1ClientQuotaUpdate) recurseRedact(v interface{}) {
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

func (o KafkaQuotasV1ClientQuotaUpdate) zeroField(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (o KafkaQuotasV1ClientQuotaUpdate) MarshalJSON() ([]byte, error) {
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
	if o.Spec != nil {
		toSerialize["spec"] = o.Spec
	}
	return json.Marshal(toSerialize)
}

type NullableKafkaQuotasV1ClientQuotaUpdate struct {
	value *KafkaQuotasV1ClientQuotaUpdate
	isSet bool
}

func (v NullableKafkaQuotasV1ClientQuotaUpdate) Get() *KafkaQuotasV1ClientQuotaUpdate {
	return v.value
}

func (v *NullableKafkaQuotasV1ClientQuotaUpdate) Set(val *KafkaQuotasV1ClientQuotaUpdate) {
	v.value = val
	v.isSet = true
}

func (v NullableKafkaQuotasV1ClientQuotaUpdate) IsSet() bool {
	return v.isSet
}

func (v *NullableKafkaQuotasV1ClientQuotaUpdate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKafkaQuotasV1ClientQuotaUpdate(val *KafkaQuotasV1ClientQuotaUpdate) *NullableKafkaQuotasV1ClientQuotaUpdate {
	return &NullableKafkaQuotasV1ClientQuotaUpdate{value: val, isSet: true}
}

func (v NullableKafkaQuotasV1ClientQuotaUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKafkaQuotasV1ClientQuotaUpdate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
