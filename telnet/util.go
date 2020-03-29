package telnet

import (
	"strconv"
)

func ToString(bytes []byte) string {
	str := ""
	for _, b := range bytes {

		if str != "" {
			str = str + " "
		}

		str = str + ByteToCodeString(b)
	}

	return str
}

func ByteToCodeString(b byte) string {
	code, found := byteToCode[b]

	if !found {
		return "??(" + strconv.Itoa(int(b)) + ")"
	}

	return CodeToString(code)
}

func CodeToString(code TelnetCode) string {
	switch code {
	case NUL:
		return "NUL"
	case ECHO:
		return "ECHO"
	case SGA:
		return "SGA"
	case ST:
		return "ST"
	case TM:
		return "TM"
	case BEL:
		return "BEL"
	case BS:
		return "BS"
	case HT:
		return "HT"
	case LF:
		return "LF"
	case FF:
		return "FF"
	case CR:
		return "CR"
	case TT:
		return "TT"
	case WS:
		return "WS"
	case TS:
		return "TS"
	case RFC:
		return "RFC"
	case LM:
		return "LM"
	case EV:
		return "EV"
	case SE:
		return "SE"
	case NOP:
		return "NOP"
	case DM:
		return "DM"
	case BRK:
		return "BRK"
	case IP:
		return "IP"
	case AO:
		return "AO"
	case AYT:
		return "AYT"
	case EC:
		return "EC"
	case EL:
		return "EL"
	case GA:
		return "GA"
	case SB:
		return "SB"
	case WILL:
		return "WILL"
	case WONT:
		return "WONT"
	case DO:
		return "DO"
	case DONT:
		return "DONT"
	case IAC:
		return "IAC"
	case CMP1:
		return "CMP1"
	case CMP2:
		return "CMP2"
	case AARD:
		return "AARD"
	case ATCP:
		return "ATCP"
	case GMCP:
		return "GMCP"
	}

	return ""
}
