package rpc

import (
	"net"
	"sync"
	"time"

	"github.com/dtynn/winston/pkg/tcp"
	"google.golang.org/grpc"
)

var (
	defaultConnMgr = NewConnectionMgr()
)

// NewConn return new grpc connection
func NewConn(header byte, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialer := func(address string, timeout time.Duration) (net.Conn, error) {
		return tcp.DialWithTimeout("tcp", address, timeout, header)
	}

	opts = append([]grpc.DialOption{grpc.WithDialer(dialer), grpc.WithInsecure()}, opts...)

	cc, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, err
	}

	return cc, nil

}

// NewConnectionMgr return a new connection mgr
func NewConnectionMgr() *ConnectionMgr {
	return &ConnectionMgr{
		conns: map[byte]map[string]*grpc.ClientConn{},
	}
}

// ConnectionMgr rpc connection manager
type ConnectionMgr struct {
	mu    sync.Mutex
	conns map[byte]map[string]*grpc.ClientConn
}

// Get get a connection
func (c *ConnectionMgr) Get(header byte, target string) (*grpc.ClientConn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.conns[header]; !ok {
		c.conns[header] = map[string]*grpc.ClientConn{}
	}

	cc, ok := c.conns[header][target]
	if ok {
		return cc, nil
	}

	cc, err := NewConn(header, target)
	if err != nil {
		return nil, err
	}

	c.conns[header][target] = cc

	return cc, nil
}
