package oled

import (
	"bytes"
	"errors"
	"fmt"
	"fonts"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"i2c"
	"sync"
)

type Screen struct {
	i2c       *i2c.I2C
	buf       []byte
	png       *bytes.Buffer
	fontlarge font.Face
	fontsmall font.Face
	mux       sync.Mutex
}

const (
	COMMANDMODE   = 0x00
	DATAMODE      = 0x40
	DISPOFF       = 0xae
	DISPON        = 0xaf
	MEMMODE       = 0x20
	SETHIGHCOL    = 0x10
	SETLOWCOL     = 0x00
	SETSEGREMAP   = 0xa0
	NORMDISP      = 0xa6
	SETMULTIPLEX  = 0xa8
	DISPONRES     = 0xa4
	SETDISPOFFSET = 0xd3
	SETDISPCLKDIV = 0xd5
	SETPRECHARGE  = 0xd9
	SETCOMPINS    = 0xda
	SETVCOMDET    = 0xdb
	CHARGEPUMP    = 0x8d
	CONTRAST      = 0x81
	BASEADDR      = 0xb0
	WIDTH         = 128
	HEIGHT        = 128
	PIXELSPERBYTE = 8
)

var (
	bytesscreen   = WIDTH * HEIGHT / PIXELSPERBYTE
	numpages      = WIDTH / PIXELSPERBYTE
	pixelsperpage = WIDTH * HEIGHT / numpages
	bytesperpage  = bytesscreen / numpages
)

func (s *Screen) init() error {
	b, n, e := s.i2c.ReadRegBytes(COMMANDMODE, 1)
	if e != nil {
		return e
	}
	if n != 1 {
		return errors.New(fmt.Sprintf("Read too many bytes %d", n))
	}
	if b[0]&0x3f != 0x07 {
		return errors.New(fmt.Sprintf("Read wrong device %x", b[0]))
	}

	if e := s.cmd(DISPOFF, MEMMODE, SETHIGHCOL, 0xb0, 0xc8, SETLOWCOL, 0x10, 0x40,
		SETSEGREMAP, NORMDISP, SETMULTIPLEX, 0xff, DISPONRES, SETDISPOFFSET,
		0x02, SETDISPCLKDIV, 0xf0, SETPRECHARGE, 0x22, SETCOMPINS, 0x12,
		SETVCOMDET, 0x20, CHARGEPUMP, 0x14, CONTRAST, 0x0f); e != nil {
		return e
	}

	if e := s.Black(); e != nil {
		return e
	}

	fsmall, e := truetype.Parse(fonts.TTFSmall)
	if e != nil {
		return e
	}
	flarge, e := truetype.Parse(fonts.TTFLarge)
	if e != nil {
		return e
	}

	s.fontlarge = truetype.NewFace(flarge, &truetype.Options{65, 0, font.HintingFull, 0, 0, 0})
	s.fontsmall = truetype.NewFace(fsmall, &truetype.Options{8, 0, font.HintingFull, 0, 0, 0})

	return s.cmd(DISPON)
}
