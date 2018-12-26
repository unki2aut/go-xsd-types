package xsd

import (
	"bytes"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	msec = int64(1000000)
	sec  = 1000 * msec
	hour = 3600 * sec
	day  = 24 * hour
	year = 365 * day
)

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

func checkDurationFromString(t *testing.T, str string, expected *Duration) {
	actual, err := DurationFromString(str)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestDurationFromString(t *testing.T) {
	checkDurationFromString(t, "PT0S", &Duration{Seconds: 0})
	checkDurationFromString(t, "PT1S", &Duration{Seconds: 1})
	checkDurationFromString(t, "PT0.11S", &Duration{Nanoseconds: 110 * msec})

	checkDurationFromString(t, "PT1M", &Duration{Minutes: 1})
	checkDurationFromString(t, "PT1M1S", &Duration{Minutes: 1, Seconds: 1})
	checkDurationFromString(t, "PT1M1.1S", &Duration{Minutes: 1, Seconds: 1, Nanoseconds: 100 * msec})
	checkDurationFromString(t, "PT1H", &Duration{Hours: 1})
	checkDurationFromString(t, "PT1H1M", &Duration{Hours: 1, Minutes: 1})
	checkDurationFromString(t, "PT1H1M1S", &Duration{Hours: 1, Minutes: 1, Seconds: 1})
	checkDurationFromString(t, "P1D", &Duration{Days: 1})
	checkDurationFromString(t, "P1DT1H1M1S", &Duration{Days: 1, Hours: 1, Minutes: 1, Seconds: 1})
	checkDurationFromString(t, "P1Y", &Duration{Years: 1})
	checkDurationFromString(t, "P1Y1DT1H1M1S", &Duration{Years: 1, Days: 1, Hours: 1, Minutes: 1, Seconds: 1})

	checkDurationFromString(t, "-PT1S", &Duration{Seconds: 1, Negative: true})
	checkDurationFromString(t, "-P1Y1DT1H1M1S", &Duration{Years: 1, Days: 1, Hours: 1, Minutes: 1, Seconds: 1, Negative: true})

	_, err := DurationFromString("PT")
	assert.Equal(t, invalidFormatError, err)

	_, err = DurationFromString("P1.")
	assert.Equal(t, invalidFormatError, err)

	_, err = DurationFromString("PT1.")
	assert.Equal(t, invalidFormatError, err)

	_, err = DurationFromString("PT1.S")
	assert.Equal(t, invalidFormatError, err)

	_, err = DurationFromString("xPT1S")
	assert.Equal(t, invalidFormatError, err)

	_, err = DurationFromString("PT1Sx")
	assert.Equal(t, invalidFormatError, err)
}

func checkDurationToNanoseconds(t *testing.T, expected int64, dur *Duration) {
	actual, err := dur.ToNanoseconds()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestDuration_ToNanoseconds(t *testing.T) {
	checkDurationToNanoseconds(t, 0, &Duration{Seconds: 0})
	checkDurationToNanoseconds(t, sec, &Duration{Seconds: 1})
	checkDurationToNanoseconds(t, 110*msec, &Duration{Nanoseconds: 110 * msec})

	checkDurationToNanoseconds(t, 60*sec, &Duration{Minutes: 1})
	checkDurationToNanoseconds(t, 61*sec, &Duration{Minutes: 1, Seconds: 1})
	checkDurationToNanoseconds(t, 61*sec+100*msec, &Duration{Minutes: 1, Seconds: 1, Nanoseconds: 100 * msec})
	checkDurationToNanoseconds(t, hour, &Duration{Hours: 1})
	checkDurationToNanoseconds(t, hour+60*sec, &Duration{Hours: 1, Minutes: 1})
	checkDurationToNanoseconds(t, hour+61*sec, &Duration{Hours: 1, Minutes: 1, Seconds: 1})
	checkDurationToNanoseconds(t, day, &Duration{Days: 1})
	checkDurationToNanoseconds(t, day+hour+61*sec, &Duration{Days: 1, Hours: 1, Minutes: 1, Seconds: 1})
	checkDurationToNanoseconds(t, year, &Duration{Years: 1})
	checkDurationToNanoseconds(t, year+day+hour+61*sec, &Duration{Years: 1, Days: 1, Hours: 1, Minutes: 1, Seconds: 1})

	checkDurationToNanoseconds(t, -sec, &Duration{Seconds: 1, Negative: true})
	checkDurationToNanoseconds(t, -(year + day + hour + 61*sec), &Duration{Years: 1, Days: 1, Hours: 1, Minutes: 1, Seconds: 1, Negative: true})
}

func checkDurationToSeconds(t *testing.T, expected float64, dur *Duration) {
	actual, err := dur.ToSeconds()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestDuration_ToSeconds(t *testing.T) {
	var (
		hour = float64(3600)
		day  = 24 * hour
		year = 365 * day
	)

	checkDurationToSeconds(t, 1, &Duration{Seconds: 1})
	checkDurationToSeconds(t, 0.11, &Duration{Nanoseconds: 110 * msec})
	checkDurationToSeconds(t, year+day+hour+61.11, &Duration{Years: 1, Days: 1, Hours: 1, Minutes: 1, Seconds: 1, Nanoseconds: 110 * msec})

	checkDurationToSeconds(t, -1, &Duration{Seconds: 1, Negative: true})
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
