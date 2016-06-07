package golagraphite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/marpaia/graphite-golang"
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
		} else {
			log.Println(err)
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

	if len(q.Tsql_row) != 0 && len(q.Tsql_row) != 0 {
		return nil, errors.New("Both SQL Metric types provided, seems you have an error in config")
	}

	// FIXME: there should be a strategy
	if len(q.Tsql_row) != 0 {

		rows, err := db.Query(q.Tsql_row)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		columnNames, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		values := make([]interface{}, len(columnNames))
		valuePointers := make([]interface{}, len(columnNames))
		for i := 0; i < len(columnNames); i++ {
			valuePointers[i] = &values[i]
		}

		timestamp := time.Now().Unix()

		for rows.Next() {
			if err = rows.Scan(valuePointers...); err != nil {
				return nil, err
			}

			if err = rows.Err(); err != nil {
				return nil, err
			}

			for i := 0; i < len(columnNames); i++ {
				if columnNames[i] == q.Timestamp {
					timestamp = values[i].(time.Time).Unix()
				}
			}

			for i := 0; i < len(columnNames); i++ {
				metrics = append(metrics, graphite.Metric{
					Name: fmt.Sprintf("%s.%s", q.Metric_prefix, NormalizeMetricName(columnNames[i])),
					Value: fmt.Sprintf("%v", values[i]),
					Timestamp: timestamp,
				})
			}
		}

		return metrics, nil

	} else if len(q.Tsql_table) != 0 {

		rows, err := db.Query(q.Tsql_table)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		data := make(map[string]string)
		var key, value string

		for rows.Next() {
			err := rows.Scan(&key, &value)
			if err != nil {
				return nil, err
			}
			data[key] = value
		}

		err = rows.Err()
		if err != nil {
			return nil, err
		}

		for key, value := range data {
			metrics = append(metrics, graphite.Metric{
				fmt.Sprintf("%s.%s", q.Metric_prefix, NormalizeMetricName(key)),
				fmt.Sprintf("%v", value),
				time.Now().Unix(),
			})
		}

		return metrics, nil

	}

	return nil, errors.New("Unknown SQL Metric type")
}
