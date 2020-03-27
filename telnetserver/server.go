package telnetserver

import (
	"log"
	"net"

	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/telnet"
	"github.com/yamamushi/kmud-2020/utils"
)

type Server struct {
	listener net.Listener
	config   *config.Config
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

func (s *Server) Listen(runner func(c *ConnectionHandler)) {
	for {
		conn, err := s.listener.Accept()
		utils.HandleError(err)
		log.Println("Client connected:", conn.RemoteAddr())
		t := telnet.NewTelnet(conn)

		wc := utils.NewWatchableReadWriter(t)

		ch := ConnectionHandler{
			config: s.config,
			conn:   &WrappedConnection{Telnet: *t, watcher: wc},
		}

		ch.Handle(runner)
	}
}

func (s *Server) Run(runner func(c *ConnectionHandler)) (err error) {
	log.Println("Starting Service")
	err = s.Setup()
	if err != nil {
		return err
	}

	log.Println("Listening for Connections...")
	s.Listen(runner)
	return nil
}
