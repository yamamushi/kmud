package main

import (
	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/telnetserver"
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
)

func mainMenu(c *telnetserver.ConnectionHandler, conf *config.Config) {
	// Menu is a helper set of utilities
	// For drawing an interactive menuing system
	utils.ExecMenu(
		conf.Login.Title,
		c,
		func(menu *utils.Menu) {
			menu.AddAction("l", "Login", func() {

			})

			menu.AddAction("n", "New user", func() {

			})

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