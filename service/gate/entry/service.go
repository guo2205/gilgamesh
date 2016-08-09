// service
package entry

import (
	"context"
	"errors"
	"gilgamesh/utility/socket"
	"io"
	"time"
)

var (
	ErrOverFlow error = errors.New("over flow")
)

type _ConnService struct {
	entry    *GateEntry
	session  uint64
	conn     *socket.Socket
	ctx      context.Context
	cancel   context.CancelFunc
	writerCh chan []byte
}

func newConnService(entry *GateEntry, session uint64, conn *socket.Socket) *_ConnService {
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
	defer c.entry.f.PostMail("gate", 1, "entry", c.session, nil)

	ch := make(chan _ReadEvent, 30)

	go c.reader(ch)
	go c.writer()

	seq := 0

	for {
		select {
		case ev := <-ch:
			if ev.Error == io.EOF {
				return nil
			}
			if ev.Error != nil {
				return ev.Error
			}
			c.entry.f.PostMail("gate", 0, "entry", c.session, ev.Data)
			seq++
		case <-c.ctx.Done():
			return c.ctx.Err()
		case <-time.After(time.Second * 5):
			if seq > int(c.entry.gateOption.PerSecondMaxPacket*5) {
				return ErrOverFlow
			}
			seq = 0
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
		err := c.conn.WritePacketOvertime(d, c.entry.gateOption.Timeout)
		if err != nil {
			return err
		}
	}
	return nil
}
