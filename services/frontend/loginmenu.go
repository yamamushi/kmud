package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/badoux/checkmail"
	"github.com/yamamushi/kmud-2020/color"
	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/crypt"
	"github.com/yamamushi/kmud-2020/telnet"
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func mainMenu(c *telnet.ConnectionHandler, term *telnet.Terminal, conf *config.Config) {
	// Menu is a helper set of utilities
	// For drawing an interactive menuing system
	utils.ExecMenu(
		conf.Frontend.Title,
		c,
		func(menu *utils.Menu) {
			menu.AddAction("l", "Login", func() {
				_, err := loginUserHandler(c.GetConn(), conf)
				if err == nil {

				}
			})

			menu.AddAction("n", "New user", func() {
				_ = registerUserHandler(c.GetConn(), conf)
			})

			menu.AddAction("t", "Testing", func() {
				term.ClearScreen()
				term.MoveCursor(0, (term.Columns/2)-10)
				term.LowIntensity()
				term.Reset()

			})

			menu.AddAction("q", "Disconnect", func() {
				menu.Exit()
				return
			})

			menu.OnExit(func() {
				// Of note here is c.GetConn() which will return a wrapped connection object
				// Note that this
				utils.WriteLine(c.GetConn(), color.Colorize(color.Red, "Come back soon!"), color.ModeDark)
				c.Close()
				return
			})
		})
}

// Login Menu
func loginUserHandler(wc *telnet.WrappedConnection, conf *config.Config) (auth types.AuthResponse, err error) {
	for {
		username := utils.GetUserInput(wc, "Username: ", color.ModeNone)

		if username == "" {
			return types.AuthResponse{}, errors.New("no username provided")
		}

		attempts := 1
		wc.WillEcho()
		for {
			password := utils.GetRawUserInputSuffix(wc, "Password: ", "\r\n", color.ModeNone)
			auth, ok := crypt.GetAuthToken(username, password, conf)
			if !ok {
				utils.WriteLine(wc, "Invalid password", color.ModeNone)
			} else {
				wc.WontEcho()
				//utils.WriteLine(wc, "Welcome "+username+" to "+conf.Game.ServerName, types.ModeNone)
				return auth, nil
			}

			if attempts >= 3 {
				utils.WriteLine(wc, "Too many failed loginUserHandler attempts", color.ModeNone)
				_ = wc.Close()
				log.Println("User booted user due to too many failed logins (" + username + ")")
			}
			attempts++
			time.Sleep(2 * time.Second)
		}
	}
}

// User Registrations
func registerUserHandler(wc *telnet.WrappedConnection, conf *config.Config) (err error) {
	for {
		var username, password, email string
		for {
			username = utils.GetUserInput(wc, "Desired username: ", color.ModeNone)
			if username == "" {
				utils.WriteLine(wc, "Exiting user registration due to empty username", color.ModeNone)
				return nil
			}
			if err := utils.ValidateName(username); err != nil {
				utils.WriteLine(wc, err.Error(), color.ModeNone)
			} else {
				break
			}
		}

		wc.WillEcho()
		for {
			pass1 := utils.GetRawUserInputSuffix(wc, "Desired password: ", "\r\n", color.ModeNone)
			if pass1 == "" {
				utils.WriteLine(wc, "Exiting user registration due to empty password", color.ModeNone)
				return nil
			}
			if len(pass1) < 7 {
				utils.WriteLine(wc, "Passwords must be at least 7 letters in length", color.ModeNone)
				continue
			}

			pass2 := utils.GetRawUserInputSuffix(wc, "Confirm password: ", "\r\n", color.ModeNone)
			if pass1 != pass2 {
				utils.WriteLine(wc, "Passwords do not match", color.ModeNone)
				continue
			}

			password = pass1
			break
		}
		wc.WontEcho()

		for {
			reademail := utils.GetUserInput(wc, "Enter your email: ", color.ModeNone)
			if reademail == "" {
				utils.WriteLine(wc, "Exiting user registration due to empty email", color.ModeNone)
				return nil
			}

			err = checkmail.ValidateFormat(reademail)
			if err != nil {
				utils.WriteLine(wc, "Invalid email format", color.ModeNone)
				continue
			}

			err = checkmail.ValidateHost(reademail)
			if err != nil {
				utils.WriteLine(wc, "Invalid email format", color.ModeNone)
				continue
			}

			email = reademail
			break
		}

		sha := crypt.Sha256Sum(password)
		hashedPass := string(sha)

		jsonData := types.AccountRegistrationRequest{
			Secret:     conf.Crypt.AccountManagerSecret,
			Username:   username,
			Email:      email,
			HashedPass: hashedPass,
		}

		jsonValue, _ := json.Marshal(jsonData)
		response, err := http.Post("http://"+conf.Cluster.AccountManagerHostname+"/registerUserHandler", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			utils.WriteLine(wc, "Unexpected Error: Please notify a Developer", color.ModeNone)
			return nil
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			output := types.AuthResponse{}
			err = json.Unmarshal(data, &output)
			if err != nil {
				utils.WriteLine(wc, "Unexpected Error: Please notify a Developer", color.ModeNone)
				log.Println("Error: GetAuthToken unmarshal failed with error: " + err.Error())
				return nil
			}
			if output.Err != "" {
				utils.WriteLine(wc, "Error: "+output.Err, color.ModeNone)
			}
			utils.WriteLine(wc, "Account Registered, you may now loginUserHandler.", color.ModeNone)
			return nil
		}
	}
}
