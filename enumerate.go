package fidoctl

import (
	"fmt"

	"github.com/buglloc/usbhid"
)

const (
	YubicoVID = 0x1050
)

func Enumerate() ([]Device, error) {
	devices, err := usbhid.Enumerate(
		usbhid.WithVidFilter(YubicoVID),
		usbhid.WithDeviceFilterFunc(isSuitableDev),
	)
	if err != nil {
		return nil, fmt.Errorf("enumerate HID devices: %w", err)
	}

	out := make([]Device, len(devices))
	for i, d := range devices {
		out[i] = Device{
			cid: initialCID,
			dev: d,
		}
	}

	return out, nil
}

func isSuitableDev(device *usbhid.Device) bool {
	switch device.ProductId() {
	case 0x0113, 0x0402:
		// FIDO
	case 0x0115, 0x0406:
		// FIDO+CCID
	case 0x0116, 0x0407:
		// OTP+FIDO+CCID
	default:
		return false
	}

	return device.UsagePage() == 0xF1D0 && device.Usage() == 0x1
}
