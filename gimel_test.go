package gimel

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var prec = big.NewInt(5)

func gen(neg bool, d, e int64) Gimel {
	return G(neg, big.NewInt(d), big.NewInt(e), prec)
}

func cmp(n bool, d, e int64, n2 bool, d2, e2 int64) int {
	return gen(n, d, e).Cmp(gen(n2, d2, e2))
}

func TestGimelConstantsE(t *testing.T) {
	assert.Equal(t, _EulerDigits[:1]+"."+_EulerDigits[1:]+"e0", Euler.String())
	assert.Equal(t, _PiDigits[:1]+"."+_PiDigits[1:]+"e0", Pi.String())
	assert.Equal(t, _Ln2Digits[:1]+"."+_Ln2Digits[1:]+"e0", Ln2.String())
}

func TestGimelConstantsNum(t *testing.T) {
	assert.Equal(t, _EulerDigits[:1]+"."+_EulerDigits[1:], Euler.Text(','))
	assert.Equal(t, _PiDigits[:1]+"."+_PiDigits[1:], Pi.Text(','))
	assert.Equal(t, _Ln2Digits[:1]+"."+_Ln2Digits[1:], Ln2.Text(','))
}

func TestGimel_Clone(t *testing.T) {
	g := gen(false, 1, 2)
	h := g.Clone()
	// verify big.Int pointers are different
	assert.False(t, g.digits == h.digits)
	assert.False(t, g.exp == h.exp)
}

func TestGimel_Abs(t *testing.T) {
	g := gen(true, 1, 2)
	assert.False(t, g.Abs().neg)
}

func TestGimel_Neg(t *testing.T) {
	g := gen(true, 1, 2)
	assert.False(t, g.Neg().neg)
	g = gen(false, 1, 2)
	assert.True(t, g.Neg().neg)
}

func TestGimel_Cmp(t *testing.T) {
	// test 1e1 and 1e2
	assert.Equal(t, 1, cmp(false, 1, 2, false, 1, 1))  // 1e2 > 1e1
	assert.Equal(t, -1, cmp(false, 1, 1, false, 1, 2)) // 1e1 < 1e2
	assert.Equal(t, 0, cmp(false, 1, 2, false, 1, 2))  // 1e2 == 1e2

	// test negative
	assert.Equal(t, 1, cmp(false, 1, 2, true, 1, 2))  // 1e2 > -1e2
	assert.Equal(t, -1, cmp(true, 1, 2, false, 1, 2)) // -1e2 < 1e2
	assert.Equal(t, 0, cmp(true, 1, 2, true, 1, 2))   // -1e2 == -1e2

	// test digit changes
	assert.Equal(t, 1, cmp(false, 5, 2, false, 1, 2))  // 5e2 > 1e2
	assert.Equal(t, -1, cmp(false, 1, 2, false, 5, 2)) // 1e2 < 5e2
	assert.Equal(t, 0, cmp(false, 5, 2, false, 5, 2))  // 5e2 == 5e2

	// test bigger digits
	assert.Equal(t, 1, cmp(false, 456, 2, false, 123, 2))  // 456e2 > 123e2
	assert.Equal(t, -1, cmp(false, 123, 2, false, 456, 2)) // 123e2 < 456e2
	assert.Equal(t, 0, cmp(false, 456, 2, false, 456, 2))  // 456e2 == 456e2
}

func TestGimel_Gt(t *testing.T) {
	assert.True(t, gen(false, 1, 2).Gt(gen(false, 1, 1)))  // 1e2 > 1e1
	assert.False(t, gen(false, 1, 2).Gt(gen(false, 1, 2))) // 1e2 == 1e1
	assert.False(t, gen(false, 1, 1).Gt(gen(false, 1, 2))) // 1e1 < 1e1
}

func TestGimel_Gte(t *testing.T) {
	assert.True(t, gen(false, 1, 2).Gte(gen(false, 1, 1)))  // 1e2 > 1e1
	assert.True(t, gen(false, 1, 2).Gte(gen(false, 1, 2)))  // 1e2 == 1e2
	assert.False(t, gen(false, 1, 1).Gte(gen(false, 1, 2))) // 1e1 < 1e1
}

func TestGimel_Lt(t *testing.T) {
	assert.True(t, gen(false, 1, 1).Lt(gen(false, 1, 2)))  // 1e1 < 1e1
	assert.False(t, gen(false, 1, 2).Lt(gen(false, 1, 2))) // 1e2 == 1e1
	assert.False(t, gen(false, 1, 2).Lt(gen(false, 1, 1))) // 1e2 > 1e1
}

func TestGimel_Lte(t *testing.T) {
	assert.True(t, gen(false, 1, 1).Lte(gen(false, 1, 2)))  // 1e1 < 1e1
	assert.True(t, gen(false, 1, 2).Lte(gen(false, 1, 2)))  // 1e2 == 1e1
	assert.False(t, gen(false, 1, 2).Lte(gen(false, 1, 1))) // 1e2 > 1e1
}

func TestGimel_Min(t *testing.T) {
	assert.Equal(t, gen(false, 1, 1), gen(false, 1, 1).Min(gen(false, 1, 2)))
	assert.Equal(t, gen(false, 1, 1), gen(false, 1, 2).Min(gen(false, 1, 1)))
	assert.Equal(t, gen(false, 1, 1), gen(false, 1, 1).Min(gen(false, 1, 1)))
}

func TestGimel_Max(t *testing.T) {
	assert.Equal(t, gen(false, 1, 2), gen(false, 1, 1).Max(gen(false, 1, 2)))
	assert.Equal(t, gen(false, 1, 2), gen(false, 1, 2).Max(gen(false, 1, 1)))
	assert.Equal(t, gen(false, 1, 1), gen(false, 1, 1).Max(gen(false, 1, 1)))
}

func TestGimel_IsPos(t *testing.T) {
	assert.True(t, gen(false, 1, 1).IsPos())
	assert.False(t, gen(true, 1, 1).IsPos())
}

func TestGimel_IsNeg(t *testing.T) {
	assert.True(t, gen(true, 1, 1).IsNeg())
	assert.False(t, gen(false, 1, 1).IsNeg())
}

func TestGimel_Add(t *testing.T) {
	assert.Equal(t, gen(false, 1, 10), gen(false, 1, 10).Add(gen(false, 1, 1)))
	assert.Equal(t, gen(false, 223, 10), gen(false, 123, 10).Add(gen(false, 1, 10)))
	assert.Equal(t, gen(false, 1, 11), gen(false, 5, 10).Add(gen(false, 5, 10)))
	assert.Equal(t, gen(false, 5, 10), gen(false, 1, 11).Add(gen(true, 5, 10)))
	assert.Equal(t, gen(false, 1, 10), gen(true, 1, 10).Add(gen(false, 2, 10)))
}

func TestGimel_Sub(t *testing.T) {
	assert.Equal(t, gen(false, 1, 10), gen(false, 1, 10).Sub(gen(false, 1, 1)))
	assert.Equal(t, gen(false, 123, 10), gen(false, 223, 10).Sub(gen(false, 1, 10)))
	assert.Equal(t, gen(false, 5, 10), gen(false, 1, 11).Sub(gen(false, 5, 10)))
	assert.Equal(t, gen(true, 3, 10), gen(true, 1, 10).Sub(gen(false, 2, 10)))
}

func TestGimel_Mul(t *testing.T) {
	assert.Equal(t, gen(false, 15, 16), gen(false, 3, 10).Mul(gen(false, 5, 5)))
	assert.Equal(t, gen(true, 182, 17), gen(false, 7, 10).Mul(gen(true, 26, 6)))
	assert.Equal(t, gen(false, 2, 100), gen(false, 1, 100).Mul(gen(false, 2, 0)))
}

func TestGimel_Div(t *testing.T) {
	assert.Equal(t, gen(false, 3, 10), gen(false, 15, 16).Div(gen(false, 5, 5)))
	assert.Equal(t, gen(true, 7, 10), gen(false, 182, 17).Div(gen(true, 26, 6)))
	assert.Equal(t, gen(false, 1, 100), gen(false, 2, 100).Div(gen(false, 2, 0)))
}

func TestGimel_Log10(t *testing.T) {
	assert.Equal(t, "1", gen(false, 10, 0).Log10().Text(0))
	assert.Equal(t, "2", gen(false, 100, 0).Log10().Text(0))
	assert.Equal(t, "3", gen(false, 1000, 0).Log10().Text(0))
}

func TestGimel_Exp(t *testing.T) {
	assert.Equal(t, gen(false, 1, 0), gen(false, 0, 0).Exp())
	assert.Equal(t, gen(false, Euler.digits.Int64(), 0), gen(false, 1, 0).Exp()) // there's a chance this will fail. It depends on how it handles e.
} // if this test succeeds, Pow() also works. Therefore, it does not need a test.
