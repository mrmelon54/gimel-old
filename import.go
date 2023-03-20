package gimel

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strings"
)

type byteRuneScanner interface {
	io.ByteScanner
	io.RuneScanner
}

type scannerCallback func(r byteRuneScanner, p *big.Int) (*Gimel, error)

type formatDetectList []struct {
	Format
	*regexp.Regexp
}

type Format uint

const (
	Auto Format = iota
	Numeric
	Scientific
)

var (
	// formatDetect contains a format, regex pair to autodetect formats
	formatDetect = formatDetectList{
		{Numeric, regexp.MustCompile(`[+-]?\d+`)},
		{Scientific, regexp.MustCompile(`[+-]?\d+e\d+`)},
	}
	// formatMap contains a map between Format and scannerCallback functions
	formatMap = map[Format]scannerCallback{
		Numeric:    scanNumeric,
		Scientific: scanScientific,
	}

	errInvalidDecimalDigit       = fmt.Errorf("invalid decimal digit")
	errInvalidScientificNotation = fmt.Errorf("invalid scientific notation")
)

// FromBigInt returns the Gimel number from a big.Int with a precision
func FromBigInt(a *big.Int, prec *big.Int) (Gimel, bool) {
	fmt.Println(a.String())
	return FromString(a.String(), Numeric, prec)
}

// FromString returns the Gimal number from a string, Format and precision
func FromString(s string, f Format, prec *big.Int) (Gimel, bool) {
	if f == Auto {
		for _, i := range formatDetect {
			if i.MatchString(s) {
				f = i.Format
				break
			}
		}
		if f == Auto {
			return Gimel{}, false
		}
	}
	if fn, ok := formatMap[f]; ok {
		return setFromScanner(strings.NewReader(s), prec, fn)
	}
	return Gimel{}, false
}

// setFromScanner implements FromString given an io.ByteScanner
func setFromScanner(r byteRuneScanner, prec *big.Int, scan scannerCallback) (Gimel, bool) {
	g, err := scan(r, prec)
	if err != nil {
		return Gimel{}, false
	}
	// entire content must have been consumed
	if _, err := r.ReadByte(); err != io.EOF {
		return Gimel{}, false
	}
	return *g, true // err == io.EOF => scan consumed all content of r
}

func scanSign(r io.ByteScanner) (neg bool, err error) {
	var ch byte
	if ch, err = r.ReadByte(); err != nil {
		return false, err
	}
	switch ch {
	case '-':
		neg = true
	case '+':
		// do nothing
	default:
		_ = r.UnreadByte()
	}
	return
}

func scanDecimalDigit(r io.ByteScanner) (n int, err error) {
	var ch byte
	if ch, err = r.ReadByte(); err != nil {
		return 0, err
	}
	if ch >= '0' && ch <= '9' {
		return int(ch - '0'), nil
	}
	_ = r.UnreadByte()
	return 0, errInvalidDecimalDigit
}

func scanDecimalDigitsLimit(r io.ByteScanner, b, p *big.Int) (n *big.Int, err error) {
	push := b != nil
	n = big.NewInt(0)
	if p == nil {
		for {
			if err = scanDecimalDigitAppender(r, b, push); err != nil {
				if errors.Is(err, errInvalidDecimalDigit) || errors.Is(err, io.EOF) {
					return n, nil
				}
				return new(big.Int).Set(zeroValue), err
			}
			n.Add(n, oneValue)
		}
	} else {
		for ; n.Cmp(p) < 0; n.Add(n, oneValue) {
			if err = scanDecimalDigitAppender(r, b, push); err != nil {
				if errors.Is(err, errInvalidDecimalDigit) || errors.Is(err, io.EOF) {
					return n, nil
				}
				return new(big.Int).Set(zeroValue), err
			}
		}
	}
	return
}

func scanDecimalDigitAppender(r io.ByteScanner, b *big.Int, push bool) error {
	ch, err := scanDecimalDigit(r)
	if err != nil {
		return err
	}
	if push {
		b.Mul(b, tenValue)
		b.Add(b, big.NewInt(int64(ch)))
	}
	return nil
}

func scanNumeric(r byteRuneScanner, p *big.Int) (*Gimel, error) {
	neg, err := scanSign(r)
	if err != nil {
		return nil, err
	}

	var b big.Int
	_, err = scanDecimalDigitsLimit(r, &b, p)
	if err != nil {
		return nil, err
	}

	n, err := scanDecimalDigitsLimit(r, nil, nil)
	if err != nil {
		return nil, err
	}

	g := G(neg, &b, n, p)
	return &g, nil
}

func scanScientific(r byteRuneScanner, p *big.Int) (*Gimel, error) {
	neg, err := scanSign(r)
	if err != nil {
		return nil, err
	}

	var b big.Int
	_, err = scanDecimalDigitsLimit(r, &b, p)
	if err != nil {
		return nil, err
	}

	ch, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if ch != 'e' {
		return nil, errInvalidScientificNotation
	}

	var b2 big.Int
	_, err = scanDecimalDigitsLimit(r, &b, nil)
	if err != nil {
		return nil, err
	}

	g := G(neg, &b, &b2, p)
	return &g, nil
}
