package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	cmd, err := parseCommandLine()
	if err != nil {
		os.Exit(-1)
	}

	formatBytes, err := ioutil.ReadFile(cmd.Format)

	if err != nil {
		panic(err)
	}

	format, err := getFormatDescriptor(formatBytes)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Format: %s", format.recordDelimiterPattern)
}
