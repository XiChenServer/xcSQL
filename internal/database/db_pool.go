package database

import (
	"fmt"
	"sync"
	"time"
)

// ConnectionPool DB 代表数据库实例
type ConnectionPool struct {
	// 最大空闲连接数. 若设置为 0，则取默认值 2；若设置为负值，则取 0，代表不启用连接池
	maxIdleCount int
	// 最多可以打开的连接数. 若设为非正值，则代表不作限制
	maxOpen       int
	maxLifetime   time.Duration
	mu            sync.Mutex
	freeConn      []*driverConn
	numOpen       int
	cleanupTicker *time.Ticker // 定时清理的定时器
	stopCleanup   chan bool    // 停止清理的信号
}

type driverConn struct {
	dbName      string
	db          *XcDB
	createdAt   time.Time
	ci          Conn
	closed      bool
	inUse       bool
	returnedAt  time.Time
	finalClosed bool

	sync.Mutex
}

// Conn 代表真实的数据库连接
type Conn interface {
	Close() error
	Ping() error
}

// NewConnectionPool 创建一个新的连接池
func NewConnectionPool(maxIdleCount, maxOpen int, maxLifetime time.Duration) *ConnectionPool {
	pool := &ConnectionPool{
		maxIdleCount: maxIdleCount,
		maxOpen:      maxOpen,
		maxLifetime:  maxLifetime,
		freeConn:     make([]*driverConn, 0),
		numOpen:      0,
		stopCleanup:  make(chan bool),
	}

	// 启动定期清理任务
	go pool.startCleanup()

	return pool
}

// NewDriverConn 创建一个新的数据库连接
func NewDriverConn(name string) *driverConn {
	db, err := NewXcDB(name)
	if err != nil {
		return nil
	}
	return &driverConn{
		dbName:      name,
		db:          db,
		createdAt:   time.Now(),
		ci:          nil, // 这里需要在实际使用时初始化
		closed:      false,
		inUse:       false,
		finalClosed: false,
	}
}

// GetConnection 从连接池获取一个连接
func (pool *ConnectionPool) GetConnection(dbName string) (*driverConn, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	for _, conn := range pool.freeConn {
		if !conn.closed && conn.dbName == dbName && conn.inUse == false {
			conn.inUse = true
			return conn, nil
		}
	}

	if pool.numOpen >= pool.maxOpen {
		return nil, fmt.Errorf("connection pool exhausted")
	}
	conn := NewDriverConn(dbName) // 这里需要在实际使用时初始化 db
	pool.numOpen++
	pool.freeConn = append(pool.freeConn, conn)
	conn.inUse = true
	return conn, nil
}

// ReleaseConnection 将连接返回给连接池或关闭它
func (pool *ConnectionPool) ReleaseConnection(conn *driverConn) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// 检查连接是否已被关闭或最终关闭
	if conn.closed || conn.finalClosed {
		if !conn.finalClosed {
			// 关闭连接并减少计数
			conn.Close()
			pool.numOpen--
		}
		return
	}

	// 将连接标记为非活动状态
	conn.inUse = false

	// 检查连接池是否已达到最大打开连接数
	if pool.numOpen < pool.maxOpen {
		// 连接池未满，将连接放回空闲列表
		pool.freeConn = append(pool.freeConn, conn)
	} else {
		// 连接池已满，关闭连接而不将其放回连接池
		conn.Close()
		pool.numOpen-- // 减少计数，因为连接已被关闭
	}
}

// CloseConnection 关闭一个连接
func (conn *driverConn) CloseConnection() {
	conn.Lock()
	defer conn.Unlock()

	if !conn.closed {
		conn.ci.Close()
		conn.closed = true
	}
}

func (conn *driverConn) Close() {

}

// startCleanup 开始定期清理任务
func (pool *ConnectionPool) startCleanup() {
	pool.cleanupTicker = time.NewTicker(pool.maxLifetime) // 每次清理间隔为 maxLifetime
	defer pool.cleanupTicker.Stop()

	for {
		select {
		case <-pool.cleanupTicker.C:
			pool.cleanupStaleConnections()
		case <-pool.stopCleanup:
			return
		}
	}
}

// StopCleanup 停止定期清理任务
func (pool *ConnectionPool) StopCleanup() {
	pool.stopCleanup <- true
}

// cleanupStaleConnections 清理达到最大生命周期的连接
func (pool *ConnectionPool) cleanupStaleConnections() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	now := time.Now()
	for i := 0; i < len(pool.freeConn); {
		conn := pool.freeConn[i]
		if now.Sub(conn.createdAt) > pool.maxLifetime {
			conn.Close() // 关闭连接
			conn.closed = true
			conn.finalClosed = true
			pool.freeConn = append(pool.freeConn[:i], pool.freeConn[i+1:]...)
			continue
		}
		i++
	}
}
