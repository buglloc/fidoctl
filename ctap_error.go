package fidoctl

import (
	"fmt"
)

type CTAPError struct {
	Code uint8
}

func (e *CTAPError) Error() string {
	return fmt.Sprintf("ctap command: %s", e.errString())
}

func (e *CTAPError) errString() string {
	switch e.Code {
	case 0x01:
		return "Invalid command"
	case 0x02:
		return "Invalid parameter"
	case 0x03:
		return "Invalid length"
	case 0x04:
		return "Invalid sequence"
	case 0x05:
		return "Msg timeout"
	case 0x06:
		return "Channel busy"
	case 0x0A:
		return "Command requires channel lock"
	case 0x0B:
		return "Invalid channel"
	case 0x7F:
		return "unspecific error"
	default:
		return fmt.Sprintf("err 0x%x", e.Code)
	}
}
