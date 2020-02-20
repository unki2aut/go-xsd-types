package xsd

import (
	"encoding/xml"
	"time"
)

type DateTime time.Time

const (
	DateTimeFormat     = "2006-01-02T15:04:05.999999999-07:00"
	DateTimeNoTimezone = "2006-01-02T15:04:05.999999999"
)

func DateTimeFromString(str string) (*DateTime, error) {
	var (
		val time.Time
		err error
		z   = str[len(str)-6]
	)

	if z == '+' || z == '-' {
		val, err = time.Parse(DateTimeFormat, str)
	} else {
		val, err = time.Parse(DateTimeNoTimezone, str)
	}

	if err != nil {
		return nil, err
	}

	return (*DateTime)(&val), nil
}

func (dt *DateTime) String() string {
	_, offset := (*time.Time)(dt).Zone()
	if offset != 0 {
		return (*time.Time)(dt).Format(DateTimeFormat)
	} else {
		return (*time.Time)(dt).Format(DateTimeNoTimezone)
	}
}

func (dt *DateTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if dt == nil {
		return xml.Attr{}, nil
	}
	return xml.Attr{Name: name, Value: dt.String()}, nil
}

func (dt *DateTime) UnmarshalXMLAttr(attr xml.Attr) error {
	val, err := DateTimeFromString(attr.Value)
	if err != nil {
		return err
	}

	*dt = *val
	return nil
}

// check interfaces
var (
	dt                     = DateTime{}
	_  xml.MarshalerAttr   = &dt
	_  xml.UnmarshalerAttr = &dt
)
