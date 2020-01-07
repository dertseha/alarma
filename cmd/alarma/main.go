package main

import (
	"flag"
	"time"

	"github.com/dertseha/alarma/internal/config"
	"github.com/dertseha/alarma/internal/core"
	"github.com/dertseha/alarma/internal/ui"
)

func main() {
	sampleConfig := flag.Bool("sampleconfig", false, "whether the app should print a sample config and exit")
	configFilename := flag.String("config", "config.json", "the configuration file to use")
	flag.Parse()

	if *sampleConfig {
		cmdSampleConfig(*configFilename)
	} else {
		cmdRun(*configFilename)
	}
}

func startTicker(update func(time.Time)) func() {
	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for now := range ticker.C {
			update(now)
		}
	}()

	return ticker.Stop
}

func cmdRun(configFilename string) {
	var configuration config.Instance
	runner := core.NewRunner()
	appWindow := ui.NewApplicationWindow(func(newConfiguration config.Instance) {
		_ = config.ToFile(configFilename, newConfiguration)
	})
	tickerStopper := startTicker(func(now time.Time) {
		newConfiguration, err := config.FromFile(configFilename)

		if err == nil {
			configuration = newConfiguration
		}
		runner.Update(configuration)
		appWindow.Update(now, configuration)
	})
	defer func() {
		tickerStopper()
	}()

	appWindow.Show()
}

func cmdSampleConfig(configFilename string) {
	configuration := config.Example()
	_ = config.ToFile(configFilename, configuration)
}
