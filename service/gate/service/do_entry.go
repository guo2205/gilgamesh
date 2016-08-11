// do_entry
package service

import (
	"errors"
	"gilgamesh/protos"
	"time"

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
		obj := protos.Public_Cts_Login{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Login(session, data, &obj)
	case proto.MessageName((*protos.Public_Cts_Resource_GetAvatar)(nil)):
		obj := protos.Public_Cts_Resource_GetAvatar{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Resource_GetAvatar(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Resource_GetLFList)(nil)):
		obj := protos.Public_Cts_Resource_GetLFList{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Resource_GetLFList(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Resource_GetLFListData)(nil)):
		obj := protos.Public_Cts_Resource_GetLFListData{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Resource_GetLFListData(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Resource_UploadAvatar)(nil)):
		obj := protos.Public_Cts_Resource_UploadAvatar{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Resource_UploadAvatar(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Player_Create)(nil)):
		obj := protos.Public_Cts_Player_Create{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Player_Create(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Player_Modify)(nil)):
		obj := protos.Public_Cts_Player_Modify{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Player_Modify(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Player_Query)(nil)):
		obj := protos.Public_Cts_Player_Query{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Player_Query(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Videotape_Get)(nil)):
		obj := protos.Public_Cts_Videotape_Get{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Videotape_Get(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Videotape_QueryList)(nil)):
		obj := protos.Public_Cts_Videotape_QueryList{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Videotape_QueryList(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Deck_Download)(nil)):
		obj := protos.Public_Cts_Deck_Download{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Deck_Download(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Deck_Query)(nil)):
		obj := protos.Public_Cts_Deck_Query{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Deck_Query(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Deck_Remove)(nil)):
		obj := protos.Public_Cts_Deck_Remove{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Deck_Remove(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Deck_Upload)(nil)):
		obj := protos.Public_Cts_Deck_Upload{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Deck_Upload(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Hall_CreateRoom)(nil)):
		obj := protos.Public_Cts_Hall_CreateRoom{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Hall_CreateRoom(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Hall_EnterRoom)(nil)):
		obj := protos.Public_Cts_Hall_EnterRoom{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Hall_EnterRoom(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Hall_Room_ChangeCamp)(nil)):
		obj := protos.Public_Cts_Hall_Room_ChangeCamp{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Hall_Room_ChangeCamp(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Hall_Room_ChangeMaster)(nil)):
		obj := protos.Public_Cts_Hall_Room_ChangeMaster{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Hall_Room_ChangeMaster(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Hall_Room_ChangeReady)(nil)):
		obj := protos.Public_Cts_Hall_Room_ChangeReady{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Hall_Room_ChangeReady(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Hall_Room_ChangeUsedeck)(nil)):
		obj := protos.Public_Cts_Hall_Room_ChangeUsedeck{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Hall_Room_ChangeUsedeck(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Hall_Room_Kick)(nil)):
		obj := protos.Public_Cts_Hall_Room_Kick{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Hall_Room_Kick(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Hall_Room_Leave)(nil)):
		obj := protos.Public_Cts_Hall_Room_Leave{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Hall_Room_Leave(session, &obj)
	case proto.MessageName((*protos.Public_Cts_Hall_Room_StartDuel)(nil)):
		obj := protos.Public_Cts_Hall_Room_StartDuel{}
		err := proto.Unmarshal(data, &obj)
		if err != nil {
			return err
		}
		return c.do_Public_Cts_Hall_Room_StartDuel(session, &obj)
	default:
		return ErrUnknownProtoType
	}
}

func (c *GateService) doEntryOffline(session uint64) error {
	c.closer(session)
	delete(c.clients, session)
	return nil
}

func (c *GateService) do_Public_Cts_Login(session uint64, data []byte, obj *protos.Public_Cts_Login) error {
	go func() {
		ret, _, err := c.f.SendMail("auth@public.auth", 0, "gate", session, data, time.Second*3)
		if err != nil {
			c.closer(session)
			return
		}
		loginResponse := protos.Public_Stc_LoginResponse{}
		err = proto.Unmarshal(ret, &loginResponse)
		if err != nil {
			c.closer(session)
			return
		}

		err = c.writer(session, ret)
		if err != nil {
			c.closer(session)
			return
		}

		if !loginResponse.State {
			c.closer(session)
			return
		}

		locker := protos.Internal_AcquireLocker{
			Account: obj.Account,
			Lock:    true,
		}
		lockData, _ := proto.Marshal(&locker)
		locker.Lock = false
		unlockData, _ := proto.Marshal(&locker)

		_, _, err = c.f.SendMail("locker@public.global", 0, "gate", session, lockData, time.Second*12)
		if err != nil {
			c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			c.closer(session)
			return
		}

		queryOnline := protos.Internal_QueryOnline{
			Account: obj.Account,
		}
		d, _ := proto.Marshal(&queryOnline)

		ret, _, err = c.f.SendMail("online@public.global", 0, "gate", session, d, time.Second*12)
		if err != nil {
			c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			c.closer(session)
			return
		}

		queryOnlineResponse := protos.Internal_QueryOnlineResponse{}
		err = proto.Unmarshal(ret, &queryOnlineResponse)
		if err != nil {
			c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			c.closer(session)
			return
		}

		if queryOnlineResponse.State {
			kick := protos.Internal_Kick{
				Session: queryOnlineResponse.Session,
			}
			d, _ := proto.Marshal(&kick)
			_, _, err := c.f.SendMail(queryOnlineResponse.Where, 0, "gate", session, d, time.Second*6)
			if err != nil {
				c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
				c.closer(session)
				return
			}
		}

		setOnline := protos.Internal_SetOnline{
			Account: obj.Account,
		}
		d, _ = proto.Marshal(&setOnline)
		_, _, err = c.f.SendMail("online@public.global", 0, "gate", session, d, time.Second*6)
		if err != nil {
			c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			c.closer(session)
			return
		}

		c.f.InsertEvent("gate", func() {
			client := &_Client{
				Session: session,
				Account: obj.Account,
			}
			c.clients[session] = client

			go func() {
				enterHall := protos.Internal_EnterHall{
					Account: obj.Account,
				}
				d, _ = proto.Marshal(&enterHall)
				_, _, err = c.f.SendMail("hall@ygo.hall", 0, "gate", session, d, time.Second*6)
				if err != nil {
					c.f.InsertEvent("gate", func() {
						c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
						delete(c.clients, session)
						c.closer(session)
					})
					return
				}
				c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			}()
		})
	}()
	return nil
}

func (c *GateService) do_Public_Cts_Resource_GetAvatar(session uint64, obj *protos.Public_Cts_Resource_GetAvatar) error {
	return nil
}

func (c *GateService) do_Public_Cts_Resource_GetLFList(session uint64, obj *protos.Public_Cts_Resource_GetLFList) error {
	return nil
}

func (c *GateService) do_Public_Cts_Resource_GetLFListData(session uint64, obj *protos.Public_Cts_Resource_GetLFListData) error {
	return nil
}

func (c *GateService) do_Public_Cts_Resource_UploadAvatar(session uint64, obj *protos.Public_Cts_Resource_UploadAvatar) error {
	return nil
}

func (c *GateService) do_Public_Cts_Player_Create(session uint64, obj *protos.Public_Cts_Player_Create) error {
	return nil
}

func (c *GateService) do_Public_Cts_Player_Modify(session uint64, obj *protos.Public_Cts_Player_Modify) error {
	return nil
}

func (c *GateService) do_Public_Cts_Player_Query(session uint64, obj *protos.Public_Cts_Player_Query) error {
	return nil
}

func (c *GateService) do_Public_Cts_Videotape_Get(session uint64, obj *protos.Public_Cts_Videotape_Get) error {
	return nil
}

func (c *GateService) do_Public_Cts_Videotape_QueryList(session uint64, obj *protos.Public_Cts_Videotape_QueryList) error {
	return nil
}

func (c *GateService) do_Public_Cts_Deck_Download(session uint64, obj *protos.Public_Cts_Deck_Download) error {
	return nil
}

func (c *GateService) do_Public_Cts_Deck_Query(session uint64, obj *protos.Public_Cts_Deck_Query) error {
	return nil
}

func (c *GateService) do_Public_Cts_Deck_Remove(session uint64, obj *protos.Public_Cts_Deck_Remove) error {
	return nil
}

func (c *GateService) do_Public_Cts_Deck_Upload(session uint64, obj *protos.Public_Cts_Deck_Upload) error {
	return nil
}

func (c *GateService) do_Public_Cts_Hall_CreateRoom(session uint64, obj *protos.Public_Cts_Hall_CreateRoom) error {
	return nil
}

func (c *GateService) do_Public_Cts_Hall_EnterRoom(session uint64, obj *protos.Public_Cts_Hall_EnterRoom) error {
	return nil
}

func (c *GateService) do_Public_Cts_Hall_Room_ChangeCamp(session uint64, obj *protos.Public_Cts_Hall_Room_ChangeCamp) error {
	return nil
}

func (c *GateService) do_Public_Cts_Hall_Room_ChangeMaster(session uint64, obj *protos.Public_Cts_Hall_Room_ChangeMaster) error {
	return nil
}

func (c *GateService) do_Public_Cts_Hall_Room_ChangeReady(session uint64, obj *protos.Public_Cts_Hall_Room_ChangeReady) error {
	return nil
}

func (c *GateService) do_Public_Cts_Hall_Room_ChangeUsedeck(session uint64, obj *protos.Public_Cts_Hall_Room_ChangeUsedeck) error {
	return nil
}

func (c *GateService) do_Public_Cts_Hall_Room_Kick(session uint64, obj *protos.Public_Cts_Hall_Room_Kick) error {
	return nil
}

func (c *GateService) do_Public_Cts_Hall_Room_Leave(session uint64, obj *protos.Public_Cts_Hall_Room_Leave) error {
	return nil
}

func (c *GateService) do_Public_Cts_Hall_Room_StartDuel(session uint64, obj *protos.Public_Cts_Hall_Room_StartDuel) error {
	return nil
}
