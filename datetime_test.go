package xsd

import (
	"bytes"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDateTime(t *testing.T) {
	var (
		val *DateTime
		err error
	)

	val, err = DateTimeFromString("2015-09-07T05:45:54+01:00")
	assert.Nil(t, err)
	assert.Equal(t, "2015-09-07T05:45:54+01:00", val.String())

	val, err = DateTimeFromString("2015-09-07T05:45:54-01:00")
	assert.Nil(t, err)
	assert.Equal(t, "2015-09-07T05:45:54-01:00", val.String())

	val, err = DateTimeFromString("2015-09-07T05:45:54")
	assert.Nil(t, err)
	assert.Equal(t, "2015-09-07T05:45:54", val.String())

	val, err = DateTimeFromString("2015-09-07T05:45:54Z")
	assert.Nil(t, err)
	assert.Equal(t, "2015-09-07T05:45:54", val.String())
}

type DateTimeAttr struct {
	DateTime *DateTime `xml:"dateTime,attr"`
}

func TestDateTime_UnmarshalXMLAttr(t *testing.T) {
	dta := DateTimeAttr{}
	err := xml.Unmarshal([]byte(`<foo dateTime="2015-09-07T05:45:54+00:00"></foo>`), &dta)
	assert.Nil(t, err)
	assert.NotNil(t, dta.DateTime)
	assert.Equal(t, "2015-09-07T05:45:54", (*dta.DateTime).String())
}

func TestDateTime_MarshalXMLAttr(t *testing.T) {
	val, err := DateTimeFromString("2015-09-07T05:45:54")
	assert.Nil(t, err)
	dur := DateTimeAttr{DateTime: val}

	b := new(bytes.Buffer)
	e := xml.NewEncoder(b)
	err = e.Encode(dur)

	assert.Nil(t, err)
	assert.Equal(t, `<DateTimeAttr dateTime="2015-09-07T05:45:54"></DateTimeAttr>`, b.String())
}

func BenchmarkDateTimeFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DateTimeFromString("2017-08-16T13:07:00.09251+02:00")
	}
}

func BenchmarkDateTime_String(b *testing.B) {
	dateTime, _ := DateTimeFromString("2017-08-16T13:07:00.09251+02:00")

	for i := 0; i < b.N; i++ {
		_ = dateTime.String()
	}
}
