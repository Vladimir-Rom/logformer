package main

import (
	"io/ioutil"
	"regexp"
	"testing"
)

func TestParallelParseLogs(t *testing.T) {
	delimeterPattern, _ := regexp.Compile(`(?:\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d)`)
	messagePattern, _ := regexp.Compile(`(?P<time>\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d\d\d)\s+(?:[^ ]+)\s+(?P<threadId>\d+)\s+(?:[^ ]+)\s+(?:[^ ]+)(?s)(?P<msg>.*)`)
	format := formatDescriptor{
		recordDelimiterPattern: delimeterPattern,
		recordPattern:          messagePattern,
		variables: []variableDescriptor{
			{
				name:    "time",
				varType: timeVarType,
				layout:  "2006-01-02 15:04:05.000",
			},
			{
				name:    "threadId",
				varType: intVarType,
			},
			{
				name:    "msg",
				varType: stringVarType,
			},
		},
	}

	logs, err := ioutil.ReadFile("testdata/logfile.txt")

	if err != nil {
		t.Fatal(err)
	}

	result, err := parallelParseLogs(4, logs, format)

	if err != nil {
		t.Fatal(err)
	}

	resultLen := len(result)

	t.Logf("Log items count: %v", resultLen)
}

func TestSplitLogs(t *testing.T) {
	logs := "aaaa---bbb---ccc---ddd"
	delimiter := regexp.MustCompile("---")
	result := splitLogs([]byte(logs), 4, delimiter)
	if l := len(result); l != 3 {
		t.Errorf("Unexpected result len: %v", l)
	}

	for i, chunk := range result {
		if i == 0 {
			continue
		}

		if logs[chunk] != '-' || logs[chunk-1] == '-' {
			t.Errorf("invalid chunk #%v: %v", i, chunk)
		}

		if result[i-1] >= chunk {
			t.Errorf("Chunk was not increased. Current chunk: %v, previous: %v", chunk, result[i-1])
		}
	}
}
