package golagraphite

import (
	"io/ioutil"

	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Graphite_settings    Graphite_settings
	Performance_counters Performance_counters
	Sql_metrics          Sql_metrics
}

type Graphite_settings struct {
	Udp   bool
	Hosts []string
}

type Performance_counters struct {
	Metric_prefix string
	Interval      int
	Counters      []string
}

type Sql_metrics struct {
	Sql_servers []Sql_server
}

type Sql_server struct {
	Connection_string string
	Queries           []Query
}

type Query struct {
	Interval      int
	Tsql_table    string
	Tsql_row      string
	Timestamp     string
	Metric_prefix string
}

func NewConfig(config_path string) (config Config) {

	source, readfile_err := ioutil.ReadFile(config_path)
	if readfile_err != nil {
		log.Fatal(readfile_err)
	}

	readyaml_err := yaml.Unmarshal(source, &config)
	if readyaml_err != nil {
		log.Fatal(readyaml_err)
	}

	return
}
