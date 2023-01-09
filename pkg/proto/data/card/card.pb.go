// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.17.3
// source: pkg/proto/data/card/card.proto

package card_store

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_data_card_card_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_data_card_card_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_pkg_proto_data_card_card_proto_rawDescGZIP(), []int{0}
}

type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetaInfo string `protobuf:"bytes,1,opt,name=metaInfo,proto3" json:"metaInfo,omitempty"`
	Number   string `protobuf:"bytes,2,opt,name=number,proto3" json:"number,omitempty"`
	Period   string `protobuf:"bytes,3,opt,name=period,proto3" json:"period,omitempty"`
	CVV      string `protobuf:"bytes,4,opt,name=CVV,proto3" json:"CVV,omitempty"`
	FullName string `protobuf:"bytes,5,opt,name=fullName,proto3" json:"fullName,omitempty"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_data_card_card_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_data_card_card_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_data_card_card_proto_rawDescGZIP(), []int{1}
}

func (x *CreateRequest) GetMetaInfo() string {
	if x != nil {
		return x.MetaInfo
	}
	return ""
}

func (x *CreateRequest) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *CreateRequest) GetPeriod() string {
	if x != nil {
		return x.Period
	}
	return ""
}

func (x *CreateRequest) GetCVV() string {
	if x != nil {
		return x.CVV
	}
	return ""
}

func (x *CreateRequest) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

type ChangeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetaInfo string `protobuf:"bytes,1,opt,name=metaInfo,proto3" json:"metaInfo,omitempty"`
	Number   string `protobuf:"bytes,2,opt,name=number,proto3" json:"number,omitempty"`
	Period   string `protobuf:"bytes,3,opt,name=period,proto3" json:"period,omitempty"`
	CVV      string `protobuf:"bytes,4,opt,name=CVV,proto3" json:"CVV,omitempty"`
	FullName string `protobuf:"bytes,5,opt,name=fullName,proto3" json:"fullName,omitempty"`
}

func (x *ChangeRequest) Reset() {
	*x = ChangeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_data_card_card_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeRequest) ProtoMessage() {}

func (x *ChangeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_data_card_card_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeRequest.ProtoReflect.Descriptor instead.
func (*ChangeRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_data_card_card_proto_rawDescGZIP(), []int{2}
}

func (x *ChangeRequest) GetMetaInfo() string {
	if x != nil {
		return x.MetaInfo
	}
	return ""
}

func (x *ChangeRequest) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *ChangeRequest) GetPeriod() string {
	if x != nil {
		return x.Period
	}
	return ""
}

func (x *ChangeRequest) GetCVV() string {
	if x != nil {
		return x.CVV
	}
	return ""
}

func (x *ChangeRequest) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetaInfo string `protobuf:"bytes,1,opt,name=metaInfo,proto3" json:"metaInfo,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_data_card_card_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_data_card_card_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_data_card_card_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteRequest) GetMetaInfo() string {
	if x != nil {
		return x.MetaInfo
	}
	return ""
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetaInfo string `protobuf:"bytes,1,opt,name=metaInfo,proto3" json:"metaInfo,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_data_card_card_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_data_card_card_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_data_card_card_proto_rawDescGZIP(), []int{4}
}

func (x *GetRequest) GetMetaInfo() string {
	if x != nil {
		return x.MetaInfo
	}
	return ""
}

type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number   string `protobuf:"bytes,2,opt,name=number,proto3" json:"number,omitempty"`
	Period   string `protobuf:"bytes,3,opt,name=period,proto3" json:"period,omitempty"`
	CVV      string `protobuf:"bytes,4,opt,name=CVV,proto3" json:"CVV,omitempty"`
	FullName string `protobuf:"bytes,5,opt,name=fullName,proto3" json:"fullName,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_data_card_card_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_data_card_card_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_data_card_card_proto_rawDescGZIP(), []int{5}
}

func (x *GetResponse) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *GetResponse) GetPeriod() string {
	if x != nil {
		return x.Period
	}
	return ""
}

func (x *GetResponse) GetCVV() string {
	if x != nil {
		return x.CVV
	}
	return ""
}

func (x *GetResponse) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

var File_pkg_proto_data_card_card_proto protoreflect.FileDescriptor

var file_pkg_proto_data_card_card_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61,
	0x2f, 0x63, 0x61, 0x72, 0x64, 0x2f, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x04, 0x63, 0x61, 0x72, 0x64, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x89, 0x01, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a,
	0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x43, 0x56, 0x56, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x43, 0x56, 0x56, 0x12,
	0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x89, 0x01, 0x0a, 0x0d,
	0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x43, 0x56, 0x56,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x43, 0x56, 0x56, 0x12, 0x1a, 0x0a, 0x08, 0x66,
	0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x2b, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x49, 0x6e, 0x66, 0x6f, 0x22, 0x28, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x6b,
	0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x43, 0x56, 0x56, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x43, 0x56, 0x56, 0x12,
	0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x32, 0xbd, 0x01, 0x0a, 0x0b,
	0x43, 0x61, 0x72, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x13, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b, 0x2e, 0x63, 0x61, 0x72,
	0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x2a, 0x0a, 0x06, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x12, 0x13, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x12, 0x2a, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x13, 0x2e,
	0x63, 0x61, 0x72, 0x64, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0b, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x2a, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x10, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x14, 0x5a, 0x12, 0x2e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_data_card_card_proto_rawDescOnce sync.Once
	file_pkg_proto_data_card_card_proto_rawDescData = file_pkg_proto_data_card_card_proto_rawDesc
)

func file_pkg_proto_data_card_card_proto_rawDescGZIP() []byte {
	file_pkg_proto_data_card_card_proto_rawDescOnce.Do(func() {
		file_pkg_proto_data_card_card_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_data_card_card_proto_rawDescData)
	})
	return file_pkg_proto_data_card_card_proto_rawDescData
}

var file_pkg_proto_data_card_card_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pkg_proto_data_card_card_proto_goTypes = []interface{}{
	(*Empty)(nil),         // 0: card.Empty
	(*CreateRequest)(nil), // 1: card.CreateRequest
	(*ChangeRequest)(nil), // 2: card.ChangeRequest
	(*DeleteRequest)(nil), // 3: card.DeleteRequest
	(*GetRequest)(nil),    // 4: card.GetRequest
	(*GetResponse)(nil),   // 5: card.GetResponse
}
var file_pkg_proto_data_card_card_proto_depIdxs = []int32{
	1, // 0: card.CardService.Create:input_type -> card.CreateRequest
	2, // 1: card.CardService.Change:input_type -> card.ChangeRequest
	3, // 2: card.CardService.Delete:input_type -> card.DeleteRequest
	4, // 3: card.CardService.Get:input_type -> card.GetRequest
	0, // 4: card.CardService.Create:output_type -> card.Empty
	0, // 5: card.CardService.Change:output_type -> card.Empty
	0, // 6: card.CardService.Delete:output_type -> card.Empty
	5, // 7: card.CardService.Get:output_type -> card.GetResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_proto_data_card_card_proto_init() }
func file_pkg_proto_data_card_card_proto_init() {
	if File_pkg_proto_data_card_card_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_data_card_card_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_pkg_proto_data_card_card_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRequest); i {
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
		file_pkg_proto_data_card_card_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeRequest); i {
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
		file_pkg_proto_data_card_card_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
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
		file_pkg_proto_data_card_card_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
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
		file_pkg_proto_data_card_card_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_proto_data_card_card_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_proto_data_card_card_proto_goTypes,
		DependencyIndexes: file_pkg_proto_data_card_card_proto_depIdxs,
		MessageInfos:      file_pkg_proto_data_card_card_proto_msgTypes,
	}.Build()
	File_pkg_proto_data_card_card_proto = out.File
	file_pkg_proto_data_card_card_proto_rawDesc = nil
	file_pkg_proto_data_card_card_proto_goTypes = nil
	file_pkg_proto_data_card_card_proto_depIdxs = nil
}
