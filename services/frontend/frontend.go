package main

// Default necessary imports from kmud-2020 libraries
import (
	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/telnetserver"
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
)

func main() {
	// utils.DefaultInit will set the random seed and scale max processes based on CPU count.
	go utils.DefaultInit()

	// config.GetConfig will pull a config either by the argument given from the command line flag -c
	// or by the default value passed in here.
	conf, err := config.GetConfig("frontend.conf")
	if err != nil {
		utils.HandleError(err)
	}

	// Here we create our server object using the provided configuration file.
	s := telnetserver.NewServer(conf)

	// We execute the server using a func(c *telnetserver.ConnectionHandler) function
	// The provided function will run in a goroutine and is expected to handle
	// All connections (the functionality will vary depending on the service)
	s.Run(mainMenu)
}

func mainMenu(c *telnetserver.ConnectionHandler) {
	// Menu is a helper set of utilities
	// For drawing an interactive menuing system
	utils.ExecMenu(
		"MUD",
		c,
		func(menu *utils.Menu) {
			menu.AddAction("q", "Disconnect", func() {
				menu.Exit()
				return
			})

			menu.OnExit(func() {
				// Of note here is c.GetConn() which will return a wrapped connection object
				// Note that this
				utils.WriteLine(c.GetConn(), "Come back soon!", types.ColorModeNone)
				c.Close()
				return
			})
		})
}
