//async packet socket test
package socket

import (
	"net"
	"testing"
	"time"
)

func TestAsyncSocket(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "localhost:60004")
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

		c := NewAsyncSocket(NewSocket(conn))
		defer c.Close()

		wr := <-c.WritePacket([]byte("hello,world!"))
		if wr.Err != nil {
			t.Fatalf("write failed : %v\n", wr.Err)
		}

		rr := <-c.ReadPacket()
		if rr.Err != nil {
			t.Fatalf("read failed : %v\n", rr.Err)
		}

		if string(rr.Data) != "nihao,shijie!" {
			t.Fatalf("read failed : data incorrect\n")
		}
	}()

	conn, err := net.Dial("tcp", "localhost:60004")
	if err != nil {
		t.Fatalf("dial failed : %v\n", err)
	}

	c := NewAsyncSocket(NewSocket(conn))
	defer c.Close()

	rr := <-c.ReadPacket()
	if rr.Err != nil {
		t.Fatalf("read failed : %v\n", rr.Err)
	}

	if string(rr.Data) != "hello,world!" {
		t.Fatalf("read failed : data incorrect\n")
	}

	wr := <-c.WritePacket([]byte("nihao,shijie!"))
	if wr.Err != nil {
		t.Fatalf("write failed : %v\n", wr.Err)
	}

	<-singal
}

func TestAsyncSocketOvertime(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "localhost:60005")
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

		c := NewAsyncSocket(NewSocket(conn))
		defer c.Close()

		time.Sleep(time.Second * 1)

		wr := <-c.WritePacket([]byte("hello,world!"))
		if wr.Err != nil {
			t.Fatalf("write failed : %v\n", wr.Err)
		}

		rr := <-c.ReadPacketOvertime(time.Second * 3)
		if rr.Err != nil {
			t.Fatalf("read failed : %v\n", rr.Err)
		}

		if string(rr.Data) != "nihao,shijie!" {
			t.Fatalf("read failed : data incorrect\n")
		}
	}()

	conn, err := net.Dial("tcp", "localhost:60005")
	if err != nil {
		t.Fatalf("dial failed : %v\n", err)
	}

	c := NewAsyncSocket(NewSocket(conn))
	defer c.Close()

	rr := <-c.ReadPacketOvertime(time.Second * 3)
	if rr.Err != nil {
		t.Fatalf("read failed : %v\n", rr.Err)
	}

	if string(rr.Data) != "hello,world!" {
		t.Fatalf("read failed : data incorrect\n")
	}

	time.Sleep(time.Second * 2)

	wr := <-c.WritePacket([]byte("nihao,shijie!"))
	if wr.Err != nil {
		t.Fatalf("write failed : %v\n", wr.Err)
	}

	<-singal
}

func TestAsyncSocketOvertimeFailed(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "localhost:60006")
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

		c := NewAsyncSocket(NewSocket(conn))
		defer c.Close()

		time.Sleep(time.Second * 3)

		rr := <-c.ReadPacketOvertime(time.Second * 1)
		if rr.Err == nil {
			t.Fatalf("overtime not work\n")
		}
	}()

	conn, err := net.Dial("tcp", "localhost:60006")
	if err != nil {
		t.Fatalf("dial failed : %v\n", err)
	}

	c := NewAsyncSocket(NewSocket(conn))
	defer c.Close()

	rr := <-c.ReadPacketOvertime(time.Second * 2)
	if rr.Err == nil {
		t.Fatalf("overtime not work\n")
	}

	time.Sleep(time.Second * 2)

	wr := <-c.WritePacket([]byte("nihao,shijie!"))
	if wr.Err != nil {
		t.Fatalf("write failed : %v\n", wr.Err)
	}

	<-singal
}

func TestAsyncSocketClose(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "localhost:60007")
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

		c := NewAsyncSocket(NewSocket(conn))
		defer c.Close()

		<-c.ReadPacketOvertime(time.Second * 60)
	}()

	conn, err := net.Dial("tcp", "localhost:60007")
	if err != nil {
		t.Fatalf("dial failed : %v\n", err)
	}

	c := NewAsyncSocket(NewSocket(conn))
	defer c.Close()

	go func() {
		time.Sleep(time.Second * 3)
		conn.Close()
	}()

	<-c.ReadPacketOvertime(time.Second * 30)

	<-singal
}
