package xsd

import (
	"bytes"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

var msec = int64(1000000)

func TestDuration_String(t *testing.T) {
	assert.Equal(t, "PT1S", Duration{Seconds: 1}.String())
	assert.Equal(t, "PT0.11S", Duration{Nanoseconds: 110 * msec}.String())
	assert.Equal(t, "PT1M", Duration{Minutes: 1}.String())
	assert.Equal(t, "PT1M1S", Duration{Minutes: 1, Seconds: 1}.String())
	assert.Equal(t, "PT1M1.1S", Duration{Minutes: 1, Seconds: 1, Nanoseconds: 100 * msec}.String())
	assert.Equal(t, "PT1H", Duration{Hours: 1}.String())
	assert.Equal(t, "PT1H1M", Duration{Hours: 1, Minutes: 1}.String())
	assert.Equal(t, "PT1H1M1S", Duration{Hours: 1, Minutes: 1, Seconds: 1}.String())
	assert.Equal(t, "P1D", Duration{Days: 1}.String())
	assert.Equal(t, "P1DT1H1M1S", Duration{Days: 1, Hours: 1, Minutes: 1, Seconds: 1}.String())
	assert.Equal(t, "P1Y", Duration{Years: 1}.String())
	assert.Equal(t, "P1Y1DT1H1M1S", Duration{Years: 1, Days: 1, Hours: 1, Minutes: 1, Seconds: 1}.String())

	assert.Equal(t, "-PT1S", Duration{Seconds: 1, Negative: true}.String())
	assert.Equal(t, "-P1Y1DT1H1M1S", Duration{Years: 1, Days: 1, Hours: 1, Minutes: 1, Seconds: 1, Negative: true}.String())
}

type DurationFromStringType func(string) (*Duration, error)

func checkDurationFromString(t *testing.T, fromString DurationFromStringType, str string, expected *Duration) {
	actual, err := fromString(str)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func checkDurationWithMethod(t *testing.T, fromString DurationFromStringType) {
	checkDurationFromString(t, fromString, "PT0S", &Duration{Seconds: 0})
	checkDurationFromString(t, fromString, "PT1S", &Duration{Seconds: 1})
	checkDurationFromString(t, fromString, "PT0.11S", &Duration{Nanoseconds: 110 * msec})

	checkDurationFromString(t, fromString, "PT1M", &Duration{Minutes: 1})
	checkDurationFromString(t, fromString, "PT1M1S", &Duration{Minutes: 1, Seconds: 1})
	checkDurationFromString(t, fromString, "PT1M1.1S", &Duration{Minutes: 1, Seconds: 1, Nanoseconds: 100 * msec})
	checkDurationFromString(t, fromString, "PT1H", &Duration{Hours: 1})
	checkDurationFromString(t, fromString, "PT1H1M", &Duration{Hours: 1, Minutes: 1})
	checkDurationFromString(t, fromString, "PT1H1M1S", &Duration{Hours: 1, Minutes: 1, Seconds: 1})
	checkDurationFromString(t, fromString, "P1D", &Duration{Days: 1})
	checkDurationFromString(t, fromString, "P1DT1H1M1S", &Duration{Days: 1, Hours: 1, Minutes: 1, Seconds: 1})
	checkDurationFromString(t, fromString, "P1Y", &Duration{Years: 1})
	checkDurationFromString(t, fromString, "P1Y1DT1H1M1S", &Duration{Years: 1, Days: 1, Hours: 1, Minutes: 1, Seconds: 1})

	checkDurationFromString(t, fromString, "-PT1S", &Duration{Seconds: 1, Negative: true})
	checkDurationFromString(t, fromString, "-P1Y1DT1H1M1S", &Duration{Years: 1, Days: 1, Hours: 1, Minutes: 1, Seconds: 1, Negative: true})

	_, err := fromString("PT")
	assert.Equal(t, invalidFormatError, err)

	_, err = fromString("P1.")
	assert.Equal(t, invalidFormatError, err)

	_, err = fromString("PT1.")
	assert.Equal(t, invalidFormatError, err)

	_, err = fromString("PT1.S")
	assert.Equal(t, invalidFormatError, err)

	_, err = fromString("xPT1S")
	assert.Equal(t, invalidFormatError, err)

	_, err = fromString("PT1Sx")
	assert.Equal(t, invalidFormatError, err)
}

func TestDurationFromString(t *testing.T) {
	checkDurationWithMethod(t, DurationFromString)
}

type DurationAttr struct {
	Duration *Duration `xml:"duration,attr"`
}

func TestDuration_UnmarshalXMLAttr(t *testing.T) {
	dur := DurationAttr{}
	err := xml.Unmarshal([]byte(`<foo duration="PT1S"></foo>`), &dur)
	assert.Nil(t, err)
	assert.NotNil(t, dur.Duration)
	assert.Equal(t, Duration{Seconds: 1}, *dur.Duration)
}

func TestDuration_MarshalXMLAttr(t *testing.T) {
	val := Duration{Seconds: 2}
	dur := DurationAttr{Duration: &val}

	b := new(bytes.Buffer)
	e := xml.NewEncoder(b)
	err := e.Encode(dur)

	assert.Nil(t, err)
	assert.Equal(t, `<DurationAttr duration="PT2S"></DurationAttr>`, b.String())
}

func BenchmarkDuration_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Duration{Years: 1, Days: 1, Hours: 1, Minutes: 1, Seconds: 1}.String()
	}
}

func BenchmarkDurationFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = DurationFromString("P1Y2DT3H4M5.67S")
	}
}
