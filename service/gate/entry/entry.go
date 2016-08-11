// entry
package entry

import (
	"errors"
	"fractal/fractal"
	"gilgamesh/utility/config"
	"gilgamesh/utility/socket"
	"gilgamesh/utility/utils"
	"net"
	"sync"
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

type GateEntry struct {
	f          *fractal.Fractal
	gateOption *config.GateOption

	connLock sync.Mutex
	connMap  map[uint64]*_ConnService
}

func NewGateEntry(f *fractal.Fractal, option *config.GateOption) *GateEntry {
	return &GateEntry{
		f:          f,
		gateOption: option,
		connMap:    make(map[uint64]*_ConnService, 500),
	}
}

func (c *GateEntry) Run() error {
	l, err := c.listen()
	if err != nil {
		return err
	}
	go c.listener(l)

	return nil
}

func (c *GateEntry) WritePacket(session uint64, d []byte) error {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	s, ok := c.connMap[session]
	if !ok {
		return ErrConnNotExist
	}

	s.WritePacket(d)

	return nil
}

func (c *GateEntry) CloseConn(session uint64) {
	c.deleteConn(session)
}

func (c *GateEntry) listen() (*net.TCPListener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", c.gateOption.LocalAddr)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (c *GateEntry) listener(l *net.TCPListener) {
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

func (c *GateEntry) serviceConn(cc net.Conn) error {
	//	var err error
	//	cc, err = ssocket.NewSSocket(cc)
	//	if err != nil {
	//		cc.Close()
	//		return err
	//	}

	err := utils.Handshake(cc, c.gateOption.Cookie, c.gateOption.Timeout)
	if err != nil {
		cc.Close()
		return err
	}

	s := newConnService(c, c.f.GenerateSession(), socket.NewSocket(cc))

	c.insertConn(s)
	defer c.deleteConn(s.session)

	return s.RunService()
}

func (c *GateEntry) insertConn(s *_ConnService) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connMap[s.session] = s
}

func (c *GateEntry) deleteConn(session uint64) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	s, ok := c.connMap[session]
	if !ok {
		return
	}
	delete(c.connMap, session)
	s.Close()
}
