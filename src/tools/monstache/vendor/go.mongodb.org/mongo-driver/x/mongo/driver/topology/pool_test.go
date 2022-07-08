package topology

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/internal/testutil/assert"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

func TestPool(t *testing.T) {
	t.Run("newPool", func(t *testing.T) {
		t.Run("should be connected", func(t *testing.T) {
			pc := poolConfig{
				Address: address.Address(""),
			}
			p, err := newPool(pc)
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			if p.connected != connected {
				t.Errorf("Expected new pool to be connected. got %v; want %v", p.connected, connected)
			}
		})
	})
	t.Run("closeConnection", func(t *testing.T) {
		t.Run("can't put connection from different pool", func(t *testing.T) {
			pc1 := poolConfig{
				Address: address.Address(""),
			}
			p1, err := newPool(pc1)
			noerr(t, err)
			err = p1.connect()
			noerr(t, err)

			pc2 := poolConfig{
				Address: address.Address(""),
			}
			p2, err := newPool(pc2)
			noerr(t, err)
			err = p2.connect()
			noerr(t, err)

			c1 := &connection{pool: p1}
			want := ErrWrongPool
			got := p2.closeConnection(c1)
			if got != want {
				t.Errorf("Errors do not match. got %v; want %v", got, want)
			}
		})
	})
	t.Run("disconnect", func(t *testing.T) {
		t.Run("cannot close twice", func(t *testing.T) {
			pc := poolConfig{
				Address: address.Address(""),
			}
			p, err := newPool(pc)
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			err = p.disconnect(context.Background())
			noerr(t, err)
			err = p.disconnect(context.Background())
			if err != ErrPoolDisconnected {
				t.Errorf("Should not be able to call disconnect twice. got %v; want %v", err, ErrPoolDisconnected)
			}
		})
		t.Run("closes idle connections", func(t *testing.T) {
			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address:     address.Address(addr.String()),
				MaxIdleTime: 100 * time.Millisecond,
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)

			err = p.connect()
			noerr(t, err)
			conns := [3]*connection{}
			for idx := range [3]struct{}{} {
				conns[idx], err = p.get(context.Background())
				noerr(t, err)
			}
			for idx := range [3]struct{}{} {
				err = p.put(conns[idx])
				noerr(t, err)
			}
			if d.lenopened() != 3 {
				t.Errorf("Should have opened 3 connections, but didn't. got %d; want %d", d.lenopened(), 3)
			}
			if p.conns.totalSize != 3 {
				t.Errorf("Pool should have 3 total connections. got %d; want %d", p.conns.totalSize, 3)
			}
			err = p.disconnect(context.Background())
			time.Sleep(time.Second)

			noerr(t, err)
			if d.lenclosed() != 3 {
				t.Errorf("Should have closed 3 connections, but didn't. got %d; want %d", d.lenclosed(), 3)
			}
			if p.conns.totalSize != 0 {
				t.Errorf("Pool should have 0 total connections. got %d; want %d", p.conns.totalSize, 0)
			}
			close(cleanup)
		})
		t.Run("closes all connections currently in pool and closes all remaining connections", func(t *testing.T) {
			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			conns := [3]*connection{}
			for idx := range [3]struct{}{} {
				conns[idx], err = p.get(context.Background())
				noerr(t, err)
			}
			for idx := range [2]struct{}{} {
				err = p.put(conns[idx])
				noerr(t, err)
			}
			if d.lenopened() != 3 {
				t.Errorf("Should have opened 3 connections, but didn't. got %d; want %d", d.lenopened(), 3)
			}
			if p.conns.totalSize != 3 {
				t.Errorf("Pool should have 3 total connections. got %d; want %d", p.conns.totalSize, 3)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Microsecond)
			defer cancel()
			err = p.disconnect(ctx)
			noerr(t, err)

			assertConnectionsClosed(t, d, 3)
			assert.Nil(t, err, "error running callback: %s", err)
			if p.conns.totalSize != 0 {
				t.Errorf("Pool should have 0 total connections. got %d; want %d", p.conns.totalSize, 0)
			}
			close(cleanup)
		})
		t.Run("properly sets the connection state on return", func(t *testing.T) {
			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address:     address.Address(addr.String()),
				MinPoolSize: 0,
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			c, err := p.get(context.Background())
			noerr(t, err)
			if p.conns.totalSize != 1 {
				t.Errorf("Pool should have 1 total connection. got %d; want %d", p.conns.totalSize, 1)
			}
			err = p.closeConnection(c)
			noerr(t, err)
			err = p.put(c)
			noerr(t, err)
			if p.conns.totalSize != 0 {
				t.Errorf("Pool should have 0 total connections. got %d; want %d", p.conns.totalSize, 0)
			}
			if d.lenopened() != 1 {
				t.Errorf("Should have opened 1 connections, but didn't. got %d; want %d", d.lenopened(), 1)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Microsecond)
			defer cancel()
			err = p.disconnect(ctx)
			noerr(t, err)
			if d.lenclosed() != 1 {
				t.Errorf("Should have closed 1 connections, but didn't. got %d; want %d", d.lenclosed(), 1)
			}
			close(cleanup)
			state := atomic.LoadInt32(&p.connected)
			if state != disconnected {
				t.Errorf("Should have set the connection state on return. got %d; want %d", state, disconnected)
			}
		})
		t.Run("no race if connections are also connecting", func(t *testing.T) {
			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			getDone := make(chan struct{})
			disconnectDone := make(chan struct{})
			_, err = p.get(context.Background())
			noerr(t, err)
			getCtx, getCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer getCancel()
			go func() {
				defer close(getDone)
				for {
					select {
					case <-disconnectDone:
						return
					default:
						loopCtx, loopCancel := context.WithTimeout(getCtx, 3*time.Second)
						c, err := p.get(loopCtx)
						loopCancel()
						if err == nil {
							_ = p.put(c)
						}
						time.Sleep(time.Microsecond)
					}
				}
			}()
			go func() {
				defer close(disconnectDone)
				_, err := p.get(getCtx)
				noerr(t, err)
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Microsecond)
				defer cancel()
				err = p.disconnect(ctx)
				noerr(t, err)
			}()
			<-getDone
			close(cleanup)
		})
	})
	t.Run("connect", func(t *testing.T) {
		t.Run("can reconnect a disconnected pool", func(t *testing.T) {
			assertGenerationMapState := func(t *testing.T, p *pool, expectedState int32) {
				t.Helper()

				actualState := atomic.LoadInt32(&p.generation.state)
				assert.Equal(t, expectedState, actualState, "expected generation map state %d, got %d", expectedState, actualState)
			}

			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			assertGenerationMapState(t, p, connected)
			c, err := p.get(context.Background())
			noerr(t, err)
			gen := c.generation
			if gen != 0 {
				t.Errorf("Connection should have a newer generation. got %d; want %d", gen, 0)
			}
			err = p.put(c)
			noerr(t, err)
			if d.lenopened() != 1 {
				t.Errorf("Should have opened 1 connections, but didn't. got %d; want %d", d.lenopened(), 1)
			}
			if p.conns.totalSize != 1 {
				t.Errorf("Pool should have 1 total connection. got %d; want %d", p.conns.totalSize, 1)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			err = p.disconnect(ctx)
			noerr(t, err)
			assertGenerationMapState(t, p, disconnected)

			assertConnectionsClosed(t, d, 1)
			if p.conns.totalSize != 0 {
				t.Errorf("Pool should have 0 total connections. got %d; want %d", p.conns.totalSize, 0)
			}
			close(cleanup)
			state := atomic.LoadInt32(&p.connected)
			if state != disconnected {
				t.Errorf("Should have set the connection state on return. got %d; want %d", state, disconnected)
			}
			err = p.connect()
			noerr(t, err)
			assertGenerationMapState(t, p, connected)

			c, err = p.get(context.Background())
			noerr(t, err)
			err = p.put(c)
			noerr(t, err)
			if d.lenopened() != 2 {
				t.Errorf("Should have opened 3 connections, but didn't. got %d; want %d", d.lenopened(), 2)
			}
			if p.conns.totalSize != 1 {
				t.Errorf("Pool should have 1 total connection. got %d; want %d", p.conns.totalSize, 1)
			}
		})
		t.Run("cannot connect multiple times without disconnect", func(t *testing.T) {
			pc := poolConfig{
				Address: "",
			}
			p, err := newPool(pc)
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			err = p.connect()
			if err != ErrPoolConnected {
				t.Errorf("Shouldn't be able to connect to already connected pool. got %v; want %v", err, ErrPoolConnected)
			}
			err = p.connect()
			if err != ErrPoolConnected {
				t.Errorf("Shouldn't be able to connect to already connected pool. got %v; want %v", err, ErrPoolConnected)
			}
			err = p.disconnect(context.Background())
			noerr(t, err)
			err = p.connect()
			if err != nil {
				t.Errorf("Should be able to connect to pool after disconnect. got %v; want <nil>", err)
			}
		})
		t.Run("can disconnect and reconnect multiple times", func(t *testing.T) {
			pc := poolConfig{
				Address: address.Address(""),
			}
			p, err := newPool(pc)
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			err = p.disconnect(context.Background())
			noerr(t, err)
			err = p.connect()
			if err != nil {
				t.Errorf("Should be able to connect to disconnected pool. got %v; want <nil>", err)
			}
			err = p.disconnect(context.Background())
			noerr(t, err)
			err = p.connect()
			if err != nil {
				t.Errorf("Should be able to connect to disconnected pool. got %v; want <nil>", err)
			}
			err = p.disconnect(context.Background())
			noerr(t, err)
			err = p.connect()
			if err != nil {
				t.Errorf("Should be able to connect to pool after disconnect. got %v; want <nil>", err)
			}
		})
	})
	t.Run("get", func(t *testing.T) {
		t.Run("return context error when already cancelled", func(t *testing.T) {
			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			cancel()
			_, err = p.get(ctx)
			if err != context.Canceled {
				t.Errorf("Should return context error when already cancelled. got %v; want %v", err, context.Canceled)
			}
			if p.conns.totalSize != 0 {
				t.Errorf("Pool should have 0 total connections. got %d; want %d", p.conns.totalSize, 0)
			}
			close(cleanup)
		})
		t.Run("return error when attempting to create new connection", func(t *testing.T) {
			wanterr := errors.New("create new connection error")
			var want error = ConnectionError{Wrapped: wanterr, init: true}
			var dialer DialerFunc = func(context.Context, string, string) (net.Conn, error) { return nil, wanterr }
			pc := poolConfig{
				Address: address.Address(""),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return dialer }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			_, got := p.get(context.Background())
			if got != want {
				t.Errorf("Should return error from calling New. got %v; want %v", got, want)
			}
			if p.conns.totalSize != 0 {
				t.Errorf("Pool should have 0 total connections. got %d; want %d", p.conns.totalSize, 0)
			}
		})
		t.Run("adds connection to inflight pool", func(t *testing.T) {
			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 1, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			c, err := p.get(ctx)
			noerr(t, err)
			inflight := len(p.opened)
			if inflight != 1 {
				t.Errorf("Incorrect number of inlight connections. got %d; want %d", inflight, 1)
			}
			if p.conns.totalSize != 1 {
				t.Errorf("Pool should have 1 total connection. got %d; want %d", p.conns.totalSize, 1)
			}
			err = p.closeConnection(c)
			noerr(t, err)
			close(cleanup)
		})
		t.Run("closes stale connections", func(t *testing.T) {
			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 2, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			closedChan := make(chan struct{}, 1)
			d.closeCallBack = func() {
				closedChan <- struct{}{}
			}
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(
				pc,
				WithDialer(func(Dialer) Dialer { return d }),
				WithIdleTimeout(func(time.Duration) time.Duration { return 10 * time.Millisecond }),
			)
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			c, err := p.get(ctx)
			noerr(t, err)
			if d.lenopened() != 1 {
				t.Errorf("Should have opened 1 connection, but didn't. got %d; want %d", d.lenopened(), 1)
			}
			if p.conns.totalSize != 1 {
				t.Errorf("Pool should have 1 total connection. got %d; want %d", p.conns.totalSize, 1)
			}
			time.Sleep(15 * time.Millisecond)
			err = p.put(c)
			noerr(t, err)
			<-closedChan
			if d.lenclosed() != 1 {
				t.Errorf("Should have closed 1 connections, but didn't. got %d; want %d", d.lenclosed(), 1)
			}
			if p.conns.totalSize != 0 {
				t.Errorf("Pool should have 0 total connections. got %d; want %d", p.conns.totalSize, 0)
			}
			c, err = p.get(ctx)
			noerr(t, err)
			if d.lenopened() != 2 {
				t.Errorf("Should have opened 2 connections, but didn't. got %d; want %d", d.lenopened(), 2)
			}
			if d.lenclosed() != 1 {
				t.Errorf("Should have closed 1 connection, but didn't. got %d; want %d", d.lenclosed(), 1)
			}
			if p.conns.totalSize != 1 {
				t.Errorf("Pool should have 1 total connection. got %d; want %d", p.conns.totalSize, 1)
			}
			close(cleanup)
		})
		t.Run("recycles connections", func(t *testing.T) {
			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			for range [3]struct{}{} {
				c, err := p.get(context.Background())
				noerr(t, err)
				err = p.put(c)
				noerr(t, err)
				if d.lenopened() != 1 {
					t.Errorf("Should have opened 1 connection, but didn't. got %d; want %d", d.lenopened(), 1)
				}
			}
			close(cleanup)
		})
		t.Run("cannot get from disconnected pool", func(t *testing.T) {
			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Microsecond)
			defer cancel()
			err = p.disconnect(ctx)
			noerr(t, err)
			_, err = p.get(context.Background())
			if err != ErrPoolDisconnected {
				t.Errorf("Should get error from disconnected pool. got %v; want %v", err, ErrPoolDisconnected)
			}
			close(cleanup)
		})
		t.Run("pool closes excess connections when returned", func(t *testing.T) {
			cleanup := make(chan struct{})
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			conns := [3]*connection{}
			for idx := range [3]struct{}{} {
				conns[idx], err = p.get(context.Background())
				noerr(t, err)
			}
			err = p.disconnect(context.Background())
			noerr(t, err)
			for idx := range [3]struct{}{} {
				err = p.put(conns[idx])
				noerr(t, err)
			}
			if d.lenopened() != 3 {
				t.Errorf("Should have opened 3 connections, but didn't. got %d; want %d", d.lenopened(), 3)
			}
			if d.lenclosed() != 3 {
				t.Errorf("Should have closed 3 connections, but didn't. got %d; want %d", d.lenclosed(), 3)
			}
			close(cleanup)
		})
		t.Run("handshaker i/o fails", func(t *testing.T) {
			want := "unable to write wire message to network: Write error"

			pc := poolConfig{
				Address: address.Address(""),
			}
			p, err := newPool(pc, WithHandshaker(func(Handshaker) Handshaker {
				return operation.NewIsMaster()
			}),
				WithDialer(func(Dialer) Dialer {
					return DialerFunc(func(context.Context, string, string) (net.Conn, error) {
						return &writeFailConn{&net.TCPConn{}}, nil
					})
				}),
			)
			noerr(t, err)
			err = p.connect()
			noerr(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err = p.get(ctx)
			connErr := err.(ConnectionError)
			if !strings.Contains(connErr.Error(), want) {
				t.Errorf("Incorrect error. got %v; error should contain %v", connErr.Wrapped, want)
			}
			if p.conns.totalSize != 0 {
				t.Errorf("Pool should have 0 total connection. got %d; want %d", p.conns.totalSize, 0)
			}
		})
	})
	t.Run("Connection", func(t *testing.T) {
		t.Run("Connection Close Does Not Error After Pool Is Disconnected", func(t *testing.T) {
			cleanup := make(chan struct{})
			defer close(cleanup)
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			c, err := p.get(context.Background())
			noerr(t, err)
			c1 := &Connection{connection: c}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err = p.disconnect(ctx)
			noerr(t, err)
			err = c1.Close()
			if err != nil {
				t.Errorf("Connection close should not error after Pool is Disconnected, but got error: %v", err)
			}
		})
		t.Run("Does not return to pool twice", func(t *testing.T) {
			cleanup := make(chan struct{})
			defer close(cleanup)
			addr := bootstrapConnections(t, 1, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newdialer(&net.Dialer{})
			pc := poolConfig{
				Address: address.Address(addr.String()),
			}
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			noerr(t, err)
			err = p.connect()
			noerr(t, err)
			c, err := p.get(context.Background())
			c1 := &Connection{connection: c}
			noerr(t, err)
			if p.conns.size != 0 {
				t.Errorf("Should be no connections in pool. got %d; want %d", p.conns.size, 0)
			}
			if p.conns.totalSize != 1 {
				t.Errorf("Pool should have 1 total connection. got %d; want %d", p.conns.totalSize, 1)
			}
			err = c1.Close()
			noerr(t, err)
			err = c1.Close()
			noerr(t, err)
			if p.conns.size != 1 {
				t.Errorf("Should not return connection to pool twice. got %d; want %d", p.conns.size, 1)
			}
			if p.conns.totalSize != 1 {
				t.Errorf("Pool should have 1 total connection. got %d; want %d", p.conns.totalSize, 1)
			}
		})
		t.Run("close does not panic if expires before connected", func(t *testing.T) {
			cleanup := make(chan struct{})
			defer close(cleanup)
			addr := bootstrapConnections(t, 3, func(nc net.Conn) {
				<-cleanup
				_ = nc.Close()
			})
			d := newSleepDialer(&net.Dialer{})
			pc := poolConfig{
				Address:     address.Address(addr.String()),
				MinPoolSize: 1,
			}
			maintainInterval = time.Second
			p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
			maintainInterval = time.Minute
			noerr(t, err)
			err = p.connect()
			noerr(t, err)

			// Increment the pool's generation number so the connection will be considered stale and will be closed by
			// get().
			p.clear(nil)
			_, err = p.get(context.Background())
			noerr(t, err)
		})
	})
	t.Run("wait queue timeout error", func(t *testing.T) {
		cleanup := make(chan struct{})
		addr := bootstrapConnections(t, 1, func(nc net.Conn) {
			<-cleanup
			_ = nc.Close()
		})
		d := newdialer(&net.Dialer{})
		pc := poolConfig{
			Address:     address.Address(addr.String()),
			MaxPoolSize: 1,
		}
		p, err := newPool(pc, WithDialer(func(Dialer) Dialer { return d }))
		noerr(t, err)
		err = p.connect()
		noerr(t, err)

		// get first connection.
		_, err = p.get(context.Background())
		noerr(t, err)

		// Set a short timeout and get again.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		_, err = p.get(ctx)
		assert.NotNil(t, err, "expected a WaitQueueTimeout; got nil")

		// Assert that error received is WaitQueueTimeoutError with context deadline exceeded.
		wqtErr, ok := err.(WaitQueueTimeoutError)
		assert.True(t, ok, "expected a WaitQueueTimeoutError; got %v", err)
		assert.True(t, wqtErr.Unwrap() == context.DeadlineExceeded,
			"expected a timeout error; got %v", wqtErr)

		close(cleanup)
	})
}

type sleepDialer struct {
	Dialer
}

func newSleepDialer(d Dialer) *sleepDialer {
	return &sleepDialer{d}
}

func (d *sleepDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	time.Sleep(5 * time.Second)
	return d.Dialer.DialContext(ctx, network, address)
}

func assertConnectionsClosed(t *testing.T, dialer *dialer, expectedClosedCount int) {
	t.Helper()

	callback := func() {
		for {
			if dialer.lenclosed() == expectedClosedCount {
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	assert.Soon(t, callback, 3*time.Second)
}
