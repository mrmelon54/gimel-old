package gimel

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
)

func init() {
	SetGimelPrecision(5) // makes tests easier to write
}

func gen(neg bool, d, e int64) Gimel {
	return Gimel{neg, big.NewInt(d), big.NewInt(e)}.Norm()
}

func cmp(n bool, d, e int64, n2 bool, d2, e2 int64) int {
	return gen(n, d, e).Cmp(gen(n2, d2, e2))
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
}

func TestGimel_Sub(t *testing.T) {
	assert.Equal(t, gen(false, 1, 10), gen(false, 1, 10).Sub(gen(false, 1, 1)))
}

func TestGimel_Mul(t *testing.T) {
}

func TestGimel_Div(t *testing.T) {
}

func TestGimel_TextE(t *testing.T) {
	assert.Equal(t, "1.23e6", gen(false, 123, 6).TextE())
	assert.Equal(t, "-3.456e15", gen(true, 3456, 15).TextE())
}

func TestGimel_Text(t *testing.T) {
	assert.Equal(t, "1230000", gen(false, 123, 6).Text(0))
	assert.Equal(t, "-3456"+strings.Repeat("0", 12), gen(true, 3456, 15).Text(0))

	assert.Equal(t, "1,230,000", gen(false, 123, 6).Text(','))
	assert.Equal(t, "-3,456"+strings.Repeat(",000", 4), gen(true, 3456, 15).Text(','))
	assert.Equal(t, "-456"+strings.Repeat(",000", 4), gen(true, 456, 14).Text(','))
	assert.Equal(t, "45,600"+strings.Repeat(",000", 4), gen(false, 456, 16).Text(','))
}
