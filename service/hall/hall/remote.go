// remote
package hall

import (
	"bufio"
	"encoding/base64"
	"gilgamesh/protos"
	"gilgamesh/utility/utils"
	"io"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
)

func (c *Service) service_Room(r io.ReadCloser, id uint64) {
	defer c.f.PostEvent("hall", func() {
		room, ok := c.rooms[id]
		if ok {
			room.R.Close()
			room.W.Close()
			delete(c.rooms, id)

			for session, client := range c.clients {
				protos.New_GameGateService_ServiceClient(c.f, client.Where).Call_PassThroughRequest("hall", session, &protos.Internal_GameGate_PassThroughRequest{
					Data: utils.Marshal(&protos.Hall_RoomDead{
						Id: id,
					}),
				})
			}
		}
	})

	rb := bufio.NewReader(r)
	for {
		line, err := rb.ReadString('\n')
		if err != nil {
			return
		}
		lines := strings.Fields(line)

		if len(lines) == 3 && lines[0] == "out" {
			c.f.PostEvent("hall", func() {
				session, err := strconv.ParseUint(lines[1], 10, 64)
				if err != nil {
					c.logger.Warningf("[%d] out session [%s] incorrect :%v\n", id, lines[1], err)
					return
				}
				client, ok := c.clients[session]
				if !ok {
					c.logger.Warningf("[%d] not found client session\n", session)
					return
				}

				d, err := base64.StdEncoding.DecodeString(lines[2])
				if err != nil {
					c.logger.Warningf("[%d] out data [%s] incorrect :%v\n", id, lines[2], err)
					return
				}

				protos.New_GameGateService_ServiceClient(c.f, client.Where).Call_PassThroughRequest("hall", session, &protos.Internal_GameGate_PassThroughRequest{
					Data: utils.Marshal(&protos.Duel_DataTransfer{
						Data: d,
					}),
				})
			})
		} else if len(lines) == 2 && (lines[0] == "ES_created" || lines[0] == "ES_changed") {
			c.f.PostEvent("hall", func() {
				d, err := base64.StdEncoding.DecodeString(lines[1])
				if err != nil {
					c.logger.Warningf("[%d] data [%s] incorrect :%v\n", id, lines[1], err)
					return
				}

				rc := protos.Hall_RoomCreated{}
				err = proto.Unmarshal(d, &rc)
				if err != nil {
					c.logger.Warningf("[%d] unmarshal protobuf failed :%v\n", id, err)
					return
				}

				room, ok := c.rooms[rc.Room.Id]
				if !ok {
					c.logger.Warningf("[%d] not found room id :%v\n", id, rc.Room.Id)
					return
				}

				room.Room = rc.Room

				if lines[0] == "ES_created" {
					d = utils.Marshal(&protos.Hall_RoomCreated{
						Room: rc.Room,
					})
				} else {
					d = utils.Marshal(&protos.Hall_RoomStateChanged{
						Room: rc.Room,
					})
				}

				for session, client := range c.clients {
					protos.New_GameGateService_ServiceClient(c.f, client.Where).Call_PassThroughRequest("hall", session, &protos.Internal_GameGate_PassThroughRequest{
						Data: d,
					})
				}
			})
		} else if len(lines) == 2 && lines[0] == "ES_userexit" {
			c.f.PostEvent("hall", func() {
				session, err := strconv.ParseUint(lines[1], 10, 64)
				if err != nil {
					c.logger.Warningf("[%d] offline session [%s] incorrect :%v\n", id, lines[1], err)
					return
				}

				client, ok := c.clients[session]
				if !ok {
					c.logger.Warningf("[%d] not found client session\n", session)
					return
				}

				client.Room = 0
			})
		}
	}
}
