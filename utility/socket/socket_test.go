//packet socket test
package socket

import (
	"net"
	"testing"
	"time"
)

func TestSocket(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "localhost:60001")
	if err != nil {
		t.Fatalf("start listen failed : %v\n", err)
	}

	singal := make(chan struct{})

	go func() {
		defer l.Close()
		defer close(singal)

		conn, err := l.Accept()
		if err != nil {
			t.Fatalf("accept failed : %v\n", err)
		}

		c := NewSocket(conn)
		defer c.Close()

		err = c.WritePacket([]byte("hello,world!"))
		if err != nil {
			t.Fatalf("write failed : %v\n", err)
		}

		d, err := c.ReadPacket()
		if err != nil {
			t.Fatalf("read failed : %v\n", err)
		}

		if string(d) != "nihao,shijie!" {
			t.Fatalf("read failed : data incorrect\n")
		}
	}()

	conn, err := net.Dial("tcp", "localhost:60001")
	if err != nil {
		t.Fatalf("dial failed : %v\n", err)
	}

	c := NewSocket(conn)
	defer c.Close()

	d, err := c.ReadPacket()
	if err != nil {
		t.Fatalf("read failed : %v\n", err)
	}

	if string(d) != "hello,world!" {
		t.Fatalf("read failed : data incorrect\n")
	}

	err = c.WritePacket([]byte("nihao,shijie!"))
	if err != nil {
		t.Fatalf("write failed : %v\n", err)
	}

	<-singal
}

func TestSocketOvertime(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "localhost:60002")
	if err != nil {
		t.Fatalf("start listen failed : %v\n", err)
	}

	singal := make(chan struct{})

	go func() {
		defer l.Close()
		defer close(singal)

		conn, err := l.Accept()
		if err != nil {
			t.Fatalf("accept failed : %v\n", err)
		}

		c := NewSocket(conn)
		defer c.Close()

		time.Sleep(time.Second * 1)

		err = c.WritePacket([]byte("hello,world!"))
		if err != nil {
			t.Fatalf("write failed : %v\n", err)
		}

		d, err := c.ReadPacketOvertime(time.Second * 3)
		if err != nil {
			t.Fatalf("read failed : %v\n", err)
		}

		if string(d) != "nihao,shijie!" {
			t.Fatalf("read failed : data incorrect\n")
		}
	}()

	conn, err := net.Dial("tcp", "localhost:60002")
	if err != nil {
		t.Fatalf("dial failed : %v\n", err)
	}

	c := NewSocket(conn)
	defer c.Close()

	d, err := c.ReadPacketOvertime(time.Second * 3)
	if err != nil {
		t.Fatalf("read failed : %v\n", err)
	}

	if string(d) != "hello,world!" {
		t.Fatalf("read failed : data incorrect\n")
	}

	time.Sleep(time.Second * 2)

	err = c.WritePacket([]byte("nihao,shijie!"))
	if err != nil {
		t.Fatalf("write failed : %v\n", err)
	}

	<-singal
}

func TestSocketOvertimeFailed(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "localhost:60003")
	if err != nil {
		t.Fatalf("start listen failed : %v\n", err)
	}

	singal := make(chan struct{})

	go func() {
		defer l.Close()
		defer close(singal)

		conn, err := l.Accept()
		if err != nil {
			t.Fatalf("accept failed : %v\n", err)
		}

		c := NewSocket(conn)
		defer c.Close()

		time.Sleep(time.Second * 3)

		_, err = c.ReadPacketOvertime(time.Second * 1)
		if err == nil || err != ErrIoTimeout {
			t.Fatalf("overtime not work\n")
		}
	}()

	conn, err := net.Dial("tcp", "localhost:60003")
	if err != nil {
		t.Fatalf("dial failed : %v\n", err)
	}

	c := NewSocket(conn)
	defer c.Close()

	_, err = c.ReadPacketOvertime(time.Second * 2)
	if err == nil || err != ErrIoTimeout {
		t.Fatalf("overtime not work\n")
	}

	time.Sleep(time.Second * 2)

	err = c.WritePacket([]byte("nihao,shijie!"))
	if err != nil {
		t.Fatalf("write failed : %v\n", err)
	}

	<-singal
}
