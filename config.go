package fidoctl

import (
	"bytes"
	"encoding/asn1"
	"encoding/binary"
	"fmt"
)

const (
	ConfigTagUsbSupported     = 0x01
	ConfigTagSerial           = 0x02
	ConfigTagUsbEnabled       = 0x03
	ConfigTagFormFactor       = 0x04
	ConfigTagVersion          = 0x05
	ConfigTagAutoEjectTimeout = 0x06
	ConfigTagChalrespTimeout  = 0x07
	ConfigTagDeviceFlags      = 0x08
	ConfigTagAppVersions      = 0x09
	ConfigTagConfigLock       = 0x0A
	ConfigTagUnlock           = 0x0B
	ConfigTagReboot           = 0x0C
	ConfigTagNfcSupported     = 0x0D
	ConfigTagNfcEnabled       = 0x0E
	ConfigTagIapDetection     = 0x0F
	ConfigTagMoreData         = 0x10
	ConfigTagFreeForm         = 0x11
	ConfigTagHidInitDelay     = 0x12
	ConfigTagPartNumber       = 0x13
	ConfigTagFipsCapable      = 0x14
	ConfigTagFipsApproved     = 0x15
	ConfigTagPinComplexity    = 0x16
	ConfigTagNfcRestricted    = 0x17
	ConfigTagResetBlocked     = 0x18
	ConfigTagFpsVersion       = 0x20
	ConfigTagStmVersion       = 0x21
)

type YubiConfig struct {
	tags map[int][]byte
}

func (c *YubiConfig) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(0xff)
	for tag, data := range c.tags {
		valBytes, err := asn1.Marshal(asn1.RawValue{
			Tag:   tag,
			Bytes: data,
		})
		if err != nil {
			return nil, fmt.Errorf("marshal cfg tag 0x%x: %w", tag, err)
		}
		buf.Write(valBytes)
	}

	out := buf.Bytes()
	if len(out) > 0xff {
		return nil, fmt.Errorf("cfg is too big: %d", len(out))
	}

	out[0] = byte(len(out) - 1)
	return out, nil
}

func (c *YubiConfig) Unmarshal(data []byte) error {
	c.tags = make(map[int][]byte)

	var cfg asn1.RawValue
	var err error

	l := int(data[0])
	if len(data) < l+1 {
		return fmt.Errorf("invalid cfg length: %d", len(data))
	}

	rest := data[1 : l+1]
	for len(rest) > 0 {
		rest, err = asn1.Unmarshal(rest, &cfg)
		if err != nil {
			return fmt.Errorf("parse config: %w", err)
		}

		c.tags[cfg.Tag] = cfg.Bytes
	}

	return nil
}

func (c *YubiConfig) Has(tag int) bool {
	_, ok := c.tags[tag]
	return ok
}

func (c *YubiConfig) Get(tag int) []byte {
	return c.tags[tag]
}

func (c *YubiConfig) Set(tag int, data []byte) *YubiConfig {
	if c.tags == nil {
		c.tags = make(map[int][]byte)
	}

	c.tags[tag] = data

	return c
}

func (c *YubiConfig) Clear() *YubiConfig {
	c.tags = make(map[int][]byte)

	return c
}

func (c *YubiConfig) Version() Version {
	rawVersion := c.tags[ConfigTagVersion]
	if len(rawVersion) != 3 {
		return Version{}
	}

	return Version{
		Major: int(rawVersion[0]),
		Minor: int(rawVersion[1]),
		Patch: int(rawVersion[2]),
	}
}

func (c *YubiConfig) Serial() uint32 {
	rawSerial := c.tags[ConfigTagSerial]
	if len(rawSerial) == 0 {
		return 0
	}

	if len(rawSerial) < 3 {
		return 0
	}

	return binary.BigEndian.Uint32(rawSerial)
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
