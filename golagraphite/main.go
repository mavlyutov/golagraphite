package main

import (
	"github.com/marpaia/graphite-golang"
	"github.yandex-team.ru/mavlyutov/golagraphite"
	"gopkg.in/alecthomas/kingpin.v2"
)

var config = kingpin.Flag("config", "Path to yaml formatted config file.").Short('c').Required().ExistingFile()

func main() {

	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	config := golagraphite.NewConfig(*config)

	metrics_channel := make(chan []graphite.Metric, 1024)  // buffer of 1024 messages to survive graphite connection outages

	go golagraphite.SendPerfCounters(config, metrics_channel)
	go golagraphite.SendSQLStatements(config, metrics_channel)

	for metrics := range metrics_channel {
		golagraphite.SendMetrics(config.Graphite_settings.Hosts, metrics)
	}

}
