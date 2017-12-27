// Package tcp provides a simple multiplexer over TCP.
// implemented by influxdb: https://github.com/influxdata/influxdb/tree/master/tcp
package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	// DefaultTimeout is the default length of time to wait for first byte.
	DefaultTimeout = 30 * time.Second
)

// Mux multiplexes a network connection.
type Mux struct {
	ln net.Listener
	m  map[byte]*listener

	defaultListener *listener

	wg sync.WaitGroup

	// The amount of time to wait for the first header byte.
	Timeout time.Duration

	// Out-of-band error logger
	logger *zap.SugaredLogger
}

type replayConn struct {
	net.Conn
	firstByte     byte
	readFirstbyte bool
}

func (rc *replayConn) Read(b []byte) (int, error) {
	if rc.readFirstbyte {
		return rc.Conn.Read(b)
	}

	if len(b) == 0 {
		return 0, nil
	}

	b[0] = rc.firstByte
	rc.readFirstbyte = true
	return 1, nil
}

// NewMux returns a new instance of Mux.
func NewMux(ln net.Listener) *Mux {
	return &Mux{
		ln:      ln,
		m:       make(map[byte]*listener),
		Timeout: DefaultTimeout,
		logger:  zap.NewNop().Sugar(),
	}
}

// Start handles connections from ln and multiplexes then across registered listeners.
func (mux *Mux) Start() error {
	for {
		// Wait for the next connection.
		// If it returns a temporary error then simply retry.
		// If it returns any other error then exit immediately.
		conn, err := mux.ln.Accept()
		if err, ok := err.(interface {
			Temporary() bool
		}); ok && err.Temporary() {
			continue
		}
		if err != nil {
			// Wait for all connections to be demux
			mux.wg.Wait()
			for _, ln := range mux.m {
				close(ln.c)
			}

			if mux.defaultListener != nil {
				close(mux.defaultListener.c)
			}

			return err
		}

		// Demux in a goroutine to
		mux.wg.Add(1)
		go mux.handleConn(conn)
	}
}

// Close should directly close the original net listener
func (mux *Mux) Close() error {
	err := mux.ln.Close()
	if err != nil && strings.HasSuffix(err.Error(), "use of closed network connection") {
		err = nil
	}

	return err
}

// WithLogger use customed logger
func (mux *Mux) WithLogger(logger *zap.Logger) {
	if logger != nil {
		mux.logger = logger.With(zap.String("pkg", "tcpmux")).Sugar()
	}
}

func (mux *Mux) handleConn(conn net.Conn) {
	defer mux.wg.Done()
	// Set a read deadline so connections with no data don't timeout.
	if err := conn.SetReadDeadline(time.Now().Add(mux.Timeout)); err != nil {
		conn.Close()
		mux.logger.Warnf("cannot set read deadline: %s", err)
		return
	}

	// Read first byte from connection to determine handler.
	var typ [1]byte
	if _, err := io.ReadFull(conn, typ[:]); err != nil {
		conn.Close()
		mux.logger.Warnf("cannot read header byte: %s", err)
		return
	}

	// Reset read deadline and let the listener handle that.
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		conn.Close()
		mux.logger.Warnf("cannot reset set read deadline: %s", err)
		return
	}

	// Retrieve handler based on first byte.
	handler := mux.m[typ[0]]
	if handler == nil {
		if mux.defaultListener == nil {
			conn.Close()
			mux.logger.Warnf("handler not registered: %d. Connection from %s closed", typ[0], conn.RemoteAddr())
			return
		}

		conn = &replayConn{
			Conn:      conn,
			firstByte: typ[0],
		}
		handler = mux.defaultListener
	}

	// Send connection to handler.  The handler is responsible for closing the connection.
	timer := time.NewTimer(mux.Timeout)
	defer timer.Stop()

	select {
	case handler.c <- conn:

	case <-timer.C:
		conn.Close()
		mux.logger.Warnf("handler not ready: %d. Connection from %s closed", typ[0], conn.RemoteAddr())
		return
	}
}

// Listen returns a listener identified by header.
// Any connection accepted by mux is multiplexed based on the initial header byte.
func (mux *Mux) Listen(header byte) net.Listener {
	// Ensure two listeners are not created for the same header byte.
	if _, ok := mux.m[header]; ok {
		panic(fmt.Sprintf("listener already registered under header byte: %d", header))
	}

	// Create a new listener and assign it.
	ln := &listener{
		c:   make(chan net.Conn),
		mux: mux,
	}
	mux.m[header] = ln

	return ln
}

// DefaultListener will return a net.Listener that will pass-through any
// connections with non-registered values for the first byte of the connection.
// The connections returned from this listener's Accept() method will replay the
// first byte of the connection as a short first Read().
//
// This can be used to pass to an HTTP server, so long as there are no conflicts
// with registered listener bytes and the first character of the HTTP request:
// 71 ('G') for GET, etc.
func (mux *Mux) DefaultListener() net.Listener {
	if mux.defaultListener == nil {
		mux.defaultListener = &listener{
			c:   make(chan net.Conn),
			mux: mux,
		}
	}

	return mux.defaultListener
}

// listener is a receiver for connections received by Mux.
type listener struct {
	c   chan net.Conn
	mux *Mux
}

// Accept waits for and returns the next connection to the listener.
func (ln *listener) Accept() (c net.Conn, err error) {
	conn, ok := <-ln.c
	if !ok {
		return nil, errors.New("network connection closed")
	}
	return conn, nil
}

// Close is a no-op. The mux's listener should be closed instead.
func (ln *listener) Close() error { return nil }

// Addr returns the Addr of the listener
func (ln *listener) Addr() net.Addr {
	if ln.mux == nil || ln.mux.ln == nil {
		return nil
	}

	return ln.mux.ln.Addr()
}

// Dial connects to a remote mux listener with a given header byte.
func Dial(network, address string, header byte) (net.Conn, error) {
	return DialWithTimeout(network, address, 0, header)
}

// DialWithTimeout connects to a remote mux listener with a given header byte and timeout
func DialWithTimeout(network, address string, timeout time.Duration, header byte) (net.Conn, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}

	if _, err := conn.Write([]byte{header}); err != nil {
		return nil, fmt.Errorf("write mux header: %s", err)
	}

	return conn, nil
}
