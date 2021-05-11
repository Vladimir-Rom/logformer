package main

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestParseLogRecordRaw(t *testing.T) {
	r, err := regexp.Compile(`(?P<digits>\d+)(?s)(?P<letters>.*)`)
	if err != nil {
		t.Fatal(err)
	}

	multiLineString := strings.TrimSpace(`
a
b
c`)

	result := parseLogRecordRaw([]byte("123"+multiLineString), *r)

	if result == nil {
		t.Fatal("string was not parsed")
	}

	if d := result["digits"]; d != "123" {
		t.Errorf("Unexpected digits group value %s", d)
	}

	if l := strings.TrimSpace(result["letters"]); l != multiLineString {
		t.Errorf("Unexpected letters group value %s", l)
	}
}

func TestConstructValue(t *testing.T) {
	v, err := constructValue("string", variableDescriptor{name: "v1", varType: stringVarType})
	if err != nil {
		t.Error(err)
	} else if s := v.value.(string); s != "string" {
		t.Errorf("Unexpected constructed value: %s", s)
	}

	v, err = constructValue("true", variableDescriptor{name: "v2", varType: boolVarType})
	if err != nil {
		t.Error(err)
	} else if b := v.value.(bool); !b {
		t.Errorf("Unexpected constructed value: %v", b)
	}

	v, err = constructValue("123", variableDescriptor{name: "v3", varType: intVarType})
	if err != nil {
		t.Error(err)
	} else if i := v.value.(int); i != 123 {
		t.Errorf("Unexpected constructed value: %v", i)
	}

	v, err = constructValue(
		"2020-09-24 05:48:54.540",
		variableDescriptor{
			name:    "v4",
			varType: timeVarType,
			layout:  "2006-01-02 15:04:05.000"})
	if err != nil {
		t.Error(err)
	} else if tm := v.value.(time.Time); tm.Year() != 2020 || tm.Month() != 9 || tm.Hour() != 5 || tm.Nanosecond()/1000000 != 540 {
		t.Errorf("Unexpected constructed value: %v", tm)
	}
}

func TestParseLogs(t *testing.T) {
	logs := `
2020-09-24 05:48:54.540 first single line
2020-09-24 05:48:54.623 multi
line
2020-09-24 05:48:54.649 second single line
`
	secondLine := strings.TrimSpace(`
multi
line
`)
	parsingResult := make(chan map[string]varValue, 3)
	err := parseLogs([]byte(logs), createFormatDescriptor(), parsingResult)
	if err != nil {
		t.Fatal(err)
	}

	if l := len(parsingResult); l != 3 {
		t.Fatalf("Unexpected log items count: %v", l)
	}

	res := <-parsingResult

	if s := res["msg"].value.(string); strings.TrimSpace(s) != "first single line" {
		t.Errorf("Unexpected first message: %s", s)
	}

	if tm := res["time"].value.(time.Time); tm.Month() != 9 {
		t.Errorf("Wrong time: %v", tm)
	}

	res = <-parsingResult

	if s := res["msg"].value.(string); strings.TrimSpace(s) != secondLine {
		t.Errorf("Unexpected second message: %s", s)
	}
}

func createFormatDescriptor() formatDescriptor {
	delimeterPattern, _ := regexp.Compile(`\n(?:\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d)`)
	messagePattern, _ := regexp.Compile(`(?P<time>\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d)(?s)(?P<msg>.*)`)
	return formatDescriptor{
		recordDelimiterPattern: *delimeterPattern,
		recordPattern:          *messagePattern,
		variables: []variableDescriptor{
			{
				name:    "time",
				varType: timeVarType,
				layout:  "2006-01-02 15:04:05.000",
			},
			{
				name:    "msg",
				varType: stringVarType,
			},
		},
	}
}
