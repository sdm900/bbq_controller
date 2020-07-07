package oled

import (
	"bytes"
	"errors"
	"fmt"
	"i2c"
	"image"
	"image/png"
	"outputs"
)

func NewSH1107(addr uint8, bus int) (*Screen, error) {
	i2c, e := i2c.NewI2C(addr, bus)
	if e != nil {
		return nil, e
	}

	s := Screen{
		i2c:       i2c,
		buf:       nil,
		png:       new(bytes.Buffer),
		fontlarge: nil,
		fontsmall: nil,
	}

	e = s.init()
	if e != nil {
		return nil, e
	}
	return &s, nil
}

func (s *Screen) writei2c(cmd byte, buf []byte) error {
	b := make([]byte, len(buf)+1)
	b[0] = cmd
	for i := 0; i < len(buf); i++ {
		b[i+1] = buf[i]
	}

	_, e := s.i2c.WriteBytes(b)
	return e
}

func (s *Screen) cmd(bytes ...byte) error {
	if len(bytes) == 0 {
		return nil
	}

	if len(bytes) > 32 {
		return errors.New("Can only send 32 commands to I2C")
	}

	return s.writei2c(COMMANDMODE, bytes)
}

func (s *Screen) data(buf []byte) error {
	if len(buf) > 4096 {
		return errors.New(fmt.Sprintf("Screen buffer too large: %d", len(buf)))
	}

	return s.writei2c(DATAMODE, buf)
}

func (s *Screen) Close() {
	s.cmd(DISPOFF)
	s.i2c.Close()
}

func (s *Screen) display(buf []byte) error {
	if len(buf) != bytesscreen {
		return errors.New(fmt.Sprintf("Not 128x128=16384 bits: %d", len(buf)*8))
	}

	for i := 0; i < numpages; i++ {
		if e := s.cmd(byte(BASEADDR+i), 0x02, 0x10); e != nil {
			return e
		}

		if e := s.data(buf[i*bytesperpage : (i+1)*bytesperpage]); e != nil {
			return e
		}
	}

	return nil
}

func (s *Screen) White() error {
	buf := make([]byte, bytesscreen, bytesscreen)
	for i := 0; i < bytesscreen; i++ {
		buf[i] = 0xff
	}

	return s.display(buf)
}

func (s *Screen) Black() error {
	buf := make([]byte, bytesscreen, bytesscreen)
	for i := 0; i < bytesscreen; i++ {
		buf[i] = 0x00
	}

	return s.display(buf)
}

func (s *Screen) SetPNG(img *image.Gray) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.png.Reset()
	if e := png.Encode(s.png, img); e != nil {
		outputs.Msg("Failed to render image: ", e)
	}
}

func (s *Screen) GetPNG() []byte {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.png.Bytes()
}

func (s *Screen) Render(img *image.Gray) error {
	s.SetPNG(img)

	buf := make([]byte, bytesscreen, bytesscreen)
	for p := 0; p < numpages; p++ {
		for x := 0; x < WIDTH; x++ {
			b := byte(0x00)
			i := byte(0x01)
			for y := p * PIXELSPERBYTE; y < (p+1)*PIXELSPERBYTE; y++ {
				if img.GrayAt(x, y).Y > 117 {
					b = b | i
				}
				i = i << 1
			}

			// this seems a little complicated
			// it matches our screen setup with SETSEGREMAP
			// the paging and adressing of the screen is a little complex
			buf[p*WIDTH+(WIDTH-x-1)] = b
		}
	}
	return s.display(buf)
}
