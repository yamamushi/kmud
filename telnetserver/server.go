package telnetserver

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/telnet"
	"github.com/yamamushi/kmud-2020/utils"
)

type Server struct {
	listener net.Listener
	config   *config.Config
	pool     *ConnectionPool
}

func NewServer(config *config.Config) (s *Server) {
	s = &Server{config: config}
	return s
}

func (s *Server) Setup() (err error) {
	address := s.config.Server.Interface + ":" + s.config.Server.Port
	log.Println("Establishing Connection on " + address)
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Bootstrap() {

}

func (s *Server) Listen(runner func(c *ConnectionHandler, conf *config.Config), conf *config.Config) {
	for {
		conn, err := s.listener.Accept()
		utils.HandleError(err)
		log.Println("Client connected:", conn.RemoteAddr())
		t := telnet.NewTelnet(conn)

		wc := utils.NewWatchableReadWriter(t)

		id, err := utils.GetUUID()
		if err != nil {
			utils.HandleError(errors.New("utils.GetUUID - " + err.Error()))
			return
		}

		ch := ConnectionHandler{
			id:     id,
			config: s.config,
			conn:   &WrappedConnection{Telnet: *t, watcher: wc},
			pool:   s.pool.messages,
		}
		err = s.pool.AddToPool(&ch)
		if err != nil {
			utils.Error("server listen() add to pool failure: " + err.Error())
		} else {
			ch.Handle(runner, conf)
		}
	}
}

func (s *Server) CreateConnectionPool() {
	s.pool = NewConnectionPool()
	go s.pool.Run()
}

func (s *Server) Run(runner func(c *ConnectionHandler, conf *config.Config), conf *config.Config) (err error) {
	log.Println("Starting Service")
	err = s.Setup()
	if err != nil {
		return err
	}

	log.Println("Creating Connection Pool")
	s.CreateConnectionPool()

	go s.TestBroadcastLoop()

	log.Println("Listening for Connections...")
	s.Listen(runner, conf)
	return nil
}

func (s *Server) TestBroadcastLoop() {
	for {
		time.Sleep(time.Second * 5)
		log.Println("Broadcasting Message")
		s.pool.messages <- PoolMessage{Type: "broadcast", Args: []string{"\nThis is a test of the emergency broadcast system"}}
	}
}
