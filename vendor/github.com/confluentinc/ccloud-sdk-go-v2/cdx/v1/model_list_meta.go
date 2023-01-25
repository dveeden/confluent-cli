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

// ListMeta ListMeta describes metadata that resource collections may have
type ListMeta struct {
	// A link to the first page of results. If a response does not contain a first link, then direct navigation to the first page is not supported.
	First NullableString `json:"first,omitempty"`
	// A link to the last page of results. If a response does not contain a last link, then direct navigation to the last page is not supported.
	Last NullableString `json:"last,omitempty"`
	// A link to the previous page of results. If a response does not contain a prev link, then either there is no previous data or backwards traversal through the result set is not supported.
	Prev NullableString `json:"prev,omitempty"`
	// A link to the next page of results. If a response does not contain a next link, then there is no more data available.
	Next NullableString `json:"next,omitempty"`
	// Number of records in the full result set. This response may be paginated and have a smaller number of records.
	TotalSize *int32 `json:"total_size,omitempty"`
}

// NewListMeta instantiates a new ListMeta object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListMeta() *ListMeta {
	this := ListMeta{}
	return &this
}

// NewListMetaWithDefaults instantiates a new ListMeta object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListMetaWithDefaults() *ListMeta {
	this := ListMeta{}
	return &this
}

// GetFirst returns the First field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ListMeta) GetFirst() string {
	if o == nil || o.First.Get() == nil {
		var ret string
		return ret
	}
	return *o.First.Get()
}

// GetFirstOk returns a tuple with the First field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ListMeta) GetFirstOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.First.Get(), o.First.IsSet()
}

// HasFirst returns a boolean if a field has been set.
func (o *ListMeta) HasFirst() bool {
	if o != nil && o.First.IsSet() {
		return true
	}

	return false
}

// SetFirst gets a reference to the given NullableString and assigns it to the First field.
func (o *ListMeta) SetFirst(v string) {
	o.First.Set(&v)
}

// SetFirstNil sets the value for First to be an explicit nil
func (o *ListMeta) SetFirstNil() {
	o.First.Set(nil)
}

// UnsetFirst ensures that no value is present for First, not even an explicit nil
func (o *ListMeta) UnsetFirst() {
	o.First.Unset()
}

// GetLast returns the Last field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ListMeta) GetLast() string {
	if o == nil || o.Last.Get() == nil {
		var ret string
		return ret
	}
	return *o.Last.Get()
}

// GetLastOk returns a tuple with the Last field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ListMeta) GetLastOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Last.Get(), o.Last.IsSet()
}

// HasLast returns a boolean if a field has been set.
func (o *ListMeta) HasLast() bool {
	if o != nil && o.Last.IsSet() {
		return true
	}

	return false
}

// SetLast gets a reference to the given NullableString and assigns it to the Last field.
func (o *ListMeta) SetLast(v string) {
	o.Last.Set(&v)
}

// SetLastNil sets the value for Last to be an explicit nil
func (o *ListMeta) SetLastNil() {
	o.Last.Set(nil)
}

// UnsetLast ensures that no value is present for Last, not even an explicit nil
func (o *ListMeta) UnsetLast() {
	o.Last.Unset()
}

// GetPrev returns the Prev field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ListMeta) GetPrev() string {
	if o == nil || o.Prev.Get() == nil {
		var ret string
		return ret
	}
	return *o.Prev.Get()
}

// GetPrevOk returns a tuple with the Prev field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ListMeta) GetPrevOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Prev.Get(), o.Prev.IsSet()
}

// HasPrev returns a boolean if a field has been set.
func (o *ListMeta) HasPrev() bool {
	if o != nil && o.Prev.IsSet() {
		return true
	}

	return false
}

// SetPrev gets a reference to the given NullableString and assigns it to the Prev field.
func (o *ListMeta) SetPrev(v string) {
	o.Prev.Set(&v)
}

// SetPrevNil sets the value for Prev to be an explicit nil
func (o *ListMeta) SetPrevNil() {
	o.Prev.Set(nil)
}

// UnsetPrev ensures that no value is present for Prev, not even an explicit nil
func (o *ListMeta) UnsetPrev() {
	o.Prev.Unset()
}

// GetNext returns the Next field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ListMeta) GetNext() string {
	if o == nil || o.Next.Get() == nil {
		var ret string
		return ret
	}
	return *o.Next.Get()
}

// GetNextOk returns a tuple with the Next field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ListMeta) GetNextOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Next.Get(), o.Next.IsSet()
}

// HasNext returns a boolean if a field has been set.
func (o *ListMeta) HasNext() bool {
	if o != nil && o.Next.IsSet() {
		return true
	}

	return false
}

// SetNext gets a reference to the given NullableString and assigns it to the Next field.
func (o *ListMeta) SetNext(v string) {
	o.Next.Set(&v)
}

// SetNextNil sets the value for Next to be an explicit nil
func (o *ListMeta) SetNextNil() {
	o.Next.Set(nil)
}

// UnsetNext ensures that no value is present for Next, not even an explicit nil
func (o *ListMeta) UnsetNext() {
	o.Next.Unset()
}

// GetTotalSize returns the TotalSize field value if set, zero value otherwise.
func (o *ListMeta) GetTotalSize() int32 {
	if o == nil || o.TotalSize == nil {
		var ret int32
		return ret
	}
	return *o.TotalSize
}

// GetTotalSizeOk returns a tuple with the TotalSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListMeta) GetTotalSizeOk() (*int32, bool) {
	if o == nil || o.TotalSize == nil {
		return nil, false
	}
	return o.TotalSize, true
}

// HasTotalSize returns a boolean if a field has been set.
func (o *ListMeta) HasTotalSize() bool {
	if o != nil && o.TotalSize != nil {
		return true
	}

	return false
}

// SetTotalSize gets a reference to the given int32 and assigns it to the TotalSize field.
func (o *ListMeta) SetTotalSize(v int32) {
	o.TotalSize = &v
}

// Redact resets all sensitive fields to their zero value.
func (o *ListMeta) Redact() {
	o.recurseRedact(o.First)
	o.recurseRedact(o.Last)
	o.recurseRedact(o.Prev)
	o.recurseRedact(o.Next)
	o.recurseRedact(o.TotalSize)
}

func (o *ListMeta) recurseRedact(v interface{}) {
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

func (o ListMeta) zeroField(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (o ListMeta) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.First.IsSet() {
		toSerialize["first"] = o.First.Get()
	}
	if o.Last.IsSet() {
		toSerialize["last"] = o.Last.Get()
	}
	if o.Prev.IsSet() {
		toSerialize["prev"] = o.Prev.Get()
	}
	if o.Next.IsSet() {
		toSerialize["next"] = o.Next.Get()
	}
	if o.TotalSize != nil {
		toSerialize["total_size"] = o.TotalSize
	}
	return json.Marshal(toSerialize)
}

type NullableListMeta struct {
	value *ListMeta
	isSet bool
}

func (v NullableListMeta) Get() *ListMeta {
	return v.value
}

func (v *NullableListMeta) Set(val *ListMeta) {
	v.value = val
	v.isSet = true
}

func (v NullableListMeta) IsSet() bool {
	return v.isSet
}

func (v *NullableListMeta) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListMeta(val *ListMeta) *NullableListMeta {
	return &NullableListMeta{value: val, isSet: true}
}

func (v NullableListMeta) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListMeta) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
