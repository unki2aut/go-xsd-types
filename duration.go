package xsd

import (
	"encoding/xml"
	"errors"
	"math"
)

type Duration struct {
	Years       int64
	Months      int64
	Days        int64
	Hours       int64
	Minutes     int64
	Seconds     int64
	Nanoseconds int64
	Negative    bool
}

var (
	invalidFormatError = errors.New("format of string no valid duration")
	errNoMonth         = errors.New("non-zero value for months is not allowed")
)

func (d Duration) String() string {
	var buf [32]byte
	w := len(buf)

	if d.Seconds > 0 || d.Nanoseconds > 0 {
		w--
		buf[w] = 'S'
	}

	if d.Nanoseconds > 0 {
		w = fmtNano(buf[:w], d.Nanoseconds)
		w--
		buf[w] = '.'

		if d.Seconds == 0 {
			w--
			buf[w] = '0'
		}
	}

	if d.Seconds > 0 {
		w = fmtInt(buf[:w], d.Seconds)
	}

	if d.Minutes > 0 {
		w--
		buf[w] = 'M'
		w = fmtInt(buf[:w], d.Minutes)
	}

	if d.Hours > 0 {
		w--
		buf[w] = 'H'
		w = fmtInt(buf[:w], d.Hours)
	}

	if w != len(buf) {
		w--
		buf[w] = 'T'
	}

	if d.Days > 0 {
		w--
		buf[w] = 'D'
		w = fmtInt(buf[:w], d.Days)
	}

	if d.Years > 0 {
		w--
		buf[w] = 'Y'
		w = fmtInt(buf[:w], d.Years)
	}

	w--
	buf[w] = 'P'

	if d.Negative {
		w--
		buf[w] = '-'
	}

	return string(buf[w:])
}

func fmtNano(buf []byte, v int64) int {
	prec := 9
	zeroes := true

	w := len(buf)
	if v == 0 {
		w++
	} else {
		for v > 0 && prec > 0 {
			digit := v % 10
			if digit != 0 && zeroes {
				zeroes = false
			}
			if !zeroes {
				w--
				buf[w] = byte(digit) + '0'
			}
			v /= 10
			prec--
		}
	}
	return w
}

// from time.go
func fmtInt(buf []byte, v int64) int {
	w := len(buf)
	if v == 0 {
		w++
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

// BenchmarkDurationFromString-8		 	 5000000	       241 ns/op
// BenchmarkDurationFromStringRegex-8   	 1000000	      1924 ns/op
func DurationFromString(str string) (*Duration, error) {
	var (
		dur = Duration{}
		buf = []byte(str)
		c   = len(buf) - 1

		timeParsed = false
	)

	parsePartial := func(indicator byte, field *int64) error {
		if buf[c] != indicator {
			return nil
		}

		c--

		i := lookupInt(&buf, c)
		if i == c+1 {
			return nil
		}

		val, err := atoi(buf[i : c+1])
		if err != nil {
			return err
		}

		*field = val
		c = i - 1

		timeParsed = true

		return nil
	}

	if buf[c] == 'S' {
		c--
		s := lookupInt(&buf, c)
		if s == c+1 {
			return nil, invalidFormatError
		}

		if buf[s-1] == '.' {
			nsVal, err := atoi(buf[s : c+1])
			if err != nil {
				return nil, err
			}
			remainingZeroes := 8 - (c - s)
			dur.Nanoseconds = nsVal * int64(math.Pow10(remainingZeroes))
			c = s - 2

			s = lookupInt(&buf, c)

			if s == c+1 {
				return nil, invalidFormatError
			}
		}

		sVal, err := atoi(buf[s : c+1])
		if err != nil {
			return nil, err
		}
		dur.Seconds = sVal
		c = s - 1

		timeParsed = true
	}

	if err := parsePartial('M', &dur.Minutes); err != nil {
		return nil, err
	}

	if err := parsePartial('H', &dur.Hours); err != nil {
		return nil, err
	}

	if timeParsed {
		if buf[c] != 'T' {
			return nil, invalidFormatError
		}

		c--
	}

	if err := parsePartial('D', &dur.Days); err != nil {
		return nil, err
	}

	if buf[c] == 'M' && buf[c+1] != '0' {
		return nil, errNoMonth
	}

	if err := parsePartial('Y', &dur.Years); err != nil {
		return nil, err
	}

	if buf[c] == 'P' {
		if c == 1 && buf[c-1] == '-' {
			dur.Negative = true
		} else if c != 0 {
			return nil, invalidFormatError
		}
	} else {
		return nil, invalidFormatError
	}

	return (*Duration)(&dur), nil
}

func atoi(buf []byte) (x int64, err error) {
	for i := 0; i < len(buf); i++ {
		c := buf[i]
		if x > (1<<63-1)/10 {
			// overflow
			return 0, invalidFormatError
		}
		x = x*10 + int64(c) - '0'
		if x < 0 {
			// overflow
			return 0, invalidFormatError
		}
	}

	return x, err
}

func lookupInt(buf *[]byte, i int) int {
	for ; i >= 0; i-- {
		c := (*buf)[i]
		if c < '0' || c > '9' {
			return i + 1
		}
	}

	return i
}

func (d *Duration) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if d == nil {
		return xml.Attr{}, nil
	}
	return xml.Attr{Name: name, Value: d.String()}, nil
}

func (d *Duration) UnmarshalXMLAttr(attr xml.Attr) error {
	dur, err := DurationFromString(attr.Value)
	if err != nil {
		return err
	}

	*d = *dur
	return nil
}

// check interfaces
var (
	dur                     = Duration{Seconds: 0}
	_   xml.MarshalerAttr   = &dur
	_   xml.UnmarshalerAttr = &dur
)
