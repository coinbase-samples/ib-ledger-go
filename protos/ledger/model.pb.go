//*
// Copyright 2022 Coinbase Global, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.4
// source: protos/ledger/model.proto

package ledger

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TransactionType int32

const (
	TransactionType_TRANSFER TransactionType = 0
	TransactionType_ORDER    TransactionType = 1
	TransactionType_CONVERT  TransactionType = 2
)

// Enum value maps for TransactionType.
var (
	TransactionType_name = map[int32]string{
		0: "TRANSFER",
		1: "ORDER",
		2: "CONVERT",
	}
	TransactionType_value = map[string]int32{
		"TRANSFER": 0,
		"ORDER":    1,
		"CONVERT":  2,
	}
)

func (x TransactionType) Enum() *TransactionType {
	p := new(TransactionType)
	*p = x
	return p
}

func (x TransactionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TransactionType) Descriptor() protoreflect.EnumDescriptor {
	return file_protos_ledger_model_proto_enumTypes[0].Descriptor()
}

func (TransactionType) Type() protoreflect.EnumType {
	return &file_protos_ledger_model_proto_enumTypes[0]
}

func (x TransactionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TransactionType.Descriptor instead.
func (TransactionType) EnumDescriptor() ([]byte, []int) {
	return file_protos_ledger_model_proto_rawDescGZIP(), []int{0}
}

type TransactionStatus int32

const (
	TransactionStatus_COMPLETE TransactionStatus = 0
	TransactionStatus_FAILED   TransactionStatus = 1
	TransactionStatus_CANCELED TransactionStatus = 2
	TransactionStatus_PENDING  TransactionStatus = 3
)

// Enum value maps for TransactionStatus.
var (
	TransactionStatus_name = map[int32]string{
		0: "COMPLETE",
		1: "FAILED",
		2: "CANCELED",
		3: "PENDING",
	}
	TransactionStatus_value = map[string]int32{
		"COMPLETE": 0,
		"FAILED":   1,
		"CANCELED": 2,
		"PENDING":  3,
	}
)

func (x TransactionStatus) Enum() *TransactionStatus {
	p := new(TransactionStatus)
	*p = x
	return p
}

func (x TransactionStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TransactionStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_protos_ledger_model_proto_enumTypes[1].Descriptor()
}

func (TransactionStatus) Type() protoreflect.EnumType {
	return &file_protos_ledger_model_proto_enumTypes[1]
}

func (x TransactionStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TransactionStatus.Descriptor instead.
func (TransactionStatus) EnumDescriptor() ([]byte, []int) {
	return file_protos_ledger_model_proto_rawDescGZIP(), []int{1}
}

type Account struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	PortfolioId *string                `protobuf:"bytes,2,opt,name=portfolioId,proto3,oneof" json:"portfolioId,omitempty"`
	UserId      string                 `protobuf:"bytes,3,opt,name=userId,proto3" json:"userId,omitempty"`
	Currency    string                 `protobuf:"bytes,4,opt,name=currency,proto3" json:"currency,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_ledger_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_protos_ledger_model_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_protos_ledger_model_proto_rawDescGZIP(), []int{0}
}

func (x *Account) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Account) GetPortfolioId() string {
	if x != nil && x.PortfolioId != nil {
		return *x.PortfolioId
	}
	return ""
}

func (x *Account) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Account) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *Account) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type AccountBalance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        *string `protobuf:"bytes,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	Balance   string  `protobuf:"bytes,2,opt,name=balance,proto3" json:"balance,omitempty"`
	Hold      string  `protobuf:"bytes,3,opt,name=hold,proto3" json:"hold,omitempty"`
	Available string  `protobuf:"bytes,4,opt,name=available,proto3" json:"available,omitempty"`
	AccountId *string `protobuf:"bytes,5,opt,name=account_id,json=accountId,proto3,oneof" json:"account_id,omitempty"`
}

func (x *AccountBalance) Reset() {
	*x = AccountBalance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_ledger_model_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountBalance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountBalance) ProtoMessage() {}

func (x *AccountBalance) ProtoReflect() protoreflect.Message {
	mi := &file_protos_ledger_model_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountBalance.ProtoReflect.Descriptor instead.
func (*AccountBalance) Descriptor() ([]byte, []int) {
	return file_protos_ledger_model_proto_rawDescGZIP(), []int{1}
}

func (x *AccountBalance) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *AccountBalance) GetBalance() string {
	if x != nil {
		return x.Balance
	}
	return ""
}

func (x *AccountBalance) GetHold() string {
	if x != nil {
		return x.Hold
	}
	return ""
}

func (x *AccountBalance) GetAvailable() string {
	if x != nil {
		return x.Available
	}
	return ""
}

func (x *AccountBalance) GetAccountId() string {
	if x != nil && x.AccountId != nil {
		return *x.AccountId
	}
	return ""
}

type AccountAndBalance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId string                 `protobuf:"bytes,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
	Currency  string                 `protobuf:"bytes,2,opt,name=currency,proto3" json:"currency,omitempty"`
	Balance   string                 `protobuf:"bytes,3,opt,name=balance,proto3" json:"balance,omitempty"`
	Hold      string                 `protobuf:"bytes,4,opt,name=hold,proto3" json:"hold,omitempty"`
	Available string                 `protobuf:"bytes,5,opt,name=available,proto3" json:"available,omitempty"`
	BalanceAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=balanceAt,proto3" json:"balanceAt,omitempty"`
}

func (x *AccountAndBalance) Reset() {
	*x = AccountAndBalance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_ledger_model_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountAndBalance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountAndBalance) ProtoMessage() {}

func (x *AccountAndBalance) ProtoReflect() protoreflect.Message {
	mi := &file_protos_ledger_model_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountAndBalance.ProtoReflect.Descriptor instead.
func (*AccountAndBalance) Descriptor() ([]byte, []int) {
	return file_protos_ledger_model_proto_rawDescGZIP(), []int{2}
}

func (x *AccountAndBalance) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *AccountAndBalance) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *AccountAndBalance) GetBalance() string {
	if x != nil {
		return x.Balance
	}
	return ""
}

func (x *AccountAndBalance) GetHold() string {
	if x != nil {
		return x.Hold
	}
	return ""
}

func (x *AccountAndBalance) GetAvailable() string {
	if x != nil {
		return x.Available
	}
	return ""
}

func (x *AccountAndBalance) GetBalanceAt() *timestamppb.Timestamp {
	if x != nil {
		return x.BalanceAt
	}
	return nil
}

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	SenderId          string                 `protobuf:"bytes,2,opt,name=senderId,proto3" json:"senderId,omitempty"`
	ReceiverId        string                 `protobuf:"bytes,3,opt,name=receiverId,proto3" json:"receiverId,omitempty"`
	RequestId         string                 `protobuf:"bytes,4,opt,name=requestId,proto3" json:"requestId,omitempty"`
	CreatedAt         *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	TransactionType   TransactionType        `protobuf:"varint,6,opt,name=transactionType,proto3,enum=ledger.TransactionType" json:"transactionType,omitempty"`
	TransactionStatus TransactionStatus      `protobuf:"varint,7,opt,name=transactionStatus,proto3,enum=ledger.TransactionStatus" json:"transactionStatus,omitempty"`
	FinalizedAt       *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=finalized_at,json=finalizedAt,proto3,oneof" json:"finalized_at,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_ledger_model_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_protos_ledger_model_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_protos_ledger_model_proto_rawDescGZIP(), []int{3}
}

func (x *Transaction) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Transaction) GetSenderId() string {
	if x != nil {
		return x.SenderId
	}
	return ""
}

func (x *Transaction) GetReceiverId() string {
	if x != nil {
		return x.ReceiverId
	}
	return ""
}

func (x *Transaction) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *Transaction) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Transaction) GetTransactionType() TransactionType {
	if x != nil {
		return x.TransactionType
	}
	return TransactionType_TRANSFER
}

func (x *Transaction) GetTransactionStatus() TransactionStatus {
	if x != nil {
		return x.TransactionStatus
	}
	return TransactionStatus_COMPLETE
}

func (x *Transaction) GetFinalizedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.FinalizedAt
	}
	return nil
}

type Hold struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AccountId     string                 `protobuf:"bytes,2,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	TransactionId string                 `protobuf:"bytes,3,opt,name=transaction_id,json=transactionId,proto3" json:"transaction_id,omitempty"`
	Amount        string                 `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Direction     string                 `protobuf:"bytes,5,opt,name=direction,proto3" json:"direction,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	ReleasedAt    *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=released_at,json=releasedAt,proto3" json:"released_at,omitempty"`
}

func (x *Hold) Reset() {
	*x = Hold{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_ledger_model_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Hold) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Hold) ProtoMessage() {}

func (x *Hold) ProtoReflect() protoreflect.Message {
	mi := &file_protos_ledger_model_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Hold.ProtoReflect.Descriptor instead.
func (*Hold) Descriptor() ([]byte, []int) {
	return file_protos_ledger_model_proto_rawDescGZIP(), []int{4}
}

func (x *Hold) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Hold) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *Hold) GetTransactionId() string {
	if x != nil {
		return x.TransactionId
	}
	return ""
}

func (x *Hold) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *Hold) GetDirection() string {
	if x != nil {
		return x.Direction
	}
	return ""
}

func (x *Hold) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Hold) GetReleasedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ReleasedAt
	}
	return nil
}

type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// UUID
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Parent Portfolio Account UUID
	AccountId string `protobuf:"bytes,2,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// Transaction UUID
	TransactionId string `protobuf:"bytes,3,opt,name=transaction_id,json=transactionId,proto3" json:"transaction_id,omitempty"`
	// Amount in asset's smallest unit e.g. Satoshi for Bitcoin
	Amount string `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount,omitempty"`
	// CREDIT or DEBIT
	Direction string                 `protobuf:"bytes,5,opt,name=direction,proto3" json:"direction,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_ledger_model_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_protos_ledger_model_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_protos_ledger_model_proto_rawDescGZIP(), []int{5}
}

func (x *Entry) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Entry) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *Entry) GetTransactionId() string {
	if x != nil {
		return x.TransactionId
	}
	return ""
}

func (x *Entry) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *Entry) GetDirection() string {
	if x != nil {
		return x.Direction
	}
	return ""
}

func (x *Entry) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

var File_protos_ledger_model_proto protoreflect.FileDescriptor

var file_protos_ledger_model_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbf, 0x01, 0x0a, 0x07, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x25, 0x0a, 0x0b, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x49, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c,
	0x69, 0x6f, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x39, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x66,
	0x6f, 0x6c, 0x69, 0x6f, 0x49, 0x64, 0x22, 0xab, 0x01, 0x0a, 0x0e, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x02, 0x69, 0x64, 0x88, 0x01, 0x01, 0x12, 0x18,
	0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x6c, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x1c, 0x0a, 0x09,
	0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x22, 0x0a, 0x0a, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01,
	0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42, 0x05,
	0x0a, 0x03, 0x5f, 0x69, 0x64, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x22, 0xd3, 0x01, 0x0a, 0x11, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x41, 0x6e, 0x64, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f,
	0x6c, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65,
	0x12, 0x38, 0x0a, 0x09, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x41, 0x74, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x41, 0x74, 0x22, 0x93, 0x03, 0x0a, 0x0b, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x65,
	0x69, 0x76, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x41, 0x0a, 0x0f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65,
	0x72, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x0f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x47, 0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e,
	0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x42, 0x0a, 0x0c, 0x66,
	0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52,
	0x0b, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x42,
	0x0f, 0x0a, 0x0d, 0x5f, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x22, 0x8a, 0x02, 0x0a, 0x04, 0x48, 0x6f, 0x6c, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x3b, 0x0a, 0x0b, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x0a, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x64, 0x41, 0x74, 0x22, 0xce, 0x01,
	0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2a, 0x37,
	0x0a, 0x0f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0c, 0x0a, 0x08, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x46, 0x45, 0x52, 0x10, 0x00, 0x12,
	0x09, 0x0a, 0x05, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f,
	0x4e, 0x56, 0x45, 0x52, 0x54, 0x10, 0x02, 0x2a, 0x48, 0x0a, 0x11, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0c, 0x0a, 0x08,
	0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41,
	0x49, 0x4c, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c,
	0x45, 0x44, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10,
	0x03, 0x42, 0x19, 0x5a, 0x17, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x41, 0x70, 0x70, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_ledger_model_proto_rawDescOnce sync.Once
	file_protos_ledger_model_proto_rawDescData = file_protos_ledger_model_proto_rawDesc
)

func file_protos_ledger_model_proto_rawDescGZIP() []byte {
	file_protos_ledger_model_proto_rawDescOnce.Do(func() {
		file_protos_ledger_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_ledger_model_proto_rawDescData)
	})
	return file_protos_ledger_model_proto_rawDescData
}

var file_protos_ledger_model_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_protos_ledger_model_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_protos_ledger_model_proto_goTypes = []interface{}{
	(TransactionType)(0),          // 0: ledger.TransactionType
	(TransactionStatus)(0),        // 1: ledger.TransactionStatus
	(*Account)(nil),               // 2: ledger.Account
	(*AccountBalance)(nil),        // 3: ledger.AccountBalance
	(*AccountAndBalance)(nil),     // 4: ledger.AccountAndBalance
	(*Transaction)(nil),           // 5: ledger.Transaction
	(*Hold)(nil),                  // 6: ledger.Hold
	(*Entry)(nil),                 // 7: ledger.Entry
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
}
var file_protos_ledger_model_proto_depIdxs = []int32{
	8, // 0: ledger.Account.created_at:type_name -> google.protobuf.Timestamp
	8, // 1: ledger.AccountAndBalance.balanceAt:type_name -> google.protobuf.Timestamp
	8, // 2: ledger.Transaction.created_at:type_name -> google.protobuf.Timestamp
	0, // 3: ledger.Transaction.transactionType:type_name -> ledger.TransactionType
	1, // 4: ledger.Transaction.transactionStatus:type_name -> ledger.TransactionStatus
	8, // 5: ledger.Transaction.finalized_at:type_name -> google.protobuf.Timestamp
	8, // 6: ledger.Hold.created_at:type_name -> google.protobuf.Timestamp
	8, // 7: ledger.Hold.released_at:type_name -> google.protobuf.Timestamp
	8, // 8: ledger.Entry.created_at:type_name -> google.protobuf.Timestamp
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_protos_ledger_model_proto_init() }
func file_protos_ledger_model_proto_init() {
	if File_protos_ledger_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_ledger_model_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Account); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_ledger_model_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountBalance); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_ledger_model_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountAndBalance); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_ledger_model_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_ledger_model_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Hold); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_ledger_model_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_protos_ledger_model_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_protos_ledger_model_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_protos_ledger_model_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_ledger_model_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protos_ledger_model_proto_goTypes,
		DependencyIndexes: file_protos_ledger_model_proto_depIdxs,
		EnumInfos:         file_protos_ledger_model_proto_enumTypes,
		MessageInfos:      file_protos_ledger_model_proto_msgTypes,
	}.Build()
	File_protos_ledger_model_proto = out.File
	file_protos_ledger_model_proto_rawDesc = nil
	file_protos_ledger_model_proto_goTypes = nil
	file_protos_ledger_model_proto_depIdxs = nil
}
