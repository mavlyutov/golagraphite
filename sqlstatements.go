package golagraphite

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"errors"

	"github.com/marpaia/graphite-golang"
	_ "github.com/denisenkom/go-mssqldb"
)

const (
	keyValueType string = "key-value"
	singleRowType string = "single-row"
)

func SendSQLStatements(config Config, metrics_channel chan []graphite.Metric) {
	for _, server := range config.Sql_metrics.Sql_servers {
		go SendSQLStatementsPerServer(server, metrics_channel)
	}
}

func SendSQLStatementsPerServer(server Sql_server, metrics_channel chan []graphite.Metric) {
	for _, query := range server.Queries {
		go SendSQLQuery(server, query, metrics_channel)
	}
}

func SendSQLQuery(server Sql_server, query Query, metrics_channel chan []graphite.Metric) {
	for {
		metrics, err := getSQLMetrics(server, query)
		if err == nil {
			metrics_channel <- metrics
		}
		time.Sleep(time.Duration(query.Interval) * time.Second)
	}
}

func getSQLMetrics(s Sql_server, q Query) (metrics []graphite.Metric, err error) {

	db, err := sql.Open("mssql", s.Connection_string)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(q.Tsql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	if q.Type == singleRowType || q.Type == nil {	// nil for backward compatibility

		columnNames, err := rows.Columns()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		values := make([]interface{}, len(columnNames))
		valuePointers := make([]interface{}, len(columnNames))
		for i := 0; i < len(columnNames); i++ {
			valuePointers[i] = &values[i]
		}

		rows.Next()
		if err = rows.Scan(valuePointers...); err != nil {
			log.Println(err)
			return nil, err
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
			return nil, err
		}

		for i := 0; i < len(columnNames); i++ {
			metrics = append(metrics, graphite.Metric{
				fmt.Sprintf("%s.%s", s.Metric_prefix, NormalizeMetricName(columnNames[i])),
				fmt.Sprintf("%v", values[i]),
				time.Now().Unix(),
			})
		}

		return metrics, nil

	} else if q.Type == keyValueType{

		return nil, errors.New("Unsupported SQL Metric type")
	}

	return nil, errors.New("Unknown SQL Metric type")
}
