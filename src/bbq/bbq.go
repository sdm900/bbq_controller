package bbq

import (
	"oled"
	"os"
	"os/signal"
	"outputs"
	"pwm"
	"sync"
	"syscall"
	"tc"
)

type BBQ struct {
	servo  *pwm.PWM
	tc     *tc.TC
	screen *oled.Screen
	setT   float32
	mux    sync.Mutex
}

func (b *BBQ) Finalise() {
	outputs.Msg("Cleaning up")
	if b.servo != nil {
		b.servo.Close()
	}

	if b.tc != nil {
		b.tc.Close()
	}

	if b.screen != nil {
		b.screen.Close()
	}
}

func (b *BBQ) signal(c chan os.Signal) {
	<-c
	b.Finalise()
	os.Exit(1)
}

func (b *BBQ) SetupCtrlC() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go b.signal(c)
}

func (b *BBQ) GetT() float32 {
	b.mux.Lock()
	defer b.mux.Unlock()

	return b.setT
}

func (b *BBQ) SetT(t float32) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.setT = t
}

func (b *BBQ) ScreenPNG() []byte {
	return b.screen.GetPNG()
}

func (b *BBQ) ProbeT() float32 {
	t, err := b.tc.ProbeT()
	if err != nil {
		outputs.Msg("Error reading probe temperature: ", err)
		t = 0
	}
	return t
}

func (b *BBQ) AmbientT() float32 {
	t, err := b.tc.AmbientT()
	if err != nil {
		outputs.Msg("Error reading ambient temperature: ", err)
		t = 0
	}
	return t
}

func (b *BBQ) ServoDuty(d float32) {
	if e := b.servo.Duty(d); e != nil {
		outputs.Msg("Could not set duty of servo: ", e)
	}
}

func (b *BBQ) Text(st, pt, ptmm, at, duty float32) {
	err := b.screen.Text(st, pt, ptmm, at, duty)
	if err != nil {
		outputs.Msg("Error displaying text: ", err)

	}
}

func New(servo *pwm.PWM, tc *tc.TC, screen *oled.Screen) *BBQ {
	return &BBQ{
		servo:  servo,
		tc:     tc,
		screen: screen,
		setT:   110.0,
	}
}
