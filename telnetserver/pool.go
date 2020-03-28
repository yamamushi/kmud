package telnetserver

import (
	"errors"
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
	"log"
	"strings"
	"sync"
)

type ConnectionPool struct {
	pool     []*ConnectionHandler
	messages chan PoolMessage
	locker   sync.Mutex
}

type PoolMessage struct {
	TargetID string
	Type     string
	Command  string
	Args     []string
}

func NewConnectionPool() (pool *ConnectionPool) {
	pool = &ConnectionPool{}
	pool.messages = make(chan PoolMessage)
	return pool
}

func (p *ConnectionPool) AddToPool(c *ConnectionHandler) error {
	p.locker.Lock()
	defer p.locker.Unlock()

	for _, conn := range p.pool {
		if conn == c {
			return errors.New("connection already in pool")
		}
	}
	p.pool = append(p.pool, c)
	return nil
}

func (p *ConnectionPool) RemoveFromPool(c *ConnectionHandler) error {
	p.locker.Lock()
	defer p.locker.Unlock()

	removeWrappedConnection(p.pool, c)
	return nil
}

func (p *ConnectionPool) CloseConnection(id string) error {
	p.locker.Lock()
	defer p.locker.Unlock()

	return nil
}

func removeWrappedConnection(s []*ConnectionHandler, r *ConnectionHandler) []*ConnectionHandler {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (p *ConnectionPool) BroadcastMessage(message string, c *ConnectionHandler) error {

	return nil
}

func (p *ConnectionPool) Run() {
	for {
		select {
		case msg := <-p.messages:
			p.ParseMessage(msg)
		}
	}
}

func (p *ConnectionPool) CleanUp() {
	p.locker.Lock()
	defer p.locker.Unlock()

	for _, conn := range p.pool {
		if conn == nil {
			_ = p.RemoveFromPool(conn)
		}
	}
}

func (p *ConnectionPool) ParseMessage(message PoolMessage) {
	message.Type = strings.ToLower(message.Type)
	if message.Type == "disconnected" {

	}
	if message.Type == "broadcast" {
		if len(message.Args) > 0 {
			if len(p.pool) > 0 {
				for _, conn := range p.pool {
					if conn != nil {
						err := utils.WriteLine(conn.conn, message.Args[0], types.ColorModeNone)
						p.HandlePoolError(conn, err)
						err = utils.Write(conn.conn, "> ", types.ColorModeNone)
						p.HandlePoolError(conn, err)
					}
				}
			}
		}
	}
}

func (p *ConnectionPool) HandlePoolError(conn *ConnectionHandler, err error) {
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			_ = p.RemoveFromPool(conn)
		} else {
			log.Println("Pool Error: " + err.Error())
		}
	}
}
