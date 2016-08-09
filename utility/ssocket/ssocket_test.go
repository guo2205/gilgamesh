// ssocket_test
package ssocket

import (
	"fractal/utility/socket"
	"net"
	"testing"
)

func TestSSocket(t *testing.T) {
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

		sconn, err := NewSSocket(conn)
		if err != nil {
			t.Fatalf("ssocket new failed : %v\n", err)
		}

		c := socket.NewSocket(sconn, 4096)
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

	sconn, err := NewSSocket(conn)
	if err != nil {
		t.Fatalf("ssocket new failed : %v\n", err)
	}

	c := socket.NewSocket(sconn, 4096)
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
