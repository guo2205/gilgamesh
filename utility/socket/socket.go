// packet socket
package socket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

const (
	_READ_BUFFER_SIZE  = 4096
	_WRITE_BUFFER_SIZE = 4096
)

var (
	idPool uint64 = 0

	ErrIoTimeout                 error = errors.New("i/o timeout")
	ErrPacketBodyLengthIncorrect error = errors.New("packet body length incorrect")
	ErrPacketBodyTooLarge        error = errors.New("packet body too large")
)

type Socket struct {
	conn      net.Conn
	id        uint64
	maxPktLen int

	pktHdrBuf []byte
	readBuf   []byte
	writeBuf  *bytes.Buffer
}

func NewSocket(conn net.Conn) *Socket {
	id := atomic.AddUint64(&idPool, 1)

	socket := &Socket{
		conn:      conn,
		id:        id,
		maxPktLen: 1024 * 1024 * 1024 * 4,
		pktHdrBuf: make([]byte, 4),
		readBuf:   make([]byte, _READ_BUFFER_SIZE),
		writeBuf:  bytes.NewBuffer(make([]byte, 0, _WRITE_BUFFER_SIZE)),
	}

	return socket
}

func NewSocketWithMaxLength(conn net.Conn, maxPacketLength int) *Socket {
	id := atomic.AddUint64(&idPool, 1)

	socket := &Socket{
		conn:      conn,
		id:        id,
		maxPktLen: maxPacketLength,
		pktHdrBuf: make([]byte, 4),
		readBuf:   make([]byte, _READ_BUFFER_SIZE),
		writeBuf:  bytes.NewBuffer(make([]byte, 0, _WRITE_BUFFER_SIZE)),
	}

	return socket
}

func (c *Socket) Id() uint64 {
	return c.id
}

func (c *Socket) ReadPacketOvertime(t time.Duration) ([]byte, error) {
	d, err := c.readPacket(time.Now().Add(t))
	if err != nil {
		return nil, convertIoTimeout(err)
	}
	return d, nil
}

func (c *Socket) ReadPacket() ([]byte, error) {
	d, err := c.readPacket(time.Time{})
	if err != nil {
		return nil, convertIoTimeout(err)
	}
	return d, nil
}

func (c *Socket) WritePacketOvertime(d []byte, t time.Duration) error {
	return convertIoTimeout(c.writePacket(d, time.Now().Add(t)))
}

func (c *Socket) WritePacket(d []byte) error {
	return convertIoTimeout(c.writePacket(d, time.Time{}))
}

func (c *Socket) Close() error {
	return c.conn.Close()
}

func (c *Socket) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Socket) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Socket) readPacket(t time.Time) ([]byte, error) {
	err := c.conn.SetReadDeadline(t)
	if err != nil {
		return nil, err
	}

	_, err = io.ReadFull(c.conn, c.pktHdrBuf)
	if err != nil {
		return nil, err
	}

	pktBodyLen := uint32(0)
	err = binary.Read(bytes.NewReader(c.pktHdrBuf), binary.LittleEndian, &pktBodyLen)
	if err != nil {
		return nil, err
	}

	if int(pktBodyLen) > c.maxPktLen {
		return nil, ErrPacketBodyTooLarge
	}

	if int(pktBodyLen) > len(c.readBuf) {
		c.readBuf = make([]byte, int(pktBodyLen))
	}

	_, err = io.ReadFull(c.conn, c.readBuf[:int(pktBodyLen)])
	if err != nil {
		return nil, err
	}

	cp := make([]byte, int(pktBodyLen))
	copy(cp, c.readBuf[:int(pktBodyLen)])

	return cp, nil
}

func (c *Socket) writePacket(d []byte, t time.Time) error {
	err := c.conn.SetWriteDeadline(t)
	if err != nil {
		return err
	}

	c.writeBuf.Reset()

	binary.Write(c.writeBuf, binary.LittleEndian, uint32(len(d)))
	c.writeBuf.Write(d)

	_, err = c.conn.Write(c.writeBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func convertIoTimeout(err error) error {
	if err != nil {
		if !strings.Contains(err.Error(), "i/o timeout") {
			return err
		} else {
			return ErrIoTimeout
		}
	}
	return nil
}
