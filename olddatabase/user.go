package olddatabase

import (
	"crypto/sha1"
	"fmt"
	"github.com/yamamushi/kmud-2020/color"
	"io"
	"net"
	"reflect"

	"github.com/yamamushi/kmud-2020/utils"
)

type User struct {
	DbObject `bson:",inline"`

	Name      string
	ColorMode color.ColorMode
	Password  []byte
	Admin     bool

	online       bool
	conn         net.Conn
	windowWidth  int
	windowHeight int
	terminalType string
}

func NewUser(name string, password string, admin bool) *User {
	user := &User{
		Name:         utils.FormatName(name),
		Password:     hash(password),
		ColorMode:    color.ModeNone,
		Admin:        admin,
		online:       false,
		windowWidth:  80,
		windowHeight: 40,
	}

	dbinit(user)
	return user
}

func (u *User) GetName() string {
	u.ReadLock()
	defer u.ReadUnlock()

	return u.Name
}

func (u *User) SetName(name string) {
	u.writeLock(func() {
		u.Name = utils.FormatName(name)
	})
}

func (u *User) SetOnline(online bool) {
	u.online = online

	if !online {
		u.conn = nil
	}
}

func (u *User) IsOnline() bool {
	return u.online
}

func (u *User) SetColorMode(cm color.ColorMode) {
	u.writeLock(func() {
		u.ColorMode = cm
	})
}

func (u *User) GetColorMode() color.ColorMode {
	u.ReadLock()
	defer u.ReadUnlock()
	return u.ColorMode
}

func hash(data string) []byte {
	h := sha1.New()
	io.WriteString(h, data)
	return h.Sum(nil)
}

// SetPassword SHA1 hashes the password before saving it to the database
func (u *User) SetPassword(password string) {
	u.writeLock(func() {
		u.Password = hash(password)
	})
}

func (u *User) VerifyPassword(password string) bool {
	hashed := hash(password)
	return reflect.DeepEqual(hashed, u.GetPassword())
}

// GetPassword returns the SHA1 of the user's password
func (u *User) GetPassword() []byte {
	u.ReadLock()
	defer u.ReadUnlock()
	return u.Password
}

func (u *User) SetConnection(conn net.Conn) {
	u.conn = conn
}

func (u *User) GetConnection() net.Conn {
	return u.conn
}

func (u *User) SetWindowSize(width int, height int) {
	u.windowWidth = width
	u.windowHeight = height
}

const MinWidth = 60
const MinHeight = 20

func (u *User) GetWindowSize() (width int, height int) {
	return utils.Max(u.windowWidth, MinWidth),
		utils.Max(u.windowHeight, MinHeight)
}

func (u *User) SetTerminalType(tt string) {
	u.terminalType = tt
}

func (u *User) GetTerminalType() string {
	return u.terminalType
}

func (u *User) GetInput(prompt string) string {
	return utils.GetUserInput(u.conn, prompt, u.GetColorMode())
}

func (u *User) WriteLine(line string, a ...interface{}) {
	utils.WriteLine(u.conn, fmt.Sprintf(line, a...), u.GetColorMode())
}

func (u *User) Write(text string) {
	utils.Write(u.conn, text, u.GetColorMode())
}

func (u *User) SetAdmin(admin bool) {
	u.writeLock(func() {
		u.Admin = admin
	})
}

func (u *User) IsAdmin() bool {
	u.ReadLock()
	defer u.ReadUnlock()
	return u.Admin
}
