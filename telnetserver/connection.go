package telnetserver

import (
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"strings"

	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/telnet"
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
)

type ConnectionHandler struct {
	user   types.User
	pc     types.PC
	conn   *WrappedConnection
	config *config.Config
}

type WrappedConnection struct {
	telnet.Telnet
	watcher *utils.WatchableReadWriter
}

func (wc *WrappedConnection) Write(p []byte) (int, error) {
	return wc.watcher.Write(p)
}

func (wc *WrappedConnection) Read(p []byte) (int, error) {
	return wc.watcher.Read(p)
}

func (c *ConnectionHandler) WriteLine(line string, a ...interface{}) {
	utils.WriteLine(c.conn, fmt.Sprintf(line, a...), types.ColorModeNone)
}

func (c *ConnectionHandler) Write(text string) {
	utils.Write(c.conn, text, types.ColorModeNone)
}

func (c *ConnectionHandler) GetInput(prompt string) string {
	return utils.GetUserInput(c.conn, prompt, types.ColorModeNone)
}

func (c *ConnectionHandler) GetWindowSize() (int, int) {
	// This is incorrect, this should ask the client for information, instead it's static for now
	return 80, 80
}

func (c *ConnectionHandler) Handle(runner func(c *ConnectionHandler, conf *config.Config), conf *config.Config) {
	go func() {
		defer c.conn.Close()
		defer c.HandleDisconnect()

		runner(c, conf)
	}()
}

func (c *ConnectionHandler) HandleDisconnect() {
	r := recover()

	if r != io.EOF && c.config.Server.Debug {
		debug.PrintStack()
	}

	if strings.ToLower(c.config.Server.LoggingLevel) == "all" {
		output := fmt.Sprintf("Lost connection to client: %v", c.conn.RemoteAddr())
		log.Println(output)
	}
}

func (c *ConnectionHandler) Close() {
	c.conn.Close()
}

func (c *ConnectionHandler) GetConn() (wc *WrappedConnection) {
	return c.conn
}
