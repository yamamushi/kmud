package server

import (
	"flag"
	"fmt"
	"github.com/yamamushi/kmud/config"
	"io"
	"log"
	"net"
	"os"
	"runtime/debug"
	"sort"
	"strconv"
	"time"

	"github.com/yamamushi/kmud/database"
	"github.com/yamamushi/kmud/engine"
	"github.com/yamamushi/kmud/model"
	"github.com/yamamushi/kmud/session"
	"github.com/yamamushi/kmud/telnet"
	"github.com/yamamushi/kmud/types"
	"github.com/yamamushi/kmud/utils"
	"gopkg.in/mgo.v2"
)

type Server struct {
	listener net.Listener
	config   config.Config
}

type connectionHandler struct {
	user   types.User
	pc     types.PC
	conn   *wrappedConnection
	config *config.Config
}

type wrappedConnection struct {
	telnet.Telnet
	watcher *utils.WatchableReadWriter
}

func (wc *wrappedConnection) Write(p []byte) (int, error) {
	return wc.watcher.Write(p)
}

func (wc *wrappedConnection) Read(p []byte) (int, error) {
	return wc.watcher.Read(p)
}

func login(wc *wrappedConnection) types.User {
	for {
		username := utils.GetUserInput(wc, "Username: ", types.ColorModeNone)

		if username == "" {
			return nil
		}

		user := model.GetUserByName(username)

		if user == nil {
			utils.WriteLine(wc, "User not found", types.ColorModeNone)
		} else if user.IsOnline() {
			utils.WriteLine(wc, "That user is already online", types.ColorModeNone)
		} else {
			attempts := 1
			wc.WillEcho()
			for {
				password := utils.GetRawUserInputSuffix(wc, "Password: ", "\r\n", types.ColorModeNone)

				// TODO - Disabling password verification to make development easier
				if user.VerifyPassword(password) || true {
					break
				}

				if attempts >= 3 {
					utils.WriteLine(wc, "Too many failed login attempts", types.ColorModeNone)
					wc.Close()
					panic("Booted user due to too many failed logins (" + user.GetName() + ")")
				}

				attempts++

				time.Sleep(2 * time.Second)
				utils.WriteLine(wc, "Invalid password", types.ColorModeNone)
			}
			wc.WontEcho()

			return user
		}
	}
}

func newUser(wc *wrappedConnection) types.User {
	for {
		name := utils.GetUserInput(wc, "Desired username: ", types.ColorModeNone)

		if name == "" {
			return nil
		}

		user := model.GetUserByName(name)
		password := ""

		if user != nil {
			utils.WriteLine(wc, "That name is unavailable", types.ColorModeNone)
		} else if err := utils.ValidateName(name); err != nil {
			utils.WriteLine(wc, err.Error(), types.ColorModeNone)
		} else {
			wc.WillEcho()
			for {
				pass1 := utils.GetRawUserInputSuffix(wc, "Desired password: ", "\r\n", types.ColorModeNone)

				if len(pass1) < 7 {
					utils.WriteLine(wc, "Passwords must be at least 7 letters in length", types.ColorModeNone)
					continue
				}

				pass2 := utils.GetRawUserInputSuffix(wc, "Confirm password: ", "\r\n", types.ColorModeNone)

				if pass1 != pass2 {
					utils.WriteLine(wc, "Passwords do not match", types.ColorModeNone)
					continue
				}

				password = pass1

				break
			}
			wc.WontEcho()

			admin := model.UserCount() == 0
			user = model.CreateUser(name, password, admin)
			return user
		}
	}
}

func (c *connectionHandler) newPlayer() types.PC {
	// TODO: character slot limit
	const SizeLimit = 12
	for {
		name := c.user.GetInput("Desired character name: ")

		if name == "" {
			return nil
		}

		char := model.GetCharacterByName(name)

		if char != nil {
			c.user.WriteLine("That name is unavailable")
		} else if err := utils.ValidateName(name); err != nil {
			c.user.WriteLine(err.Error())
		} else {
			room := model.GetRooms()[0] // TODO: Better way to pick an initial character location
			return model.CreatePlayerCharacter(name, c.user.GetId(), room)
		}
	}
}

func (c *connectionHandler) WriteLine(line string, a ...interface{}) {
	utils.WriteLine(c.conn, fmt.Sprintf(line, a...), types.ColorModeNone)
}

func (c *connectionHandler) Write(text string) {
	utils.Write(c.conn, text, types.ColorModeNone)
}

func (c *connectionHandler) GetInput(prompt string) string {
	return utils.GetUserInput(c.conn, prompt, types.ColorModeNone)
}

func (c *connectionHandler) GetWindowSize() (int, int) {
	return 80, 80
}

func (c *connectionHandler) mainMenu() {
	utils.ExecMenu(
		"MUD",
		c,
		func(menu *utils.Menu) {
			menu.AddAction("l", "Login", func() {
				c.user = login(c.conn)
				c.loggedIn()
			})

			menu.AddAction("n", "New user", func() {
				c.user = newUser(c.conn)
				c.loggedIn()
			})

			menu.AddAction("q", "Disconnect", func() {
				menu.Exit()
				return
			})

			menu.OnExit(func() {
				utils.WriteLine(c.conn, "Take luck!", types.ColorModeNone)
				c.conn.Close()
				return
			})
		})
}

func (c *connectionHandler) userMenu() {
	utils.ExecMenu(
		c.user.GetName(),
		c.user,
		func(menu *utils.Menu) {
			menu.OnExit(func() {
				c.user.SetOnline(false)
				c.user = nil
			})

			if c.user.IsAdmin() {
				menu.AddAction("a", "Admin", func() {
					c.adminMenu()
				})
			}

			menu.AddAction("n", "New character", func() {
				c.pc = c.newPlayer()
			})

			menu.AddAction("q", "logout", func() {
				menu.Exit()
				return
			})

			// TODO: Sort character list
			chars := model.GetUserCharacters(c.user.GetId())

			if len(chars) > 0 {
				menu.AddAction("d", "Delete character", func() {
					c.deleteMenu()
				})
			}

			for i, char := range chars {
				menu.AddAction(strconv.Itoa(i+1), char.GetName(), func() {
					c.pc = char
					c.launchSession()
				})
			}
		})
}

func (c *connectionHandler) deleteMenu() {
	utils.ExecMenu(
		"Delete character",
		c.user,
		func(menu *utils.Menu) {
			// TODO: Sort character list
			chars := model.GetUserCharacters(c.user.GetId())
			for i, char := range chars {
				c := char
				menu.AddAction(strconv.Itoa(i+1), char.GetName(), func() {
					// TODO: Delete confirmation
					model.DeleteCharacter(c.GetId())
				})
			}

			menu.AddAction("q", "Return to previous menu", func() {
				menu.Exit()
			})
		})
}

func (c *connectionHandler) adminMenu() {
	utils.ExecMenu(
		"Admin",
		c.user,
		func(menu *utils.Menu) {
			menu.AddAction("u", "Users", func() {
				c.userAdminMenu()
			})

			menu.AddAction("q", "Return to previous menu", func() {
				menu.Exit()
			})
		})
}

func (c *connectionHandler) userAdminMenu() {
	utils.ExecMenu("User Admin", c.user, func(menu *utils.Menu) {
		users := model.GetUsers()
		sort.Sort(users)

		for i, user := range users {
			online := ""
			if user.IsOnline() {
				online = "*"
			}

			u := user
			menu.AddAction(strconv.Itoa(i+1), user.GetName()+online, func() {
				c.specificUserMenu(u)
			})
		}

		menu.AddAction("q", "Return to previous menu", func() {
			menu.Exit()
		})
	})
}

func (c *connectionHandler) specificUserMenu(user types.User) {
	suffix := ""
	if user.IsOnline() {
		suffix = "(Online)"
	} else {
		suffix = "(Offline)"
	}

	utils.ExecMenu(
		fmt.Sprintf("User: %c %c", user.GetName(), suffix),
		c.user,
		func(menu *utils.Menu) {
			menu.AddAction("d", "Delete", func() {
				model.DeleteUser(user.GetId())
				menu.Exit()
			})

			menu.AddAction("a", fmt.Sprintf("Admin - %v", user.IsAdmin()), func() {
				u := model.GetUser(user.GetId())
				u.SetAdmin(!u.IsAdmin())
			})

			menu.AddAction("q", "Return to previous menu", func() {
				menu.Exit()
			})

			if user.IsOnline() {
				menu.AddAction("w", "Watch", func() {
					if user == c.user {
						c.user.WriteLine("You can't watch yours!")
					} else {
						userConn := user.GetConnection().(*wrappedConnection)

						userConn.watcher.AddWatcher(c.conn)
						utils.GetRawUserInput(c.conn, "Type anything to stop watching\r\n", c.user.GetColorMode())
						userConn.watcher.RemoveWatcher(c.conn)
					}
				})
			}
		})
}

func (c *connectionHandler) Handle() {
	go func() {
		defer c.conn.Close()

		defer func() {
			r := recover()

			username := ""
			charname := ""

			if c.user != nil {
				c.user.SetOnline(false)
				username = c.user.GetName()
			}

			if c.pc != nil {
				c.pc.SetOnline(false)
				charname = c.pc.GetName()
			}

			if r != io.EOF && c.config.Server.Debug {
				debug.PrintStack()
			}

			if username != "" && charname != "" {
				output := fmt.Sprintf("Lost connection to client (%v/%v): %v, %v",
					username,
					charname,
					c.conn.RemoteAddr(), r)
				log.Println(output)
			} else {
				output := fmt.Sprintf("Lost connection to client: %v", c.conn.RemoteAddr())
				log.Println(output)
			}

		}()

		c.mainMenu()
	}()
}

func (c *connectionHandler) loggedIn() {
	if c.user == nil {
		return
	}

	c.user.SetOnline(true)
	c.user.SetConnection(c.conn)

	c.conn.DoWindowSize()
	c.conn.DoTerminalType()

	c.conn.Listen(func(code telnet.TelnetCode, data []byte) {
		switch code {
		case telnet.WS:
			if len(data) != 4 {
				log.Println("Malformed window size data:", data)
				return
			}

			width := int((255 * data[0])) + int(data[1])
			height := int((255 * data[2])) + int(data[3])
			c.user.SetWindowSize(width, height)

		case telnet.TT:
			c.user.SetTerminalType(string(data))
		}
	})

	c.userMenu()
}

func (c *connectionHandler) launchSession() {
	if c.pc == nil {
		return
	}

	session := session.NewSession(c.conn, c.user, c.pc)
	session.Exec()
	c.pc = nil
}

func (s *Server) ReadConfig() {
	var confPath string
	flag.StringVar(&confPath, "c", "kmud.conf", "Path to Config File")
	flag.Parse()

	_, err := os.Stat(confPath)
	if err != nil {
		log.Fatal("Config file is missing: ", confPath)
		os.Exit(1)
	}

	s.config, err = config.ReadConfig(confPath)
	utils.HandleError(err)
}

func (s *Server) Start() {
	log.Println("Connecting to database... ")
	session, err := mgo.Dial(s.config.DB.MongoHost)
	utils.HandleError(err)
	log.Println("Database Connection Established")

	address := s.config.Server.Interface + ":" + s.config.Server.Port
	log.Println("Establishing Connection on " + address)
	s.listener, err = net.Listen("tcp", address)
	utils.HandleError(err)

	database.Init(database.NewMongoSession(session.Copy()), "mud")
	log.Println("Database Initialized")

}

func (s *Server) Bootstrap() {
	// Create the world object if necessary
	model.GetWorld()

	// If there are no rooms at all create one
	rooms := model.GetRooms()
	if len(rooms) == 0 {
		zones := model.GetZones()

		var zone types.Zone

		if len(zones) == 0 {
			zone, _ = model.CreateZone("Default")
		} else {
			zone = zones[0]
		}

		model.CreateRoom(zone, types.Coordinate{X: 0, Y: 0, Z: 0})
	}
}

func (s *Server) Listen() {
	for {
		conn, err := s.listener.Accept()
		utils.HandleError(err)
		log.Println("Client connected:", conn.RemoteAddr())
		t := telnet.NewTelnet(conn)

		wc := utils.NewWatchableReadWriter(t)

		ch := connectionHandler{
			config: &s.config,
			conn:   &wrappedConnection{Telnet: *t, watcher: wc},
		}

		ch.Handle()
	}
}

func (s *Server) Exec() {
	s.ReadConfig()
	log.Println("Starting Server Setup")
	s.Start()
	log.Println("Bootstrapping World")
	s.Bootstrap()
	log.Println("Starting World Engine")
	engine.Start()
	log.Println("Listening for Connections...")
	s.Listen()
}
