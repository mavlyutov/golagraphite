// +build linux darwin

package golagraphite

import (
	"log"

	"github.com/marpaia/graphite-golang"
)

func SendPerfCounters(c Config, metrics_channel chan []graphite.Metric) {
	log.Println("PerfCounters are not available at unix host systems")
	return
}
