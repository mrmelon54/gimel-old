package gimel

import (
	"math/big"
)

type Sign bool

const (
	SignPositive Sign = true
	SignNegative Sign = false
)

type Gimel struct {
	sign   Sign
	digits *big.Int
	exp    *big.Int
}

func (g Gimel) Clone() Gimel {
	return Gimel{
		g.sign,
		(&big.Int{}).Set(g.digits),
		(&big.Int{}).Set(g.exp),
	}
}

func (g Gimel) Abs() Gimel {
	a := g.Clone()
	a.sign = SignPositive
	return a
}
