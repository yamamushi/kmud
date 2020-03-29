package telnet

import (
	"github.com/yamamushi/kmud-2020/toys"
	"github.com/yamamushi/kmud-2020/utils"
	"strconv"
	"strings"
	"time"
)

type Terminal struct {
	telnet  *Telnet
	Type    string
	Columns string
	ColI    int
	Rows    string
	RowI    int
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
	t.Columns = strconv.Itoa(x)
	t.ColI = x
	t.Rows = strconv.Itoa(y)
	t.RowI = y
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

func (t *Terminal) ClearCursorRight() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[0K"))
	}
}

func (t *Terminal) ClearCursorLeft() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[1K"))
	}
}

func (t *Terminal) ClearLine() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B[2K"))
	}
}

func (t *Terminal) Refresh() {
	if t.VT100 {
		_, _ = t.telnet.Write([]byte("\u001B8"))
		for row := 0; row <= t.RowI; row++ {
			_, _ = t.telnet.Write([]byte("\u001B[" + strconv.Itoa(row) + ";0H"))
			_, _ = t.telnet.Write([]byte("\u001B[2K"))
		}
	}
}

func (t *Terminal) Nyan() {
	if t.VT100 {
		for {
			var wait time.Duration
			wait = 100
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan0))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan1))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan2))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan3))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan4))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan5))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan6))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan7))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan8))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan9))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan10))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()

			_, _ = t.telnet.Write([]byte(toys.Nyan11))
			time.Sleep(wait * time.Millisecond)
			t.Refresh()
		}
	}
}
