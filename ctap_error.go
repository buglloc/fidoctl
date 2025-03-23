package fidoctl

import (
	"fmt"
)

type CTAPErrCode uint8

const (
	CTAPErrCodeNone                CTAPErrCode = 0x00
	CTAPErrCodeInvalidCommand      CTAPErrCode = 0x01
	CTAPErrCodeInvalidParameter    CTAPErrCode = 0x02
	CTAPErrCodeInvalidLen          CTAPErrCode = 0x03
	CTAPErrCodeInvalidSequence     CTAPErrCode = 0x04
	CTAPErrCodeMsgTimeout          CTAPErrCode = 0x05
	CTAPErrCodeChannelBusy         CTAPErrCode = 0x06
	CTAPErrCodeChannelLockRequired CTAPErrCode = 0x07
	CTAPErrCodeChannelInvalid      CTAPErrCode = 0x08
	CTAPErrCodeChannelUnknown      CTAPErrCode = 0x09
)

type CTAPError struct {
	Code CTAPErrCode
}

func (e *CTAPError) Error() string {
	return fmt.Sprintf("ctap command: %s", e.errString())
}

func (e *CTAPError) errString() string {
	switch e.Code {
	case CTAPErrCodeInvalidCommand:
		return "Invalid command"
	case CTAPErrCodeInvalidParameter:
		return "Invalid parameter"
	case CTAPErrCodeInvalidLen:
		return "Invalid length"
	case CTAPErrCodeInvalidSequence:
		return "Invalid sequence"
	case CTAPErrCodeMsgTimeout:
		return "Msg timeout"
	case CTAPErrCodeChannelBusy:
		return "Channel busy"
	case CTAPErrCodeChannelLockRequired:
		return "Command requires channel lock"
	case CTAPErrCodeChannelInvalid:
		return "Invalid channel"
	case CTAPErrCodeChannelUnknown:
		return "unspecific error"
	default:
		return fmt.Sprintf("err 0x%x", e.Code)
	}
}
