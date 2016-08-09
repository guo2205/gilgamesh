// async_socket
package socket

import (
	"errors"
	"net"
	"sync/atomic"
	"time"
)

type AsyncSocket struct {
	conn *Socket

	closed uint32

	shutdownSingal chan *struct{}
	readerSignal   chan *struct{}
	writerSignal   chan *struct{}
	readerCh       chan _ReadEvent
	writerCh       chan _WriterEvent
}

var (
	ErrChannelUnexpectedClosed error = errors.New("channel unexpected closed")
)

type _ReadEvent struct {
	T      time.Duration
	Result chan _ReadResult
}

type _WriterEvent struct {
	Data   []byte
	T      time.Duration
	Result chan _WriterResult
}

type _ReadResult struct {
	Data []byte
	Err  error
}

type _WriterResult struct {
	Err error
}

func NewAsyncSocket(conn *Socket) *AsyncSocket {
	c := &AsyncSocket{
		conn:           conn,
		closed:         0,
		shutdownSingal: make(chan *struct{}),
		readerSignal:   make(chan *struct{}, 2),
		writerSignal:   make(chan *struct{}, 2),
		readerCh:       make(chan _ReadEvent, 64),
		writerCh:       make(chan _WriterEvent, 64),
	}
	go c.reader()
	go c.writer()
	return c
}

func (c *AsyncSocket) ReadPacketOvertime(t time.Duration) chan _ReadResult {
	r := make(chan _ReadResult)
	c.readerCh <- _ReadEvent{
		T:      t,
		Result: r,
	}
	return r
}

func (c *AsyncSocket) ReadPacket() chan _ReadResult {
	return c.ReadPacketOvertime(0)
}

func (c *AsyncSocket) WritePacketOvertime(d []byte, t time.Duration) chan _WriterResult {
	r := make(chan _WriterResult)
	c.writerCh <- _WriterEvent{
		Data:   d,
		T:      t,
		Result: r,
	}
	return r
}

func (c *AsyncSocket) WritePacket(d []byte) chan _WriterResult {
	return c.WritePacketOvertime(d, 0)
}

func (c *AsyncSocket) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *AsyncSocket) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *AsyncSocket) Close() {
	if atomic.SwapUint32(&c.closed, 1) != 0 {
		return
	}

	c.readerSignal <- nil
	c.writerSignal <- nil
	for i := 0; i < 2; i++ {
		<-c.shutdownSingal
	}
	c.conn.Close()
	close(c.readerSignal)
	close(c.writerSignal)
	close(c.shutdownSingal)
	close(c.readerCh)
	close(c.writerCh)
}

func (c *AsyncSocket) reader() {
	defer func() {
		recover()
		c.writerSignal <- nil
		c.shutdownSingal <- nil
	}()

	for {
		select {
		case ev, open := <-c.readerCh:
			if !open {
				ev.Result <- _ReadResult{
					Err: ErrChannelUnexpectedClosed,
				}
				close(ev.Result)
				return
			}
			select {
			case r := <-c.readWorker(ev.T):
				ev.Result <- r
				close(ev.Result)
			case <-c.readerSignal:
				return
			}
		case <-c.readerSignal:
			return
		}
	}
}

func (c *AsyncSocket) writer() {
	defer func() {
		recover()
		c.readerSignal <- nil
		c.shutdownSingal <- nil
	}()

	for {
		select {
		case ev, open := <-c.writerCh:
			if !open {
				ev.Result <- _WriterResult{
					Err: ErrChannelUnexpectedClosed,
				}
				close(ev.Result)
				return
			}
			select {
			case r := <-c.writeWorker(ev.Data, ev.T):
				ev.Result <- r
				close(ev.Result)
			case <-c.writerSignal:
				return
			}
		case <-c.writerSignal:
			return
		}
	}
}

func (c *AsyncSocket) readWorker(t time.Duration) chan _ReadResult {
	ch := make(chan _ReadResult)

	go func() {
		defer func() {
			recover()
			close(ch)
		}()

		var (
			d   []byte
			err error
		)

		if t == 0 {
			d, err = c.conn.ReadPacket()
		} else {
			d, err = c.conn.ReadPacketOvertime(t)
		}

		ch <- _ReadResult{
			Data: d,
			Err:  err,
		}
	}()

	return ch
}

func (c *AsyncSocket) writeWorker(d []byte, t time.Duration) chan _WriterResult {
	ch := make(chan _WriterResult)

	go func() {
		defer func() {
			recover()
			close(ch)
		}()

		var (
			err error
		)

		if t == 0 {
			err = c.conn.WritePacket(d)
		} else {
			err = c.conn.WritePacketOvertime(d, t)
		}

		ch <- _WriterResult{
			Err: err,
		}
	}()

	return ch
}
