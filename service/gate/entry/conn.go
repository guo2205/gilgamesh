// conn
package entry

import (
	"errors"
	"time"

	"github.com/liuhanlcj/fractal/fractal/utility/socket"
)

var (
	ErrOutOfLimit error = errors.New("out of limit")
)

type _Conn struct {
	entry   *GateEntry
	session uint64
	r       socket.Reader
	w       socket.Writer
	closer  socket.Closer
}

func (c *_Conn) start() error {
	defer c.entry.deleteConn(c.session)

	qos := 0

	for {
		select {
		case ev := <-c.r:
			if ev.Error != nil {
				return ev.Error
			}

			qos++
			if qos > int(c.entry.option.PerSecondMaxPacket) {
				return ErrOutOfLimit
			}

			if err := c.entry.onData(c.session, ev.Data); err != nil {
				return err
			}
		case <-time.After(time.Second):
			qos = 0
		}
	}

	return nil
}
