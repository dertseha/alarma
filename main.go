package main

import (
	"time"

	"github.com/docopt/docopt-go"

	"github.com/dertseha/alarma/config"
	"github.com/dertseha/alarma/core"
	"github.com/dertseha/alarma/ui"
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
	parser := docopt.Parser{
		OptionsFirst:  true,
		SkipHelpFlags: false,
	}
	arguments, _ := parser.ParseArgs(usage(), nil, Title)
	configFilename := arguments["--config"].(string)

	if arguments["sampleconfig"].(bool) {
		cmdSampleConfig(configFilename)
	} else {
		cmdRun(configFilename)
	}
}

func startTicker(deferrer chan<- func(), update func(time.Time)) func() {
	ticker := time.NewTicker(100 * time.Millisecond)
	postUpdate := func(now time.Time) {
		deferrer <- func() { update(now) }
	}

	deferrer <- func() {
		go func() {
			for now := range ticker.C {
				postUpdate(now)
			}
		}()
	}

	return ticker.Stop
}

func cmdRun(configFilename string) {
	var configuration config.Instance
	runner := core.NewRunner()
	appWindow := ui.NewApplicationWindow(func(newConfiguration config.Instance) {
		_ = config.ToFile(configFilename, newConfiguration)
	})
	deferrer := make(chan func(), 100)
	tickerStopper := startTicker(deferrer, func(now time.Time) {
		newConfiguration, err := config.FromFile(configFilename)

		if err == nil {
			configuration = newConfiguration
		}
		runner.Update(configuration)
		appWindow.Update(now, configuration)
	})
	defer func() {
		tickerStopper()
		close(deferrer)
	}()

	appWindow.Show(deferrer)
}

func cmdSampleConfig(configFilename string) {
	var configuration config.Instance

	configuration.TimeSpansActive = true
	configuration.TimeSpans = []config.TimeSpan{
		{
			ID:      "sample-entry",
			Enabled: true,
			From:    "08:00",
			To:      "08:30",
			Path:    ".",
		},
	}
	_ = config.ToFile(configFilename, configuration)
}
