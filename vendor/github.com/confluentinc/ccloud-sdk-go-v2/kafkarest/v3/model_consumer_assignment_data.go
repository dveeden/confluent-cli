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

// ConsumerAssignmentData struct for ConsumerAssignmentData
type ConsumerAssignmentData struct {
	Kind            string           `json:"kind"`
	Metadata        ResourceMetadata `json:"metadata"`
	ClusterId       string           `json:"cluster_id"`
	ConsumerGroupId string           `json:"consumer_group_id"`
	ConsumerId      string           `json:"consumer_id"`
	TopicName       string           `json:"topic_name"`
	PartitionId     int32            `json:"partition_id"`
	Partition       Relationship     `json:"partition"`
	Lag             Relationship     `json:"lag"`
}

// NewConsumerAssignmentData instantiates a new ConsumerAssignmentData object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConsumerAssignmentData(kind string, metadata ResourceMetadata, clusterId string, consumerGroupId string, consumerId string, topicName string, partitionId int32, partition Relationship, lag Relationship) *ConsumerAssignmentData {
	this := ConsumerAssignmentData{}
	this.Kind = kind
	this.Metadata = metadata
	this.ClusterId = clusterId
	this.ConsumerGroupId = consumerGroupId
	this.ConsumerId = consumerId
	this.TopicName = topicName
	this.PartitionId = partitionId
	this.Partition = partition
	this.Lag = lag
	return &this
}

// NewConsumerAssignmentDataWithDefaults instantiates a new ConsumerAssignmentData object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConsumerAssignmentDataWithDefaults() *ConsumerAssignmentData {
	this := ConsumerAssignmentData{}
	return &this
}

// GetKind returns the Kind field value
func (o *ConsumerAssignmentData) GetKind() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Kind
}

// GetKindOk returns a tuple with the Kind field value
// and a boolean to check if the value has been set.
func (o *ConsumerAssignmentData) GetKindOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Kind, true
}

// SetKind sets field value
func (o *ConsumerAssignmentData) SetKind(v string) {
	o.Kind = v
}

// GetMetadata returns the Metadata field value
func (o *ConsumerAssignmentData) GetMetadata() ResourceMetadata {
	if o == nil {
		var ret ResourceMetadata
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *ConsumerAssignmentData) GetMetadataOk() (*ResourceMetadata, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *ConsumerAssignmentData) SetMetadata(v ResourceMetadata) {
	o.Metadata = v
}

// GetClusterId returns the ClusterId field value
func (o *ConsumerAssignmentData) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *ConsumerAssignmentData) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *ConsumerAssignmentData) SetClusterId(v string) {
	o.ClusterId = v
}

// GetConsumerGroupId returns the ConsumerGroupId field value
func (o *ConsumerAssignmentData) GetConsumerGroupId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ConsumerGroupId
}

// GetConsumerGroupIdOk returns a tuple with the ConsumerGroupId field value
// and a boolean to check if the value has been set.
func (o *ConsumerAssignmentData) GetConsumerGroupIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ConsumerGroupId, true
}

// SetConsumerGroupId sets field value
func (o *ConsumerAssignmentData) SetConsumerGroupId(v string) {
	o.ConsumerGroupId = v
}

// GetConsumerId returns the ConsumerId field value
func (o *ConsumerAssignmentData) GetConsumerId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ConsumerId
}

// GetConsumerIdOk returns a tuple with the ConsumerId field value
// and a boolean to check if the value has been set.
func (o *ConsumerAssignmentData) GetConsumerIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ConsumerId, true
}

// SetConsumerId sets field value
func (o *ConsumerAssignmentData) SetConsumerId(v string) {
	o.ConsumerId = v
}

// GetTopicName returns the TopicName field value
func (o *ConsumerAssignmentData) GetTopicName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TopicName
}

// GetTopicNameOk returns a tuple with the TopicName field value
// and a boolean to check if the value has been set.
func (o *ConsumerAssignmentData) GetTopicNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TopicName, true
}

// SetTopicName sets field value
func (o *ConsumerAssignmentData) SetTopicName(v string) {
	o.TopicName = v
}

// GetPartitionId returns the PartitionId field value
func (o *ConsumerAssignmentData) GetPartitionId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.PartitionId
}

// GetPartitionIdOk returns a tuple with the PartitionId field value
// and a boolean to check if the value has been set.
func (o *ConsumerAssignmentData) GetPartitionIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PartitionId, true
}

// SetPartitionId sets field value
func (o *ConsumerAssignmentData) SetPartitionId(v int32) {
	o.PartitionId = v
}

// GetPartition returns the Partition field value
func (o *ConsumerAssignmentData) GetPartition() Relationship {
	if o == nil {
		var ret Relationship
		return ret
	}

	return o.Partition
}

// GetPartitionOk returns a tuple with the Partition field value
// and a boolean to check if the value has been set.
func (o *ConsumerAssignmentData) GetPartitionOk() (*Relationship, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Partition, true
}

// SetPartition sets field value
func (o *ConsumerAssignmentData) SetPartition(v Relationship) {
	o.Partition = v
}

// GetLag returns the Lag field value
func (o *ConsumerAssignmentData) GetLag() Relationship {
	if o == nil {
		var ret Relationship
		return ret
	}

	return o.Lag
}

// GetLagOk returns a tuple with the Lag field value
// and a boolean to check if the value has been set.
func (o *ConsumerAssignmentData) GetLagOk() (*Relationship, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Lag, true
}

// SetLag sets field value
func (o *ConsumerAssignmentData) SetLag(v Relationship) {
	o.Lag = v
}

// Redact resets all sensitive fields to their zero value.
func (o *ConsumerAssignmentData) Redact() {
	o.recurseRedact(&o.Kind)
	o.recurseRedact(&o.Metadata)
	o.recurseRedact(&o.ClusterId)
	o.recurseRedact(&o.ConsumerGroupId)
	o.recurseRedact(&o.ConsumerId)
	o.recurseRedact(&o.TopicName)
	o.recurseRedact(&o.PartitionId)
	o.recurseRedact(&o.Partition)
	o.recurseRedact(&o.Lag)
}

func (o *ConsumerAssignmentData) recurseRedact(v interface{}) {
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

func (o ConsumerAssignmentData) zeroField(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (o ConsumerAssignmentData) MarshalJSON() ([]byte, error) {
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
		toSerialize["consumer_group_id"] = o.ConsumerGroupId
	}
	if true {
		toSerialize["consumer_id"] = o.ConsumerId
	}
	if true {
		toSerialize["topic_name"] = o.TopicName
	}
	if true {
		toSerialize["partition_id"] = o.PartitionId
	}
	if true {
		toSerialize["partition"] = o.Partition
	}
	if true {
		toSerialize["lag"] = o.Lag
	}
	return json.Marshal(toSerialize)
}

type NullableConsumerAssignmentData struct {
	value *ConsumerAssignmentData
	isSet bool
}

func (v NullableConsumerAssignmentData) Get() *ConsumerAssignmentData {
	return v.value
}

func (v *NullableConsumerAssignmentData) Set(val *ConsumerAssignmentData) {
	v.value = val
	v.isSet = true
}

func (v NullableConsumerAssignmentData) IsSet() bool {
	return v.isSet
}

func (v *NullableConsumerAssignmentData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConsumerAssignmentData(val *ConsumerAssignmentData) *NullableConsumerAssignmentData {
	return &NullableConsumerAssignmentData{value: val, isSet: true}
}

func (v NullableConsumerAssignmentData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConsumerAssignmentData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
