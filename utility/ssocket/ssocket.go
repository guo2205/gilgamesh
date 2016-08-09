// ssocket
package ssocket

import (
	"crypto/rc4"
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"time"
)

const (
	_SALT = "ab75kjg&&54#$"
)

type SSocket struct {
	conn    net.Conn
	encoder *rc4.Cipher
	decoder *rc4.Cipher
}

func NewSSocket(conn net.Conn) (*SSocket, error) {
	X1, E1 := _DHExchange()
	X2, E2 := _DHExchange()

	err := binary.Write(conn, binary.LittleEndian, E1.Int64())
	if err != nil {
		return nil, err
	}

	err = binary.Write(conn, binary.LittleEndian, E2.Int64())
	if err != nil {
		return nil, err
	}

	var (
		S1 int64
		S2 int64
	)

	err = binary.Read(conn, binary.LittleEndian, &S1)
	if err != nil {
		return nil, err
	}
	err = binary.Read(conn, binary.LittleEndian, &S2)
	if err != nil {
		return nil, err
	}

	key1 := _DHKey(X1, big.NewInt(S1))
	key2 := _DHKey(X2, big.NewInt(S2))

	encoder, err := rc4.NewCipher([]byte(fmt.Sprintf("%v%v", _SALT, key2)))
	if err != nil {
		return nil, err
	}
	decoder, err := rc4.NewCipher([]byte(fmt.Sprintf("%v%v", _SALT, key1)))
	if err != nil {
		return nil, err
	}

	return &SSocket{
		conn:    conn,
		encoder: encoder,
		decoder: decoder,
	}, nil
}

func (c *SSocket) Read(b []byte) (n int, err error) {
	rlen, err := c.conn.Read(b)
	if err != nil {
		return -1, err
	}

	d := b[:rlen]
	c.decoder.XORKeyStream(d, d)

	return rlen, nil
}

func (c *SSocket) Write(b []byte) (n int, err error) {
	d := make([]byte, len(b))
	copy(d, b)

	c.decoder.XORKeyStream(d, d)

	return c.conn.Write(d)
}

func (c *SSocket) Close() error {
	return c.conn.Close()
}

func (c *SSocket) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *SSocket) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *SSocket) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *SSocket) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *SSocket) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
