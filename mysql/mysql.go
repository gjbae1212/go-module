package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/gjbae1212/go-module/util"
	_ "github.com/go-sql-driver/mysql"
)

type Connector interface {
	DB() *sql.DB
	DSN() string
	Connect(tries uint) error
	Close()
}

type Connection struct {
	dsn             string
	conn            *sql.DB
	backoff         *util.Backoff
	connectedAmount uint
	connectionMux   *sync.Mutex
}

func NewWithDSN(dsn string, tries uint) (*Connection, error) {
	conn := &Connection{
		dsn:           dsn,
		backoff:       util.NewBackoff(20 * time.Second),
		connectionMux: new(sync.Mutex),
	}
	if err := conn.Connect(tries); err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Connection) DB() *sql.DB {
	return c.conn
}

func (c *Connection) DSN() string {
	return c.dsn
}

func (c *Connection) Connect(tries uint) error {
	if tries == 0 {
		return nil
	}
	c.connectionMux.Lock()
	defer c.connectionMux.Unlock()
	if c.connectedAmount > 0 {
		// already have opened connection
		c.connectedAmount++
		return nil
	}
	var err error
	var db *sql.DB
	for i := tries; i > 0; i-- {
		// Wait before attempt.
		time.Sleep(c.backoff.Wait())

		// Open connection to MySQL but...
		db, err = sql.Open("mysql", c.dsn)
		if err != nil {
			continue
		}

		// ...try to use the connection for real.
		if err = db.Ping(); err != nil { // Connection failed.  Wrong username or password?
			db.Close()
			continue
		}

		// Connected
		c.conn = db
		c.backoff.Success()
		c.connectedAmount++
		return nil
	}

	return fmt.Errorf("[Err] Failed to connect to MySQL after %d tries (%v)", tries, err)
}

func (c *Connection) Close() {
	c.connectionMux.Lock()
	defer c.connectionMux.Unlock()
	if c.connectedAmount == 0 { // connection closed already
		return
	}
	c.connectedAmount--
	if c.connectedAmount == 0 && c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
