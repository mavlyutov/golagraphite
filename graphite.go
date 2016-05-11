package golagraphite

import (
	"log"
	"net"
	"strconv"
	"strings"
	"os"
	"fmt"

	"github.com/marpaia/graphite-golang"
)

func sendMetricRoutine(host string, port int, metrics []graphite.Metric) {
	err := sendMetric(host, port, metrics)
	if err != nil {
		for _, metric := range metrics {
			log.Println(fmt.Sprintf("Unable to sent metric '[%v]' to graphite '%v:%v': %s", metric, host, port, err))
		}
	} else {
		for _, metric := range metrics {
			log.Printf("Sent metric '[%v]' to graphite '%v:%v'", metric, host, port)
		}
	}
}

func sendMetric(host string, port int, metrics []graphite.Metric) error {
	for i, _ := range metrics {
		metrics[i].Name = replaceHostnameStub(metrics[i].Name)
	}
	Graphite, conn_err := graphite.NewGraphite(host, port)
	if conn_err != nil {
		return conn_err
	}
	defer Graphite.Disconnect()
	send_err := Graphite.SendMetrics(metrics)
	if send_err != nil {
		return send_err
	}
	return nil
}

func SendMetrics(hosts []string, metrics []graphite.Metric) {
	for _, v := range hosts {
		host, port_string, _ := net.SplitHostPort(v)
		port, _ := strconv.Atoi(port_string)
		go sendMetricRoutine(host, port, metrics)
	}
}

func replaceHostnameStub(string_with_stub string) string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Cannot detect hostname")
		return string_with_stub
	}
	string_with_stub = strings.Replace(string_with_stub, `%hostname%`, NormalizeMetricName(hostname), -1)
	return string_with_stub
}

func NormalizeMetricName(rawName string) (normalizedName string) {

	normalizedName = strings.Replace(rawName, `.`, `_`, -1)
	normalizedName = strings.Replace(normalizedName, ` `, `_`, -1)
	normalizedName = strings.Replace(normalizedName, `\`, `.`, -1)

	normalizedName = strings.Replace(normalizedName, `:`, `.`, -1)
	normalizedName = strings.Replace(normalizedName, `/`, ``, -1)
	normalizedName = strings.Replace(normalizedName, `(`, ``, -1)
	normalizedName = strings.Replace(normalizedName, `)`, ``, -1)
	normalizedName = strings.Replace(normalizedName, `[`, ``, -1)
	normalizedName = strings.Replace(normalizedName, `]`, ``, -1)
	normalizedName = strings.Replace(normalizedName, `*`, ``, -1)
	normalizedName = strings.Replace(normalizedName, `%`, ``, -1)
	normalizedName = strings.Replace(normalizedName, `#`, ``, -1)
	normalizedName = strings.Replace(normalizedName, `-`, ``, -1)

	normalizedName = strings.ToLower(normalizedName)

	normalizedName = strings.TrimPrefix(normalizedName, "_")
	normalizedName = strings.TrimPrefix(normalizedName, ".")
	normalizedName = strings.TrimSuffix(normalizedName, "_")
	normalizedName = strings.TrimSuffix(normalizedName, ".")

	return
}
