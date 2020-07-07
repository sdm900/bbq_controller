// allows use of a MCP9600 thermocouple controller

package tc

import (
	"errors"
	"fmt"
	"i2c"
)

type TC struct {
	i2c *i2c.I2C
}

const (
	CHIPID    = 0x20
	HOT       = 0x00
	DELTA     = 0x01
	COLD      = 0x02
	TCCONFIG  = 0x05
	DEVCONFIG = 0x6
)

func New(addr uint8, bus int) (*TC, error) {
	i2c, e := i2c.NewI2C(addr, bus)
	if e != nil {
		return nil, e
	}

	t := TC{i2c}

	chipid, e := t.i2c.ReadRegU8(CHIPID)
	if e != nil {
		return nil, e
	}

	if chipid != 64 {
		return nil, errors.New(fmt.Sprintf("Thermocouple chipid %08b not $08b", chipid, 64))
	}

	tcconfig, e := t.i2c.ReadRegU8(TCCONFIG)
	if e != nil {
		return nil, e
	}

	if tcconfig != 0 {
		return nil, errors.New(fmt.Sprintf("Thermocouple %08b not K %08b", tcconfig, 0))
	}

	devconfig, e := t.i2c.ReadRegU8(DEVCONFIG)
	if e != nil {
		return nil, e
	}

	if devconfig != 0 {
		return nil, errors.New(fmt.Sprintf("Device %08b not K %08b", devconfig, 0))
	}

	return &t, nil
}

func reg2temp(t []byte) float32 {
	if t[0]&0x80 == 0x80 {
		return float32(t[0]<<4) + float32(t[1])/16.0 - 4096.0
	}
	return float32(t[0]<<4) + float32(t[1])/16.0
}

func (tc *TC) ProbeT() (float32, error) {
	t, n, e := tc.i2c.ReadRegBytes(HOT, 2)
	if n != 2 || e != nil {
		return 0, errors.New(fmt.Sprintf("Could not read probe temp %d %s", n, e))
	}

	return reg2temp(t), nil
}

func (tc *TC) AmbientT() (float32, error) {
	t, n, e := tc.i2c.ReadRegBytes(COLD, 2)
	if n != 2 || e != nil {
		return 0, errors.New(fmt.Sprintf("Could not read ambient temp %d %s", n, e))
	}

	return reg2temp(t), nil
}

func (tc *TC) Close() {
	tc.i2c.Close()
}
