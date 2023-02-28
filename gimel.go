package gimel

import (
	"fmt"
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

func G(neg bool, digits, exp *big.Int) Gimel {
	return Gimel{neg, digits, exp}.normInit()
}

func SetGimelPrecision(digits int64) {
	defaultPrecision = big.NewInt(digits)
}

func (g Gimel) normInit() Gimel {
	gl := len(g.digits.String())
	var a, b big.Int
	a.Sub(defaultPrecision, big.NewInt(int64(gl))).Int64()
	b.Exp(tenValue, &a, nil)
	g.digits.Mul(g.digits, &b)
	return g
}

func (g Gimel) normShift() Gimel {
	// if the sign is negative then set the negative flag and only store absolute values
	if g.digits.Sign() == -1 {
		g.neg = !g.neg
		g.digits.Abs(g.digits)
	}

	// get the length of digits
	gl := len(g.digits.String())
	var a, b big.Int

	// shift the exponent to match the precision
	a.Sub(defaultPrecision, big.NewInt(int64(gl)))
	g.exp.Sub(g.exp, &a)

	switch a.Sign() {
	case 1:
		// if the current digits are too short then multiply to line up
		b.Exp(tenValue, &a, nil)
		g.digits.Mul(g.digits, &b)
	case -1:
		// if the current digits are too short then divide to line up
		a.Abs(&a)
		b.Exp(tenValue, &a, nil)
		g.digits.Div(g.digits, &b)
	}
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

// maxMin outputs the maximum, minimum and a boolean defining if the values were swapped
func (g Gimel) maxMin(o Gimel) (Gimel, Gimel, bool) {
	if g.Gt(o) {
		return g.Clone(), o.Clone(), false
	} else {
		return o.Clone(), g.Clone(), true
	}
}

func (g Gimel) IsPos() bool { return !g.neg }
func (g Gimel) IsNeg() bool { return g.neg }

func (g Gimel) shiftToLineUpDigits(o Gimel) (d1, d2, exp *big.Int) {
	// find the difference between the exponents in max - min order
	m1, m2, swapped := g.maxMin(o)
	var a big.Int
	a.Sub(m1.exp, m2.exp)
	if a.CmpAbs(defaultPrecision) == 1 {
		return m1.digits, zeroValue, m1.exp
	}

	// make pow10 multiplier to shift bigger number left to align digits
	var a3 big.Int
	a3.Exp(tenValue, &a, nil)

	// perform pow10 multiply
	d1, d2 = new(big.Int), new(big.Int)
	d1.Mul(m1.digits, &a3)
	d2.Set(m2.digits)

	// flip digits to negative
	if m1.neg {
		d1.Neg(d1)
	}
	if m2.neg {
		d2.Neg(d2)
	}

	// swap the numbers back to the original order
	// this makes sure subtraction is called correctly
	if swapped {
		d1, d2 = d2, d1
	}

	// add to the exponent for normShift to calculate later
	exp = new(big.Int)
	exp.Set(m2.exp)
	return
}

func (g Gimel) Add(o Gimel) Gimel {
	var a big.Int
	d1, d2, exp := g.shiftToLineUpDigits(o)
	if d2.Sign() == 0 {
		return Gimel{false, d1, exp}.normShift()
	}
	a.Add(d1, d2)
	return Gimel{false, &a, exp}.normShift()
}

func (g Gimel) Sub(o Gimel) Gimel {
	return g.Add(o.Neg()) // yes this works
}

func (g Gimel) Mul(o Gimel) Gimel {
	// multiply the digits
	var a big.Int
	a.Mul(g.digits, o.digits)

	// sum exponents
	var b big.Int
	b.Add(g.exp, o.exp)

	// shift the exponent to account for the weird shift of the digits
	b.Sub(&b, defaultPrecision)
	b.Add(&b, oneValue)
	return Gimel{g.neg != o.neg, &a, &b}.normShift()
}

func (g Gimel) Div(o Gimel) Gimel {
	// TODO: fix division
	// it doesn't work due to the weird way I store decimal numbers as integers
	fmt.Println(g, "/", o)
	var a big.Int
	a.Div(g.digits, o.digits)
	var b big.Int
	b.Sub(g.exp, o.exp)
	fmt.Println(a, b)
	return Gimel{g.neg != o.neg, &a, &b}.normShift()
}

func (g Gimel) BigInt() *big.Int {
	if g.digits.Sign() == 0 {
		return big.NewInt(0)
	}
	var c big.Int
	c.Sub(g.exp, defaultPrecision)
	c.Add(&c, oneValue)
	var d big.Int
	d.Exp(tenValue, &c, nil)
	d.Mul(&d, g.digits)
	if g.neg {
		d.Neg(&d)
	}
	return &d
}

// String is just an alias for TextE for the Stringer interface
func (g Gimel) String() string { return g.TextE() }

func (g Gimel) TextE() string {
	var b strings.Builder
	if g.neg {
		b.WriteByte('-')
	}
	a := strings.TrimRight(g.digits.String(), "0")
	switch len(a) {
	case 0:
		return "0" // end early
	case 1:
		b.WriteByte(a[0])
	default:
		b.WriteByte(a[0])
		b.WriteByte('.')
		b.WriteString(a[1:])
	}
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
