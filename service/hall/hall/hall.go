// hall
package hall

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"fractal/fractal"
	"gilgamesh/protos"
	"gilgamesh/utility/config"
	"gilgamesh/utility/utils"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/liuhanlcj/mylog"
)

type _Client struct {
	Where     string
	Account   string
	RoomWhere uint64
}

type _Room struct {
	Room *protos.Public_Room
	R    io.ReadCloser
	W    io.WriteCloser
}

type Service struct {
	fractal.DefaultServiceProvider
	logger mylog.Logger
	f      *fractal.Fractal
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
	f *fractal.Fractal,
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
	case proto.MessageName((*protos.Public_Cts_Hall_CreateRoom)(nil)):
		return c.do_Public_Cts_Hall_CreateRoom(session, obj.(*protos.Public_Cts_Hall_CreateRoom))
	case proto.MessageName((*protos.Public_Cts_Hall_EnterRoom)(nil)):
		return c.do_Public_Cts_Hall_EnterRoom(session, obj.(*protos.Public_Cts_Hall_EnterRoom))

	case proto.MessageName((*protos.Public_Stc_Hall_RoomDestoried)(nil)):
		return c.do_Public_Stc_Hall_RoomDestoried(session, data, obj.(*protos.Public_Stc_Hall_RoomDestoried))
	case proto.MessageName((*protos.Public_Stc_Hall_RoomStateChanged)(nil)):
		return c.do_Public_Stc_Hall_RoomStateChanged(session, data, obj.(*protos.Public_Stc_Hall_RoomStateChanged))
	case proto.MessageName((*protos.Public_Cts_Duel)(nil)):
		return c.do_Public_Cts_Duel(session, data, obj.(*protos.Public_Cts_Duel))
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
		if client.RoomWhere != 0 {
			room, ok := c.rooms[client.RoomWhere]
			if ok {
				room.W.Write([]byte(fmt.Sprintf("offline %d\n", session)))
			}
			client.RoomWhere = 0
		} else {
			delete(c.clients, session)
		}
	}
	c.f.PostMail(caller, 0, "hall", session, utils.Marshal(&protos.Public_Stc_Hall_YouLeaveHall{}))
	return []byte{}, nil
}

func (c *Service) do_Internal_Hall_Room_PlayerLeave(caller string, session uint64, obj *protos.Internal_Hall_Room_PlayerLeave) ([]byte, error) {
	client, ok := c.clients[session]
	if ok {
		client.RoomWhere = 0
	}
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Hall_CreateRoom(session uint64, obj *protos.Public_Cts_Hall_CreateRoom) ([]byte, error) {
	client, ok := c.clients[session]
	if !ok {
		return []byte{}, nil
	}

	cmd := exec.Command(c.hall.RoomExe)
	pwd, _ := filepath.Split(c.hall.RoomExe)
	cmd.Dir = pwd

	r, err := cmd.StdoutPipe()
	if err != nil {
		c.logger.Debug("create room r pipe failed :", err)
		return []byte{}, nil
	}
	w, err := cmd.StdinPipe()
	if err != nil {
		c.logger.Debug("create room w pipe failed :", err)
		return []byte{}, nil
	}

	id := c.roomIdPool
	c.roomIdPool++
	c.rooms[id] = &_Room{
		Room: &protos.Public_Room{
			Id:     id,
			Option: obj.Option,
			State:  protos.Public_Init,
		},
		R: r,
		W: w,
	}

	err = cmd.Start()
	if err != nil {
		c.logger.Debug("start room failed :", err)
		return []byte{}, nil
	}

	go func() {
		rb := bufio.NewReader(r)
		shutdown := false
		for {
			line, err := rb.ReadString('\n')
			if err != nil {
				shutdown = true
			}
			lines := strings.Fields(line)
			switch len(lines) {
			case 3:
				if lines[0] == "out" {
					c.f.InsertEvent("hall", func() {
						session, err := strconv.ParseUint(lines[1], 10, 64)
						if err != nil {
							c.logger.Warning("room in session incorrect :", lines[1], err)
							return
						}
						client, ok := c.clients[session]
						if ok {
							d, err := base64.StdEncoding.DecodeString(lines[2])
							if err != nil {
								c.logger.Warning("room in data incorrect :", lines[2], err)
								return
							}
							c.logger.Debug("out", session, lines[2], d)
							duel := protos.Public_Stc_Duel{
								Data: d,
							}
							d, err = proto.Marshal(&duel)
							if err != nil {
								c.logger.Warning("marshal duel failed :", err)
								return
							}
							obj := protos.Gilgamesh{
								Type: proto.MessageName(&duel),
								Data: d,
							}
							d, err = proto.Marshal(&obj)
							if err != nil {
								c.logger.Warning("marshal obj failed :", err)
								return
							}
							c.f.PostMail(client.Where, 0, "hall", session, d)
						} else {
							c.logger.Warning("unknown session :", session)
						}
					})
				} else if lines[0] == "ES_created" {
					c.f.InsertEvent("hall", func() {
						d, err := base64.StdEncoding.DecodeString(lines[2])
						if err != nil {
							return
						}
						rc := protos.Public_Stc_Hall_RoomCreated{}
						err = proto.Unmarshal(d, &rc)
						if err != nil {
							return
						}

						room, ok := c.rooms[rc.Room.Id]
						if !ok {
							return
						}

						room.Room = rc.Room
						for s, v := range c.clients {
							c.f.PostMail(v.Where, 0, "hall", s, d)
						}

						c.logger.Debug("room created :", rc.Room.Id, rc.Room.Option.Name)
					})
				}
			}

			if shutdown {
				c.f.InsertEvent("hall", func() {
					room, ok := c.rooms[id]
					if ok {
						room.R.Close()
						room.W.Close()
						delete(c.rooms, id)
					}
				})
				return
			}
		}
	}()

	d := utils.Marshal(&protos.Internal_Hall_CreateRoom{
		Account: client.Account,
		Option:  obj.Option,
	})

	client.RoomWhere = id

	tmp := fmt.Sprintf("createroom %d %d %s\n", id, session, base64.StdEncoding.EncodeToString(d))
	c.logger.Debugf("createroom %d %d %v\n", id, session, d)
	w.Write([]byte(tmp))

	tmp = fmt.Sprintf("new %d\n", session)
	c.logger.Debug(tmp)
	w.Write([]byte(tmp))

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

	if room.Room.Option.Password != obj.Password {
		return []byte{}, nil
	}

	client.RoomWhere = obj.Id

	room.W.Write([]byte(fmt.Sprintf("new %d\n", session)))

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

func (c *Service) do_Public_Cts_Duel(session uint64, data []byte, obj *protos.Public_Cts_Duel) ([]byte, error) {
	client, ok := c.clients[session]
	if ok {
		if client.RoomWhere != 0 {
			room, ok := c.rooms[client.RoomWhere]
			if ok {
				c.logger.Debugf("in %d %v\n", session, obj.Data)
				room.W.Write([]byte(fmt.Sprintf("in %d %s\n", session, base64.StdEncoding.EncodeToString(obj.Data))))
			}
		}
	}
	return []byte{}, nil
}
