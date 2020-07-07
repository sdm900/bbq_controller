package oled

import (
	"bytes"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"io/ioutil"
	"network"
	"strconv"
	"strings"
	"time"
)

func left(img *image.Gray, s string, strn, y int, f font.Face) {
	d := &font.Drawer{img, image.White, f, fixed.Point26_6{0, fixed.Int26_6(y << 6)}}
	d.DrawString(s)
}

func right(img *image.Gray, s string, strn, y int, f font.Face) {
	d := &font.Drawer{img, image.White, f, fixed.Point26_6{0, 0}}

	var l fixed.Int26_6
	if strn == 0 {
		l = d.MeasureString(strings.Repeat(s, strn))
	} else {
		l = d.MeasureString(strings.Repeat("W", strn))
	}
	d.Dot = fixed.Point26_6{fixed.Int26_6((WIDTH-1)<<6) - l, fixed.Int26_6(y << 6)}
	d.DrawString(s)
}

func centre(img *image.Gray, s string, y int, f font.Face) {
	d := &font.Drawer{img, image.White, f, fixed.Point26_6{0, 0}}
	l := d.MeasureString(s)
	d.Dot = fixed.Point26_6{(fixed.Int26_6((WIDTH-1)<<6) - l) / 2, fixed.Int26_6(y << 6)}
	d.DrawString(s)
}

func tempformat(t float32) string {
	if t < 100 {
		return fmt.Sprintf("%4.1f", t)
	}
	return fmt.Sprintf("%3.0f", t)
}

func systemp() float32 {
	var t string
	if b, e := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp"); e != nil {
		t = "0"
	} else {
		t = string(bytes.TrimSpace(b))
	}
	st, e := strconv.ParseFloat(t, 32)
	if e != nil {
		st = 0
	}

	return float32(st / 1000.0)
}

func sysload() string {
	var t string
	if b, e := ioutil.ReadFile("/proc/loadavg"); e != nil {
		t = "Unknown"
	} else {
		bb := bytes.Fields(b)
		t = string(bb[0]) + " " + string(bb[1]) + " " + string(bb[2])
	}

	return t
}

func servoarrow(img *image.Gray, duty, y int, f font.Face) {
	d := &font.Drawer{img, image.White, f, fixed.Point26_6{fixed.Int26_6(0), fixed.Int26_6(y << 6)}}
	s := fmt.Sprintf("%d", duty)
	d.DrawString(s)

	iini := 7
	iend := WIDTH - 1
	imid := (iini + iend) / 2

	for j := y - 1; j >= y-7; j-- {
		img.SetGray(imid, j, color.Gray{0xff})
	}

	y = y - 4
	if duty < 50 {
		i := imid - 1
		for ; i > imid-(50-duty)*(imid-iini)/50+3; i-- {
			img.SetGray(i, y, color.Gray{0xff})
		}

		for ii := 2; ii >= 0; ii-- {
			for j := y - ii; j <= y+ii; j++ {
				img.SetGray(i, j, color.Gray{0xff})
			}
			i--
		}
	} else if duty > 50 {
		i := imid + 1
		for ; i < imid+(duty-50)*(iend-imid)/50-3; i++ {
			img.SetGray(i, y, color.Gray{0xff})
		}

		for ii := 2; ii >= 0; ii-- {
			for j := y - ii; j <= y+ii; j++ {
				img.SetGray(i, j, color.Gray{0xff})
			}
			i++
		}
	}
}

func (s *Screen) Text(setpoint, temp, mm, amb, duty float32) error {
	img := image.NewGray(image.Rect(0, 0, WIDTH-1, HEIGHT-1))

	left(img, time.Now().Format("Mon 2 Jan"), 10, 7, s.fontsmall)
	right(img, time.Now().Format("15:04:05"), 8, 7, s.fontsmall)
	left(img, network.Hostname()+": "+network.GetLocalIP(), 0, 18, s.fontsmall)
	left(img, "Load: "+sysload(), 0, 29, s.fontsmall)
	servoarrow(img, int(duty), 40, s.fontsmall)
	left(img, "SetT: "+tempformat(setpoint), 10, 56, s.fontsmall)
	right(img, "SysT: "+tempformat(systemp()), 10, 56, s.fontsmall)
	centre(img, tempformat(temp), HEIGHT-17, s.fontlarge)
	left(img, "AveT: "+tempformat(mm), 10, HEIGHT-1, s.fontsmall)
	right(img, "AmbT: "+tempformat(amb), 10, HEIGHT-1, s.fontsmall)

	return s.Render(img)
}
