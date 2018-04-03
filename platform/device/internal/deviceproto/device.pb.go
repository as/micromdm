// Code generated by protoc-gen-go.
// source: device.proto
// DO NOT EDIT!

/*
Package deviceproto is a generated protocol buffer package.

It is generated from these files:
	device.proto

It has these top-level messages:
	Device
*/
package deviceproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Device struct {
	Uuid                   string `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
	Udid                   string `protobuf:"bytes,2,opt,name=udid" json:"udid,omitempty"`
	SerialNumber           string `protobuf:"bytes,3,opt,name=serial_number,json=serialNumber" json:"serial_number,omitempty"`
	OsVersion              string `protobuf:"bytes,4,opt,name=os_version,json=osVersion" json:"os_version,omitempty"`
	BuildVersion           string `protobuf:"bytes,5,opt,name=build_version,json=buildVersion" json:"build_version,omitempty"`
	ProductName            string `protobuf:"bytes,6,opt,name=product_name,json=productName" json:"product_name,omitempty"`
	Imei                   string `protobuf:"bytes,7,opt,name=imei" json:"imei,omitempty"`
	Meid                   string `protobuf:"bytes,8,opt,name=meid" json:"meid,omitempty"`
	Token                  string `protobuf:"bytes,9,opt,name=token" json:"token,omitempty"`
	PushMagic              string `protobuf:"bytes,10,opt,name=push_magic,json=pushMagic" json:"push_magic,omitempty"`
	MdmTopic               string `protobuf:"bytes,11,opt,name=mdm_topic,json=mdmTopic" json:"mdm_topic,omitempty"`
	UnlockToken            string `protobuf:"bytes,12,opt,name=unlock_token,json=unlockToken" json:"unlock_token,omitempty"`
	Enrolled               bool   `protobuf:"varint,13,opt,name=enrolled" json:"enrolled,omitempty"`
	AwaitingConfiguration  bool   `protobuf:"varint,14,opt,name=awaiting_configuration,json=awaitingConfiguration" json:"awaiting_configuration,omitempty"`
	DeviceName             string `protobuf:"bytes,15,opt,name=device_name,json=deviceName" json:"device_name,omitempty"`
	Model                  string `protobuf:"bytes,16,opt,name=model" json:"model,omitempty"`
	ModelName              string `protobuf:"bytes,17,opt,name=model_name,json=modelName" json:"model_name,omitempty"`
	Description            string `protobuf:"bytes,18,opt,name=description" json:"description,omitempty"`
	Color                  string `protobuf:"bytes,19,opt,name=color" json:"color,omitempty"`
	AssetTag               string `protobuf:"bytes,20,opt,name=asset_tag,json=assetTag" json:"asset_tag,omitempty"`
	DepDevice              bool   `protobuf:"varint,21,opt,name=dep_device,json=depDevice" json:"dep_device,omitempty"`
	DepProfileStatus       string `protobuf:"bytes,22,opt,name=dep_profile_status,json=depProfileStatus" json:"dep_profile_status,omitempty"`
	DepProfileUuid         string `protobuf:"bytes,23,opt,name=dep_profile_uuid,json=depProfileUuid" json:"dep_profile_uuid,omitempty"`
	DepProfileAssignTime   int64  `protobuf:"varint,24,opt,name=dep_profile_assign_time,json=depProfileAssignTime" json:"dep_profile_assign_time,omitempty"`
	DepProfilePushTime     int64  `protobuf:"varint,25,opt,name=dep_profile_push_time,json=depProfilePushTime" json:"dep_profile_push_time,omitempty"`
	DepProfileAssignedDate int64  `protobuf:"varint,26,opt,name=dep_profile_assigned_date,json=depProfileAssignedDate" json:"dep_profile_assigned_date,omitempty"`
	DepProfileAssignedBy   string `protobuf:"bytes,27,opt,name=dep_profile_assigned_by,json=depProfileAssignedBy" json:"dep_profile_assigned_by,omitempty"`
	LastCheckIn            int64  `protobuf:"varint,28,opt,name=last_check_in,json=lastCheckIn" json:"last_check_in,omitempty"`
	LastQueryResponse      []byte `protobuf:"bytes,29,opt,name=last_query_response,json=lastQueryResponse,proto3" json:"last_query_response,omitempty"`
}

func (m *Device) Reset()                    { *m = Device{} }
func (m *Device) String() string            { return proto.CompactTextString(m) }
func (*Device) ProtoMessage()               {}
func (*Device) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Device) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Device) GetUdid() string {
	if m != nil {
		return m.Udid
	}
	return ""
}

func (m *Device) GetSerialNumber() string {
	if m != nil {
		return m.SerialNumber
	}
	return ""
}

func (m *Device) GetOsVersion() string {
	if m != nil {
		return m.OsVersion
	}
	return ""
}

func (m *Device) GetBuildVersion() string {
	if m != nil {
		return m.BuildVersion
	}
	return ""
}

func (m *Device) GetProductName() string {
	if m != nil {
		return m.ProductName
	}
	return ""
}

func (m *Device) GetImei() string {
	if m != nil {
		return m.Imei
	}
	return ""
}

func (m *Device) GetMeid() string {
	if m != nil {
		return m.Meid
	}
	return ""
}

func (m *Device) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *Device) GetPushMagic() string {
	if m != nil {
		return m.PushMagic
	}
	return ""
}

func (m *Device) GetMdmTopic() string {
	if m != nil {
		return m.MdmTopic
	}
	return ""
}

func (m *Device) GetUnlockToken() string {
	if m != nil {
		return m.UnlockToken
	}
	return ""
}

func (m *Device) GetEnrolled() bool {
	if m != nil {
		return m.Enrolled
	}
	return false
}

func (m *Device) GetAwaitingConfiguration() bool {
	if m != nil {
		return m.AwaitingConfiguration
	}
	return false
}

func (m *Device) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *Device) GetModel() string {
	if m != nil {
		return m.Model
	}
	return ""
}

func (m *Device) GetModelName() string {
	if m != nil {
		return m.ModelName
	}
	return ""
}

func (m *Device) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Device) GetColor() string {
	if m != nil {
		return m.Color
	}
	return ""
}

func (m *Device) GetAssetTag() string {
	if m != nil {
		return m.AssetTag
	}
	return ""
}

func (m *Device) GetDepDevice() bool {
	if m != nil {
		return m.DepDevice
	}
	return false
}

func (m *Device) GetDepProfileStatus() string {
	if m != nil {
		return m.DepProfileStatus
	}
	return ""
}

func (m *Device) GetDepProfileUuid() string {
	if m != nil {
		return m.DepProfileUuid
	}
	return ""
}

func (m *Device) GetDepProfileAssignTime() int64 {
	if m != nil {
		return m.DepProfileAssignTime
	}
	return 0
}

func (m *Device) GetDepProfilePushTime() int64 {
	if m != nil {
		return m.DepProfilePushTime
	}
	return 0
}

func (m *Device) GetDepProfileAssignedDate() int64 {
	if m != nil {
		return m.DepProfileAssignedDate
	}
	return 0
}

func (m *Device) GetDepProfileAssignedBy() string {
	if m != nil {
		return m.DepProfileAssignedBy
	}
	return ""
}

func (m *Device) GetLastCheckIn() int64 {
	if m != nil {
		return m.LastCheckIn
	}
	return 0
}

func (m *Device) GetLastQueryResponse() []byte {
	if m != nil {
		return m.LastQueryResponse
	}
	return nil
}

func init() {
	proto.RegisterType((*Device)(nil), "deviceproto.Device")
}

func init() { proto.RegisterFile("device.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 565 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x94, 0x4f, 0x6f, 0x13, 0x31,
	0x10, 0xc5, 0x15, 0xda, 0xa6, 0x89, 0x93, 0x94, 0xd4, 0x4d, 0x52, 0xb7, 0xa5, 0x22, 0x94, 0x4b,
	0x0e, 0xa8, 0x12, 0x42, 0x3d, 0x70, 0x84, 0xf6, 0xc2, 0x81, 0xaa, 0x84, 0xc0, 0xd5, 0x72, 0xd6,
	0xd3, 0x8d, 0x95, 0xb5, 0xbd, 0xf8, 0x4f, 0x51, 0xbe, 0x3c, 0x42, 0x1e, 0xa7, 0x4d, 0x44, 0xb9,
	0x79, 0x7e, 0xf3, 0xde, 0x78, 0x66, 0xd6, 0x5a, 0xd2, 0x95, 0xf0, 0xa0, 0x0a, 0xb8, 0xac, 0x9d,
	0x0d, 0x96, 0x76, 0x72, 0x84, 0xc1, 0xc5, 0x9f, 0x7d, 0xd2, 0xbc, 0xc1, 0x98, 0x52, 0xb2, 0x1b,
	0xa3, 0x92, 0xac, 0x31, 0x6e, 0x4c, 0xda, 0x53, 0x3c, 0x23, 0x93, 0x4a, 0xb2, 0x17, 0x6b, 0x26,
	0x95, 0xa4, 0x6f, 0x49, 0xcf, 0x83, 0x53, 0xa2, 0xe2, 0x26, 0xea, 0x39, 0x38, 0xb6, 0x83, 0xc9,
	0x6e, 0x86, 0xb7, 0xc8, 0xe8, 0x39, 0x21, 0xd6, 0xf3, 0x07, 0x70, 0x5e, 0x59, 0xc3, 0x76, 0x51,
	0xd1, 0xb6, 0xfe, 0x67, 0x06, 0xa9, 0xc6, 0x3c, 0xaa, 0x4a, 0x3e, 0x29, 0xf6, 0x72, 0x0d, 0x84,
	0x8f, 0xa2, 0x37, 0xa4, 0x5b, 0x3b, 0x2b, 0x63, 0x11, 0xb8, 0x11, 0x1a, 0x58, 0x13, 0x35, 0x9d,
	0x35, 0xbb, 0x15, 0x1a, 0x7b, 0x56, 0x1a, 0x14, 0xdb, 0xcf, 0xfd, 0xa5, 0x73, 0x62, 0x1a, 0x94,
	0x64, 0xad, 0xcc, 0xd2, 0x99, 0x0e, 0xc8, 0x5e, 0xb0, 0x4b, 0x30, 0xac, 0x8d, 0x30, 0x07, 0xa9,
	0xc9, 0x3a, 0xfa, 0x05, 0xd7, 0xa2, 0x54, 0x05, 0x23, 0xb9, 0xc9, 0x44, 0xbe, 0x26, 0x40, 0xcf,
	0x48, 0x5b, 0x4b, 0xcd, 0x83, 0xad, 0x55, 0xc1, 0x3a, 0x98, 0x6d, 0x69, 0xa9, 0x67, 0x29, 0x4e,
	0xcd, 0x45, 0x53, 0xd9, 0x62, 0xc9, 0x73, 0xe1, 0x6e, 0x6e, 0x2e, 0xb3, 0x19, 0x96, 0x3f, 0x25,
	0x2d, 0x30, 0xce, 0x56, 0x15, 0x48, 0xd6, 0x1b, 0x37, 0x26, 0xad, 0xe9, 0x53, 0x4c, 0xaf, 0xc8,
	0x48, 0xfc, 0x16, 0x2a, 0x28, 0x53, 0xf2, 0xc2, 0x9a, 0x7b, 0x55, 0x46, 0x27, 0x42, 0xda, 0xc4,
	0x01, 0x2a, 0x87, 0x8f, 0xd9, 0xeb, 0xed, 0x24, 0x7d, 0x4d, 0xd6, 0x5f, 0x2f, 0x6f, 0xe4, 0x25,
	0x5e, 0x4a, 0x32, 0xc2, 0x85, 0x0c, 0xc8, 0x9e, 0xb6, 0x12, 0x2a, 0xd6, 0xcf, 0x83, 0x62, 0x90,
	0x06, 0xc5, 0x43, 0x76, 0x1d, 0xe6, 0x41, 0x91, 0xa0, 0x69, 0x9c, 0xaa, 0xfa, 0xc2, 0xa9, 0x1a,
	0x3b, 0xa0, 0x79, 0x94, 0x2d, 0x94, 0xca, 0x16, 0xb6, 0xb2, 0x8e, 0x1d, 0xe5, 0xb2, 0x18, 0xa4,
	0x05, 0x09, 0xef, 0x21, 0xf0, 0x20, 0x4a, 0x36, 0xc8, 0x0b, 0x42, 0x30, 0x13, 0x65, 0xba, 0x53,
	0x42, 0xcd, 0x73, 0x6f, 0x6c, 0x88, 0x53, 0xb5, 0x25, 0xd4, 0xeb, 0xd7, 0xf6, 0x8e, 0xd0, 0x94,
	0xae, 0x9d, 0xbd, 0x57, 0x15, 0x70, 0x1f, 0x44, 0x88, 0x9e, 0x8d, 0xb0, 0x48, 0x5f, 0x42, 0x7d,
	0x97, 0x13, 0xdf, 0x91, 0xd3, 0x09, 0xe9, 0x6f, 0xab, 0xf1, 0x9d, 0x1e, 0xa3, 0xf6, 0x60, 0xa3,
	0xfd, 0x91, 0x5e, 0xec, 0x15, 0x39, 0xde, 0x56, 0x0a, 0xef, 0x55, 0x69, 0x78, 0x50, 0x1a, 0x18,
	0x1b, 0x37, 0x26, 0x3b, 0xd3, 0xc1, 0xc6, 0xf0, 0x09, 0x93, 0x33, 0xa5, 0x81, 0xbe, 0x27, 0xc3,
	0x6d, 0x1b, 0x3e, 0x0b, 0x34, 0x9d, 0xa0, 0x89, 0x6e, 0x4c, 0x77, 0xd1, 0x2f, 0xd0, 0xf2, 0x91,
	0x9c, 0x3c, 0xbf, 0x09, 0x24, 0x97, 0x22, 0x00, 0x3b, 0x45, 0xdb, 0xe8, 0xdf, 0xbb, 0x40, 0xde,
	0x88, 0x00, 0xff, 0x6f, 0x12, 0x24, 0x9f, 0xaf, 0xd8, 0x19, 0x4e, 0x35, 0x78, 0x6e, 0xfc, 0xbc,
	0xa2, 0x17, 0xa4, 0x57, 0x09, 0x1f, 0x78, 0xb1, 0x80, 0x62, 0xc9, 0x95, 0x61, 0xaf, 0xf0, 0x96,
	0x4e, 0x82, 0xd7, 0x89, 0x7d, 0x31, 0xf4, 0x92, 0x1c, 0xa1, 0xe6, 0x57, 0x04, 0xb7, 0xe2, 0x0e,
	0x7c, 0x6d, 0x8d, 0x07, 0x76, 0x3e, 0x6e, 0x4c, 0xba, 0xd3, 0xc3, 0x94, 0xfa, 0x96, 0x32, 0xd3,
	0x75, 0x62, 0xde, 0xc4, 0xff, 0xc0, 0x87, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xd4, 0x53, 0xe5,
	0xaa, 0x24, 0x04, 0x00, 0x00,
}
