package main

import (
	"fmt"
	"time"

	"github.com/docopt/docopt-go"

	"github.com/dertseha/alarma/config"
	"github.com/dertseha/alarma/core"
)

const (
	// Version contains the current version number
	Version = "0.0.1"
	// Name is the name of the application
	Name = "alarma"
	// Title contains a combined string of name and version
	Title = Name + " v" + Version
)

func usage() string {
	return Title + `

Usage:
   alarma --config=<filename>
   alarma sampleconfig --config=<filename>
   alarma -h | --help
   alarma --version

Options:
   -h --help              Show this screen.
   --version              Show version.
   --config=<filename>    Filename of the configuration file
`
}

func main() {
	arguments, _ := docopt.Parse(usage(), nil, true, Title, false)
	configFilename := arguments["--config"].(string)

	if arguments["sampleconfig"].(bool) {
		cmdSampleConfig(configFilename)
	} else {
		cmdRun(configFilename)
	}
}

func cmdRun(configFilename string) {
	var configuration config.Instance
	runner := core.NewRunner()

	for true {
		newConfiguration, err := config.FromFile(configFilename)

		if err == nil {
			configuration = newConfiguration
		} else {
			fmt.Printf("\r failed config: %v", err)
		}
		runner.Update(configuration)
		time.Sleep(time.Millisecond * 100)
	}
}

func cmdSampleConfig(configFilename string) {
	var configuration config.Instance

	configuration.TimeSpansActive = true
	configuration.TimeSpans = []config.TimeSpan{config.TimeSpan{
		ID:      "sample-entry",
		Enabled: true,
		From:    "08:00",
		To:      "08:30",
		Path:    "."}}
	config.ToFile(configFilename, configuration)
}
