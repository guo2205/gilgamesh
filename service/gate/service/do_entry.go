// do_entry
package service

import (
	"errors"
	"gilgamesh/protos"

	"github.com/golang/protobuf/proto"
)

var (
	ErrUnknownProtoType error = errors.New("unknown proto type")
)

func (c *GateService) doEntry(_type uint32, session uint64, data []byte) ([]byte, error) {
	switch _type {
	case 0:
		return []byte{}, c.doEntryNormal(session, data)
	case 1:
		return []byte{}, c.doEntryOffline(session)
	default:
		return []byte{}, ErrUnknownMailType
	}
}

func (c *GateService) doEntryNormal(session uint64, data []byte) error {
	gl := protos.Gilgamesh{}
	err := proto.Unmarshal(data, &gl)
	if err != nil {
		return err
	}

	switch gl.Type {
	case proto.MessageName((*protos.Public_Cts_Login)(nil)):
	default:
		return ErrUnknownProtoType
	}
}

func (c *GateService) doEntryOffline(session uint64) error {
	return nil
}
