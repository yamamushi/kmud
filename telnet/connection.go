package telnet

import (
	"fmt"
	"github.com/yamamushi/kmud-2020/color"
	"io"
	"log"
	"runtime/debug"
	"strings"

	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/utils"
)

type ConnectionHandler struct {
	id        string
	AuthToken string
	pool      chan PoolMessage
	conn      *WrappedConnection
	config    *config.Config
}

type WrappedConnection struct {
	Telnet
	watcher *utils.WatchableReadWriter
}

// Write a raw byte to the connection rather than through the io.Writer (the wc.watcher writer)
func (wc *WrappedConnection) RawWrite(p []byte) (int, error) {
	return wc.Telnet.Write(p)
}

// Raw Read byte to the connection rather than through the io.Reader (the wc.watcher reader)
func (wc *WrappedConnection) RawRead(p []byte) (int, error) {
	return wc.Telnet.Read(p)
}

func (wc *WrappedConnection) Write(p []byte) (int, error) {
	return wc.watcher.Write(p)
}

func (wc *WrappedConnection) Read(p []byte) (int, error) {
	return wc.watcher.Read(p)
}

func (c *ConnectionHandler) WriteLine(line string, a ...interface{}) {
	utils.WriteLine(c.conn, fmt.Sprintf(line, a...), color.ModeNone)
}

func (c *ConnectionHandler) Write(text string) {
	utils.Write(c.conn, text, color.ModeNone)
}

func (c *ConnectionHandler) GetInput(prompt string) string {
	return utils.GetUserInput(c.conn, prompt, color.ModeNone)
}

func (c *ConnectionHandler) GetWindowSize() (int, int) {
	// This is incorrect, this should ask the client for information, instead it's static for now
	return 80, 80
}

func (c *ConnectionHandler) Handle(runner func(c *ConnectionHandler, term *Terminal, conf *config.Config), term *Terminal, conf *config.Config) {
	go func() {
		defer c.conn.Close()
		defer c.HandleDisconnect()

		runner(c, term, conf)
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

	c.pool <- PoolMessage{TargetID: c.id, Type: "disconnected"}
}

func (c *ConnectionHandler) Close() {
	c.conn.Close()
}

func (c *ConnectionHandler) GetConn() (wc *WrappedConnection) {
	return c.conn
}
