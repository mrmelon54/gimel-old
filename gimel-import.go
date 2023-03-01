package gimel

import (
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
		Numeric: scanNumeric,
	}

	errInvalidDecimalDigit = fmt.Errorf("invalid decimal digit")
)

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
	return 0, errInvalidDecimalDigit
}

func scanDecimalDigitsLimit(r io.ByteScanner, b, p *big.Int) (err error) {
	var ch int

	for i := big.NewInt(0); i.Cmp(p) < 0; i.Add(i, oneValue) {
		if ch, err = scanDecimalDigit(r); err != nil {
			return err
		}
		b.Mul(b, tenValue)
		b.Add(b, big.NewInt(int64(ch)))
	}
	return
}

func scanNumeric(r byteRuneScanner, p *big.Int) (*Gimel, error) {
	neg, err := scanSign(r)
	if err != nil {
		return nil, err
	}

	var b big.Int
	if scanDecimalDigitsLimit(r, &b, p) != nil {
		return nil, err
	}

	return
}

func scanScientific(r io.RuneScanner, p *big.Int) (*Gimel, error) {

}
