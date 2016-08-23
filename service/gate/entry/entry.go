// entry
package entry

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/liuhanlcj/fractal/fractal/utility/socket"
)

type SessionGenerator func() (uint64, error)
type OnData func(uint64, []byte) error

var (
	ErrSessionNotExist error = errors.New("session not exist")
)

type GateEntry struct {
	laddr            string
	sessionGenerator SessionGenerator
	onData           OnData
	option           *Option

	lock     sync.Mutex
	connList map[uint64]*_Conn
}

func NewEntry(laddr string,
	sessionGenerator SessionGenerator,
	onData OnData,
	option *Option) *GateEntry {
	return &GateEntry{
		laddr:            laddr,
		sessionGenerator: sessionGenerator,
		onData:           onData,
		option:           option,
		connList:         make(map[uint64]*_Conn, 1000),
	}
}

func (c *GateEntry) Start() error {
	laddr, err := net.ResolveTCPAddr("tcp", c.laddr)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}

	go c.listener(l)

	return nil
}

func (c *GateEntry) WriteConn(session uint64, data []byte) error {
	return c.writeConn(session, data)
}

func (c *GateEntry) CloseConn(session uint64) error {
	return c.deleteConn(session)
}

func (c *GateEntry) listener(l *net.TCPListener) {
	defer l.Close()

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			return
		}

		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(time.Second * socket.TCP_KEEPALIVE_PERIOD)
		conn.SetReadBuffer(socket.TCP_READ_BUFFER_SIZE)
		conn.SetWriteBuffer(socket.TCP_WRITE_BUFFER_SIZE)

		go c.srvConn(conn)
	}
}

func (c *GateEntry) srvConn(cc net.Conn) error {
	r, w, closer := socket.NewChanSocket(socket.NewSocket(cc), 0, time.Second*30)

	session, err := c.sessionGenerator()
	if err != nil {
		closer()
		return err
	}

	conn := &_Conn{
		entry:   c,
		session: session,
		r:       r,
		w:       w,
		closer:  closer,
	}

	c.addConn(session, conn)
	defer c.deleteConn(session)

	return conn.start()
}

func (c *GateEntry) writeConn(session uint64, data []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	conn, ok := c.connList[session]
	if !ok {
		return ErrSessionNotExist
	}

	conn.w <- data

	return nil
}

func (c *GateEntry) addConn(session uint64, conn *_Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.connList[session] = conn
}

func (c *GateEntry) deleteConn(session uint64) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	conn, ok := c.connList[session]
	if !ok {
		return ErrSessionNotExist
	}

	delete(c.connList, session)

	conn.closer()

	return nil
}
