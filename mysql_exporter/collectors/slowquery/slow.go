package slowquery

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type SlowQueryCollector struct {
	Desc *prometheus.Desc
	Db   *sql.DB
}

func NewSlowQueryCollector(db *sql.DB, mysqladdr string) *SlowQueryCollector {
	return &SlowQueryCollector{
		Desc: prometheus.NewDesc("Mysql_Status_SlowQuery", "help mysql status slowQuery", nil, prometheus.Labels{"MysqlAddr": mysqladdr}),
		Db:   db,
	}

}

func (s *SlowQueryCollector) Describe(desc chan<- *prometheus.Desc) {
	desc <- s.Desc
}

func (s *SlowQueryCollector) Collect(metric chan<- prometheus.Metric) {
	var (
		name  string
		count float64
	)

	s.Db.QueryRow("show global status where variable_name=?", "Slow_queries").Scan(&name, &count)
	metric <- prometheus.MustNewConstMetric(s.Desc, prometheus.CounterValue, count)

}
