package main

import (
	"github.com/jessevdk/go-flags"
)

type cmdLineOption struct {
	Format string `short:"f" long:"format" description:"path to yaml file with log format description"`
}

func parseCommandLine() (cmdLineOption, error) {
	var result cmdLineOption
	_, err := flags.Parse(&result)
	return result, err
}
