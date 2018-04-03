package depsync

import (
	"time"

	"github.com/as/micromdm/dep"
	"github.com/gogo/protobuf/proto"
	uuid "github.com/satori/go.uuid"

	"github.com/as/micromdm/dep/depsync/internal/depsyncproto"
)

type Event struct {
	ID      string
	Time    time.Time
	Devices []dep.Device
}

func NewEvent(devices []dep.Device) *Event {
	event := Event{
		ID:      uuid.NewV4().String(),
		Time:    time.Now().UTC(),
		Devices: devices,
	}
	return &event
}

// MarshalEvent serializes an event to a protocol buffer wire format.
func MarshalEvent(e *Event) ([]byte, error) {
	var devices []*depsyncproto.Device
	for _, d := range e.Devices {
		devices = append(devices, &depsyncproto.Device{
			SerialNumber:       d.SerialNumber,
			Model:              d.Model,
			Description:        d.Description,
			Color:              d.Color,
			AssetTag:           d.AssetTag,
			ProfileUuid:        d.ProfileUUID,
			ProfileAssignTime:  d.ProfileAssignTime.UnixNano(),
			ProfilePushTime:    d.ProfilePushTime.UnixNano(),
			DeviceAssignedDate: d.DeviceAssignedDate.UnixNano(),
			DeviceAssignedBy:   d.DeviceAssignedBy,
			OpType:             d.OpType,
			OpDate:             d.OpDate.UnixNano(),
		})
	}
	return proto.Marshal(&depsyncproto.Event{
		Id:      e.ID,
		Time:    e.Time.UnixNano(),
		Devices: devices,
	})
}

// UnmarshalEvent parses a protocol buffer representation of data into
// the Event.
func UnmarshalEvent(data []byte, e *Event) error {
	var pb depsyncproto.Event
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}
	e.ID = pb.GetId()
	e.Time = time.Unix(0, pb.GetTime()).UTC()
	protodev := pb.GetDevices()
	var devices []dep.Device
	for _, d := range protodev {
		devices = append(devices, dep.Device{
			SerialNumber:       d.GetSerialNumber(),
			Model:              d.GetModel(),
			Description:        d.GetDescription(),
			Color:              d.GetColor(),
			AssetTag:           d.GetAssetTag(),
			ProfileUUID:        d.GetProfileUuid(),
			ProfileAssignTime:  time.Unix(0, d.GetProfileAssignTime()).UTC(),
			ProfilePushTime:    time.Unix(0, d.GetProfilePushTime()).UTC(),
			DeviceAssignedDate: time.Unix(0, d.GetDeviceAssignedDate()).UTC(),
			DeviceAssignedBy:   d.GetDeviceAssignedBy(),
			OpType:             d.GetOpType(),
			OpDate:             time.Unix(0, d.GetOpDate()).UTC(),
		})
	}
	e.Devices = devices
	return nil
}
