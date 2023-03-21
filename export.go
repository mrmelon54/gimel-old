package gimel

import (
	"fmt"
	"math/big"
	"strings"
)

// BigInt returns the big.Int representing the full Gimel number
func (g Gimel) BigInt() *big.Int {
	var a big.Int
	g.BigFloat().Int(&a)
	return &a
}

// BigFloat returns the big.Float representing the full Gimel number
func (g Gimel) BigFloat() *big.Float {
	// shortcut for 0^±x
	if g.digits.Sign() == 0 {
		return big.NewFloat(0)
	}

	// shortcut for ±x^0
	if g.exp.Sign() == 0 {
		if g.neg {
			return big.NewFloat(-1)
		}
		return big.NewFloat(1)
	}

	shiftDigits := g.getFloatShiftedDigits()
	fmt.Println("shiftDigits:", shiftDigits.String(), "orig:", g.digits.String(), "prec:", g.prec)
	{
		// shift by exponent
		var c big.Int
		c.Exp(tenValue, g.exp, nil)
		var d big.Float
		d.SetInt(&c)
		fmt.Println(d.String())
		switch d.Sign() {
		case 1:
			shiftDigits.Mul(shiftDigits, &d)
		case -1:
			shiftDigits.Mul(shiftDigits, &d)
		}
	}
	if g.neg {
		shiftDigits.Neg(shiftDigits)
	}
	return shiftDigits
}

// getFloatShiftedDigits converts to digits 12300 (5 prec) format => 1.23 (as *big.Float)
func (g Gimel) getFloatShiftedDigits() *big.Float {
	var (
		shiftDigits, d, e big.Float
		c1, c2            big.Int
	)
	c1.Set(g.prec)
	c1.Sub(&c1, oneValue)
	c2.Exp(tenValue, &c1, nil)
	d.SetInt(&c2)
	e.SetInt(g.digits)
	shiftDigits.Quo(&e, &d)
	return &shiftDigits
}

// String is just an alias for TextE for the Stringer interface
func (g Gimel) String() string { return g.TextE() }

// TextE returns the scientific representation of the Gimel number
// For example: 1.23e15
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

// Text returns the decimal text representation of the Gimel number
// If the sep parameter is set to 0 then no separator is used
// For example: 1,230,000,000,000,000
func (g Gimel) Text(sep rune) string {
	var b strings.Builder
	if g.neg {
		b.WriteByte('-')
	}

	if sep == 0 {
		g.writeFullDigits(&b)
	} else {
		var b2 strings.Builder
		dec := g.writeFullDigits(&b2)
		a := b2.String()
		l := int64(len(a))

		// decimal point check
		var dec2 int64
		if dec != nil {
			dec2 = dec.Int64()
			l += dec2
		}

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
		if dec2 < 0 {
			b.WriteRune('.')
			b.WriteString(a[l:])
		}
	}
	return b.String()
}

// writeFullDigits is an internal function to write the full digits of a Gimel number
// this is equivalent to running Text(0) but missing the sign
func (g Gimel) writeFullDigits(b *strings.Builder) *big.Int {
	var c big.Int
	c.Sub(g.exp, g.prec)
	c.Add(&c, oneValue)
	if c.Sign() == -1 {
		b.WriteString(g.digits.String())
		return &c
	}
	b.WriteString(g.digits.String())
	for i := new(big.Int); i.Cmp(&c) < 0; i.Add(i, oneValue) {
		b.WriteByte('0')
	}
	return nil
}
