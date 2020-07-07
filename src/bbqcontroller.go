package main

import (
	"bbq"
	"filter"
	"oled"
	"outputs"
	"pwm"
	"tc"
	"time"
	"webserver"
)

func min32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func main() {

	servo, e := pwm.New(25, 0.02, 0.0007, 0.0021)
	if e != nil {
		outputs.Err(e)
	}

	TC, e := tc.New(0x66, 1)
	if e != nil {
		outputs.Err(e)
	}

	screen, e := oled.NewSH1107(0x3c, 1)
	if e != nil {
		outputs.Err(e)
	}

	b := bbq.New(servo, TC, screen)
	defer b.Finalise()
	b.SetupCtrlC()
	go webserver.Serve(b)

	var duty float32
	duty = 50
	mm := filter.NewMM(300, 70)

	// Going to assume that a duty of 0 is closed and 100 is open

	for i := 0; ; i++ {
		pt := b.ProbeT()
		ptmm := mm.Add(pt)
		sett := b.GetT()

		if ptmm < sett-0.5 {
			duty = min32(duty+0.5, 100.0)
		}

		if ptmm > sett+0.5 {
			duty = max32(duty-0.5, 0)
		}

		b.ServoDuty(duty)
		b.Text(sett, pt, ptmm, b.AmbientT(), duty)
		time.Sleep(1 * time.Second)
	}

	time.Sleep(5 * time.Second)
}
