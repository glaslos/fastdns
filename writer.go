package fastdns

import (
	"net"
	"sync"
)

// A ResponseWriter interface is used by an DNS handler to construct an DNS response.
type ResponseWriter interface {
	// LocalAddr returns the net.Addr of the server
	LocalAddr() net.Addr

	// RemoteAddr returns the net.Addr of the client that sent the current request.
	RemoteAddr() net.Addr

	// Write writes a raw buffer back to the client.
	Write([]byte) (int, error)
}

type memResponseWriter struct {
	data  []byte
	raddr net.Addr
	laddr net.Addr
}

func (rw *memResponseWriter) RemoteAddr() net.Addr {
	return rw.raddr
}

func (rw *memResponseWriter) LocalAddr() net.Addr {
	return rw.laddr
}

func (rw *memResponseWriter) Write(p []byte) (n int, err error) {
	rw.data = append(rw.data, p...)
	n = len(p)
	return
}

type udpResponseWriter struct {
	rbuf []byte
	conn *net.UDPConn
	addr *net.UDPAddr
}

func (rw *udpResponseWriter) RemoteAddr() net.Addr {
	return rw.addr
}

func (rw *udpResponseWriter) LocalAddr() net.Addr {
	return rw.conn.LocalAddr()
}

func (rw *udpResponseWriter) Write(p []byte) (n int, err error) {
	n, _, err = rw.conn.WriteMsgUDP(p, nil, rw.addr)
	return
}

var udpResponseWriterPool = sync.Pool{
	New: func() interface{} {
		return &udpResponseWriter{
			rbuf: make([]byte, 0, 1024),
		}
	},
}
