// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: proto/rpc.proto

package proto

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

type Rpc struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Handlers []*Handler `protobuf:"bytes,1,rep,name=handlers,proto3" json:"handlers,omitempty"`
}

func (x *Rpc) Reset() {
	*x = Rpc{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rpc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rpc) ProtoMessage() {}

func (x *Rpc) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rpc.ProtoReflect.Descriptor instead.
func (*Rpc) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *Rpc) GetHandlers() []*Handler {
	if x != nil {
		return x.Handlers
	}
	return nil
}

type Handler struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Method string `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Route  string `protobuf:"bytes,3,opt,name=route,proto3" json:"route,omitempty"`
}

func (x *Handler) Reset() {
	*x = Handler{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Handler) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Handler) ProtoMessage() {}

func (x *Handler) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Handler.ProtoReflect.Descriptor instead.
func (*Handler) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *Handler) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Handler) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Handler) GetRoute() string {
	if x != nil {
		return x.Route
	}
	return ""
}

type DummyMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DummyField string `protobuf:"bytes,1,opt,name=dummy_field,json=dummyField,proto3" json:"dummy_field,omitempty"`
}

func (x *DummyMessage) Reset() {
	*x = DummyMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DummyMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DummyMessage) ProtoMessage() {}

func (x *DummyMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DummyMessage.ProtoReflect.Descriptor instead.
func (*DummyMessage) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *DummyMessage) GetDummyField() string {
	if x != nil {
		return x.DummyField
	}
	return ""
}

var File_proto_rpc_proto protoreflect.FileDescriptor

var file_proto_rpc_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x06, 0x62, 0x6f, 0x69, 0x6c, 0x65, 0x72, 0x22, 0x32, 0x0a, 0x03, 0x52, 0x70, 0x63,
	0x12, 0x2b, 0x0a, 0x08, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x62, 0x6f, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x52, 0x08, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x73, 0x22, 0x4b, 0x0a,
	0x07, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x22, 0x2f, 0x0a, 0x0c, 0x44, 0x75,
	0x6d, 0x6d, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x75,
	0x6d, 0x6d, 0x79, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x64, 0x75, 0x6d, 0x6d, 0x79, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x42, 0x0b, 0x5a, 0x09, 0x70,
	0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_rpc_proto_rawDescOnce sync.Once
	file_proto_rpc_proto_rawDescData = file_proto_rpc_proto_rawDesc
)

func file_proto_rpc_proto_rawDescGZIP() []byte {
	file_proto_rpc_proto_rawDescOnce.Do(func() {
		file_proto_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_rpc_proto_rawDescData)
	})
	return file_proto_rpc_proto_rawDescData
}

var file_proto_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_rpc_proto_goTypes = []interface{}{
	(*Rpc)(nil),          // 0: boiler.Rpc
	(*Handler)(nil),      // 1: boiler.Handler
	(*DummyMessage)(nil), // 2: boiler.DummyMessage
}
var file_proto_rpc_proto_depIdxs = []int32{
	1, // 0: boiler.Rpc.handlers:type_name -> boiler.Handler
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_rpc_proto_init() }
func file_proto_rpc_proto_init() {
	if File_proto_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rpc); i {
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
		file_proto_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Handler); i {
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
		file_proto_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DummyMessage); i {
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
			RawDescriptor: file_proto_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_rpc_proto_goTypes,
		DependencyIndexes: file_proto_rpc_proto_depIdxs,
		MessageInfos:      file_proto_rpc_proto_msgTypes,
	}.Build()
	File_proto_rpc_proto = out.File
	file_proto_rpc_proto_rawDesc = nil
	file_proto_rpc_proto_goTypes = nil
	file_proto_rpc_proto_depIdxs = nil
}
