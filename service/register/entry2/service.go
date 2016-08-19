// service
package entry2

import (
	"context"
	"encoding/hex"
	"errors"
	"gilgamesh/protos"
	"gilgamesh/utility/models"
	"gilgamesh/utility/socket"
	"io"
	"log"

	"github.com/golang/protobuf/proto"
)

var (
	ErrOverFlow error = errors.New("over flow")
)

type _ConnService struct {
	entry    *GateEntry2
	session  uint64
	conn     *socket.Socket
	ctx      context.Context
	cancel   context.CancelFunc
	writerCh chan []byte
}

func newConnService(entry *GateEntry2, session uint64, conn *socket.Socket) *_ConnService {
	ctx, cancel := context.WithCancel(context.Background())
	return &_ConnService{
		entry:    entry,
		session:  session,
		conn:     conn,
		ctx:      ctx,
		cancel:   cancel,
		writerCh: make(chan []byte, 30),
	}
}

func (c *_ConnService) Close() {
	c.cancel()
	close(c.writerCh)
}

func (c *_ConnService) WritePacket(d []byte) {
	c.writerCh <- d
}

func (c *_ConnService) RunService() error {
	defer c.conn.Close()

	ch := make(chan _ReadEvent, 30)

	go c.reader(ch)
	go c.writer()

	for {
		select {
		case ev := <-ch:
			if ev.Error == io.EOF {
				return nil
			}
			if ev.Error != nil {
				return ev.Error
			}

			register := protos.Public_Cts_Register{}
			err := proto.Unmarshal(ev.Data, &register)
			if err != nil {
				log.Println("unmarshal failed :", err)
				return err
			}

			log.Println("register :", register)

			response := protos.Public_Stc_RegisterResponse{}

			reason, ok, err := models.Register(register.Account, hex.EncodeToString(register.Password))
			if err == nil && ok {
				response.State = true
			} else {
				response.State = false
				response.Reason = reason
			}
			log.Println("response :", response)
			d, _ := proto.Marshal(&response)
			c.conn.WritePacket(d)

		case <-c.ctx.Done():
			return c.ctx.Err()
		}
	}
}

type _ReadEvent struct {
	Data  []byte
	Error error
}

func (c *_ConnService) reader(ch chan _ReadEvent) {
	defer close(ch)

	for {
		d, err := c.conn.ReadPacket()
		if err != nil {
			ch <- _ReadEvent{
				Error: err,
			}
			return
		}
		ch <- _ReadEvent{
			Data: d,
		}
	}
}

func (c *_ConnService) writer() error {
	for d := range c.writerCh {
		err := c.conn.WritePacketOvertime(d, c.entry.option.Timeout)
		if err != nil {
			return err
		}
	}
	return nil
}
