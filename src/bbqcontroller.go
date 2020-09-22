package main

import (
	"bbq"
	"filter"
	"math/rand"
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
	mm := filter.NewMM(100, 80)
	mm1 := filter.NewMM(10, 60)
	mm2 := filter.NewMM(10, 60)

	// Going to assume that a duty of 0 is closed and 100 is open

	for i := 0; i <= 100; i++ {
		b.ServoDuty(float32(i))
		time.Sleep(30 * time.Millisecond)
	}
	for i := 100; i >= 0; i-- {
		b.ServoDuty(float32(i))
		time.Sleep(30 * time.Millisecond)
	}

	for i := 0; ; i++ {
		pt := b.ProbeT()
		ptmm := mm.Add(pt)
		curpt := mm1.Add(pt)
		dispat := mm2.Add(b.AmbientT())
		sett := b.GetT()

		if ((curpt < ptmm+0.1) && (ptmm < sett+0.5)) || (ptmm < sett-10.0) {
			duty = min32(duty+0.5, 100.0)
		}

		if ((curpt > ptmm-0.1) && (ptmm > sett-0.5)) || (ptmm > sett+10.0) {
			duty = max32(duty-0.5, 0)
		}

		b.ServoDuty(duty)
		b.Text(sett, curpt, ptmm, dispat, duty)
		time.Sleep(time.Duration(1000+rand.Int63n(1000)) * time.Millisecond)
	}

	time.Sleep(5 * time.Second)
}
