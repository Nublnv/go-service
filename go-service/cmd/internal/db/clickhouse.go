package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

var ErrPoolClosed = fmt.Errorf("Error. Pool is closed")

type Pool struct {
	opts clickhouse.Options

	maxOpen       int
	maxIdle       int
	idleTTL       time.Duration
	pingOnAcquire bool

	idleCh chan *pooledConn

	openCount int32
	closed    int32

	mu sync.Mutex
}

type pooledConn struct {
	conn     clickhouse.Conn
	created  time.Time
	lastUsed time.Time
	bad      bool
}

type Conn struct {
	conn *pooledConn
	p    *Pool
}

func GetClickConnPull(ctx context.Context) *Pool {

	dbHost := os.Getenv("CLICKHOUSE_HOST")
	if dbHost == "" {
		panic("CLICKHOUSE_HOST env variable is not set")
	}
	dbPort := os.Getenv("CLICKHOUSE_PORT")
	if dbPort == "" {
		dbPort = "9000"
	}
	dbUser := os.Getenv("CLICKHOUSE_USER")
	if dbUser == "" {
		panic("CLICKHOUSE_USER env variable is not set")
	}
	dbPass := os.Getenv("CLICKHOUSE_PASSWORD")
	if dbPass == "" {
		panic("CLICKHOUSE_PASSWORD env variable is not set")
	}
	dbName := os.Getenv("CLICKHOUSE_DB")
	if dbName == "" {
		panic("CLICKHOUSE_DB env variable is not set")
	}
	return new_(clickhouse.Options{
		Addr: []string{fmt.Sprint(dbHost, ":", dbPort)},
		Auth: clickhouse.Auth{
			Username: dbUser,
			Database: dbName,
			Password: dbPass,
		},
	},
		0,
		0,
		30*time.Minute,
		false,
	)
}

// If maxOpen or maxIdle <= 0 it will set as 10
func new_(opts clickhouse.Options, maxOpen, maxIdle int, idleTTL time.Duration, pingOnAcquire bool) *Pool {
	if maxIdle <= 0 {
		maxIdle = 10
	}
	if maxOpen <= 0 {
		maxOpen = 10
	}
	if maxIdle > maxOpen {
		maxIdle = maxOpen
	}
	return &Pool{
		opts:          opts,
		maxOpen:       maxOpen,
		maxIdle:       maxIdle,
		idleTTL:       idleTTL,
		pingOnAcquire: pingOnAcquire,
		closed:        0,
		idleCh:        make(chan *pooledConn, maxIdle),
	}
}

func (p *Pool) Close() error {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return nil
	}
	close(p.idleCh)
	for pc := range p.idleCh {
		pc.conn.Close()
		atomic.AddInt32(&p.openCount, -1)
	}
	return nil
}

func (p *Pool) Acquire(ctx context.Context) (*Conn, error) {
	if atomic.LoadInt32(&p.closed) == 1 {
		return nil, ErrPoolClosed
	}

	select {
	case pc := <-p.idleCh:
		return prepareConn(ctx, p, pc)
	default:

	}

	for {
		cur := atomic.LoadInt32(&p.openCount)
		if int(cur) >= p.maxOpen {
			break
		}
		if atomic.CompareAndSwapInt32(&p.openCount, cur, cur+1) {
			conn, err := clickhouse.Open(&p.opts)
			if err != nil {
				atomic.AddInt32(&p.openCount, -1)
				return nil, err
			}
			pc := &pooledConn{
				conn:     conn,
				created:  time.Now(),
				lastUsed: time.Now(),
			}
			if p.pingOnAcquire {
				if err := conn.Ping(ctx); err != nil {
					_ = conn.Close()
					atomic.AddInt32(&p.openCount, -1)
					return nil, err
				}
			}
			return &Conn{
				p:    p,
				conn: pc,
			}, nil
		}
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case pc := <-p.idleCh:
		return prepareConn(ctx, p, pc)
	}
}

func prepareConn(ctx context.Context, p *Pool, pc *pooledConn) (*Conn, error) {
	if pc == nil {
		return nil, ErrPoolClosed
	}
	if p.shouldDrop(pc) {
		_ = pc.conn.Close()
		atomic.AddInt32(&p.openCount, -1)
		return p.Acquire(ctx)
	}
	pc.lastUsed = time.Now()
	if p.pingOnAcquire {
		err := pc.conn.Ping(ctx)
		if err != nil {
			_ = pc.conn.Close()
			atomic.AddInt32(&p.openCount, -1)
			return p.Acquire(ctx)
		}
	}
	return &Conn{
		p:    p,
		conn: pc,
	}, nil
}

func (p *Pool) shouldDrop(pc *pooledConn) bool {
	if pc.bad {
		return true
	}
	if p.idleTTL > 0 && time.Since(pc.lastUsed) > p.idleTTL {
		return true
	}
	return false
}

func (c *Conn) Realese() {
	p := c.p
	if p == nil || c.conn == nil {
		return
	}
	pc := c.conn
	c.p = nil
	c.conn = nil

	if atomic.LoadInt32(&p.closed) == 1 || p.shouldDrop(pc) {
		_ = pc.conn.Close()
		atomic.AddInt32(&p.openCount, -1)
		return
	}
	select {
	case p.idleCh <- pc:
	default:
		_ = pc.conn.Close()
		atomic.AddInt32(&p.openCount, -1)
	}
}

func (c *Conn) Conn() clickhouse.Conn {
	return c.conn.conn
}
