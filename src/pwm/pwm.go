package pwm

import (
	"errors"
	"fmt"
	"os"
)

type PWM struct {
	pin     uint
	fd      *os.File
	minDuty float32
	maxDuty float32
}

func (p *PWM) write(s string) error {
	if _, e := p.fd.WriteString(s); e != nil {
		return e
	}

	return nil
}

func (p *PWM) setpin(f float32) error {
	// f is a percentage 0-100
	// what we write to the pin is the fraction of duty 0-1.0
	s := fmt.Sprintf("%d=%0.3f\n", p.pin, f/100.0)

	if e := p.write(s); e != nil {
		return e
	}

	return nil
}

func New(pin uint, period, minPulse, maxPulse float32) (*PWM, error) {
	// period and pulse widths are in seconds (s)
	fd, e := os.OpenFile("/dev/pi-blaster", os.O_APPEND|os.O_WRONLY, 0600)
	if e != nil {
		return nil, e
	}

	p := PWM{pin, fd, minPulse / period * 100.0, maxPulse / period * 100.0}
	if e := p.write("\n"); e != nil {
		return nil, e
	}

	if e := p.setpin(0); e != nil {
		return nil, e
	}

	return &p, nil
}

func (p *PWM) Duty(d float32) error {
	// d is a percentage eg. 0-100

	if d < 0.0 || d > 100.0 {
		return errors.New(fmt.Sprintf("Illegal value for duty cycle: %f", d))
	}

	d = p.minDuty + (p.maxDuty-p.minDuty)*d/100.0
	if e := p.setpin(d); e != nil {
		return errors.New("Could not write to /dev/pi-blaster: " + e.Error())
	}

	return nil
}

func (p *PWM) Close() {
	p.write("\n")
	p.setpin(0)
	p.fd.Close()
}
