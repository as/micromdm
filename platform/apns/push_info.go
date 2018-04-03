package apns

import (
	"github.com/pkg/errors"

	"github.com/as/micromdm/platform/apns/internal/pushproto"
	"github.com/gogo/protobuf/proto"
)

type PushInfo struct {
	UDID      string
	PushMagic string
	Token     string
	MDMTopic  string
}

func MarshalPushInfo(p *PushInfo) ([]byte, error) {
	protopush := pushproto.PushInfo{
		Udid:      p.UDID,
		PushMagic: p.PushMagic,
		Token:     p.Token,
		MdmTopic:  p.MDMTopic,
	}
	return proto.Marshal(&protopush)
}

func UnmarshalPushInfo(data []byte, p *PushInfo) error {
	var pb pushproto.PushInfo
	if err := proto.Unmarshal(data, &pb); err != nil {
		return errors.Wrap(err, "unmarshal proto to PushInfo")
	}
	p.UDID = pb.GetUdid()
	p.Token = pb.GetToken()
	p.PushMagic = pb.GetPushMagic()
	p.MDMTopic = pb.GetMdmTopic()
	return nil
}
