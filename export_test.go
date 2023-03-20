package gimel

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
)

func TestGimel_BigInt(t *testing.T) {
	assert.Equal(t, big.NewInt(1230000), gen(false, 123, 6).BigInt())
	assert.Equal(t, big.NewInt(-3456000000), gen(true, 3456, 9).BigInt())
	assert.Equal(t, big.NewInt(0), gen(true, 0, 9).BigInt())
}

func TestGimel_TextE(t *testing.T) {
	assert.Equal(t, "1.23e6", gen(false, 123, 6).TextE())
	assert.Equal(t, "-3.456e15", gen(true, 3456, 15).TextE())
	assert.Equal(t, "-3e15", gen(true, 3, 15).TextE())
	assert.Equal(t, "0", gen(false, 0, 15).TextE())
	assert.Equal(t, "0", gen(true, 0, 15).TextE())
}

func TestGimel_Text(t *testing.T) {
	assert.Equal(t, "1230000", gen(false, 123, 6).Text(0))
	assert.Equal(t, "-3456"+strings.Repeat("0", 12), gen(true, 3456, 15).Text(0))

	assert.Equal(t, "1,230,000", gen(false, 123, 6).Text(','))
	assert.Equal(t, "-3,456,000,000,000,000", gen(true, 3456, 15).Text(','))
	assert.Equal(t, "-456,000,000,000,000", gen(true, 456, 14).Text(','))
	assert.Equal(t, "45,600,000,000,000,000", gen(false, 456, 16).Text(','))

	assert.Equal(t, "1,234.5", gen(false, 12345, 3).Text(','))
}
