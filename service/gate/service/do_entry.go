// do_entry
package service

import (
	"errors"
	"gilgamesh/protos"
	"gilgamesh/utility/utils"
	"time"

	"github.com/golang/protobuf/proto"
)

var (
	ErrUnknownProtoType error = errors.New("unknown proto type")
	ErrNotFoundClient   error = errors.New("not dound client")
)

func (c *Service) doEntry(_type uint32, session uint64, data []byte) ([]byte, error) {
	switch _type {
	case 0:
		return []byte{}, c.doEntryNormal(session, data)
	case 1:
		return []byte{}, c.doEntryOffline(session)
	default:
		return []byte{}, ErrUnknownMailType
	}
}

func (c *Service) doEntryNormal(session uint64, data []byte) error {
	obj, _type, err := utils.Unmarshal(data)
	if err != nil {
		c.logger.Debug("doEntryNormal failed:", err)
		return err
	}

	switch _type {
	case proto.MessageName((*protos.Public_Cts_Login)(nil)):
		return c.do_Public_Cts_Login(session, data, obj.(*protos.Public_Cts_Login))
	case proto.MessageName((*protos.Public_Cts_Hall_CreateRoom)(nil)):
		return c.do_Public_Cts_Hall_CreateRoom(session, data, obj.(*protos.Public_Cts_Hall_CreateRoom))
	case proto.MessageName((*protos.Public_Cts_Hall_EnterRoom)(nil)):
		return c.do_Public_Cts_Hall_EnterRoom(session, data, obj.(*protos.Public_Cts_Hall_EnterRoom))
	case proto.MessageName((*protos.Public_Cts_Duel)(nil)):
		return c.do_Public_Cts_Duel(session, data, obj.(*protos.Public_Cts_Duel))
	case proto.MessageName((*protos.Public_Chat)(nil)):
		return c.do_Public_Chat(session, data, obj.(*protos.Public_Chat))
	default:
		return ErrUnknownProtoType
	}
}

func (c *Service) doEntryOffline(session uint64) error {
	client, ok := c.clients[session]
	if ok {
		c.f.PostMail("online@public.global", 0, "gate", session,
			utils.Marshal(&protos.Internal_Global_Online_Set{
				Account: client.Account,
				State:   false,
			}))
		c.f.PostMail("hall@ygo.hall", 0, "gate", session,
			utils.Marshal(&protos.Internal_Hall_Leave{
				Account: client.Account,
			}))
	}
	delete(c.clients, session)
	return nil
}

func (c *Service) do_Public_Cts_Login(session uint64, data []byte, obj *protos.Public_Cts_Login) error {
	c.logger.Debug("login :", *obj)
	go func() {
		ret, _, err := c.f.SendMail("auth@public.auth", 0, "gate", session, data, time.Second*8)
		if err != nil {
			c.closer(session)
			return
		}
		uret, _, err := utils.UnmarshalWithType(ret, new(protos.Public_Stc_LoginResponse))
		if err != nil {
			c.closer(session)
			return
		}
		loginResponse := uret.(*protos.Public_Stc_LoginResponse)

		err = c.writer(session, ret)
		if err != nil {
			c.closer(session)
			return
		}

		if !loginResponse.State {
			c.closer(session)
			return
		}

		lockData := utils.Marshal(&protos.Internal_Global_Locker_Acquire{
			Account: obj.Account,
			Lock:    true,
		})
		unlockData := utils.Marshal(&protos.Internal_Global_Locker_Acquire{
			Account: obj.Account,
			Lock:    false,
		})

		ret, _, err = c.f.SendMail("locker@public.global", 0, "gate", session, lockData, time.Second*12)
		if err != nil || string(ret) != "success" {
			c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			c.closer(session)
			return
		}

		ret, _, err = c.f.SendMail("online@public.global", 0, "gate", session,
			utils.Marshal(&protos.Internal_Global_Online_Query{
				Account: obj.Account,
			}), time.Second*12)
		if err != nil {
			c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			c.closer(session)
			return
		}

		uret, _, err = utils.UnmarshalWithType(ret, new(protos.Internal_Global_Online_QueryResponse))
		if err != nil {
			c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			c.closer(session)
			return
		}
		queryOnlineResponse := uret.(*protos.Internal_Global_Online_QueryResponse)

		if queryOnlineResponse.State {
			c.f.PostMail(queryOnlineResponse.Where, 1, "gate", session,
				utils.Marshal(&protos.Internal_Gate_Kick{
					Session: queryOnlineResponse.Session,
				}))
			c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			c.closer(session)
			return
		}

		_, _, err = c.f.SendMail("online@public.global", 0, "gate", session,
			utils.Marshal(&protos.Internal_Global_Online_Set{
				Account: obj.Account,
				State:   true,
			}), time.Second*6)
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
				_, _, err = c.f.SendMail("hall@ygo.hall", 0, "gate", session,
					utils.Marshal(&protos.Internal_Hall_Enter{
						Account: obj.Account,
					}), time.Second*6)
				if err != nil {
					c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
					c.closer(session)
					return
				}

				c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			}()
		})
	}()
	return nil
}

func (c *Service) do_Public_Cts_Hall_CreateRoom(session uint64, data []byte, obj *protos.Public_Cts_Hall_CreateRoom) error {
	c.f.PostMail("hall@ygo.hall", 0, "gate", session, data)
	return nil
}

func (c *Service) do_Public_Cts_Hall_EnterRoom(session uint64, data []byte, obj *protos.Public_Cts_Hall_EnterRoom) error {
	c.f.PostMail("hall@ygo.hall", 0, "gate", session, data)
	return nil
}

func (c *Service) do_Public_Cts_Duel(session uint64, data []byte, obj *protos.Public_Cts_Duel) error {
	c.f.PostMail("hall@ygo.hall", 0, "gate", session, data)
	return nil
}

func (c *Service) do_Public_Chat(session uint64, data []byte, obj *protos.Public_Chat) error {
	c.f.PostMail("chat@public.global", 0, "gate", session, data)
	return nil
}
