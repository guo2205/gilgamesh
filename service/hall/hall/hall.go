// hall
package hall

import (
	"errors"
	"fmt"
	"fractal/fractal"
	"gilgamesh/protos"
	"gilgamesh/utility/mylog"
	"gilgamesh/utility/utils"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

type _Client struct {
	Where     string
	Account   string
	RoomWhere string
}

type _Room struct {
	Where string
	Room  *protos.Public_Room
}

type Service struct {
	fractal.DefaultServiceProvider
	logger *mylog.Logger
	f      *fractal.Fractal

	roomIdPool uint64

	clients map[uint64]*_Client
	rooms   map[uint64]*_Room
}

var (
	ErrUnknownProtoType error = errors.New("unknown proto type")
	ErrNotFoundClient   error = errors.New("not found client")
	ErrNotFoundRoom     error = errors.New("not found room")
	ErrClientNotInRoom  error = errors.New("client not in room")
)

func NewService(
	logger *mylog.Logger,
	f *fractal.Fractal) *Service {
	return &Service{
		logger:     logger,
		f:          f,
		roomIdPool: 1,
		clients:    make(map[uint64]*_Client, 1000),
		rooms:      make(map[uint64]*_Room, 1000),
	}
}

func (c *Service) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	obj, ptype, err := utils.Unmarshal(data)
	if err != nil {
		return []byte{}, err
	}

	switch ptype {
	case proto.MessageName((*protos.Internal_Hall_Enter)(nil)):
		return c.do_Internal_Hall_Enter(caller, session, obj.(*protos.Internal_Hall_Enter))
	case proto.MessageName((*protos.Internal_Hall_Leave)(nil)):
		return c.do_Internal_Hall_Leave(caller, session, obj.(*protos.Internal_Hall_Leave))
	case proto.MessageName((*protos.Internal_Hall_Room_RoomInitlized)(nil)):
		return c.do_Internal_Hall_Room_RoomInitlized(caller, session, obj.(*protos.Internal_Hall_Room_RoomInitlized))
	case proto.MessageName((*protos.Public_Cts_Hall_CreateRoom)(nil)):
		return c.do_Public_Cts_Hall_CreateRoom(session, obj.(*protos.Public_Cts_Hall_CreateRoom))
	case proto.MessageName((*protos.Public_Cts_Hall_EnterRoom)(nil)):
		return c.do_Public_Cts_Hall_EnterRoom(session, obj.(*protos.Public_Cts_Hall_EnterRoom))
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
	case proto.MessageName((*protos.Public_Stc_Hall_RoomCreated)(nil)):
		return c.do_Public_Stc_Hall_RoomCreated(session, data, obj.(*protos.Public_Stc_Hall_RoomCreated))
	case proto.MessageName((*protos.Public_Stc_Hall_RoomDestoried)(nil)):
		return c.do_Public_Stc_Hall_RoomDestoried(session, data, obj.(*protos.Public_Stc_Hall_RoomDestoried))
	case proto.MessageName((*protos.Public_Stc_Hall_RoomStateChanged)(nil)):
		return c.do_Public_Stc_Hall_RoomStateChanged(session, data, obj.(*protos.Public_Stc_Hall_RoomStateChanged))
	default:
		return []byte{}, ErrUnknownProtoType
	}
}

func (c *Service) do_Internal_Hall_Enter(caller string, session uint64, obj *protos.Internal_Hall_Enter) ([]byte, error) {
	c.clients[session] = &_Client{
		Where:   caller,
		Account: obj.Account,
	}
	c.f.PostMail(caller, 0, "hall", session, utils.Marshal(&protos.Public_Stc_Hall_YouEnterHall{}))
	roomList := make([]*protos.Public_Room, 0, 100)
	for _, v := range c.rooms {
		if v.Room.State == protos.Public_Init {
			continue
		}
		roomList = append(roomList, v.Room)
	}
	c.f.PostMail(caller, 0, "hall", session, utils.Marshal(&protos.Public_Stc_Hall_RoomList{
		RoomList: roomList,
	}))
	return []byte{}, nil
}

func (c *Service) do_Internal_Hall_Leave(caller string, session uint64, obj *protos.Internal_Hall_Leave) ([]byte, error) {
	client, ok := c.clients[session]
	if ok {
		if client.RoomWhere != "" {
			c.f.PostMail(client.RoomWhere, 0, "hall", session, utils.Marshal(&protos.Public_Cts_Hall_Room_Leave{}))
		} else {
			delete(c.clients, session)
		}
	}
	c.f.PostMail(caller, 0, "hall", session, utils.Marshal(&protos.Public_Stc_Hall_YouLeaveHall{}))
	return []byte{}, nil
}

func (c *Service) do_Internal_Hall_Room_RoomInitlized(caller string, session uint64, obj *protos.Internal_Hall_Room_RoomInitlized) ([]byte, error) {
	room, ok := c.rooms[obj.Id]
	if !ok {
		return []byte{}, ErrNotFoundRoom
	}
	client, ok := c.clients[session]
	if !ok {
		return []byte{}, ErrNotFoundClient
	}

	c.logger.Debug("room initlized :", obj.Id, caller, client.Account)

	room.Where = caller
	client.RoomWhere = caller
	c.f.PostMail(caller, 0, "hall", session, utils.Marshal(&protos.Internal_Hall_CreateRoom{
		Account: client.Account,
		Option:  room.Room.Option,
	}))
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Hall_CreateRoom(session uint64, obj *protos.Public_Cts_Hall_CreateRoom) ([]byte, error) {
	id := c.roomIdPool
	c.roomIdPool++
	c.rooms[id] = &_Room{
		Room: &protos.Public_Room{
			Id:     id,
			Option: obj.Option,
			State:  protos.Public_Init,
		},
	}

	c.logger.Debug("create room :", id, *obj.Option)

	//exec.Command("room.exe", fmt.Sprint(session), fmt.Sprint(id)).Start()
	ioutil.WriteFile(`args.txt`, []byte(fmt.Sprintf(`%d %d`, session, id)), 777)
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Hall_EnterRoom(session uint64, obj *protos.Public_Cts_Hall_EnterRoom) ([]byte, error) {
	room, ok := c.rooms[obj.Id]
	if !ok {
		return []byte{}, ErrNotFoundRoom
	}
	client, ok := c.clients[session]
	if !ok {
		return []byte{}, ErrNotFoundClient
	}

	c.logger.Debug("enter room :", obj.Id, room.Room.Option.Name, client.Account)

	client.RoomWhere = room.Where
	c.f.PostMail(room.Where, 0, "hall", session, utils.Marshal(&protos.Internal_Hall_EnterRoom{
		Account:      client.Account,
		RoomPassword: obj.Password,
	}))
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Hall_Room_ChangeCamp(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_ChangeCamp) ([]byte, error) {
	client, ok := c.clients[session]
	if !ok {
		return []byte{}, ErrNotFoundClient
	}
	if client.RoomWhere == "" {
		return []byte{}, ErrClientNotInRoom
	}
	c.f.PostMail(client.RoomWhere, 0, "hall", session, data)
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Hall_Room_ChangeMaster(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_ChangeMaster) ([]byte, error) {
	client, ok := c.clients[session]
	if !ok {
		return []byte{}, ErrNotFoundClient
	}
	if client.RoomWhere == "" {
		return []byte{}, ErrClientNotInRoom
	}
	c.f.PostMail(client.RoomWhere, 0, "hall", session, data)
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Hall_Room_ChangeReady(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_ChangeReady) ([]byte, error) {
	client, ok := c.clients[session]
	if !ok {
		return []byte{}, ErrNotFoundClient
	}
	if client.RoomWhere == "" {
		return []byte{}, ErrClientNotInRoom
	}
	c.f.PostMail(client.RoomWhere, 0, "hall", session, data)
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Hall_Room_Kick(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_Kick) ([]byte, error) {
	client, ok := c.clients[session]
	if !ok {
		return []byte{}, ErrNotFoundClient
	}
	if client.RoomWhere == "" {
		return []byte{}, ErrClientNotInRoom
	}
	c.f.PostMail(client.RoomWhere, 0, "hall", session, data)
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Hall_Room_Leave(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_Leave) ([]byte, error) {
	client, ok := c.clients[session]
	if !ok {
		return []byte{}, ErrNotFoundClient
	}
	client.RoomWhere = ""
	c.f.PostMail(client.RoomWhere, 0, "hall", session, data)
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Hall_Room_StartDuel(session uint64, data []byte, obj *protos.Public_Cts_Hall_Room_StartDuel) ([]byte, error) {
	client, ok := c.clients[session]
	if !ok {
		return []byte{}, ErrNotFoundClient
	}
	if client.RoomWhere == "" {
		return []byte{}, ErrClientNotInRoom
	}
	c.f.PostMail(client.RoomWhere, 0, "hall", session, data)
	return []byte{}, nil
}

func (c *Service) do_Public_Stc_Hall_RoomCreated(session uint64, data []byte, obj *protos.Public_Stc_Hall_RoomCreated) ([]byte, error) {
	room, ok := c.rooms[obj.Room.Id]
	if !ok {
		return []byte{}, nil
	}

	room.Room = obj.Room
	for s, v := range c.clients {
		c.f.PostMail(v.Where, 0, "hall", s, data)
	}

	c.logger.Debug("room created :", obj.Room.Id, room.Room.Option.Name)

	return []byte{}, nil
}

func (c *Service) do_Public_Stc_Hall_RoomDestoried(session uint64, data []byte, obj *protos.Public_Stc_Hall_RoomDestoried) ([]byte, error) {
	_, ok := c.rooms[obj.Id]
	if !ok {
		return []byte{}, nil
	}
	delete(c.rooms, obj.Id)
	for s, v := range c.clients {
		c.f.PostMail(v.Where, 0, "hall", s, data)
	}

	c.logger.Debug("room destoried :", obj.Id)

	return []byte{}, nil
}

func (c *Service) do_Public_Stc_Hall_RoomStateChanged(session uint64, data []byte, obj *protos.Public_Stc_Hall_RoomStateChanged) ([]byte, error) {
	room, ok := c.rooms[obj.Room.Id]
	if !ok {
		return []byte{}, nil
	}
	room.Room = obj.Room
	for s, v := range c.clients {
		c.f.PostMail(v.Where, 0, "hall", s, data)
	}

	c.logger.Debug("room state changed :", obj.Room.Id)

	return []byte{}, nil
}
