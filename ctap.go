package fidoctl

import (
	"fmt"
)

const (
	initialCID uint32 = 0xFFFFFFFF
)

type CTAPPacket struct {
	CMD   uint8
	BCNTH uint8
	BCNTL uint8
	DATA  []byte
	seq   uint8
}

func (c *CTAPPacket) ToHID(cid uint32, reportInputLength uint16) [][]byte {
	var result [][]byte
	result = make([][]byte, 0)
	length := uint16(len(c.DATA))
	c.BCNTL = uint8(length & 0xff)
	c.BCNTH = uint8((length >> 8) & 0xff)

	out := make([]byte, 0, int(reportInputLength))
	out = append(out, cid2Bytes(cid)...)
	out = append(out, c.CMD)
	out = append(out, c.BCNTH, c.BCNTL)

	seq := uint8(0)
	for n := 0; n < len(c.DATA); n++ {
		out = append(out, c.DATA[n])

		if len(out) == int(reportInputLength) {
			result = append(result, out)
			out = make([]byte, 0, int(reportInputLength))
			out = append(out, cid2Bytes(cid)...)
			out = append(out, seq)
			seq++
		}
	}

	if len(out) != 5 {
		for n := len(out); n < int(reportInputLength); n++ {
			out = append(out, 0)
		}
		result = append(result, out)
	}
	return result
}

func (c *CTAPPacket) FromHID(cid uint32, in []byte) (bool, error) {
	if len(in) < 5 {
		return true, fmt.Errorf("packet too short: %d", len(in))
	}

	packetCID := bytesToCid(in[3], in[2], in[1], in[0])
	if packetCID != cid {
		return false, fmt.Errorf("CID mismatch: %x (ecpected) != %x (actual)", cid, packetCID)
	}

	switch {
	case in[4]&0x80 == 0:
		// seq packet
		seq := in[4]
		if seq != c.seq {
			return true, fmt.Errorf("sequence mismatc: %d (actual) != %d (expected)", seq, c.seq)
		}

		c.seq++
		c.DATA = append(c.DATA, in[5:]...)
	default:
		// command packet
		switch in[4] & 0x7f {
		case 0x3F:
			// error
			return true, &CTAPError{
				Code: CTAPErrCode(in[7]),
			}
		case 0x3B:
			// keep-alive
			return false, nil
		default:
			c.CMD = in[4]
			c.BCNTH = in[5]
			c.BCNTL = in[6]
			c.DATA = append(c.DATA, in[7:]...)
			c.seq = 0
		}
	}

	if uint32(len(c.DATA)) < uint32(c.BCNTL)+uint32(c.BCNTH)<<8 {
		// need more data
		return false, nil
	}

	return true, nil
}

func cid2Bytes(cid uint32) []byte {
	return []byte{
		byte(cid >> 24), byte(cid >> 16), byte(cid >> 8), byte(cid & 0xff),
	}
}

func bytesToCid(b1, b2, b3, b4 byte) uint32 {
	return uint32(b1) | uint32(b2)<<8 | uint32(b3)<<16 | uint32(b4)<<24
}
