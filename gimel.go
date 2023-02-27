package gimel

import (
	"math"
	"math/big"
	"strings"
)

var (
	zeroValue        = big.NewInt(0)
	oneValue         = big.NewInt(1)
	tenValue         = big.NewInt(10)
	defaultPrecision = big.NewInt(100)
)

type Gimel struct {
	neg    bool
	digits *big.Int
	exp    *big.Int
}

func SetGimelPrecision(digits int64) {
	defaultPrecision = big.NewInt(digits)
}

func (g Gimel) Norm() Gimel {
	gl := len(g.digits.String())
	var a big.Int
	b := a.Sub(defaultPrecision, big.NewInt(int64(gl))).Int64()
	g.digits.Mul(g.digits, big.NewInt(int64(math.Pow10(int(b)))))
	return g
}

func (g Gimel) Clone() Gimel {
	return Gimel{
		g.neg,
		(&big.Int{}).Set(g.digits),
		(&big.Int{}).Set(g.exp),
	}
}

func (g Gimel) Abs() Gimel {
	a := g.Clone()
	a.neg = false
	return a
}

func (g Gimel) Neg() Gimel {
	a := g.Clone()
	a.neg = !g.neg
	return a
}

func (g Gimel) Cmp(o Gimel) (r int) {
	switch {
	case g == o:
		// do nothing
	case g.neg == o.neg:
		r = g.exp.Cmp(o.exp)
		if r == 0 {
			r = g.digits.Cmp(o.digits)
		}
		if g.neg {
			r = -r
		}
	case g.neg:
		r = -1
	default:
		r = 1
	}
	return
}

func (g Gimel) Gt(o Gimel) bool  { return g.Cmp(o) == 1 }
func (g Gimel) Gte(o Gimel) bool { return g.Cmp(o) != -1 }
func (g Gimel) Lt(o Gimel) bool  { return g.Cmp(o) == -1 }
func (g Gimel) Lte(o Gimel) bool { return g.Cmp(o) != 1 }
func (g Gimel) Eq(o Gimel) bool  { return g.Cmp(o) == 0 }
func (g Gimel) Neq(o Gimel) bool { return g.Cmp(o) != 0 }

func (g Gimel) Min(o Gimel) Gimel {
	if g.Lt(o) {
		return g.Clone()
	} else {
		return o.Clone()
	}
}

func (g Gimel) Max(o Gimel) Gimel {
	if g.Gt(o) {
		return g.Clone()
	} else {
		return o.Clone()
	}
}

func (g Gimel) IsPos() bool { return !g.neg }
func (g Gimel) IsNeg() bool { return g.neg }

func (g Gimel) Add(o Gimel) Gimel {
	var a big.Int
	a.Sub(g.exp, o.exp)
	if a.CmpAbs(defaultPrecision) == 1 {
		return g.Max(o)
	}
	a2 := a.Sign()
	a.Abs(&a)
	var a3 big.Int
	a3.Exp(tenValue, &a, nil)
	var a4 big.Int // max exp

	// shift to match
	var b1, b2 big.Int
	switch a2 {
	case 0:
		b1.Set(g.digits)
		b2.Set(o.digits)
		a4.Set(g.exp)
	case 1:
		b1.Mul(g.digits, &a3)
		b2.Set(o.digits)
		a4.Set(g.exp)
	case -1:
		b1.Set(g.digits)
		b2.Mul(o.digits, &a3)
		a4.Set(o.exp)
	}

	// flip to negative
	if g.neg {
		b1.Neg(&b1)
	}
	if o.neg {
		b2.Neg(&b2)
	}

	var c1 big.Int
	c1.Add(&b1, &b2)
	if a2 != 0 {
		c1.Div(&c1, &a3)
	}
	c2 := c1.Sign() == -1
	c1.Abs(&c1)
	return Gimel{c2, &c1, &a4}
}

func (g Gimel) Sub(o Gimel) Gimel {
	var a big.Int
	a.Sub(g.exp, o.exp)
	if a.CmpAbs(defaultPrecision) == 1 {
		return g.Max(o)
	}
	// TODO
	return Gimel{}
}

func (g Gimel) Mul(o Gimel) Gimel {
	return Gimel{}
}

func (g Gimel) Div(o Gimel) Gimel {
	return Gimel{}
}

// String is just an alias for TextE for the Stringer interface
func (g Gimel) String() string { return g.TextE() }

func (g Gimel) TextE() string {
	var b strings.Builder
	if g.neg {
		b.WriteByte('-')
	}
	a := strings.TrimRight(g.digits.String(), "0")
	b.WriteByte(a[0])
	b.WriteByte('.')
	b.WriteString(a[1:])
	b.WriteByte('e')
	b.WriteString(g.exp.String())
	return b.String()
}

func (g Gimel) Text(sep rune) string {
	var b strings.Builder
	if g.neg {
		b.WriteByte('-')
	}

	if sep == 0 {
		g.writeFullDigits(&b)
	} else {
		var b2 strings.Builder
		g.writeFullDigits(&b2)
		a := b2.String()
		l := len(a)
		// start at digit 0th triple
		for i := -(3 - l%3); i < l; i += 3 {
			if i < 0 {
				b.WriteString(a[0 : i+3])
			} else {
				if i != 0 {
					b.WriteRune(',')
				}
				b.WriteString(a[i : i+3])
			}
		}
	}
	return b.String()
}

func (g Gimel) writeFullDigits(b *strings.Builder) {
	b.WriteString(g.digits.String())
	var c big.Int
	c.Sub(g.exp, defaultPrecision)
	c.Add(&c, oneValue)
	for i := new(big.Int); i.Cmp(&c) < 0; i.Add(i, oneValue) {
		b.WriteByte('0')
	}
}
