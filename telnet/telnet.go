package telnet

import (
	"log"
	"net"
	"time"
)

// RFC 854: http://tools.ietf.org/html/rfc854, http://support.microsoft.com/kb/231866

var byteToCode map[byte]TelnetCode
var codeToByte map[TelnetCode]byte

type TelnetCode int

func init() {
	byteToCode = map[byte]TelnetCode{}
	codeToByte = map[TelnetCode]byte{}

	codeToByte[NUL] = '\x00'
	codeToByte[ECHO] = '\x01'
	codeToByte[SGA] = '\x03'
	codeToByte[ST] = '\x05'
	codeToByte[TM] = '\x06'
	codeToByte[BEL] = '\x07'
	codeToByte[BS] = '\x08'
	codeToByte[HT] = '\x09'
	codeToByte[LF] = '\x0a'
	codeToByte[FF] = '\x0c'
	codeToByte[CR] = '\x0d'
	codeToByte[TT] = '\x18'
	codeToByte[WS] = '\x1F'
	codeToByte[TS] = '\x20'
	codeToByte[RFC] = '\x21'
	codeToByte[LM] = '\x22'
	codeToByte[EV] = '\x24'
	codeToByte[SE] = '\xf0'
	codeToByte[NOP] = '\xf1'
	codeToByte[DM] = '\xf2'
	codeToByte[BRK] = '\xf3'
	codeToByte[IP] = '\xf4'
	codeToByte[AO] = '\xf5'
	codeToByte[AYT] = '\xf6'
	codeToByte[EC] = '\xf7'
	codeToByte[EL] = '\xf8'
	codeToByte[GA] = '\xf9'
	codeToByte[SB] = '\xfa'
	codeToByte[WILL] = '\xfb'
	codeToByte[WONT] = '\xfc'
	codeToByte[DO] = '\xfd'
	codeToByte[DONT] = '\xfe'
	codeToByte[IAC] = '\xff'

	codeToByte[CMP1] = '\x55'
	codeToByte[CMP2] = '\x56'
	codeToByte[AARD] = '\x66'
	codeToByte[ATCP] = '\xc8'
	codeToByte[GMCP] = '\xc9'

	for enum, code := range codeToByte {
		byteToCode[code] = enum
	}
}

// Telnet wraps the given connection object, processing telnet codes from its byte
// stream and interpreting them as necessary, making it possible to hand the connection
// object off to other code so that it doesn't have to worry about telnet escape sequences
// being found in its data.
type Telnet struct {
	conn net.Conn
	err  error

	processor *telnetProcessor
}

func NewTelnet(conn net.Conn) *Telnet {
	var t Telnet
	t.conn = conn
	t.processor = newTelnetProcessor()
	return &t
}

func (t *Telnet) Write(p []byte) (int, error) {
	return t.conn.Write(p)
}

func (t *Telnet) Read(p []byte) (int, error) {
	for {
		t.fill()
		if t.err != nil {
			return 0, t.err
		}

		n, err := t.processor.Read(p)
		if n > 0 {
			return n, err
		}
	}
}

func (t *Telnet) Data(code TelnetCode) []byte {
	return t.processor.subdata[code]
}

func (t *Telnet) Listen(listenFunc func(TelnetCode, []byte)) {
	t.processor.listenFunc = listenFunc
}

// Idea/name for this function shamelessly stolen from bufio
func (t *Telnet) fill() {
	buf := make([]byte, 1024)
	n, err := t.conn.Read(buf)
	t.err = err
	t.processor.addBytes(buf[:n])
}

func (t *Telnet) Close() error {
	return t.conn.Close()
}

func (t *Telnet) LocalAddr() net.Addr {
	return t.conn.LocalAddr()
}

func (t *Telnet) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *Telnet) SetDeadline(dl time.Time) error {
	return t.conn.SetDeadline(dl)
}

func (t *Telnet) SetReadDeadline(dl time.Time) error {
	return t.conn.SetReadDeadline(dl)
}

func (t *Telnet) SetWriteDeadline(dl time.Time) error {
	return t.conn.SetWriteDeadline(dl)
}

type processorState int

const (
	stateBase   processorState = iota
	stateInIAC  processorState = iota
	stateInSB   processorState = iota
	stateCapSB  processorState = iota
	stateEscIAC processorState = iota
)

// telnetProcessor implements a state machine that reads input one byte at a time
// and processes it according to the telnet spec. It is designed to read a raw telnet
// stream, from which it will extract telnet escape codes and subnegotiation data.
// The processor can then be read from with all of the telnet codes removed, leaving
// the pure user input stream.
type telnetProcessor struct {
	state     processorState
	currentSB TelnetCode

	capturedBytes []byte
	subdata       map[TelnetCode][]byte
	cleanData     string
	listenFunc    func(TelnetCode, []byte)

	debug bool
}

func newTelnetProcessor() *telnetProcessor {
	var tp telnetProcessor
	tp.state = stateBase
	tp.debug = false
	tp.currentSB = NUL

	return &tp
}

func (tp *telnetProcessor) Read(p []byte) (int, error) {
	maxLen := len(p)

	n := 0

	if maxLen >= len(tp.cleanData) {
		n = len(tp.cleanData)
	} else {
		n = maxLen
	}

	for i := 0; i < n; i++ {
		p[i] = tp.cleanData[i]
	}

	tp.cleanData = tp.cleanData[n:] // TODO: Memory leak?

	return n, nil
}

func (tp *telnetProcessor) capture(b byte) {
	if tp.debug {
		log.Println("Captured:", ByteToCodeString(b))
	}

	tp.capturedBytes = append(tp.capturedBytes, b)
}

func (tp *telnetProcessor) dontCapture(b byte) {
	tp.cleanData = tp.cleanData + string(b)
}

func (tp *telnetProcessor) resetSubDataField(code TelnetCode) {
	if tp.subdata == nil {
		tp.subdata = map[TelnetCode][]byte{}
	}

	tp.subdata[code] = []byte{}
}

func (tp *telnetProcessor) captureSubData(code TelnetCode, b byte) {
	if tp.debug {
		log.Println("Captured subdata:", CodeToString(code), b)
	}

	if tp.subdata == nil {
		tp.subdata = map[TelnetCode][]byte{}
	}

	tp.subdata[code] = append(tp.subdata[code], b)
}

func (tp *telnetProcessor) addBytes(bytes []byte) {
	for _, b := range bytes {
		tp.addByte(b)
	}
}

func (tp *telnetProcessor) addByte(b byte) {
	code := byteToCode[b]

	switch tp.state {
	case stateBase:
		if code == IAC {
			tp.state = stateInIAC
			tp.capture(b)
		} else {
			tp.dontCapture(b)
		}

	case stateInIAC:
		if code == WILL || code == WONT || code == DO || code == DONT {
			// Stay in this state
		} else if code == SB {
			tp.state = stateInSB
		} else {
			tp.state = stateBase
		}
		tp.capture(b)

	case stateInSB:
		tp.capture(b)
		tp.currentSB = code
		tp.state = stateCapSB
		tp.resetSubDataField(code)

	case stateCapSB:
		if code == IAC {
			tp.state = stateEscIAC
		} else {
			tp.captureSubData(tp.currentSB, b)
		}

	case stateEscIAC:
		if code == IAC {
			tp.state = stateCapSB
			tp.captureSubData(tp.currentSB, b)
		} else {
			tp.subDataFinished(tp.currentSB)
			tp.currentSB = NUL
			tp.state = stateBase
			tp.addByte(codeToByte[IAC])
			tp.addByte(b)
		}
	}
}

func (tp *telnetProcessor) subDataFinished(code TelnetCode) {
	if tp.listenFunc != nil {
		tp.listenFunc(code, tp.subdata[code])
	}
}
