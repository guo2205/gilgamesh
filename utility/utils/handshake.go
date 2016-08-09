// handshake
package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"net"
	"time"
)

var (
	ErrHandshakeChecksumFailed error = errors.New("handshake failed : checksum incorrect")
)

func Handshake(rawConn net.Conn, cookie uint64, t time.Duration) error {
	if t != 0 {
		err := rawConn.SetDeadline(time.Now().Add(t))
		if err != nil {
			return err
		}
	}

	mySeed, err := generateSeed(rawConn)
	if err != nil {
		return err
	}
	err = handshake1(rawConn, cookie)
	if err != nil {
		return err
	}
	err = handshake2(rawConn, cookie, mySeed)
	if err != nil {
		return err
	}
	return nil
}

func generateSeed(rawConn net.Conn) (uint32, error) {
	seed := rand.New(rand.NewSource(time.Now().Unix())).Uint32()
	err := binary.Write(rawConn, binary.LittleEndian, seed)
	if err != nil {
		return 0, err
	}
	return seed, nil
}

func handshake1(rawConn net.Conn, cookie uint64) error {
	var seed uint32

	err := binary.Read(rawConn, binary.LittleEndian, &seed)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(b, binary.LittleEndian, cookie+uint64(seed))
	h := md5.New()
	h.Write(b.Bytes())
	_, err = rawConn.Write(h.Sum(nil))
	if err != nil {
		return err
	}

	return nil
}

func handshake2(rawConn net.Conn, cookie uint64, mySeed uint32) error {
	checksum := make([]byte, 16)

	_, err := io.ReadFull(rawConn, checksum)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(b, binary.LittleEndian, cookie+uint64(mySeed))
	h := md5.New()
	h.Write(b.Bytes())
	myChecksum := h.Sum(nil)

	if bytes.Compare(myChecksum, checksum) != 0 {
		return ErrHandshakeChecksumFailed
	}

	return nil
}
