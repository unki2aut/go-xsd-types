# go-xsd-types
XML Schema Definition Primitive Types for GO

## Usage
```go
type MyStruct struct {
	Duration *Duration `xml:"duration,attr"`
	DateTime *DateTime `xml:"dateTime,attr"`
}

val := MyStruct{}

err := xml.Unmarshal([]byte(`<node duration="PT1S">
  <entry dateTime="2015-09-07T05:45:54"/>
</node>`), &val)

fmt.Println(val.Duration, val.DateTime)
```
