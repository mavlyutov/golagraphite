package golagraphite

import (
	"log"
	"net"
	"strconv"
	"strings"

	"os"

	"github.com/marpaia/graphite-golang"
)

func sendMetric(host string, port int, metrics []graphite.Metric) {
	for i, _ := range metrics {
		metrics[i].Name = replaceHostnameStub(metrics[i].Name)
	}
	Graphite, conn_err := graphite.NewGraphite(host, port)
	if conn_err != nil {
		log.Println(conn_err)
		return
	}
	defer Graphite.Disconnect()
	send_err := Graphite.SendMetrics(metrics)
	if send_err != nil {
		log.Println(send_err)
		return
	}
	for _, metric := range metrics {
		log.Printf("Sent metric '[%v]' to graphite '%v:%v'\n", metric, host, port)
	}
}

func SendMetrics(hosts []string, metrics []graphite.Metric) {
	for _, v := range hosts {
		host, port_string, _ := net.SplitHostPort(v)
		port, _ := strconv.Atoi(port_string)
		go sendMetric(host, port, metrics)
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