package connect

import (
	"time"

	"github.com/gogo/protobuf/proto"
	uuid "github.com/satori/go.uuid"

	"github.com/as/micromdm/mdm"
	"github.com/as/micromdm/mdm/connect/internal/connectproto"
)

type Event struct {
	ID       string
	Time     time.Time
	Response mdm.Response
	Raw      []byte
}

func NewEvent(resp MDMConnectRequest) *Event {
	event := Event{
		ID:       uuid.NewV4().String(),
		Time:     time.Now().UTC(),
		Response: resp.MDMResponse,
		Raw:      resp.Raw,
	}
	return &event
}

func MarshalEvent(e *Event) ([]byte, error) {
	response := &connectproto.Response{
		CommandUuid: e.Response.CommandUUID,
		Udid:        e.Response.UDID,
		Status:      e.Response.Status,
		RequestType: e.Response.RequestType,
	}
	if e.Response.UserID != nil {
		response.UserId = *e.Response.UserID
	}

	return proto.Marshal(&connectproto.Event{
		Id:       e.ID,
		Time:     e.Time.UnixNano(),
		Response: response,
		Raw:      e.Raw,
	})
}

func UnmarshalEvent(data []byte, e *Event) error {
	var pb connectproto.Event
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}
	e.ID = pb.Id
	e.Time = time.Unix(0, pb.Time).UTC()
	if pb.Response == nil {
		return nil
	}
	r := pb.GetResponse()
	e.Response = mdm.Response{
		UDID:        r.GetUdid(),
		UserID:      strPtr(r.GetUserId()),
		Status:      r.GetStatus(),
		RequestType: r.GetRequestType(),
		CommandUUID: r.GetCommandUuid(),
	}
	e.Raw = pb.Raw
	return nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
