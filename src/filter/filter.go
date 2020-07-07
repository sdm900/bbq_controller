package filter

import "math"

type MM struct {
	list []float32
	sort []float32
	perc int
	mm   float32
	i    int
}

func NewMM(n, p int) *MM {
	l := make([]float32, 0, n)
	s := make([]float32, 0, n)

	if p < 0 {
		p = 0
	}
	if p > 100 {
		p = 100
	}

	mm := MM{l, s, p, 0.0, -1}
	return &mm
}

func (mm *MM) Add(a float32) float32 {
	if mm.i < 0 {
		mm.list = append(mm.list, a)
		mm.sort = append(mm.sort, a)
		mm.mm = a
		mm.i = 0

		return a
	}

	mm.i = (mm.i + 1) % cap(mm.list)
	if mm.i == len(mm.list) {
		mm.list = append(mm.list, -math.MaxFloat32)
		mm.sort = append(mm.sort, -math.MaxFloat32)
	}

	b := mm.list[mm.i]
	mm.list[mm.i] = a

	if a > b {
		for i := 0; i < len(mm.list); i++ {
			bb := mm.sort[i]

			if bb == b {
				mm.sort[i] = a
				break
			} else if bb < a {
				mm.sort[i] = a
				a = bb
			}
		}
	} else {
		for i := len(mm.list) - 1; i >= 0; i-- {
			bb := mm.sort[i]

			if bb == b {
				mm.sort[i] = a
				break
			} else if bb > a {
				mm.sort[i] = a
				a = bb
			}
		}
	}

	pn := len(mm.list) * (100 - mm.perc) / 200
	c := float32(0)
	j := int(0)
	for i := len(mm.list) - pn - 1; i >= pn; i-- {
		j++
		c = c + mm.sort[i]
	}
	mm.mm = c / float32(j)

	return mm.mm
}
