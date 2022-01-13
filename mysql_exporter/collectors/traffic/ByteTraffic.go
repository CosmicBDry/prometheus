package traffic

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type TrafficCollector struct {
	Desc *prometheus.Desc
	Db   *sql.DB
}

func NewTrafficCollector(db *sql.DB, mysqladdr string) *TrafficCollector {
	return &TrafficCollector{
		Desc: prometheus.NewDesc("Mysql_Status_ByteTraffic", "help ByteTraffic", []string{"TrafficType"}, prometheus.Labels{"MysqlAddr": mysqladdr}),
		Db:   db,
	}
}

func (t *TrafficCollector) Describe(desc chan<- *prometheus.Desc) {
	desc <- t.Desc
}

func (t *TrafficCollector) Collect(metric chan<- prometheus.Metric) {

	var (
		name     string
		quantity float64
		types    = []string{"sent", "received"}
	)

	for _, trafficType := range types {

		t.Db.QueryRow("show global status where variable_name=?", "Bytes_"+trafficType).Scan(&name, &quantity)
		metric <- prometheus.MustNewConstMetric(t.Desc, prometheus.CounterValue, quantity, trafficType)
	}

}
