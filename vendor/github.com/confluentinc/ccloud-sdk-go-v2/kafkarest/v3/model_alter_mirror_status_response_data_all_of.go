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

// AlterMirrorStatusResponseDataAllOf struct for AlterMirrorStatusResponseDataAllOf
type AlterMirrorStatusResponseDataAllOf struct {
	MirrorTopicName string         `json:"mirror_topic_name"`
	ErrorMessage    NullableString `json:"error_message"`
	ErrorCode       NullableInt32  `json:"error_code"`
	MirrorLags      MirrorLags     `json:"mirror_lags"`
}

// NewAlterMirrorStatusResponseDataAllOf instantiates a new AlterMirrorStatusResponseDataAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAlterMirrorStatusResponseDataAllOf(mirrorTopicName string, errorMessage NullableString, errorCode NullableInt32, mirrorLags MirrorLags) *AlterMirrorStatusResponseDataAllOf {
	this := AlterMirrorStatusResponseDataAllOf{}
	this.MirrorTopicName = mirrorTopicName
	this.ErrorMessage = errorMessage
	this.ErrorCode = errorCode
	this.MirrorLags = mirrorLags
	return &this
}

// NewAlterMirrorStatusResponseDataAllOfWithDefaults instantiates a new AlterMirrorStatusResponseDataAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAlterMirrorStatusResponseDataAllOfWithDefaults() *AlterMirrorStatusResponseDataAllOf {
	this := AlterMirrorStatusResponseDataAllOf{}
	return &this
}

// GetMirrorTopicName returns the MirrorTopicName field value
func (o *AlterMirrorStatusResponseDataAllOf) GetMirrorTopicName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.MirrorTopicName
}

// GetMirrorTopicNameOk returns a tuple with the MirrorTopicName field value
// and a boolean to check if the value has been set.
func (o *AlterMirrorStatusResponseDataAllOf) GetMirrorTopicNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MirrorTopicName, true
}

// SetMirrorTopicName sets field value
func (o *AlterMirrorStatusResponseDataAllOf) SetMirrorTopicName(v string) {
	o.MirrorTopicName = v
}

// GetErrorMessage returns the ErrorMessage field value
// If the value is explicit nil, the zero value for string will be returned
func (o *AlterMirrorStatusResponseDataAllOf) GetErrorMessage() string {
	if o == nil || o.ErrorMessage.Get() == nil {
		var ret string
		return ret
	}

	return *o.ErrorMessage.Get()
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AlterMirrorStatusResponseDataAllOf) GetErrorMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.ErrorMessage.Get(), o.ErrorMessage.IsSet()
}

// SetErrorMessage sets field value
func (o *AlterMirrorStatusResponseDataAllOf) SetErrorMessage(v string) {
	o.ErrorMessage.Set(&v)
}

// GetErrorCode returns the ErrorCode field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *AlterMirrorStatusResponseDataAllOf) GetErrorCode() int32 {
	if o == nil || o.ErrorCode.Get() == nil {
		var ret int32
		return ret
	}

	return *o.ErrorCode.Get()
}

// GetErrorCodeOk returns a tuple with the ErrorCode field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AlterMirrorStatusResponseDataAllOf) GetErrorCodeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ErrorCode.Get(), o.ErrorCode.IsSet()
}

// SetErrorCode sets field value
func (o *AlterMirrorStatusResponseDataAllOf) SetErrorCode(v int32) {
	o.ErrorCode.Set(&v)
}

// GetMirrorLags returns the MirrorLags field value
func (o *AlterMirrorStatusResponseDataAllOf) GetMirrorLags() MirrorLags {
	if o == nil {
		var ret MirrorLags
		return ret
	}

	return o.MirrorLags
}

// GetMirrorLagsOk returns a tuple with the MirrorLags field value
// and a boolean to check if the value has been set.
func (o *AlterMirrorStatusResponseDataAllOf) GetMirrorLagsOk() (*MirrorLags, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MirrorLags, true
}

// SetMirrorLags sets field value
func (o *AlterMirrorStatusResponseDataAllOf) SetMirrorLags(v MirrorLags) {
	o.MirrorLags = v
}

// Redact resets all sensitive fields to their zero value.
func (o *AlterMirrorStatusResponseDataAllOf) Redact() {
	o.recurseRedact(&o.MirrorTopicName)
	o.recurseRedact(&o.ErrorMessage)
	o.recurseRedact(&o.ErrorCode)
	o.recurseRedact(&o.MirrorLags)
}

func (o *AlterMirrorStatusResponseDataAllOf) recurseRedact(v interface{}) {
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

func (o AlterMirrorStatusResponseDataAllOf) zeroField(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (o AlterMirrorStatusResponseDataAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["mirror_topic_name"] = o.MirrorTopicName
	}
	if true {
		toSerialize["error_message"] = o.ErrorMessage.Get()
	}
	if true {
		toSerialize["error_code"] = o.ErrorCode.Get()
	}
	if true {
		toSerialize["mirror_lags"] = o.MirrorLags
	}
	return json.Marshal(toSerialize)
}

type NullableAlterMirrorStatusResponseDataAllOf struct {
	value *AlterMirrorStatusResponseDataAllOf
	isSet bool
}

func (v NullableAlterMirrorStatusResponseDataAllOf) Get() *AlterMirrorStatusResponseDataAllOf {
	return v.value
}

func (v *NullableAlterMirrorStatusResponseDataAllOf) Set(val *AlterMirrorStatusResponseDataAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableAlterMirrorStatusResponseDataAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableAlterMirrorStatusResponseDataAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAlterMirrorStatusResponseDataAllOf(val *AlterMirrorStatusResponseDataAllOf) *NullableAlterMirrorStatusResponseDataAllOf {
	return &NullableAlterMirrorStatusResponseDataAllOf{value: val, isSet: true}
}

func (v NullableAlterMirrorStatusResponseDataAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAlterMirrorStatusResponseDataAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
