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
	"time"
)

import (
	"reflect"
)

// AnyUnevenLoadDataAllOf struct for AnyUnevenLoadDataAllOf
type AnyUnevenLoadDataAllOf struct {
	ClusterId      string `json:"cluster_id"`
	Status         string `json:"status"`
	PreviousStatus string `json:"previous_status"`
	// The date and time at which this task was created.
	StatusUpdatedAt time.Time `json:"status_updated_at"`
	// The date and time at which this task was created.
	PreviousStatusUpdatedAt time.Time      `json:"previous_status_updated_at"`
	ErrorCode               NullableInt32  `json:"error_code,omitempty"`
	ErrorMessage            NullableString `json:"error_message,omitempty"`
	BrokerTasks             Relationship   `json:"broker_tasks"`
}

// NewAnyUnevenLoadDataAllOf instantiates a new AnyUnevenLoadDataAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAnyUnevenLoadDataAllOf(clusterId string, status string, previousStatus string, statusUpdatedAt time.Time, previousStatusUpdatedAt time.Time, brokerTasks Relationship) *AnyUnevenLoadDataAllOf {
	this := AnyUnevenLoadDataAllOf{}
	this.ClusterId = clusterId
	this.Status = status
	this.PreviousStatus = previousStatus
	this.StatusUpdatedAt = statusUpdatedAt
	this.PreviousStatusUpdatedAt = previousStatusUpdatedAt
	this.BrokerTasks = brokerTasks
	return &this
}

// NewAnyUnevenLoadDataAllOfWithDefaults instantiates a new AnyUnevenLoadDataAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAnyUnevenLoadDataAllOfWithDefaults() *AnyUnevenLoadDataAllOf {
	this := AnyUnevenLoadDataAllOf{}
	return &this
}

// GetClusterId returns the ClusterId field value
func (o *AnyUnevenLoadDataAllOf) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *AnyUnevenLoadDataAllOf) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *AnyUnevenLoadDataAllOf) SetClusterId(v string) {
	o.ClusterId = v
}

// GetStatus returns the Status field value
func (o *AnyUnevenLoadDataAllOf) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *AnyUnevenLoadDataAllOf) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *AnyUnevenLoadDataAllOf) SetStatus(v string) {
	o.Status = v
}

// GetPreviousStatus returns the PreviousStatus field value
func (o *AnyUnevenLoadDataAllOf) GetPreviousStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PreviousStatus
}

// GetPreviousStatusOk returns a tuple with the PreviousStatus field value
// and a boolean to check if the value has been set.
func (o *AnyUnevenLoadDataAllOf) GetPreviousStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PreviousStatus, true
}

// SetPreviousStatus sets field value
func (o *AnyUnevenLoadDataAllOf) SetPreviousStatus(v string) {
	o.PreviousStatus = v
}

// GetStatusUpdatedAt returns the StatusUpdatedAt field value
func (o *AnyUnevenLoadDataAllOf) GetStatusUpdatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.StatusUpdatedAt
}

// GetStatusUpdatedAtOk returns a tuple with the StatusUpdatedAt field value
// and a boolean to check if the value has been set.
func (o *AnyUnevenLoadDataAllOf) GetStatusUpdatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.StatusUpdatedAt, true
}

// SetStatusUpdatedAt sets field value
func (o *AnyUnevenLoadDataAllOf) SetStatusUpdatedAt(v time.Time) {
	o.StatusUpdatedAt = v
}

// GetPreviousStatusUpdatedAt returns the PreviousStatusUpdatedAt field value
func (o *AnyUnevenLoadDataAllOf) GetPreviousStatusUpdatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.PreviousStatusUpdatedAt
}

// GetPreviousStatusUpdatedAtOk returns a tuple with the PreviousStatusUpdatedAt field value
// and a boolean to check if the value has been set.
func (o *AnyUnevenLoadDataAllOf) GetPreviousStatusUpdatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PreviousStatusUpdatedAt, true
}

// SetPreviousStatusUpdatedAt sets field value
func (o *AnyUnevenLoadDataAllOf) SetPreviousStatusUpdatedAt(v time.Time) {
	o.PreviousStatusUpdatedAt = v
}

// GetErrorCode returns the ErrorCode field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *AnyUnevenLoadDataAllOf) GetErrorCode() int32 {
	if o == nil || o.ErrorCode.Get() == nil {
		var ret int32
		return ret
	}
	return *o.ErrorCode.Get()
}

// GetErrorCodeOk returns a tuple with the ErrorCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AnyUnevenLoadDataAllOf) GetErrorCodeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ErrorCode.Get(), o.ErrorCode.IsSet()
}

// HasErrorCode returns a boolean if a field has been set.
func (o *AnyUnevenLoadDataAllOf) HasErrorCode() bool {
	if o != nil && o.ErrorCode.IsSet() {
		return true
	}

	return false
}

// SetErrorCode gets a reference to the given NullableInt32 and assigns it to the ErrorCode field.
func (o *AnyUnevenLoadDataAllOf) SetErrorCode(v int32) {
	o.ErrorCode.Set(&v)
}

// SetErrorCodeNil sets the value for ErrorCode to be an explicit nil
func (o *AnyUnevenLoadDataAllOf) SetErrorCodeNil() {
	o.ErrorCode.Set(nil)
}

// UnsetErrorCode ensures that no value is present for ErrorCode, not even an explicit nil
func (o *AnyUnevenLoadDataAllOf) UnsetErrorCode() {
	o.ErrorCode.Unset()
}

// GetErrorMessage returns the ErrorMessage field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *AnyUnevenLoadDataAllOf) GetErrorMessage() string {
	if o == nil || o.ErrorMessage.Get() == nil {
		var ret string
		return ret
	}
	return *o.ErrorMessage.Get()
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AnyUnevenLoadDataAllOf) GetErrorMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.ErrorMessage.Get(), o.ErrorMessage.IsSet()
}

// HasErrorMessage returns a boolean if a field has been set.
func (o *AnyUnevenLoadDataAllOf) HasErrorMessage() bool {
	if o != nil && o.ErrorMessage.IsSet() {
		return true
	}

	return false
}

// SetErrorMessage gets a reference to the given NullableString and assigns it to the ErrorMessage field.
func (o *AnyUnevenLoadDataAllOf) SetErrorMessage(v string) {
	o.ErrorMessage.Set(&v)
}

// SetErrorMessageNil sets the value for ErrorMessage to be an explicit nil
func (o *AnyUnevenLoadDataAllOf) SetErrorMessageNil() {
	o.ErrorMessage.Set(nil)
}

// UnsetErrorMessage ensures that no value is present for ErrorMessage, not even an explicit nil
func (o *AnyUnevenLoadDataAllOf) UnsetErrorMessage() {
	o.ErrorMessage.Unset()
}

// GetBrokerTasks returns the BrokerTasks field value
func (o *AnyUnevenLoadDataAllOf) GetBrokerTasks() Relationship {
	if o == nil {
		var ret Relationship
		return ret
	}

	return o.BrokerTasks
}

// GetBrokerTasksOk returns a tuple with the BrokerTasks field value
// and a boolean to check if the value has been set.
func (o *AnyUnevenLoadDataAllOf) GetBrokerTasksOk() (*Relationship, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BrokerTasks, true
}

// SetBrokerTasks sets field value
func (o *AnyUnevenLoadDataAllOf) SetBrokerTasks(v Relationship) {
	o.BrokerTasks = v
}

// Redact resets all sensitive fields to their zero value.
func (o *AnyUnevenLoadDataAllOf) Redact() {
	o.recurseRedact(&o.ClusterId)
	o.recurseRedact(&o.Status)
	o.recurseRedact(&o.PreviousStatus)
	o.recurseRedact(&o.StatusUpdatedAt)
	o.recurseRedact(&o.PreviousStatusUpdatedAt)
	o.recurseRedact(o.ErrorCode)
	o.recurseRedact(o.ErrorMessage)
	o.recurseRedact(&o.BrokerTasks)
}

func (o *AnyUnevenLoadDataAllOf) recurseRedact(v interface{}) {
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

func (o AnyUnevenLoadDataAllOf) zeroField(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (o AnyUnevenLoadDataAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["cluster_id"] = o.ClusterId
	}
	if true {
		toSerialize["status"] = o.Status
	}
	if true {
		toSerialize["previous_status"] = o.PreviousStatus
	}
	if true {
		toSerialize["status_updated_at"] = o.StatusUpdatedAt
	}
	if true {
		toSerialize["previous_status_updated_at"] = o.PreviousStatusUpdatedAt
	}
	if o.ErrorCode.IsSet() {
		toSerialize["error_code"] = o.ErrorCode.Get()
	}
	if o.ErrorMessage.IsSet() {
		toSerialize["error_message"] = o.ErrorMessage.Get()
	}
	if true {
		toSerialize["broker_tasks"] = o.BrokerTasks
	}
	return json.Marshal(toSerialize)
}

type NullableAnyUnevenLoadDataAllOf struct {
	value *AnyUnevenLoadDataAllOf
	isSet bool
}

func (v NullableAnyUnevenLoadDataAllOf) Get() *AnyUnevenLoadDataAllOf {
	return v.value
}

func (v *NullableAnyUnevenLoadDataAllOf) Set(val *AnyUnevenLoadDataAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableAnyUnevenLoadDataAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableAnyUnevenLoadDataAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAnyUnevenLoadDataAllOf(val *AnyUnevenLoadDataAllOf) *NullableAnyUnevenLoadDataAllOf {
	return &NullableAnyUnevenLoadDataAllOf{value: val, isSet: true}
}

func (v NullableAnyUnevenLoadDataAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAnyUnevenLoadDataAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
