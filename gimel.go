package gimel

import (
	"math/big"
)

func strToBigInt(s string) *big.Int {
	var a big.Int
	_, ok := a.SetString(s, 10)
	if !ok {
		panic("strToBigInt failed")
	}
	return &a
}

var ( // constants
	zeroValue = big.NewInt(0)
	oneValue  = big.NewInt(1)
	tenValue  = big.NewInt(10)
	// Gimel Constants
	euler = G(false, strToBigInt("2718281828459045235360287471352662497757247093699959574966967627724076630353547594571382178525166427"), big.NewInt(-100), big.NewInt(100))
	pi    = G(false, strToBigInt("3141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117067"), big.NewInt(-100), big.NewInt(100))
	ln2   = G(false, strToBigInt("6931471805599453094172321214581765680755001343602552541206800094933936219696947156058633269964186875"), big.NewInt(-100), big.NewInt(100))
)

type Gimel struct {
	neg    bool
	digits *big.Int
	exp    *big.Int
	prec   *big.Int
	p10p   *big.Int
}

// G returns a normalised version of the Gimel struct
// Gimel{false, 123, 10} will be converted to Gimel{false, 12300, 10} with a precision value of 5
func G(neg bool, digits, exp, prec *big.Int) Gimel {
	var p, p2 big.Int
	p.Set(prec)
	p2.Exp(tenValue, prec, nil)
	return Gimel{neg, digits, exp, &p, &p2}.normPrec()
}

// g2 is an internal function to return the Gimel struct with cloned precision values
func g2(neg bool, digits, exp, prec *big.Int) Gimel {
	var p, p2 big.Int
	p.Set(prec)
	p2.Exp(tenValue, prec, nil)
	return Gimel{neg, digits, exp, &p, &p2}
}

// minBigInt is an internal function to get the minimum big int value
func minBigInt(a, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return a
	}
	return b
}

// normPrec is an internal function to return the normalised version of the Gimel struct
// Gimel{false, 123, 10} will be converted to Gimel{false, 12300, 10} with a precision value of 5
func (g Gimel) normPrec() Gimel {
	gl := len(g.digits.String())
	var a, b big.Int
	a.Sub(g.prec, big.NewInt(int64(gl))).Int64()

	switch a.Sign() {
	case 1:
		// if the current digits are too short then multiply to line up
		b.Exp(tenValue, &a, nil)
		g.digits.Mul(g.digits, &b)
	case -1:
		// if the current digits are too short then devide to line up
		a.Abs(&a)
		b.Exp(tenValue, &a, nil)
		g.digits.Div(g.digits, &b)
	}
	return g
}

// normShift is an internal function to return the normalised version of the Gimel struct
// this is equivalent to normPrec but also shifts the exponent the same amount as the digits
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
	a.Sub(g.prec, big.NewInt(int64(gl)))
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

// Norm returns the normalised version of the Gimel struct
// this is equivalent to normPrec but also shifts the exponent the same amount as the digits
func (g Gimel) Norm() Gimel {
	return g.normShift().Clone()
}

// Precision returns a new Gimel struct with a different precision value
// normPrec is called after to retain the normalised Gimel struct
func (g Gimel) Precision(prec *big.Int) Gimel {
	g.prec = new(big.Int).Set(prec)
	g.p10p = new(big.Int).Exp(tenValue, prec, nil)
	return g.normPrec()
}

// Clone returns a clone of the Gimel struct
func (g Gimel) Clone() Gimel {
	return Gimel{
		g.neg,
		(&big.Int{}).Set(g.digits),
		(&big.Int{}).Set(g.exp),
		new(big.Int).Set(g.prec),
		new(big.Int).Set(g.p10p),
	}
}

// Abs returns a clone of the Gimel struct with a positive sign
func (g Gimel) Abs() Gimel {
	a := g.Clone()
	a.neg = false
	return a
}

// Neg returns a clone of the Gimel struct with an inverted sign
func (g Gimel) Neg() Gimel {
	a := g.Clone()
	a.neg = !g.neg
	return a
}

// Cmp returns:
//
// -1 if g <  o
//
//	0 if g == o
//
// +1 if g >  o
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

// Gt is an alias for g > o
func (g Gimel) Gt(o Gimel) bool { return g.Cmp(o) == 1 }

// Gte is an alias for g >= o
func (g Gimel) Gte(o Gimel) bool { return g.Cmp(o) != -1 }

// Lt is an alias for g < o
func (g Gimel) Lt(o Gimel) bool { return g.Cmp(o) == -1 }

// Lte is an alias for g <= o
func (g Gimel) Lte(o Gimel) bool { return g.Cmp(o) != 1 }

// Eq is an alias for g == o
func (g Gimel) Eq(o Gimel) bool { return g.Cmp(o) == 0 }

// Neq is an alias for g != o
func (g Gimel) Neq(o Gimel) bool { return g.Cmp(o) != 0 }

// Min returns a clone of the minimum value
func (g Gimel) Min(o Gimel) Gimel {
	if g.Lt(o) {
		return g.Clone()
	} else {
		return o.Clone()
	}
}

// Max returns a clone of the maximum value
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

// IsPos returns true if the sign is positive
func (g Gimel) IsPos() bool { return !g.neg }

// IsNeg returns true if the sign is negative
func (g Gimel) IsNeg() bool { return g.neg }

// shiftToLineUpDigits is an internal function to shift the digits to line up for add/subtract operations
func (g Gimel) shiftToLineUpDigits(o Gimel) (d1, d2, exp, prec *big.Int) {
	prec = new(big.Int).Set(minBigInt(g.prec, o.prec))

	// find the difference between the exponents in max - min order
	m1, m2, swapped := g.maxMin(o)
	var a big.Int
	a.Sub(m1.exp, m2.exp)
	if a.CmpAbs(prec) == 1 {
		d1 = m1.digits
		d2 = zeroValue
		exp = m1.exp
		return
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

// Add returns the sum g+o
func (g Gimel) Add(o Gimel) Gimel {
	var a big.Int
	d1, d2, exp, prec := g.shiftToLineUpDigits(o)
	if d2.Sign() == 0 {
		return g2(false, d1, exp, prec).normShift()
	}
	a.Add(d1, d2)
	return g2(false, &a, exp, prec).normShift()
}

// Sub returns the difference g-o
func (g Gimel) Sub(o Gimel) Gimel {
	return g.Add(o.Neg()) // yes this works
}

// Mul returns the product g*o
func (g Gimel) Mul(o Gimel) Gimel {
	prec := new(big.Int).Set(minBigInt(g.prec, o.prec))

	// multiply the digits
	var a big.Int
	a.Mul(g.digits, o.digits)

	// sum exponents
	var b big.Int
	b.Add(g.exp, o.exp)

	// shift the exponent to account for the weird shift of the digits
	b.Sub(&b, prec)
	b.Add(&b, oneValue)
	return g2(g.neg != o.neg, &a, &b, prec).normShift()
}

// Div returns the quotient g/o
func (g Gimel) Div(o Gimel) Gimel {
	prec := new(big.Int).Set(minBigInt(g.prec, o.prec))
	p10p := new(big.Int).Set(minBigInt(g.p10p, o.p10p))

	// multiply bigger number by 10^prec to give space for full integer division
	var a big.Int
	a.Mul(g.digits, p10p)
	a.Div(&a, o.digits)

	// subtract the exponents
	var b big.Int
	b.Sub(g.exp, o.exp)
	b.Sub(&b, oneValue)
	return g2(g.neg != o.neg, &a, &b, prec).normShift()
}

// Log10 Returns the natural logarithm. (log base e)
// This uses the Taylor series expansion of ln(x). The precision of the result is the same as the precision of the input, with a max of the precision of Euler's number.
func (g Gimel) Ln() Gimel {
	if g.neg {
		panic("Cannot take ln of negative Gimel number")
	}
	// check precision is greater than e
	var epsilon *big.Int
	if g.prec.Cmp(euler.prec) == 1 {
		epsilon = new(big.Int).Set(euler.prec)
	} else {
		epsilon = new(big.Int).Set(g.prec)
	}
	var (
		y = new(big.Int).Set(g.digits)
		x = new(big.Int).Div(new(big.Int).Sub(y, oneValue), new(big.Int).Add(y, oneValue))
		z = new(big.Int).Mul(x, x)
		L = new(big.Int).Set(zeroValue)
	)
	for k := 1; x.Cmp(epsilon) == 1; k += 2 {
		t := new(big.Int).Div(new(big.Int).Mul(new(big.Int).Set(2), x), k)
		L.Add(L, t)
		x.Mul(x, z)
	}
	return g // TODO: @MrMelon54
}

// BigInt returns the big.Int representing the full Gimel number
func (g Gimel) BigInt() *big.Int {
	if g.digits.Sign() == 0 {
		return big.NewInt(0)
	}
	var c big.Int
	c.Sub(g.exp, g.prec)
	c.Add(&c, oneValue)
	var d big.Int
	d.Exp(tenValue, &c, nil)
	d.Mul(&d, g.digits)
	if g.neg {
		d.Neg(&d)
	}
	return &d
}

// Ln returns the natural logarithm. (log base e)
//
// This uses the Taylor series expansion of ln(x).
//
// The precision of the result is the same as the precision of the input,
// with a max of the precision of Euler's number.
func (g Gimel) Ln() Gimel {
	if g.neg {
		panic("Cannot take ln of negative Gimel number")
	}
	// check precision is greater than e
	var epsilon *big.Int
	if g.prec.Cmp(euler.prec) == 1 {
		epsilon = new(big.Int).Set(euler.prec)
	} else {
		epsilon = new(big.Int).Set(g.prec)
	}
	var (
		y = new(big.Int).Set(g.digits)
		x = new(big.Int).Div(new(big.Int).Sub(y, oneValue), new(big.Int).Add(y, oneValue))
		z = new(big.Int).Mul(x, x)
		L = new(big.Int).Set(zeroValue)
	)
	for k := 1; x.Cmp(epsilon) == 1; k += 2 {
		t := new(big.Int).Div(new(big.Int).Mul(new(big.Int).Set(2), x), k)
		L.Add(L, t)
		x.Mul(x, z)
	}

	a, ok := FromBigInt(L, epsilon)
	if !ok {
		panic("failed to parse big int")
	}
	return a
}
