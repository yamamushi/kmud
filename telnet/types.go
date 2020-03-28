package telnet

const (
	NUL  TelnetCode = iota // NULL, no operation
	ECHO TelnetCode = iota // Echo
	SGA  TelnetCode = iota // Suppress go ahead
	ST   TelnetCode = iota // Status
	TM   TelnetCode = iota // Timing mark
	BEL  TelnetCode = iota // Bell
	BS   TelnetCode = iota // Backspace
	HT   TelnetCode = iota // Horizontal tab
	LF   TelnetCode = iota // Line feed
	FF   TelnetCode = iota // Form feed
	CR   TelnetCode = iota // Carriage return
	TT   TelnetCode = iota // Terminal type
	WS   TelnetCode = iota // Window size
	TS   TelnetCode = iota // Terminal speed
	RFC  TelnetCode = iota // Remote flow control
	LM   TelnetCode = iota // Line mode
	EV   TelnetCode = iota // Environment variables
	SE   TelnetCode = iota // End of subnegotiation parameters.
	NOP  TelnetCode = iota // No operation.
	DM   TelnetCode = iota // Data Mark. The data stream portion of a Synch. This should always be accompanied by a TCP Urgent notification.
	BRK  TelnetCode = iota // Break. NVT character BRK.
	IP   TelnetCode = iota // Interrupt Process
	AO   TelnetCode = iota // Abort output
	AYT  TelnetCode = iota // Are you there
	EC   TelnetCode = iota // Erase character
	EL   TelnetCode = iota // Erase line
	GA   TelnetCode = iota // Go ahead signal
	SB   TelnetCode = iota // Indicates that what follows is subnegotiation of the indicated option.
	WILL TelnetCode = iota // Indicates the desire to begin performing, or confirmation that you are now performing, the indicated option.
	WONT TelnetCode = iota // Indicates the refusal to perform, or continue performing, the indicated option.
	DO   TelnetCode = iota // Indicates the request that the other party perform, or confirmation that you are expecting the other party to perform, the indicated option.
	DONT TelnetCode = iota // Indicates the demand that the other party stop performing, or confirmation that you are no longer expecting the other party to perform, the indicated option.
	IAC  TelnetCode = iota // Interpret as command

	// Non-standard codes:
	CMP1 TelnetCode = iota // MCCP Compress
	CMP2 TelnetCode = iota // MCCP Compress2
	AARD TelnetCode = iota // Aardwolf MUD out of band communication, http://www.aardwolf.com/blog/2008/07/10/telnet-negotiation-control-mud-client-interaction/
	ATCP TelnetCode = iota // Achaea Telnet Client Protocol, http://www.ironrealms.com/rapture/manual/files/FeatATCP-txt.html
	GMCP TelnetCode = iota // Generic Mud Communication Protocol
)
