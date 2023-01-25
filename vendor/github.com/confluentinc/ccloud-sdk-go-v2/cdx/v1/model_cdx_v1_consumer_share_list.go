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

// CdxV1ConsumerShareList Resources accessible by the consumer   ## The Consumer Shared Resources Model <SchemaDefinition schemaRef=\"#/components/schemas/cdx.v1.ConsumerSharedResource\" />
type CdxV1ConsumerShareList struct {
	// APIVersion defines the schema version of this representation of a resource.
	ApiVersion string `json:"api_version"`
	// Kind defines the object this REST resource represents.
	Kind     string   `json:"kind"`
	Metadata ListMeta `json:"metadata"`
	// A data property that contains an array of resource items. Each entry in the array is a separate resource.
	Data []CdxV1ConsumerShare `json:"data"`
}

// NewCdxV1ConsumerShareList instantiates a new CdxV1ConsumerShareList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCdxV1ConsumerShareList(apiVersion string, kind string, metadata ListMeta, data []CdxV1ConsumerShare) *CdxV1ConsumerShareList {
	this := CdxV1ConsumerShareList{}
	this.ApiVersion = apiVersion
	this.Kind = kind
	this.Metadata = metadata
	this.Data = data
	return &this
}

// NewCdxV1ConsumerShareListWithDefaults instantiates a new CdxV1ConsumerShareList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCdxV1ConsumerShareListWithDefaults() *CdxV1ConsumerShareList {
	this := CdxV1ConsumerShareList{}
	return &this
}

// GetApiVersion returns the ApiVersion field value
func (o *CdxV1ConsumerShareList) GetApiVersion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ApiVersion
}

// GetApiVersionOk returns a tuple with the ApiVersion field value
// and a boolean to check if the value has been set.
func (o *CdxV1ConsumerShareList) GetApiVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ApiVersion, true
}

// SetApiVersion sets field value
func (o *CdxV1ConsumerShareList) SetApiVersion(v string) {
	o.ApiVersion = v
}

// GetKind returns the Kind field value
func (o *CdxV1ConsumerShareList) GetKind() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Kind
}

// GetKindOk returns a tuple with the Kind field value
// and a boolean to check if the value has been set.
func (o *CdxV1ConsumerShareList) GetKindOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Kind, true
}

// SetKind sets field value
func (o *CdxV1ConsumerShareList) SetKind(v string) {
	o.Kind = v
}

// GetMetadata returns the Metadata field value
func (o *CdxV1ConsumerShareList) GetMetadata() ListMeta {
	if o == nil {
		var ret ListMeta
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *CdxV1ConsumerShareList) GetMetadataOk() (*ListMeta, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *CdxV1ConsumerShareList) SetMetadata(v ListMeta) {
	o.Metadata = v
}

// GetData returns the Data field value
func (o *CdxV1ConsumerShareList) GetData() []CdxV1ConsumerShare {
	if o == nil {
		var ret []CdxV1ConsumerShare
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *CdxV1ConsumerShareList) GetDataOk() (*[]CdxV1ConsumerShare, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *CdxV1ConsumerShareList) SetData(v []CdxV1ConsumerShare) {
	o.Data = v
}

// Redact resets all sensitive fields to their zero value.
func (o *CdxV1ConsumerShareList) Redact() {
	o.recurseRedact(&o.ApiVersion)
	o.recurseRedact(&o.Kind)
	o.recurseRedact(&o.Metadata)
	o.recurseRedact(&o.Data)
}

func (o *CdxV1ConsumerShareList) recurseRedact(v interface{}) {
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

func (o CdxV1ConsumerShareList) zeroField(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (o CdxV1ConsumerShareList) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["api_version"] = o.ApiVersion
	}
	if true {
		toSerialize["kind"] = o.Kind
	}
	if true {
		toSerialize["metadata"] = o.Metadata
	}
	if true {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableCdxV1ConsumerShareList struct {
	value *CdxV1ConsumerShareList
	isSet bool
}

func (v NullableCdxV1ConsumerShareList) Get() *CdxV1ConsumerShareList {
	return v.value
}

func (v *NullableCdxV1ConsumerShareList) Set(val *CdxV1ConsumerShareList) {
	v.value = val
	v.isSet = true
}

func (v NullableCdxV1ConsumerShareList) IsSet() bool {
	return v.isSet
}

func (v *NullableCdxV1ConsumerShareList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCdxV1ConsumerShareList(val *CdxV1ConsumerShareList) *NullableCdxV1ConsumerShareList {
	return &NullableCdxV1ConsumerShareList{value: val, isSet: true}
}

func (v NullableCdxV1ConsumerShareList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCdxV1ConsumerShareList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
