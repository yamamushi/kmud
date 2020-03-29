package telnet

import (
	"errors"
	"strconv"
	"strings"
)

func (t *Telnet) WillEcho() {
	t.SendCommand(WILL, ECHO)
}

func (t *Telnet) WontEcho() {
	t.SendCommand(WONT, ECHO)
}

func (t *Telnet) DoWindowSize() (int, int, error) {
	t.SendCommand(DO, WS)
	response, err := t.ReadIACResponse()
	if err != nil {
		return 0, 0, err
	}
	x := 0
	y := 0
	if strings.Contains(response, "IAC WILL WS") {
		codes := strings.Split(response, " ")
		for _, code := range codes {
			if strings.Contains(code, "??(") {
				code = strings.TrimPrefix(code, "??(")
				code = strings.TrimSuffix(code, ")")
				if x == 0 {
					x, err = strconv.Atoi(code)
					if err != nil {
						return 0, 0, errors.New("IAC window size conversion failed")
					}
				} else {
					y, err = strconv.Atoi(code)
					if err != nil {
						return 0, 0, errors.New("IAC window size conversion failed")
					}
				}
			}
		}
		return x, y, nil
	}

	return 0, 0, errors.New("no proper window size response found")
}

func (t *Telnet) DoTerminalType() (string, error) {
	// This is really supposed to be two commands, one to ask if they'll send a
	// terminal type, and another to indicate that they should send it if
	// they've expressed a "willingness" to send it. For the time being this
	// works well enough.

	// See http://tools.ietf.org/html/rfc884

	t.SendCommand(DO, TT, IAC, SB, TT, 1, IAC, SE) // 1 = SEND
	iac, err := t.ReadIACResponse()
	if err != nil {
		return "", err
	}
	iacfields := strings.Split(iac, " ")

	var term string
	for _, field := range iacfields {
		if strings.Contains(field, "??(") {
			field = strings.TrimPrefix(field, "??(")
			field = strings.TrimSuffix(field, ")")
			val, err := strconv.Atoi(field)
			if err != nil {
				return "", err
			}
			term = term + string(val)
		}
	}

	return strings.ToLower(term), err
}

func (t *Telnet) SendCommand(codes ...TelnetCode) {
	t.conn.Write(BuildCommand(codes...))
}

func (t *Telnet) ReadIACResponse() (string, error) {
	b := make([]byte, 1024)
	_, err := t.conn.Read(b)
	if err != nil {
		return "", err
	}

	fields := strings.Split(ToString(b), " ")

	for i := len(fields); i >= 0; i-- {
		if i == 1 {
			break
		}
		if fields[i-1] != "NUL" {
			break
		} else {
			fields = fields[:i-1+copy(fields[i-1:], fields[i:])]
		}
	}

	var command string
	for _, code := range fields {
		command = command + " " + code
	}
	return command, nil
}

func BuildCommand(codes ...TelnetCode) []byte {
	command := make([]byte, len(codes)+1)
	command[0] = codeToByte[IAC]

	for i, code := range codes {
		command[i+1] = codeToByte[code]
	}

	return command
}
