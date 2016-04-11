package golagraphite

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/marpaia/graphite-golang"
	_ "github.com/zensqlmonitor/go-mssqldb" // we cant use default github.com/denisenkom/go-mssqldb cause of unmerged pull/145
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

}
