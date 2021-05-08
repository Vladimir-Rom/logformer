package main

import "testing"

func TestGetFormatDescriptor(t *testing.T) {
	var yaml = `
recordDelimiterPattern: \n
recordPattern: "(?P<var1>.*)(?P<var2>.*)"
variables:
  - name: var1
    type: string
  - name: var2
    type: bool
`
	format, err := getFormatDescriptor([]byte(yaml))
	if err != nil {
		t.Fatal(err)
	}

	if varsCount := len(format.variables); varsCount != 2 {
		t.Errorf("Unexpected variables count: %v", varsCount)
	}

	if format.variables[0].name != "var1" {
		t.Errorf("Unexpected variable name: %s", format.variables[0].name)
	}

	if format.variables[0].varType != stringVarType {
		t.Errorf("Unexpected variable type: %v", format.variables[0].varType)
	}
}
