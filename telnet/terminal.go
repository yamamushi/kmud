package telnet

import (
	"github.com/yamamushi/kmud-2020/utils"
	"strconv"
	"strings"
)

type Terminal struct {
	telnet  *Telnet
	Type    string
	Columns int
	Rows    int
	VT100   bool
}

func GetTermInfo(telnet *Telnet) (*Terminal, error) {
	terminal := &Terminal{}
	return ResetSettings(telnet, terminal)
}

func ResetSettings(telnet *Telnet, t *Terminal) (*Terminal, error) {
	termtype, err := telnet.DoTerminalType()
	if err != nil {
		utils.Error("server Read IAC error: " + err.Error())
		return nil, err
	}
	var vt100 bool
	if strings.Contains(termtype, "xterm") {
		vt100 = true
	}
	//log.Println(termtype)

	x, y, err := telnet.DoWindowSize()
	if err != nil {
		utils.Error("server Read IAC error: " + err.Error())
		return nil, err
	}
	//log.Println(x, y)
	t.telnet = telnet
	t.Type = termtype
	t.Columns = x
	t.Rows = y
	t.VT100 = vt100
	return t, nil
}

// VT100 Escape Sequences
// http://ascii-table.com/ansi-escape-sequences-vt-100.php

func (t *Terminal) Reset() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001Bc"))
	}
}

func (t *Terminal) ClearScreen() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[2J"))
	}
}

// ClearDown() Clears the terminal from the cursor down
func (t *Terminal) ClearDown() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[J"))
	}
}

// ClearDown() Clears the terminal from the cursor up
func (t *Terminal) ClearUp() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[1J"))
	}
}

func (t *Terminal) LineFeedMode() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[20l"))
	}
}

func (t *Terminal) DisableCharAttr() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[0m"))
	}
}

func (t *Terminal) Bold() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[1m"))
	}
}

func (t *Terminal) LowIntensity() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[2m"))
	}
}

func (t *Terminal) Underline() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[4m"))
	}
}

func (t *Terminal) Blink() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[5m"))
	}
}

func (t *Terminal) Reverse() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[7m"))
	}
}

func (t *Terminal) Invisible() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[8m"))
	}
}

func (t *Terminal) MoveCursor(row int, column int) {

	x := strconv.Itoa(row)
	y := strconv.Itoa(column)

	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[" + x + ";" + y + "H"))
	}
}

func (t *Terminal) ResetCursor() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B8"))
	}
}
