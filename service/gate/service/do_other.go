// do_other
package service

import (
	"gilgamesh/protos"
	"gilgamesh/utility/utils"

	"github.com/golang/protobuf/proto"
)

func (c *GateService) doOther(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	switch _type {
	case 0:
		return []byte{}, c.writer(session, data)
	case 1:
		obj, _type, err := utils.Unmarshal(data)
		if err != nil {
			return []byte{}, err
		}
		switch _type {
		case proto.MessageName((*protos.Internal_Kick)(nil)):
			return c.do_Internal_Kick(session, data, obj.(*protos.Internal_Kick))
		}
	}
	return []byte{}, nil
}

func (c *GateService) do_Internal_Kick(session uint64, data []byte, obj *protos.Internal_Kick) ([]byte, error) {
	c.closer(obj.Session)
	return []byte{}, nil
}
