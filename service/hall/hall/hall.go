// hall
package hall

import (
	"encoding/base64"
	"errors"
	"fmt"
	"gilgamesh/protos"
	"gilgamesh/utility/config"
	"gilgamesh/utility/utils"
	"io"
	"os/exec"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/liuhanlcj/fractal/fractal/sdk"
	"github.com/liuhanlcj/mylog"
)

type _Client struct {
	Account string
	Where   string
	Room    uint64
}

type _Room struct {
	Room *protos.Common_Room
	R    io.ReadCloser
	W    io.WriteCloser
}

type Service struct {
	logger mylog.Logger
	f      *fsdk.Fractal
	hall   *config.HallOption

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
	logger mylog.Logger,
	f *fsdk.Fractal,
	hall *config.HallOption) *Service {
	return &Service{
		logger:     logger,
		f:          f,
		hall:       hall,
		roomIdPool: 1,
		clients:    make(map[uint64]*_Client, 1000),
		rooms:      make(map[uint64]*_Room, 1000),
	}
}

func (c *Service) On_Enter(caller string, session uint64, in *protos.Internal_Hall_EnterRequest, responser func(e error)) {
	c.clients[session] = &_Client{
		Account: in.Account,
		Where:   caller,
	}

	gateClient := protos.New_GameGateService_ServiceClient(c.f, caller)

	gateClient.Call_PassThrough("hall", session, &protos.Internal_GameGate_PassThroughRequest{
		Data: utils.Marshal(&protos.Hall_YouEnter{}),
	})

	roomList := make([]*protos.Common_Room, 0, 100)
	for _, v := range c.rooms {
		if v.Room.State == protos.Common_Init {
			continue
		}
		roomList = append(roomList, v.Room)
	}
	gateClient.Call_PassThrough("hall", session, &protos.Internal_GameGate_PassThroughRequest{
		Data: utils.Marshal(&protos.Hall_RoomList{
			List: roomList,
		}),
	})

	responser(nil)
}

func (c *Service) On_Leave(caller string, session uint64, responser func(out *protos.Internal_Hall_LeaveResponse, e error)) {
	client, ok := c.clients[session]
	if ok {
		if client.Room != 0 {
			room, ok := c.rooms[client.Room]
			if ok {
				room.W.Write([]byte(fmt.Sprintf("offline %d\n", session)))
			}
		}
		delete(c.clients, session)
	}

	gateClient := protos.New_GameGateService_ServiceClient(c.f, caller)

	gateClient.Call_PassThrough("hall", session, &protos.Internal_GameGate_PassThroughRequest{
		Data: utils.Marshal(&protos.Hall_YouLeave{}),
	})

	responser(&protos.Internal_Hall_LeaveResponse{
		Success: true,
	}, nil)
}

func (c *Service) On_CreateRoom(caller string, session uint64, in *protos.Internal_Hall_CreateRoomRequest, responser func(e error)) {
	client, ok := c.clients[session]
	if !ok {
		c.logger.Warningf("[%d] not found client session\n", session)
		responser(ErrNotFoundClient)
		return
	}

	cmd := exec.Command(c.hall.RoomExe)
	pwd, _ := filepath.Split(c.hall.RoomExe)
	cmd.Dir = pwd

	r, err := cmd.StdoutPipe()
	if err != nil {
		responser(err)
		c.logger.Warningf("[%s %d] create stdout pipe failed : %v\n", client.Account, session, err)
		return
	}
	w, err := cmd.StdinPipe()
	if err != nil {
		responser(err)
		c.logger.Warningf("[%s %d] create stdin pipe failed : %v\n", client.Account, session, err)
		return
	}

	id := c.roomIdPool
	c.roomIdPool++
	c.rooms[id] = &_Room{
		Room: &protos.Common_Room{
			Id:     id,
			Option: in.Option,
			State:  protos.Common_Init,
		},
		R: r,
		W: w,
	}

	go c.service_Room(r, id)

	err = cmd.Start()
	if err != nil {
		responser(err)
		delete(c.rooms, id)
		c.logger.Warningf("[%s %d] start room failed : %v\n", client.Account, session, err)
		return
	}

	d, _ := proto.Marshal(&protos.Internal_Hall_CreateRoomRequest{
		Option: in.Option,
	})

	client.Room = id

	w.Write([]byte(fmt.Sprintf("createroom %d %s\n", id, base64.StdEncoding.EncodeToString(d))))
	w.Write([]byte(fmt.Sprintf("new %d\n", session)))

	responser(nil)
}

func (c *Service) On_EnterRoom(caller string, session uint64, in *protos.Internal_Hall_EnterRoomRequest, responser func(e error)) {
	client, ok := c.clients[session]
	if !ok {
		responser(ErrNotFoundClient)
		c.logger.Warningf("[%d] not found client session\n", session)
		return
	}

	room, ok := c.rooms[in.Id]
	if !ok {
		responser(nil)
		c.logger.Warningf("[%s %d] enter room but not found id : %d\n", client.Account, session, in.Id)
		return
	}

	if room.Room.Option.Password != in.Password {
		responser(nil)
		return
	}

	c.logger.Debugf("[%s %d] enter room : %d %s\n", client.Account, session, in.Id, room.Room.Option.Name)

	client.Room = in.Id

	room.W.Write([]byte(fmt.Sprintf("new %d\n", session)))

	responser(nil)
}

func (c *Service) On_DataTransfer(caller string, session uint64, in *protos.Internal_Duel_DataTransfer, responser func(e error)) {
	client, ok := c.clients[session]
	if !ok {
		responser(ErrNotFoundClient)
		c.logger.Warningf("[%d] not found client session\n", session)
		return
	}

	if client.Room == 0 {
		responser(ErrNotFoundRoom)
		c.logger.Warningf("[%s %d] not in room\n", client.Account, session)
		return
	}

	room, ok := c.rooms[client.Room]
	if !ok {
		responser(ErrNotFoundRoom)
		c.logger.Warningf("[%s %d] not found room id : %d\n", client.Account, session, client.Room)
		return
	}

	room.W.Write([]byte(fmt.Sprintf("in %d %s\n", session, base64.StdEncoding.EncodeToString(in.Data))))

	c.logger.Debugf("in %d %v\n", session, in.Data)

	responser(nil)
}
