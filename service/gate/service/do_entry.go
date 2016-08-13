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
		return err
	}

	switch _type {
	case proto.MessageName((*protos.Public_Cts_Login)(nil)):
		return c.do_Public_Cts_Login(session, data, obj.(*protos.Public_Cts_Login))
	case proto.MessageName((*protos.Public_Cts_Resource_GetAvatar)(nil)):
		return c.do_Public_Cts_Resource_GetAvatar(session, data, obj.(*protos.Public_Cts_Resource_GetAvatar))
	case proto.MessageName((*protos.Public_Cts_Resource_GetLFList)(nil)):
		return c.do_Public_Cts_Resource_GetLFList(session, data, obj.(*protos.Public_Cts_Resource_GetLFList))
	case proto.MessageName((*protos.Public_Cts_Resource_GetLFListData)(nil)):
		return c.do_Public_Cts_Resource_GetLFListData(session, data, obj.(*protos.Public_Cts_Resource_GetLFListData))
	case proto.MessageName((*protos.Public_Cts_Resource_UploadAvatar)(nil)):
		return c.do_Public_Cts_Resource_UploadAvatar(session, data, obj.(*protos.Public_Cts_Resource_UploadAvatar))
	case proto.MessageName((*protos.Public_Cts_Player_Create)(nil)):
		return c.do_Public_Cts_Player_Create(session, data, obj.(*protos.Public_Cts_Player_Create))
	case proto.MessageName((*protos.Public_Cts_Player_Modify)(nil)):
		return c.do_Public_Cts_Player_Modify(session, data, obj.(*protos.Public_Cts_Player_Modify))
	case proto.MessageName((*protos.Public_Cts_Player_Query)(nil)):
		return c.do_Public_Cts_Player_Query(session, data, obj.(*protos.Public_Cts_Player_Query))
	case proto.MessageName((*protos.Public_Cts_Videotape_Get)(nil)):
		return c.do_Public_Cts_Videotape_Get(session, data, obj.(*protos.Public_Cts_Videotape_Get))
	case proto.MessageName((*protos.Public_Cts_Videotape_QueryList)(nil)):
		return c.do_Public_Cts_Videotape_QueryList(session, data, obj.(*protos.Public_Cts_Videotape_QueryList))
	case proto.MessageName((*protos.Public_Cts_Deck_Download)(nil)):
		return c.do_Public_Cts_Deck_Download(session, data, obj.(*protos.Public_Cts_Deck_Download))
	case proto.MessageName((*protos.Public_Cts_Deck_Query)(nil)):
		return c.do_Public_Cts_Deck_Query(session, data, obj.(*protos.Public_Cts_Deck_Query))
	case proto.MessageName((*protos.Public_Cts_Deck_Remove)(nil)):
		return c.do_Public_Cts_Deck_Remove(session, data, obj.(*protos.Public_Cts_Deck_Remove))
	case proto.MessageName((*protos.Public_Cts_Deck_Upload)(nil)):
		return c.do_Public_Cts_Deck_Upload(session, data, obj.(*protos.Public_Cts_Deck_Upload))
	case proto.MessageName((*protos.Public_Cts_Hall_CreateRoom)(nil)):
		return c.do_Public_Cts_Hall_CreateRoom(session, data, obj.(*protos.Public_Cts_Hall_CreateRoom))
	case proto.MessageName((*protos.Public_Cts_Hall_EnterRoom)(nil)):
		return c.do_Public_Cts_Hall_EnterRoom(session, data, obj.(*protos.Public_Cts_Hall_EnterRoom))
	case proto.MessageName((*protos.Public_Cts_Hall_Room_ChangeCamp)(nil)):
		return c.do_Public_Cts_Hall_Room_ChangeCamp(session, data, obj.(*protos.Public_Cts_Hall_Room_ChangeCamp))
	case proto.MessageName((*protos.Public_Cts_Hall_Room_ChangeMaster)(nil)):
		return c.do_Public_Cts_Hall_Room_ChangeMaster(session, data, obj.(*protos.Public_Cts_Hall_Room_ChangeMaster))
	case proto.MessageName((*protos.Public_Cts_Hall_Room_ChangeReady)(nil)):
		return c.do_Public_Cts_Hall_Room_ChangeReady(session, data, obj.(*protos.Public_Cts_Hall_Room_ChangeReady))
	case proto.MessageName((*protos.Public_Cts_Hall_Room_Kick)(nil)):
		return c.do_Public_Cts_Hall_Room_Kick(session, data, obj.(*protos.Public_Cts_Hall_Room_Kick))
	case proto.MessageName((*protos.Public_Cts_Hall_Room_Leave)(nil)):
		return c.do_Public_Cts_Hall_Room_Leave(session, data, obj.(*protos.Public_Cts_Hall_Room_Leave))
	case proto.MessageName((*protos.Public_Cts_Hall_Room_StartDuel)(nil)):
		return c.do_Public_Cts_Hall_Room_StartDuel(session, data, obj.(*protos.Public_Cts_Hall_Room_StartDuel))
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
		// TODO 通知大厅玩家离线
	}
	delete(c.clients, session)
	return nil
}

func (c *Service) do_Public_Cts_Login(session uint64, data []byte, obj *protos.Public_Cts_Login) error {
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

		queryOnlineResponse := protos.Internal_Global_Online_QueryResponse{}
		err = proto.Unmarshal(ret, &queryOnlineResponse)
		if err != nil {
			c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			c.closer(session)
			return
		}

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
				ret, _, err = c.f.SendMail("player@ygo.database", 0, "gate", session,
					utils.Marshal(&protos.Internal_Database_Player_Query{
						Account: obj.Account,
					}), time.Second*6)
				if err != nil {
					c.f.InsertEvent("gate", func() {
						c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
						c.closer(session)
					})
					return
				}

				if len(ret) == 0 {
					if err := c.writer(session,
						utils.Marshal(&protos.Public_Stc_Player_NeedCreatePlayer{})); err != nil {
						c.f.InsertEvent("gate", func() {
							c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
							c.closer(session)
						})
						return
					}
				} else {
					_, _, err = c.f.SendMail("hall@ygo.hall", 0, "gate", session,
						utils.Marshal(&protos.Internal_Hall_Enter{
							Account: obj.Account,
						}), time.Second*6)
					if err != nil {
						c.f.InsertEvent("gate", func() {
							c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
							c.closer(session)
						})
						return
					}
				}
				c.f.PostMail("locker@public.global", 0, "gate", session, unlockData)
			}()
		})
	}()
	return nil
}

func (c *Service) do_Public_Cts_Resource_GetAvatar(session uint64, data []byte, obj *protos.Public_Cts_Resource_GetAvatar) error {
	go func() {
		ret, _, err := c.f.SendMail("avatar@ygo.database", 0, "gate", session, data, time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Resource_GetAvatarResponse{
					AvatarHashId: obj.HashId,
				}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Resource_UploadAvatar(session uint64, data []byte, obj *protos.Public_Cts_Resource_UploadAvatar) error {
	go func() {
		ret, _, err := c.f.SendMail("avatar@ygo.database", 0, "gate", session, data, time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Resource_UploadAvatarResponse{}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Resource_GetLFList(session uint64, data []byte, obj *protos.Public_Cts_Resource_GetLFList) error {
	go func() {
		ret, _, err := c.f.SendMail("lflist@ygo.database", 0, "gate", session, data, time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Resource_GetLFListResponse{}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Resource_GetLFListData(session uint64, data []byte, obj *protos.Public_Cts_Resource_GetLFListData) error {
	go func() {
		ret, _, err := c.f.SendMail("lflist@ygo.database", 0, "gate", session, data, time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Resource_GetLFListDataResponse{}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Player_Create(session uint64, data []byte, obj *protos.Public_Cts_Player_Create) error {
	client, ok := c.clients[session]
	if !ok {
		c.closer(session)
		return ErrNotFoundClient
	}
	go func() {
		ret, _, err := c.f.SendMail("player@ygo.database", 0, "gate", session, utils.Marshal(&protos.Internal_Database_Player_Create{
			Player:  obj.Player,
			Account: client.Account,
		}), time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Player_CreateResponse{
					State: false,
				}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
				return
			}
			_, _, err = c.f.SendMail("hall@ygo.hall", 0, "gate", session,
				utils.Marshal(&protos.Internal_Hall_Enter{
					Account: client.Account,
				}), time.Second*6)
			if err != nil {
				c.closer(session)
				return
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Player_Modify(session uint64, data []byte, obj *protos.Public_Cts_Player_Modify) error {
	client, ok := c.clients[session]
	if !ok {
		c.closer(session)
		return ErrNotFoundClient
	}
	go func() {
		ret, _, err := c.f.SendMail("player@ygo.database", 0, "gate", session, utils.Marshal(&protos.Internal_Database_Player_Modify{
			Player:  obj.Player,
			Account: client.Account,
		}), time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Player_ModifyResponse{
					State: false,
				}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Player_Query(session uint64, data []byte, obj *protos.Public_Cts_Player_Query) error {
	go func() {
		ret, _, err := c.f.SendMail("player@ygo.database", 0, "gate", session, data, time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Player_QueryResponse{}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Videotape_Get(session uint64, data []byte, obj *protos.Public_Cts_Videotape_Get) error {
	client, ok := c.clients[session]
	if !ok {
		c.closer(session)
		return ErrNotFoundClient
	}
	go func() {
		ret, _, err := c.f.SendMail("videotape@ygo.database", 0, "gate", session, utils.Marshal(&protos.Internal_Database_Videotape_Get{
			Id:      obj.Id,
			Account: client.Account,
		}), time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Videotape_VideoTapeData{
					Id: obj.Id,
				}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Videotape_QueryList(session uint64, data []byte, obj *protos.Public_Cts_Videotape_QueryList) error {
	client, ok := c.clients[session]
	if !ok {
		c.closer(session)
		return ErrNotFoundClient
	}
	go func() {
		ret, _, err := c.f.SendMail("videotape@ygo.database", 0, "gate", session, utils.Marshal(&protos.Internal_Database_Videotape_QueryList{
			Account: client.Account,
		}), time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Videotape_VideoTapeList{}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Deck_Download(session uint64, data []byte, obj *protos.Public_Cts_Deck_Download) error {
	client, ok := c.clients[session]
	if !ok {
		c.closer(session)
		return ErrNotFoundClient
	}
	go func() {
		ret, _, err := c.f.SendMail("deck@ygo.database", 0, "gate", session, utils.Marshal(&protos.Internal_Database_Deck_Download{
			Id:      obj.Id,
			Account: client.Account,
		}), time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Deck_DownloadResponse{
					State: false,
				}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Deck_Query(session uint64, data []byte, obj *protos.Public_Cts_Deck_Query) error {
	client, ok := c.clients[session]
	if !ok {
		c.closer(session)
		return ErrNotFoundClient
	}
	go func() {
		ret, _, err := c.f.SendMail("deck@ygo.database", 0, "gate", session, utils.Marshal(&protos.Internal_Database_Deck_Query{
			Account: client.Account,
		}), time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Deck_QueryListResponse{}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Deck_Remove(session uint64, data []byte, obj *protos.Public_Cts_Deck_Remove) error {
	client, ok := c.clients[session]
	if !ok {
		c.closer(session)
		return ErrNotFoundClient
	}
	go func() {
		ret, _, err := c.f.SendMail("deck@ygo.database", 0, "gate", session, utils.Marshal(&protos.Internal_Database_Deck_Remove{
			Id:      obj.Id,
			Account: client.Account,
		}), time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Deck_RemoveResponse{
					State: false,
				}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
	}()
	return nil
}

func (c *Service) do_Public_Cts_Deck_Upload(session uint64, data []byte, obj *protos.Public_Cts_Deck_Upload) error {
	client, ok := c.clients[session]
	if !ok {
		c.closer(session)
		return ErrNotFoundClient
	}
	go func() {
		ret, _, err := c.f.SendMail("deck@ygo.database", 0, "gate", session, utils.Marshal(&protos.Internal_Database_Deck_Upload{
			Deck:    obj.Deck,
			Account: client.Account,
		}), time.Second*6)
		if err != nil {
			err := c.writer(session,
				utils.Marshal(&protos.Public_Stc_Deck_UploadResponse{
					State: false,
				}))
			if err != nil {
				c.closer(session)
			}
		} else {
			err := c.writer(session, ret)
			if err != nil {
				c.closer(session)
			}
		}
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

func (c *Service) do_Public_Cts_Hall_Room_ChangeCamp(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_ChangeCamp) error {
	c.f.PostMail("hall@ygo.hall", 0, "gate", session, data)
	return nil
}

func (c *Service) do_Public_Cts_Hall_Room_ChangeMaster(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_ChangeMaster) error {
	c.f.PostMail("hall@ygo.hall", 0, "gate", session, data)
	return nil
}

func (c *Service) do_Public_Cts_Hall_Room_ChangeReady(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_ChangeReady) error {
	c.f.PostMail("hall@ygo.hall", 0, "gate", session, data)
	return nil
}

func (c *Service) do_Public_Cts_Hall_Room_Kick(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_Kick) error {
	c.f.PostMail("hall@ygo.hall", 0, "gate", session, data)
	return nil
}

func (c *Service) do_Public_Cts_Hall_Room_Leave(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_Leave) error {
	c.f.PostMail("hall@ygo.hall", 0, "gate", session, data)
	return nil
}

func (c *Service) do_Public_Cts_Hall_Room_StartDuel(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_StartDuel) error {
	c.f.PostMail("hall@ygo.hall", 0, "gate", session, data)
	return nil
}
