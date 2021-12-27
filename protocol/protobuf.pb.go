// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf.proto

package protocol

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type LoginReq struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginReq) Reset()         { *m = LoginReq{} }
func (m *LoginReq) String() string { return proto.CompactTextString(m) }
func (*LoginReq) ProtoMessage()    {}
func (*LoginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{0}
}

func (m *LoginReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginReq.Unmarshal(m, b)
}
func (m *LoginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginReq.Marshal(b, m, deterministic)
}
func (m *LoginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginReq.Merge(m, src)
}
func (m *LoginReq) XXX_Size() int {
	return xxx_messageInfo_LoginReq.Size(m)
}
func (m *LoginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginReq.DiscardUnknown(m)
}

var xxx_messageInfo_LoginReq proto.InternalMessageInfo

func (m *LoginReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type LoginResp struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginResp) Reset()         { *m = LoginResp{} }
func (m *LoginResp) String() string { return proto.CompactTextString(m) }
func (*LoginResp) ProtoMessage()    {}
func (*LoginResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{1}
}

func (m *LoginResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginResp.Unmarshal(m, b)
}
func (m *LoginResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginResp.Marshal(b, m, deterministic)
}
func (m *LoginResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginResp.Merge(m, src)
}
func (m *LoginResp) XXX_Size() int {
	return xxx_messageInfo_LoginResp.Size(m)
}
func (m *LoginResp) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginResp.DiscardUnknown(m)
}

var xxx_messageInfo_LoginResp proto.InternalMessageInfo

func (m *LoginResp) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type StartReq struct {
	ItemID               int32    `protobuf:"varint,1,opt,name=itemID,proto3" json:"itemID,omitempty"`
	Config               []byte   `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	Args                 string   `protobuf:"bytes,3,opt,name=args,proto3" json:"args,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartReq) Reset()         { *m = StartReq{} }
func (m *StartReq) String() string { return proto.CompactTextString(m) }
func (*StartReq) ProtoMessage()    {}
func (*StartReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{2}
}

func (m *StartReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartReq.Unmarshal(m, b)
}
func (m *StartReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartReq.Marshal(b, m, deterministic)
}
func (m *StartReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartReq.Merge(m, src)
}
func (m *StartReq) XXX_Size() int {
	return xxx_messageInfo_StartReq.Size(m)
}
func (m *StartReq) XXX_DiscardUnknown() {
	xxx_messageInfo_StartReq.DiscardUnknown(m)
}

var xxx_messageInfo_StartReq proto.InternalMessageInfo

func (m *StartReq) GetItemID() int32 {
	if m != nil {
		return m.ItemID
	}
	return 0
}

func (m *StartReq) GetConfig() []byte {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *StartReq) GetArgs() string {
	if m != nil {
		return m.Args
	}
	return ""
}

type StartResp struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartResp) Reset()         { *m = StartResp{} }
func (m *StartResp) String() string { return proto.CompactTextString(m) }
func (*StartResp) ProtoMessage()    {}
func (*StartResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{3}
}

func (m *StartResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartResp.Unmarshal(m, b)
}
func (m *StartResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartResp.Marshal(b, m, deterministic)
}
func (m *StartResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartResp.Merge(m, src)
}
func (m *StartResp) XXX_Size() int {
	return xxx_messageInfo_StartResp.Size(m)
}
func (m *StartResp) XXX_DiscardUnknown() {
	xxx_messageInfo_StartResp.DiscardUnknown(m)
}

var xxx_messageInfo_StartResp proto.InternalMessageInfo

func (m *StartResp) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type SignalReq struct {
	ItemID               int32    `protobuf:"varint,1,opt,name=itemID,proto3" json:"itemID,omitempty"`
	Signal               int32    `protobuf:"varint,2,opt,name=signal,proto3" json:"signal,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignalReq) Reset()         { *m = SignalReq{} }
func (m *SignalReq) String() string { return proto.CompactTextString(m) }
func (*SignalReq) ProtoMessage()    {}
func (*SignalReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{4}
}

func (m *SignalReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignalReq.Unmarshal(m, b)
}
func (m *SignalReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignalReq.Marshal(b, m, deterministic)
}
func (m *SignalReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignalReq.Merge(m, src)
}
func (m *SignalReq) XXX_Size() int {
	return xxx_messageInfo_SignalReq.Size(m)
}
func (m *SignalReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SignalReq.DiscardUnknown(m)
}

var xxx_messageInfo_SignalReq proto.InternalMessageInfo

func (m *SignalReq) GetItemID() int32 {
	if m != nil {
		return m.ItemID
	}
	return 0
}

func (m *SignalReq) GetSignal() int32 {
	if m != nil {
		return m.Signal
	}
	return 0
}

type SignalResp struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignalResp) Reset()         { *m = SignalResp{} }
func (m *SignalResp) String() string { return proto.CompactTextString(m) }
func (*SignalResp) ProtoMessage()    {}
func (*SignalResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{5}
}

func (m *SignalResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignalResp.Unmarshal(m, b)
}
func (m *SignalResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignalResp.Marshal(b, m, deterministic)
}
func (m *SignalResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignalResp.Merge(m, src)
}
func (m *SignalResp) XXX_Size() int {
	return xxx_messageInfo_SignalResp.Size(m)
}
func (m *SignalResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SignalResp.DiscardUnknown(m)
}

var xxx_messageInfo_SignalResp proto.InternalMessageInfo

func (m *SignalResp) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type ItemStatueReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ItemStatueReq) Reset()         { *m = ItemStatueReq{} }
func (m *ItemStatueReq) String() string { return proto.CompactTextString(m) }
func (*ItemStatueReq) ProtoMessage()    {}
func (*ItemStatueReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{6}
}

func (m *ItemStatueReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ItemStatueReq.Unmarshal(m, b)
}
func (m *ItemStatueReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ItemStatueReq.Marshal(b, m, deterministic)
}
func (m *ItemStatueReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ItemStatueReq.Merge(m, src)
}
func (m *ItemStatueReq) XXX_Size() int {
	return xxx_messageInfo_ItemStatueReq.Size(m)
}
func (m *ItemStatueReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ItemStatueReq.DiscardUnknown(m)
}

var xxx_messageInfo_ItemStatueReq proto.InternalMessageInfo

type ItemStatueResp struct {
	Items                map[int32]*ItemStatue `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ItemStatueResp) Reset()         { *m = ItemStatueResp{} }
func (m *ItemStatueResp) String() string { return proto.CompactTextString(m) }
func (*ItemStatueResp) ProtoMessage()    {}
func (*ItemStatueResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{7}
}

func (m *ItemStatueResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ItemStatueResp.Unmarshal(m, b)
}
func (m *ItemStatueResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ItemStatueResp.Marshal(b, m, deterministic)
}
func (m *ItemStatueResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ItemStatueResp.Merge(m, src)
}
func (m *ItemStatueResp) XXX_Size() int {
	return xxx_messageInfo_ItemStatueResp.Size(m)
}
func (m *ItemStatueResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ItemStatueResp.DiscardUnknown(m)
}

var xxx_messageInfo_ItemStatueResp proto.InternalMessageInfo

func (m *ItemStatueResp) GetItems() map[int32]*ItemStatue {
	if m != nil {
		return m.Items
	}
	return nil
}

type ItemStatue struct {
	ItemID               int32    `protobuf:"varint,1,opt,name=itemID,proto3" json:"itemID,omitempty"`
	Pid                  int32    `protobuf:"varint,2,opt,name=pid,proto3" json:"pid,omitempty"`
	Timestamp            int64    `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	IsAlive              bool     `protobuf:"varint,4,opt,name=isAlive,proto3" json:"isAlive,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ItemStatue) Reset()         { *m = ItemStatue{} }
func (m *ItemStatue) String() string { return proto.CompactTextString(m) }
func (*ItemStatue) ProtoMessage()    {}
func (*ItemStatue) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{8}
}

func (m *ItemStatue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ItemStatue.Unmarshal(m, b)
}
func (m *ItemStatue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ItemStatue.Marshal(b, m, deterministic)
}
func (m *ItemStatue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ItemStatue.Merge(m, src)
}
func (m *ItemStatue) XXX_Size() int {
	return xxx_messageInfo_ItemStatue.Size(m)
}
func (m *ItemStatue) XXX_DiscardUnknown() {
	xxx_messageInfo_ItemStatue.DiscardUnknown(m)
}

var xxx_messageInfo_ItemStatue proto.InternalMessageInfo

func (m *ItemStatue) GetItemID() int32 {
	if m != nil {
		return m.ItemID
	}
	return 0
}

func (m *ItemStatue) GetPid() int32 {
	if m != nil {
		return m.Pid
	}
	return 0
}

func (m *ItemStatue) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *ItemStatue) GetIsAlive() bool {
	if m != nil {
		return m.IsAlive
	}
	return false
}

type PanicLogReq struct {
	ItemID               int32    `protobuf:"varint,1,opt,name=itemID,proto3" json:"itemID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PanicLogReq) Reset()         { *m = PanicLogReq{} }
func (m *PanicLogReq) String() string { return proto.CompactTextString(m) }
func (*PanicLogReq) ProtoMessage()    {}
func (*PanicLogReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{9}
}

func (m *PanicLogReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PanicLogReq.Unmarshal(m, b)
}
func (m *PanicLogReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PanicLogReq.Marshal(b, m, deterministic)
}
func (m *PanicLogReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PanicLogReq.Merge(m, src)
}
func (m *PanicLogReq) XXX_Size() int {
	return xxx_messageInfo_PanicLogReq.Size(m)
}
func (m *PanicLogReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PanicLogReq.DiscardUnknown(m)
}

var xxx_messageInfo_PanicLogReq proto.InternalMessageInfo

func (m *PanicLogReq) GetItemID() int32 {
	if m != nil {
		return m.ItemID
	}
	return 0
}

type PanicLogResp struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PanicLogResp) Reset()         { *m = PanicLogResp{} }
func (m *PanicLogResp) String() string { return proto.CompactTextString(m) }
func (*PanicLogResp) ProtoMessage()    {}
func (*PanicLogResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{10}
}

func (m *PanicLogResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PanicLogResp.Unmarshal(m, b)
}
func (m *PanicLogResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PanicLogResp.Marshal(b, m, deterministic)
}
func (m *PanicLogResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PanicLogResp.Merge(m, src)
}
func (m *PanicLogResp) XXX_Size() int {
	return xxx_messageInfo_PanicLogResp.Size(m)
}
func (m *PanicLogResp) XXX_DiscardUnknown() {
	xxx_messageInfo_PanicLogResp.DiscardUnknown(m)
}

var xxx_messageInfo_PanicLogResp proto.InternalMessageInfo

func (m *PanicLogResp) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *PanicLogResp) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*LoginReq)(nil), "loginReq")
	proto.RegisterType((*LoginResp)(nil), "loginResp")
	proto.RegisterType((*StartReq)(nil), "startReq")
	proto.RegisterType((*StartResp)(nil), "startResp")
	proto.RegisterType((*SignalReq)(nil), "signalReq")
	proto.RegisterType((*SignalResp)(nil), "signalResp")
	proto.RegisterType((*ItemStatueReq)(nil), "itemStatueReq")
	proto.RegisterType((*ItemStatueResp)(nil), "itemStatueResp")
	proto.RegisterMapType((map[int32]*ItemStatue)(nil), "itemStatueResp.ItemsEntry")
	proto.RegisterType((*ItemStatue)(nil), "itemStatue")
	proto.RegisterType((*PanicLogReq)(nil), "panicLogReq")
	proto.RegisterType((*PanicLogResp)(nil), "panicLogResp")
}

func init() { proto.RegisterFile("protobuf.proto", fileDescriptor_c77a803fcbc0c059) }

var fileDescriptor_c77a803fcbc0c059 = []byte{
	// 345 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x25, 0x4d, 0x53, 0x93, 0x49, 0xad, 0x25, 0x07, 0x09, 0x45, 0x34, 0x2e, 0x08, 0x39, 0x15,
	0xa9, 0x1e, 0x44, 0x4f, 0x8a, 0x3d, 0x14, 0xc4, 0xc3, 0x7a, 0xf3, 0xb6, 0x4d, 0xb7, 0x61, 0x31,
	0x5f, 0xcd, 0x6e, 0x0b, 0xfd, 0x09, 0xfe, 0x6b, 0x99, 0xcd, 0xd6, 0x54, 0xb0, 0x7a, 0x7b, 0xf3,
	0xde, 0xdb, 0x37, 0xb3, 0xc3, 0xc0, 0xa0, 0xaa, 0x4b, 0x55, 0xce, 0xd7, 0xcb, 0xb1, 0x06, 0xe4,
	0x1c, 0xdc, 0xac, 0x4c, 0x45, 0x41, 0xf9, 0x2a, 0x08, 0xa0, 0x5b, 0xb0, 0x9c, 0x87, 0x56, 0x64,
	0xc5, 0x1e, 0xd5, 0x98, 0x5c, 0x80, 0x67, 0x74, 0x59, 0xa1, 0x21, 0x29, 0x17, 0xdf, 0x06, 0xc4,
	0xe4, 0x15, 0x5c, 0xa9, 0x58, 0xad, 0x30, 0xe0, 0x14, 0x7a, 0x42, 0xf1, 0x7c, 0xf6, 0xac, 0x1d,
	0x0e, 0x35, 0x15, 0xf2, 0x49, 0x59, 0x2c, 0x45, 0x1a, 0x76, 0x22, 0x2b, 0xee, 0x53, 0x53, 0x61,
	0x1e, 0xab, 0x53, 0x19, 0xda, 0x4d, 0x1e, 0x62, 0x6c, 0x68, 0xf2, 0x0e, 0x34, 0x7c, 0x00, 0x4f,
	0x8a, 0xb4, 0x60, 0xd9, 0x3f, 0x1d, 0x1b, 0x93, 0xee, 0xe8, 0x50, 0x53, 0x91, 0x08, 0x60, 0xf7,
	0xf8, 0x40, 0xfc, 0x09, 0x1c, 0x63, 0xc6, 0x9b, 0x62, 0x6a, 0xcd, 0x29, 0x5f, 0x91, 0x4f, 0x0b,
	0x06, 0xfb, 0x8c, 0xac, 0x82, 0x6b, 0x70, 0x90, 0x91, 0xa1, 0x15, 0xd9, 0xb1, 0x3f, 0x19, 0x8d,
	0x7f, 0xea, 0xe3, 0x19, 0x8a, 0xd3, 0x42, 0xd5, 0x5b, 0xda, 0x18, 0x47, 0x53, 0x80, 0x96, 0x0c,
	0x86, 0x60, 0x7f, 0xf0, 0xad, 0x19, 0x19, 0x61, 0x70, 0x09, 0xce, 0x86, 0x65, 0x6b, 0xae, 0xc7,
	0xf5, 0x27, 0xfe, 0x7e, 0x62, 0xa3, 0xdc, 0x77, 0xee, 0x2c, 0x52, 0x00, 0xb4, 0xc2, 0xc1, 0xcf,
	0x0f, 0xc1, 0xae, 0xc4, 0xc2, 0xfc, 0x1c, 0x61, 0x70, 0x06, 0x9e, 0x12, 0x39, 0x97, 0x8a, 0xe5,
	0x95, 0xde, 0xb6, 0x4d, 0x5b, 0x22, 0x08, 0xe1, 0x48, 0xc8, 0xc7, 0x4c, 0x6c, 0x78, 0xd8, 0x8d,
	0xac, 0xd8, 0xa5, 0xbb, 0x92, 0x5c, 0x81, 0x5f, 0xb1, 0x42, 0x24, 0x2f, 0x65, 0xfa, 0xc7, 0xb6,
	0xc9, 0x2d, 0xf4, 0x5b, 0xdb, 0xef, 0x7b, 0xc5, 0xa1, 0x72, 0xd9, 0x1c, 0x80, 0x47, 0x11, 0x3e,
	0xc1, 0xbb, 0xab, 0x6f, 0x30, 0x29, 0xb3, 0x79, 0x4f, 0xa3, 0x9b, 0xaf, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x84, 0xef, 0xac, 0x18, 0x9f, 0x02, 0x00, 0x00,
}
