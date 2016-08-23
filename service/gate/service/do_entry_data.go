// do_entry
package service

import (
	"encoding/hex"
	"errors"
	"gilgamesh/protos"
	"gilgamesh/utility/utils"
	"time"
)

var (
	ErrNotFoundClient      error = errors.New("not found client")
	ErrAuthFailed          error = errors.New("auth failed")
	ErrLockFailed          error = errors.New("lock failed")
	ErrReleaseLockFailed   error = errors.New("release lock failed")
	ErrAlreadyOnlineFailed error = errors.New("already online")
	ErrSetOnlineFailed     error = errors.New("set online failed")
)

func (c *Service) do_Auth_AuthRequest(session uint64, obj *protos.Auth_AuthRequest, responser func(e error)) {
	go func() {
		c.logger.Debugf("[%s %d] auth : [%s]\n", obj.Account, session, hex.EncodeToString(obj.Password))

		// 创建认证客户端
		authClient := protos.New_AuthService_ServiceClient(c.f, "/public/auth/?.auth")

		// 认证账号
		authResp, _, err := authClient.Call_Auth_Sync("gate", session, time.Second*3, obj)
		if err != nil {
			c.logger.Warningf("[%s %d] auth failed :%v\n", obj.Account, session, err)
			c.closer(session)
			responser(err)
			return
		}

		// 写回结果
		err = c.writer(session, utils.Marshal(authResp))
		if err != nil {
			c.logger.Warningf("[%s %d] write auth response failed :%v\n", obj.Account, session, err)
			c.closer(session)
			responser(err)
			return
		}

		// 等待写入完毕
		time.Sleep(time.Millisecond * 100)

		// 认证失败
		if !authResp.Success {
			c.closer(session)
			responser(ErrAuthFailed)
			return
		}

		// 创建锁客户端
		lockerClient := protos.New_GlobalLockerService_ServiceClient(c.f, "/public/global.locker")

		// 获取登录锁
		lockResp, _, err := lockerClient.Call_Acquire_Sync("gate", session, time.Second*3, &protos.Internal_Global_Locker_AcquireRequest{
			Key: obj.Account,
		})

		// 获取失败
		if err != nil {
			c.logger.Warningf("[%s %d] lock failed :%v\n", obj.Account, session, err)
			c.closer(session)
			responser(err)
			return
		}

		// 加锁失败
		if !lockResp.Success {
			c.logger.Warningf("[%s %d] lock failed :unknown\n", obj.Account, session)
			c.closer(session)
			responser(ErrLockFailed)
			return
		}

		// 创建在线列表客户端
		onlineClient := protos.New_GlobalOnlineStateService_ServiceClient(c.f, "/public/global.online")

		// 查询状态
		queryResp, _, err := onlineClient.Call_Query_Sync("gate", session, time.Second*3, &protos.Internal_Global_OnlineState_QueryRequest{
			Account: obj.Account,
		})

		// 查询失败
		if err != nil {
			c.logger.Warningf("[%s %d] query online state failed :%v\n", obj.Account, session, err)

			c.closer(session)
			responser(err)

			resp, _, err := lockerClient.Call_Release_Sync("gate", session, time.Second*3, &protos.Internal_Global_Locker_ReleaseRequest{
				Key: obj.Account,
			})
			if err != nil {
				c.logger.Warningf("[%s %d] release lock failed :%v\n", obj.Account, session, err)
			} else if !resp.Success {
				c.logger.Warningf("[%s %d] release lock failed :unknown\n", obj.Account, session)
			}

			return
		}

		if queryResp.State {
			// 账号在线，踢人
			_, err := protos.New_GameGateService_ServiceClient(c.f, queryResp.Where).Call_Kick("gate", session, &protos.Internal_GameGate_KickRequest{
				Session: queryResp.Session,
			})
			// 失败
			if err != nil {
				c.logger.Warningf("[%s %d] kick failed :%v\n", obj.Account, session, err)
			}

			c.closer(session)
			responser(ErrAlreadyOnlineFailed)

			releaseResponse, _, err := lockerClient.Call_Release_Sync("gate", session, time.Second*3, &protos.Internal_Global_Locker_ReleaseRequest{
				Key: obj.Account,
			})

			if err != nil {
				c.logger.Warningf("[%s %d] release lock failed :%v\n", obj.Account, session, err)
			} else if !releaseResponse.Success {
				c.logger.Warningf("[%s %d] release lock failed :unknown\n", obj.Account, session)
			}

			return
		}

		// 设置在线状态和所在网关
		setResp, _, err := onlineClient.Call_Set_Sync("gate", session, time.Second*3, &protos.Internal_Global_OnlineState_SetRequest{
			Account: obj.Account,
			State:   true,
		})
		// 失败，解锁
		if err != nil || !setResp.Success {
			if err != nil {
				c.logger.Warningf("[%s %d] set online failed :%v\n", obj.Account, session, err)
			} else if !setResp.Success {
				c.logger.Warningf("[%s %d] set online failed :unknown\n", obj.Account, session)
			}

			c.closer(session)
			responser(ErrSetOnlineFailed)

			releaseResponse, _, err := lockerClient.Call_Release_Sync("gate", session, time.Second*3, &protos.Internal_Global_Locker_ReleaseRequest{
				Key: obj.Account,
			})

			if err != nil {
				c.logger.Warningf("[%s %d] release lock failed :%v\n", obj.Account, session, err)
			} else if !releaseResponse.Success {
				c.logger.Warningf("[%s %d] release lock failed :unknown\n", obj.Account, session)
			}

			return
		}

		// 添加新客户端
		err = c.f.SendEvent("gate", func() {
			client := &_Client{
				Session: session,
				Account: obj.Account,
			}
			c.clients[session] = client
		})
		// 失败，致命错误，直接退出
		if err != nil {
			c.logger.Fatalf("[%s %d] add client failed :%v\n", obj.Account, session, err)
		}

		// 创建大厅客户端
		hallClient := protos.New_HallService_ServiceClient(c.f, "/ygo/hall.hall")

		// 通知大厅新用户进入
		_, err = hallClient.Call_Enter("gate", session, &protos.Internal_Hall_EnterRequest{
			Account: obj.Account,
		})
		// 失败，解锁
		if err != nil {
			c.logger.Warningf("[%s %d] enter hall failed :%v\n", obj.Account, session, err)

			c.closer(session)
			responser(err)

			releaseResponse, _, err := lockerClient.Call_Release_Sync("gate", session, time.Second*3, &protos.Internal_Global_Locker_ReleaseRequest{
				Key: obj.Account,
			})

			if err != nil {
				c.logger.Warningf("[%s %d] release lock failed :%v\n", obj.Account, session, err)
			} else if !releaseResponse.Success {
				c.logger.Warningf("[%s %d] release lock failed :unknown\n", obj.Account, session)
			}

			return
		}

		// 创建聊天客户端
		chatClient := protos.New_ChatService_ServiceClient(c.f, "/ygo/chat.chat")

		// 进入大厅聊天
		_, err = chatClient.Call_EnterHall("gate", session)
		// 失败，解锁
		if err != nil {
			c.logger.Warningf("[%s %d] enter chat failed :%v\n", obj.Account, session, err)

			c.closer(session)
			responser(err)

			releaseResponse, _, err := lockerClient.Call_Release_Sync("gate", session, time.Second*3, &protos.Internal_Global_Locker_ReleaseRequest{
				Key: obj.Account,
			})

			if err != nil {
				c.logger.Warningf("[%s %d] release lock failed :%v\n", obj.Account, session, err)
			} else if !releaseResponse.Success {
				c.logger.Warningf("[%s %d] release lock failed :unknown\n", obj.Account, session)
			}

			return
		}

		// 登录结束，解锁
		releaseResponse, _, err := lockerClient.Call_Release_Sync("gate", session, time.Second*3, &protos.Internal_Global_Locker_ReleaseRequest{
			Key: obj.Account,
		})

		if err != nil || !releaseResponse.Success {
			if err != nil {
				c.logger.Warningf("[%s %d] release lock failed :%v\n", obj.Account, session, err)
			} else if !releaseResponse.Success {
				c.logger.Warningf("[%s %d] release lock failed :unknown\n", obj.Account, session)
			}

			c.closer(session)
			responser(ErrReleaseLockFailed)

			return
		}

		responser(nil)
	}()
}

func (c *Service) do_Chat_HallRequest(session uint64, obj *protos.Chat_HallRequest, responser func(e error)) {
	client, ok := c.clients[session]
	if !ok {
		responser(ErrNotFoundClient)
		return
	}

	_, err := protos.New_ChatService_ServiceClient(c.f, "/ygo/chat.chat").Call_HallChat("gate", session, &protos.Internal_Chat_HallRequest{
		Account: client.Account,
		Content: obj.Content,
	})
	if err != nil {
		responser(err)
	}
	responser(nil)
}

func (c *Service) do_Hall_EnterRoomRequest(session uint64, obj *protos.Hall_EnterRoomRequest, responser func(e error)) {
	_, err := protos.New_HallService_ServiceClient(c.f, "/ygo/hall.hall").Call_EnterRoom("gate", session, &protos.Internal_Hall_EnterRoomRequest{
		Id:       obj.Id,
		Password: obj.Password,
	})
	if err != nil {
		responser(err)
	}
	responser(nil)
}

func (c *Service) do_Hall_CreateRoomRequest(session uint64, obj *protos.Hall_CreateRoomRequest, responser func(e error)) {
	_, err := protos.New_HallService_ServiceClient(c.f, "/ygo/hall.hall").Call_CreateRoom("gate", session, &protos.Internal_Hall_CreateRoomRequest{
		Option: obj.Option,
	})
	if err != nil {
		responser(err)
	}
	responser(nil)
}

func (c *Service) do_Duel_DataTransfer(session uint64, obj *protos.Duel_DataTransfer, responser func(e error)) {
	_, err := protos.New_HallService_ServiceClient(c.f, "/ygo/hall.hall").Call_DataTransfer("gate", session, &protos.Internal_Duel_DataTransfer{
		Data: obj.Data,
	})
	if err != nil {
		responser(err)
	}
	responser(nil)
}
