package queriesAll

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type QeuryClollector struct {
	Desc *prometheus.Desc
	Db   *sql.DB
}

func NewQeuryClollector(db *sql.DB, mysqladdr string) *QeuryClollector {
	return &QeuryClollector{
		Desc: prometheus.NewDesc("Mysql_status_Queries", "hep mysql status queries", nil, prometheus.Labels{"MysqlAddr": mysqladdr}),
		Db:   db,
	}
}

func (q *QeuryClollector) Describe(desc chan<- *prometheus.Desc) {
	desc <- q.Desc
}

func (q *QeuryClollector) Collect(metric chan<- prometheus.Metric) {

	var (
		name  string
		count float64
	)

	q.Db.QueryRow("show global status where variable_name=?", "Queries").Scan(&name, &count)

	metric <- prometheus.MustNewConstMetric(q.Desc, prometheus.CounterValue, count)

}
