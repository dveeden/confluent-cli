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
Stream Designer API

# Introduction  Stream Designer API provides resources/API for defining stream processing pipelines. Each pipeline describes a set of stream processing components, including connectors, topics, streams, tables, queries and schemas. The components in a pipeline need not exist as Confluent Cloud resources until the pipeline is activated.  This API defines operations to create, list, modify, manage and delete pipelines. 

API version: 0.0.1-alpha0
Contact: stream-designer@confluent.io
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package v1

import (
	"encoding/json"
)

import (
	"reflect"
)

// TypedGlobalObjectReference ObjectReference provides information for you to locate the referred object
type TypedGlobalObjectReference struct {
	// ID of the referred resource
	Id string `json:"id"`
	// API URL for accessing or modifying the referred object
	Related string `json:"related"`
	// CRN reference to the referred resource
	ResourceName string `json:"resource_name"`
	// API group and version of the referred resource
	ApiVersion *string `json:"api_version,omitempty"`
	// Kind of the referred resource
	Kind *string `json:"kind,omitempty"`
}

// NewTypedGlobalObjectReference instantiates a new TypedGlobalObjectReference object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTypedGlobalObjectReference(id string, related string, resourceName string) *TypedGlobalObjectReference {
	this := TypedGlobalObjectReference{}
	this.Id = id
	this.Related = related
	this.ResourceName = resourceName
	return &this
}

// NewTypedGlobalObjectReferenceWithDefaults instantiates a new TypedGlobalObjectReference object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTypedGlobalObjectReferenceWithDefaults() *TypedGlobalObjectReference {
	this := TypedGlobalObjectReference{}
	return &this
}

// GetId returns the Id field value
func (o *TypedGlobalObjectReference) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *TypedGlobalObjectReference) GetIdOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *TypedGlobalObjectReference) SetId(v string) {
	o.Id = v
}

// GetRelated returns the Related field value
func (o *TypedGlobalObjectReference) GetRelated() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Related
}

// GetRelatedOk returns a tuple with the Related field value
// and a boolean to check if the value has been set.
func (o *TypedGlobalObjectReference) GetRelatedOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Related, true
}

// SetRelated sets field value
func (o *TypedGlobalObjectReference) SetRelated(v string) {
	o.Related = v
}

// GetResourceName returns the ResourceName field value
func (o *TypedGlobalObjectReference) GetResourceName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ResourceName
}

// GetResourceNameOk returns a tuple with the ResourceName field value
// and a boolean to check if the value has been set.
func (o *TypedGlobalObjectReference) GetResourceNameOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.ResourceName, true
}

// SetResourceName sets field value
func (o *TypedGlobalObjectReference) SetResourceName(v string) {
	o.ResourceName = v
}

// GetApiVersion returns the ApiVersion field value if set, zero value otherwise.
func (o *TypedGlobalObjectReference) GetApiVersion() string {
	if o == nil || o.ApiVersion == nil {
		var ret string
		return ret
	}
	return *o.ApiVersion
}

// GetApiVersionOk returns a tuple with the ApiVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TypedGlobalObjectReference) GetApiVersionOk() (*string, bool) {
	if o == nil || o.ApiVersion == nil {
		return nil, false
	}
	return o.ApiVersion, true
}

// HasApiVersion returns a boolean if a field has been set.
func (o *TypedGlobalObjectReference) HasApiVersion() bool {
	if o != nil && o.ApiVersion != nil {
		return true
	}

	return false
}

// SetApiVersion gets a reference to the given string and assigns it to the ApiVersion field.
func (o *TypedGlobalObjectReference) SetApiVersion(v string) {
	o.ApiVersion = &v
}

// GetKind returns the Kind field value if set, zero value otherwise.
func (o *TypedGlobalObjectReference) GetKind() string {
	if o == nil || o.Kind == nil {
		var ret string
		return ret
	}
	return *o.Kind
}

// GetKindOk returns a tuple with the Kind field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TypedGlobalObjectReference) GetKindOk() (*string, bool) {
	if o == nil || o.Kind == nil {
		return nil, false
	}
	return o.Kind, true
}

// HasKind returns a boolean if a field has been set.
func (o *TypedGlobalObjectReference) HasKind() bool {
	if o != nil && o.Kind != nil {
		return true
	}

	return false
}

// SetKind gets a reference to the given string and assigns it to the Kind field.
func (o *TypedGlobalObjectReference) SetKind(v string) {
	o.Kind = &v
}

// Redact resets all sensitive fields to their zero value.
func (o *TypedGlobalObjectReference) Redact() {
    o.recurseRedact(&o.Id)
    o.recurseRedact(&o.Related)
    o.recurseRedact(&o.ResourceName)
    o.recurseRedact(o.ApiVersion)
    o.recurseRedact(o.Kind)
}

func (o *TypedGlobalObjectReference) recurseRedact(v interface{}) {
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

func (o TypedGlobalObjectReference) zeroField(v interface{}) {
    p := reflect.ValueOf(v).Elem()
    p.Set(reflect.Zero(p.Type()))
}

func (o TypedGlobalObjectReference) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["id"] = o.Id
	}
	if true {
		toSerialize["related"] = o.Related
	}
	if true {
		toSerialize["resource_name"] = o.ResourceName
	}
	if o.ApiVersion != nil {
		toSerialize["api_version"] = o.ApiVersion
	}
	if o.Kind != nil {
		toSerialize["kind"] = o.Kind
	}
	return json.Marshal(toSerialize)
}

type NullableTypedGlobalObjectReference struct {
	value *TypedGlobalObjectReference
	isSet bool
}

func (v NullableTypedGlobalObjectReference) Get() *TypedGlobalObjectReference {
	return v.value
}

func (v *NullableTypedGlobalObjectReference) Set(val *TypedGlobalObjectReference) {
	v.value = val
	v.isSet = true
}

func (v NullableTypedGlobalObjectReference) IsSet() bool {
	return v.isSet
}

func (v *NullableTypedGlobalObjectReference) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTypedGlobalObjectReference(val *TypedGlobalObjectReference) *NullableTypedGlobalObjectReference {
	return &NullableTypedGlobalObjectReference{value: val, isSet: true}
}

func (v NullableTypedGlobalObjectReference) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTypedGlobalObjectReference) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


