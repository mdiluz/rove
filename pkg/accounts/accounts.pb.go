// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.6.1
// source: accounts/accounts.proto

package accounts

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// RegisterInfo contains the information needed to register an account
type RegisterInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name for the account, must be unique
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *RegisterInfo) Reset() {
	*x = RegisterInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounts_accounts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterInfo) ProtoMessage() {}

func (x *RegisterInfo) ProtoReflect() protoreflect.Message {
	mi := &file_accounts_accounts_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterInfo.ProtoReflect.Descriptor instead.
func (*RegisterInfo) Descriptor() ([]byte, []int) {
	return file_accounts_accounts_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// RegisterResponse is the response information from registering an account
type RegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounts_accounts_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_accounts_accounts_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_accounts_accounts_proto_rawDescGZIP(), []int{1}
}

// DataKeyValue represents a simple key value pair to assign to an account
type DataKeyValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The account to assign the new key value pair to
	Account string `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	// The key value pair to assign
	Key   string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *DataKeyValue) Reset() {
	*x = DataKeyValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounts_accounts_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataKeyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataKeyValue) ProtoMessage() {}

func (x *DataKeyValue) ProtoReflect() protoreflect.Message {
	mi := &file_accounts_accounts_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataKeyValue.ProtoReflect.Descriptor instead.
func (*DataKeyValue) Descriptor() ([]byte, []int) {
	return file_accounts_accounts_proto_rawDescGZIP(), []int{2}
}

func (x *DataKeyValue) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *DataKeyValue) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *DataKeyValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// DataKeyResponse is a simple response
type DataKeyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DataKeyResponse) Reset() {
	*x = DataKeyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounts_accounts_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataKeyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataKeyResponse) ProtoMessage() {}

func (x *DataKeyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_accounts_accounts_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataKeyResponse.ProtoReflect.Descriptor instead.
func (*DataKeyResponse) Descriptor() ([]byte, []int) {
	return file_accounts_accounts_proto_rawDescGZIP(), []int{3}
}

// DataKey describes a simple key value with an account, for fetching
type DataKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The account to fetch data for
	Account string `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	// The key to fetch
	Key string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DataKey) Reset() {
	*x = DataKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounts_accounts_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataKey) ProtoMessage() {}

func (x *DataKey) ProtoReflect() protoreflect.Message {
	mi := &file_accounts_accounts_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataKey.ProtoReflect.Descriptor instead.
func (*DataKey) Descriptor() ([]byte, []int) {
	return file_accounts_accounts_proto_rawDescGZIP(), []int{4}
}

func (x *DataKey) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *DataKey) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

// DataResponse describes a data fetch response
type DataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The value of the key
	Value string `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *DataResponse) Reset() {
	*x = DataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounts_accounts_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataResponse) ProtoMessage() {}

func (x *DataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_accounts_accounts_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataResponse.ProtoReflect.Descriptor instead.
func (*DataResponse) Descriptor() ([]byte, []int) {
	return file_accounts_accounts_proto_rawDescGZIP(), []int{5}
}

func (x *DataResponse) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_accounts_accounts_proto protoreflect.FileDescriptor

var file_accounts_accounts_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x73, 0x22, 0x22, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x12, 0x0a, 0x10, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x50, 0x0a, 0x0c, 0x44,
	0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x11, 0x0a,
	0x0f, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x35, 0x0a, 0x07, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x24, 0x0a, 0x0c, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, 0xcb, 0x01,
	0x0a, 0x0a, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x61, 0x6e, 0x74, 0x12, 0x40, 0x0a, 0x08,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x73, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x1a, 0x1a, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2e, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x42,
	0x0a, 0x0b, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x2e,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x19, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x37, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x11,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65,
	0x79, 0x1a, 0x16, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x25, 0x5a, 0x23, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x64, 0x69, 0x6c, 0x75, 0x7a,
	0x2f, 0x72, 0x6f, 0x76, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_accounts_accounts_proto_rawDescOnce sync.Once
	file_accounts_accounts_proto_rawDescData = file_accounts_accounts_proto_rawDesc
)

func file_accounts_accounts_proto_rawDescGZIP() []byte {
	file_accounts_accounts_proto_rawDescOnce.Do(func() {
		file_accounts_accounts_proto_rawDescData = protoimpl.X.CompressGZIP(file_accounts_accounts_proto_rawDescData)
	})
	return file_accounts_accounts_proto_rawDescData
}

var file_accounts_accounts_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_accounts_accounts_proto_goTypes = []interface{}{
	(*RegisterInfo)(nil),     // 0: accounts.RegisterInfo
	(*RegisterResponse)(nil), // 1: accounts.RegisterResponse
	(*DataKeyValue)(nil),     // 2: accounts.DataKeyValue
	(*DataKeyResponse)(nil),  // 3: accounts.DataKeyResponse
	(*DataKey)(nil),          // 4: accounts.DataKey
	(*DataResponse)(nil),     // 5: accounts.DataResponse
}
var file_accounts_accounts_proto_depIdxs = []int32{
	0, // 0: accounts.Accountant.Register:input_type -> accounts.RegisterInfo
	2, // 1: accounts.Accountant.AssignValue:input_type -> accounts.DataKeyValue
	4, // 2: accounts.Accountant.GetValue:input_type -> accounts.DataKey
	1, // 3: accounts.Accountant.Register:output_type -> accounts.RegisterResponse
	3, // 4: accounts.Accountant.AssignValue:output_type -> accounts.DataKeyResponse
	5, // 5: accounts.Accountant.GetValue:output_type -> accounts.DataResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_accounts_accounts_proto_init() }
func file_accounts_accounts_proto_init() {
	if File_accounts_accounts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_accounts_accounts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterInfo); i {
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
		file_accounts_accounts_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterResponse); i {
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
		file_accounts_accounts_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataKeyValue); i {
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
		file_accounts_accounts_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataKeyResponse); i {
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
		file_accounts_accounts_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataKey); i {
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
		file_accounts_accounts_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataResponse); i {
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
			RawDescriptor: file_accounts_accounts_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_accounts_accounts_proto_goTypes,
		DependencyIndexes: file_accounts_accounts_proto_depIdxs,
		MessageInfos:      file_accounts_accounts_proto_msgTypes,
	}.Build()
	File_accounts_accounts_proto = out.File
	file_accounts_accounts_proto_rawDesc = nil
	file_accounts_accounts_proto_goTypes = nil
	file_accounts_accounts_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AccountantClient is the client API for Accountant service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AccountantClient interface {
	// Register should create a new account in the database
	// It will return an error if the account already exists
	Register(ctx context.Context, in *RegisterInfo, opts ...grpc.CallOption) (*RegisterResponse, error)
	// AssignValue assigns a key-value pair to an account, or overwrites an existing key
	AssignValue(ctx context.Context, in *DataKeyValue, opts ...grpc.CallOption) (*DataKeyResponse, error)
	// GetValue will get the value for a key for an account
	GetValue(ctx context.Context, in *DataKey, opts ...grpc.CallOption) (*DataResponse, error)
}

type accountantClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountantClient(cc grpc.ClientConnInterface) AccountantClient {
	return &accountantClient{cc}
}

func (c *accountantClient) Register(ctx context.Context, in *RegisterInfo, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/accounts.Accountant/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountantClient) AssignValue(ctx context.Context, in *DataKeyValue, opts ...grpc.CallOption) (*DataKeyResponse, error) {
	out := new(DataKeyResponse)
	err := c.cc.Invoke(ctx, "/accounts.Accountant/AssignValue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountantClient) GetValue(ctx context.Context, in *DataKey, opts ...grpc.CallOption) (*DataResponse, error) {
	out := new(DataResponse)
	err := c.cc.Invoke(ctx, "/accounts.Accountant/GetValue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountantServer is the server API for Accountant service.
type AccountantServer interface {
	// Register should create a new account in the database
	// It will return an error if the account already exists
	Register(context.Context, *RegisterInfo) (*RegisterResponse, error)
	// AssignValue assigns a key-value pair to an account, or overwrites an existing key
	AssignValue(context.Context, *DataKeyValue) (*DataKeyResponse, error)
	// GetValue will get the value for a key for an account
	GetValue(context.Context, *DataKey) (*DataResponse, error)
}

// UnimplementedAccountantServer can be embedded to have forward compatible implementations.
type UnimplementedAccountantServer struct {
}

func (*UnimplementedAccountantServer) Register(context.Context, *RegisterInfo) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedAccountantServer) AssignValue(context.Context, *DataKeyValue) (*DataKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignValue not implemented")
}
func (*UnimplementedAccountantServer) GetValue(context.Context, *DataKey) (*DataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetValue not implemented")
}

func RegisterAccountantServer(s *grpc.Server, srv AccountantServer) {
	s.RegisterService(&_Accountant_serviceDesc, srv)
}

func _Accountant_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountantServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accounts.Accountant/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountantServer).Register(ctx, req.(*RegisterInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Accountant_AssignValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataKeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountantServer).AssignValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accounts.Accountant/AssignValue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountantServer).AssignValue(ctx, req.(*DataKeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Accountant_GetValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountantServer).GetValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accounts.Accountant/GetValue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountantServer).GetValue(ctx, req.(*DataKey))
	}
	return interceptor(ctx, in, info, handler)
}

var _Accountant_serviceDesc = grpc.ServiceDesc{
	ServiceName: "accounts.Accountant",
	HandlerType: (*AccountantServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Accountant_Register_Handler,
		},
		{
			MethodName: "AssignValue",
			Handler:    _Accountant_AssignValue_Handler,
		},
		{
			MethodName: "GetValue",
			Handler:    _Accountant_GetValue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "accounts/accounts.proto",
}
