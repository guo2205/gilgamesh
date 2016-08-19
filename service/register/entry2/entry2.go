// entry2
package entry2

import (
	"errors"
	"gilgamesh/utility/config"
	"gilgamesh/utility/socket"
	"gilgamesh/utility/utils"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	_TCP_KEEPALIVE_PERIOD  = 15
	_TCP_READ_BUFFER_SIZE  = 512
	_TCP_WRITE_BUFFER_SIZE = 512
)

var (
	ErrConnNotExist error = errors.New("conn session not exist")
)

type GateEntry2 struct {
	option *config.RegisterOption

	idPool uint64

	connLock sync.Mutex
	connMap  map[uint64]*_ConnService
}

func NewGateEntry2(option *config.RegisterOption) *GateEntry2 {
	return &GateEntry2{
		option:  option,
		connMap: make(map[uint64]*_ConnService, 500),
	}
}

func (c *GateEntry2) Run() error {
	l, err := c.listen()
	if err != nil {
		return err
	}
	c.listener(l)

	return nil
}

func (c *GateEntry2) WritePacket(session uint64, d []byte) error {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	s, ok := c.connMap[session]
	if !ok {
		return ErrConnNotExist
	}

	s.WritePacket(d)

	return nil
}

func (c *GateEntry2) CloseConn(session uint64) {
	c.deleteConn(session)
}

func (c *GateEntry2) listen() (*net.TCPListener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", c.option.LocalAddr)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (c *GateEntry2) listener(l *net.TCPListener) {
	defer l.Close()

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			return
		}

		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(time.Second * _TCP_KEEPALIVE_PERIOD)
		conn.SetReadBuffer(_TCP_READ_BUFFER_SIZE)
		conn.SetWriteBuffer(_TCP_WRITE_BUFFER_SIZE)

		go c.serviceConn(conn)
	}
}

func (c *GateEntry2) serviceConn(cc net.Conn) error {
	err := utils.Handshake(cc, c.option.Cookie, c.option.Timeout)
	if err != nil {
		cc.Close()
		return err
	}

	s := newConnService(c, atomic.AddUint64(&c.idPool, 1), socket.NewSocket(cc))

	c.insertConn(s)
	defer c.deleteConn(s.session)

	return s.RunService()
}

func (c *GateEntry2) insertConn(s *_ConnService) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connMap[s.session] = s
}

func (c *GateEntry2) deleteConn(session uint64) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	s, ok := c.connMap[session]
	if !ok {
		return
	}
	delete(c.connMap, session)
	s.Close()
}
