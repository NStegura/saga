// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.6.1
// source: productsapi.proto

package api

import (
	reflect "reflect"
	sync "sync"

	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Products struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Products []*Product `protobuf:"bytes,1,rep,name=products,proto3" json:"products,omitempty"`
}

func (x *Products) Reset() {
	*x = Products{}
	if protoimpl.UnsafeEnabled {
		mi := &file_productsapi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Products) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Products) ProtoMessage() {}

func (x *Products) ProtoReflect() protoreflect.Message {
	mi := &file_productsapi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Products.ProtoReflect.Descriptor instead.
func (*Products) Descriptor() ([]byte, []int) {
	return file_productsapi_proto_rawDescGZIP(), []int{0}
}

func (x *Products) GetProducts() []*Product {
	if x != nil {
		return x.Products
	}
	return nil
}

type ProductId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId int64 `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
}

func (x *ProductId) Reset() {
	*x = ProductId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_productsapi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProductId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductId) ProtoMessage() {}

func (x *ProductId) ProtoReflect() protoreflect.Message {
	mi := &file_productsapi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductId.ProtoReflect.Descriptor instead.
func (*ProductId) Descriptor() ([]byte, []int) {
	return file_productsapi_proto_rawDescGZIP(), []int{1}
}

func (x *ProductId) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

type Product struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId   int64  `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Category    string `protobuf:"bytes,2,opt,name=category,proto3" json:"category,omitempty"`
	Name        string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Count       int64  `protobuf:"varint,5,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *Product) Reset() {
	*x = Product{}
	if protoimpl.UnsafeEnabled {
		mi := &file_productsapi_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Product) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Product) ProtoMessage() {}

func (x *Product) ProtoReflect() protoreflect.Message {
	mi := &file_productsapi_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Product.ProtoReflect.Descriptor instead.
func (*Product) Descriptor() ([]byte, []int) {
	return file_productsapi_proto_rawDescGZIP(), []int{2}
}

func (x *Product) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *Product) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Product) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Product) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Product) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type Pong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pong bool `protobuf:"varint,1,opt,name=pong,proto3" json:"pong,omitempty"`
}

func (x *Pong) Reset() {
	*x = Pong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_productsapi_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pong) ProtoMessage() {}

func (x *Pong) ProtoReflect() protoreflect.Message {
	mi := &file_productsapi_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pong.ProtoReflect.Descriptor instead.
func (*Pong) Descriptor() ([]byte, []int) {
	return file_productsapi_proto_rawDescGZIP(), []int{3}
}

func (x *Pong) GetPong() bool {
	if x != nil {
		return x.Pong
	}
	return false
}

var File_productsapi_proto protoreflect.FileDescriptor

var file_productsapi_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x61, 0x70, 0x69,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3c, 0x0a,
	0x08, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x12, 0x30, 0x0a, 0x08, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x22, 0x2a, 0x0a, 0x09, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x22, 0x90, 0x01, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x1a, 0x0a, 0x04, 0x50, 0x6f,
	0x6e, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x32, 0xc7, 0x01, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x73, 0x41, 0x70, 0x69, 0x12, 0x3e, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x15, 0x2e,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x73, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64,
	0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x50,
	0x69, 0x6e, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x11, 0x2e, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x22, 0x00,
	0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4e,
	0x53, 0x74, 0x65, 0x67, 0x75, 0x72, 0x61, 0x2f, 0x73, 0x61, 0x67, 0x61, 0x2f, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_productsapi_proto_rawDescOnce sync.Once
	file_productsapi_proto_rawDescData = file_productsapi_proto_rawDesc
)

func file_productsapi_proto_rawDescGZIP() []byte {
	file_productsapi_proto_rawDescOnce.Do(func() {
		file_productsapi_proto_rawDescData = protoimpl.X.CompressGZIP(file_productsapi_proto_rawDescData)
	})
	return file_productsapi_proto_rawDescData
}

var file_productsapi_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_productsapi_proto_goTypes = []any{
	(*Products)(nil),    // 0: productsapi.Products
	(*ProductId)(nil),   // 1: productsapi.ProductId
	(*Product)(nil),     // 2: productsapi.Product
	(*Pong)(nil),        // 3: productsapi.Pong
	(*empty.Empty)(nil), // 4: google.protobuf.Empty
}
var file_productsapi_proto_depIdxs = []int32{
	2, // 0: productsapi.Products.products:type_name -> productsapi.Product
	4, // 1: productsapi.ProductsApi.GetProducts:input_type -> google.protobuf.Empty
	1, // 2: productsapi.ProductsApi.GetProductInfo:input_type -> productsapi.ProductId
	4, // 3: productsapi.ProductsApi.GetPing:input_type -> google.protobuf.Empty
	0, // 4: productsapi.ProductsApi.GetProducts:output_type -> productsapi.Products
	2, // 5: productsapi.ProductsApi.GetProductInfo:output_type -> productsapi.Product
	3, // 6: productsapi.ProductsApi.GetPing:output_type -> productsapi.Pong
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_productsapi_proto_init() }
func file_productsapi_proto_init() {
	if File_productsapi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_productsapi_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Products); i {
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
		file_productsapi_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ProductId); i {
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
		file_productsapi_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Product); i {
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
		file_productsapi_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Pong); i {
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
			RawDescriptor: file_productsapi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_productsapi_proto_goTypes,
		DependencyIndexes: file_productsapi_proto_depIdxs,
		MessageInfos:      file_productsapi_proto_msgTypes,
	}.Build()
	File_productsapi_proto = out.File
	file_productsapi_proto_rawDesc = nil
	file_productsapi_proto_goTypes = nil
	file_productsapi_proto_depIdxs = nil
}
