package fidoctl

import (
	"crypto/rand"
	"fmt"

	"github.com/buglloc/usbhid"
)

const (
	typeInit      = 0x80
	cmdInit       = 0x06
	cmdReadConfig = 0xC2
	cmdInsConfig  = 0xC3
)

type Device struct {
	cid uint32
	dev *usbhid.Device
}

func (d *Device) Open() error {
	if err := d.dev.Open(true); err != nil {
		return fmt.Errorf("opening device: %w", err)
	}

	return d.init()
}

func (d *Device) OneShot(fn func(d *Device) error) error {
	if !d.dev.IsOpen() {
		if err := d.Open(); err != nil {
			return fmt.Errorf("opening device: %w", err)
		}
		defer func() { _ = d.Close() }()
	}

	return fn(d)
}

func (d *Device) Reboot() error {
	return d.OneShot(func(d *Device) error {
		var cfg YubiConfig
		data, err := cfg.Set(ConfigTagReboot, nil).Marshal()
		if err != nil {
			return fmt.Errorf("marshal config: %w", err)
		}

		_, err = d.SendAndReceive(cmdInsConfig, data)
		return err
	})
}

func (d *Device) YubiConfig() (*YubiConfig, error) {
	var cfg YubiConfig
	err := d.OneShot(func(d *Device) error {
		data, err := d.SendAndReceive(cmdReadConfig, []byte{0x00})
		if err != nil {
			return fmt.Errorf("reading config: %w", err)
		}

		if err := cfg.Unmarshal(data); err != nil {
			return fmt.Errorf("unmarshaling config: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (d *Device) SendAndReceive(cmd uint8, data []byte) ([]byte, error) {
	req := CTAPPacket{
		CMD:  cmd | typeInit,
		DATA: data,
	}

	if err := d.send(&req); err != nil {
		return nil, fmt.Errorf("send packet: %w", err)
	}

	rsp, err := d.recv()
	if err != nil {
		return nil, fmt.Errorf("recv packet: %w", err)
	}

	if req.CMD != rsp.CMD {
		return nil, fmt.Errorf("unexpected command response: %x (got) != %x (expected)", req.CMD, rsp.CMD)
	}

	return rsp.DATA, nil
}

func (d *Device) Path() string {
	return d.dev.Path()
}

func (d *Device) String() string {
	return d.dev.String()
}

func (d *Device) send(req *CTAPPacket) error {
	for i, report := range req.ToHID(d.cid, d.dev.GetOutputReportLength()) {
		if err := d.dev.SetOutputReport(0, report); err != nil {
			return fmt.Errorf("send report[%d]: %w", i, err)
		}
	}

	return nil
}

func (d *Device) recv() (*CTAPPacket, error) {
	var packet CTAPPacket
	for {
		id, data, err := d.dev.GetInputReport()
		if err != nil {
			return nil, fmt.Errorf("read report: %w", err)
		}

		if id != 0 {
			return nil, fmt.Errorf("invalid report id: %x (expected) != %x (actual)", 0, id)
		}

		done, err := packet.FromHID(d.cid, data)
		if err != nil {
			return nil, fmt.Errorf("parse report: %w", err)
		}

		if done {
			break
		}
	}

	return &packet, nil
}

func (d *Device) init() error {
	nonce := make([]byte, 8)
	_, err := rand.Read(nonce)
	if err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	rsp, err := d.SendAndReceive(cmdInit, nonce)
	if err != nil {
		return fmt.Errorf("send INIT command: %w", err)
	}

	if len(rsp) < 17 {
		return fmt.Errorf("invalid INIT response: too short: %d (actual) < 17 (expected)", len(rsp))
	}

	d.cid = bytesToCid(rsp[8], rsp[9], rsp[10], rsp[11])
	return nil
}

func (d *Device) Close() error {
	return d.dev.Close()
}
