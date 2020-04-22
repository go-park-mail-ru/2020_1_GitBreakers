// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: user.proto

package usergrpc

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
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

//использует модель юзера
type UserAvatarModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chunk    []byte     `protobuf:"bytes,1,opt,name=Chunk,proto3" json:"Chunk,omitempty"`
	User     *UserModel `protobuf:"bytes,2,opt,name=User,proto3" json:"User,omitempty"`
	FileName string     `protobuf:"bytes,3,opt,name=FileName,proto3" json:"FileName,omitempty"`
}

func (x *UserAvatarModel) Reset() {
	*x = UserAvatarModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserAvatarModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserAvatarModel) ProtoMessage() {}

func (x *UserAvatarModel) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserAvatarModel.ProtoReflect.Descriptor instead.
func (*UserAvatarModel) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *UserAvatarModel) GetChunk() []byte {
	if x != nil {
		return x.Chunk
	}
	return nil
}

func (x *UserAvatarModel) GetUser() *UserModel {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *UserAvatarModel) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

type UserUpdateModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID   *UserIDModel `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
	UserData *UserModel   `protobuf:"bytes,2,opt,name=UserData,proto3" json:"UserData,omitempty"`
}

func (x *UserUpdateModel) Reset() {
	*x = UserUpdateModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserUpdateModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserUpdateModel) ProtoMessage() {}

func (x *UserUpdateModel) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserUpdateModel.ProtoReflect.Descriptor instead.
func (*UserUpdateModel) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *UserUpdateModel) GetUserID() *UserIDModel {
	if x != nil {
		return x.UserID
	}
	return nil
}

func (x *UserUpdateModel) GetUserData() *UserModel {
	if x != nil {
		return x.UserData
	}
	return nil
}

type UserIDModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID int64 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
}

func (x *UserIDModel) Reset() {
	*x = UserIDModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserIDModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserIDModel) ProtoMessage() {}

func (x *UserIDModel) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserIDModel.ProtoReflect.Descriptor instead.
func (*UserIDModel) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *UserIDModel) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

type CheckPassResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsCorrect bool `protobuf:"varint,1,opt,name=is_correct,json=isCorrect,proto3" json:"is_correct,omitempty"`
}

func (x *CheckPassResp) Reset() {
	*x = CheckPassResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckPassResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckPassResp) ProtoMessage() {}

func (x *CheckPassResp) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckPassResp.ProtoReflect.Descriptor instead.
func (*CheckPassResp) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{3}
}

func (x *CheckPassResp) GetIsCorrect() bool {
	if x != nil {
		return x.IsCorrect
	}
	return false
}

type CheckPassModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login *LoginModel `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Pass  string      `protobuf:"bytes,2,opt,name=pass,proto3" json:"pass,omitempty"`
}

func (x *CheckPassModel) Reset() {
	*x = CheckPassModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckPassModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckPassModel) ProtoMessage() {}

func (x *CheckPassModel) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckPassModel.ProtoReflect.Descriptor instead.
func (*CheckPassModel) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{4}
}

func (x *CheckPassModel) GetLogin() *LoginModel {
	if x != nil {
		return x.Login
	}
	return nil
}

func (x *CheckPassModel) GetPass() string {
	if x != nil {
		return x.Pass
	}
	return ""
}

type LoginModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
}

func (x *LoginModel) Reset() {
	*x = LoginModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginModel) ProtoMessage() {}

func (x *LoginModel) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginModel.ProtoReflect.Descriptor instead.
func (*LoginModel) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{5}
}

func (x *LoginModel) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

type UserModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Name     string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Login    string `protobuf:"bytes,4,opt,name=login,proto3" json:"login,omitempty"`
	Image    string `protobuf:"bytes,5,opt,name=image,proto3" json:"image,omitempty"`
	Email    string `protobuf:"bytes,6,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *UserModel) Reset() {
	*x = UserModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserModel) ProtoMessage() {}

func (x *UserModel) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserModel.ProtoReflect.Descriptor instead.
func (*UserModel) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{6}
}

func (x *UserModel) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserModel) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *UserModel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserModel) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *UserModel) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

func (x *UserModel) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x6c, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x41, 0x76, 0x61, 0x74, 0x61,
	0x72, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x27, 0x0a, 0x04,
	0x55, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52,
	0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x22, 0x71, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d,
	0x6f, 0x64, 0x65, 0x6c, 0x12, 0x2d, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x44, 0x12, 0x2f, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52, 0x08, 0x55, 0x73, 0x65, 0x72,
	0x44, 0x61, 0x74, 0x61, 0x22, 0x25, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x4d, 0x6f,
	0x64, 0x65, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x22, 0x2e, 0x0a, 0x0d, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x50, 0x61, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x12, 0x1d, 0x0a, 0x0a,
	0x69, 0x73, 0x5f, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x69, 0x73, 0x43, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x22, 0x50, 0x0a, 0x0e, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x50, 0x61, 0x73, 0x73, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x2a, 0x0a,
	0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x4d, 0x6f, 0x64,
	0x65, 0x6c, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x73,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x73, 0x73, 0x22, 0x22, 0x0a,
	0x0a, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x6c,
	0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x22, 0x8d, 0x01, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x32, 0x83, 0x03, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x47, 0x72, 0x70, 0x63, 0x12, 0x37,
	0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x13, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x42, 0x79,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x14, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x13, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x6c,
	0x22, 0x00, 0x12, 0x41, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72,
	0x12, 0x19, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x09, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x61,
	0x73, 0x73, 0x12, 0x18, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x50, 0x61, 0x73, 0x73, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x17, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x61, 0x73,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x42, 0x79,
	0x49, 0x44, 0x12, 0x15, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x44, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x13, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x22, 0x00,
	0x12, 0x45, 0x0a, 0x0c, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72,
	0x12, 0x19, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x28, 0x01, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x3b, 0x75, 0x73, 0x65,
	0x72, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData = file_user_proto_rawDesc
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_proto_rawDescData)
	})
	return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_user_proto_goTypes = []interface{}{
	(*UserAvatarModel)(nil), // 0: usergrpc.UserAvatarModel
	(*UserUpdateModel)(nil), // 1: usergrpc.UserUpdateModel
	(*UserIDModel)(nil),     // 2: usergrpc.UserIDModel
	(*CheckPassResp)(nil),   // 3: usergrpc.CheckPassResp
	(*CheckPassModel)(nil),  // 4: usergrpc.CheckPassModel
	(*LoginModel)(nil),      // 5: usergrpc.LoginModel
	(*UserModel)(nil),       // 6: usergrpc.UserModel
	(*empty.Empty)(nil),     // 7: google.protobuf.Empty
}
var file_user_proto_depIdxs = []int32{
	6,  // 0: usergrpc.UserAvatarModel.User:type_name -> usergrpc.UserModel
	2,  // 1: usergrpc.UserUpdateModel.userID:type_name -> usergrpc.UserIDModel
	6,  // 2: usergrpc.UserUpdateModel.UserData:type_name -> usergrpc.UserModel
	5,  // 3: usergrpc.CheckPassModel.login:type_name -> usergrpc.LoginModel
	6,  // 4: usergrpc.UserGrpc.Create:input_type -> usergrpc.UserModel
	5,  // 5: usergrpc.UserGrpc.GetByLogin:input_type -> usergrpc.LoginModel
	1,  // 6: usergrpc.UserGrpc.UpdateUser:input_type -> usergrpc.UserUpdateModel
	4,  // 7: usergrpc.UserGrpc.CheckPass:input_type -> usergrpc.CheckPassModel
	2,  // 8: usergrpc.UserGrpc.GetByID:input_type -> usergrpc.UserIDModel
	1,  // 9: usergrpc.UserGrpc.UploadAvatar:input_type -> usergrpc.UserUpdateModel
	7,  // 10: usergrpc.UserGrpc.Create:output_type -> google.protobuf.Empty
	6,  // 11: usergrpc.UserGrpc.GetByLogin:output_type -> usergrpc.UserModel
	7,  // 12: usergrpc.UserGrpc.UpdateUser:output_type -> google.protobuf.Empty
	3,  // 13: usergrpc.UserGrpc.CheckPass:output_type -> usergrpc.CheckPassResp
	6,  // 14: usergrpc.UserGrpc.GetByID:output_type -> usergrpc.UserModel
	7,  // 15: usergrpc.UserGrpc.UploadAvatar:output_type -> google.protobuf.Empty
	10, // [10:16] is the sub-list for method output_type
	4,  // [4:10] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserAvatarModel); i {
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
		file_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserUpdateModel); i {
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
		file_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserIDModel); i {
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
		file_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckPassResp); i {
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
		file_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckPassModel); i {
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
		file_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginModel); i {
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
		file_user_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserModel); i {
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
			RawDescriptor: file_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_rawDesc = nil
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UserGrpcClient is the client API for UserGrpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserGrpcClient interface {
	Create(ctx context.Context, in *UserModel, opts ...grpc.CallOption) (*empty.Empty, error)
	GetByLogin(ctx context.Context, in *LoginModel, opts ...grpc.CallOption) (*UserModel, error)
	UpdateUser(ctx context.Context, in *UserUpdateModel, opts ...grpc.CallOption) (*empty.Empty, error)
	CheckPass(ctx context.Context, in *CheckPassModel, opts ...grpc.CallOption) (*CheckPassResp, error)
	GetByID(ctx context.Context, in *UserIDModel, opts ...grpc.CallOption) (*UserModel, error)
	UploadAvatar(ctx context.Context, opts ...grpc.CallOption) (UserGrpc_UploadAvatarClient, error)
}

type userGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewUserGrpcClient(cc grpc.ClientConnInterface) UserGrpcClient {
	return &userGrpcClient{cc}
}

func (c *userGrpcClient) Create(ctx context.Context, in *UserModel, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/usergrpc.UserGrpc/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGrpcClient) GetByLogin(ctx context.Context, in *LoginModel, opts ...grpc.CallOption) (*UserModel, error) {
	out := new(UserModel)
	err := c.cc.Invoke(ctx, "/usergrpc.UserGrpc/GetByLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGrpcClient) UpdateUser(ctx context.Context, in *UserUpdateModel, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/usergrpc.UserGrpc/UpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGrpcClient) CheckPass(ctx context.Context, in *CheckPassModel, opts ...grpc.CallOption) (*CheckPassResp, error) {
	out := new(CheckPassResp)
	err := c.cc.Invoke(ctx, "/usergrpc.UserGrpc/CheckPass", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGrpcClient) GetByID(ctx context.Context, in *UserIDModel, opts ...grpc.CallOption) (*UserModel, error) {
	out := new(UserModel)
	err := c.cc.Invoke(ctx, "/usergrpc.UserGrpc/GetByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userGrpcClient) UploadAvatar(ctx context.Context, opts ...grpc.CallOption) (UserGrpc_UploadAvatarClient, error) {
	stream, err := c.cc.NewStream(ctx, &_UserGrpc_serviceDesc.Streams[0], "/usergrpc.UserGrpc/UploadAvatar", opts...)
	if err != nil {
		return nil, err
	}
	x := &userGrpcUploadAvatarClient{stream}
	return x, nil
}

type UserGrpc_UploadAvatarClient interface {
	Send(*UserUpdateModel) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type userGrpcUploadAvatarClient struct {
	grpc.ClientStream
}

func (x *userGrpcUploadAvatarClient) Send(m *UserUpdateModel) error {
	return x.ClientStream.SendMsg(m)
}

func (x *userGrpcUploadAvatarClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UserGrpcServer is the server API for UserGrpc service.
type UserGrpcServer interface {
	Create(context.Context, *UserModel) (*empty.Empty, error)
	GetByLogin(context.Context, *LoginModel) (*UserModel, error)
	UpdateUser(context.Context, *UserUpdateModel) (*empty.Empty, error)
	CheckPass(context.Context, *CheckPassModel) (*CheckPassResp, error)
	GetByID(context.Context, *UserIDModel) (*UserModel, error)
	UploadAvatar(UserGrpc_UploadAvatarServer) error
}

// UnimplementedUserGrpcServer can be embedded to have forward compatible implementations.
type UnimplementedUserGrpcServer struct {
}

func (*UnimplementedUserGrpcServer) Create(context.Context, *UserModel) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedUserGrpcServer) GetByLogin(context.Context, *LoginModel) (*UserModel, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByLogin not implemented")
}
func (*UnimplementedUserGrpcServer) UpdateUser(context.Context, *UserUpdateModel) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (*UnimplementedUserGrpcServer) CheckPass(context.Context, *CheckPassModel) (*CheckPassResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPass not implemented")
}
func (*UnimplementedUserGrpcServer) GetByID(context.Context, *UserIDModel) (*UserModel, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (*UnimplementedUserGrpcServer) UploadAvatar(UserGrpc_UploadAvatarServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadAvatar not implemented")
}

func RegisterUserGrpcServer(s *grpc.Server, srv UserGrpcServer) {
	s.RegisterService(&_UserGrpc_serviceDesc, srv)
}

func _UserGrpc_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserModel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGrpcServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usergrpc.UserGrpc/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGrpcServer).Create(ctx, req.(*UserModel))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGrpc_GetByLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginModel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGrpcServer).GetByLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usergrpc.UserGrpc/GetByLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGrpcServer).GetByLogin(ctx, req.(*LoginModel))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGrpc_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserUpdateModel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGrpcServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usergrpc.UserGrpc/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGrpcServer).UpdateUser(ctx, req.(*UserUpdateModel))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGrpc_CheckPass_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckPassModel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGrpcServer).CheckPass(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usergrpc.UserGrpc/CheckPass",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGrpcServer).CheckPass(ctx, req.(*CheckPassModel))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGrpc_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIDModel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserGrpcServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usergrpc.UserGrpc/GetByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserGrpcServer).GetByID(ctx, req.(*UserIDModel))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserGrpc_UploadAvatar_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UserGrpcServer).UploadAvatar(&userGrpcUploadAvatarServer{stream})
}

type UserGrpc_UploadAvatarServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*UserUpdateModel, error)
	grpc.ServerStream
}

type userGrpcUploadAvatarServer struct {
	grpc.ServerStream
}

func (x *userGrpcUploadAvatarServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *userGrpcUploadAvatarServer) Recv() (*UserUpdateModel, error) {
	m := new(UserUpdateModel)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _UserGrpc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "usergrpc.UserGrpc",
	HandlerType: (*UserGrpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _UserGrpc_Create_Handler,
		},
		{
			MethodName: "GetByLogin",
			Handler:    _UserGrpc_GetByLogin_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserGrpc_UpdateUser_Handler,
		},
		{
			MethodName: "CheckPass",
			Handler:    _UserGrpc_CheckPass_Handler,
		},
		{
			MethodName: "GetByID",
			Handler:    _UserGrpc_GetByID_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadAvatar",
			Handler:       _UserGrpc_UploadAvatar_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "user.proto",
}
