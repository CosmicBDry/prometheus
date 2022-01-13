package connect

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type ConnectCollector struct {
	MaxConnDesc   *prometheus.Desc
	ConnectedDesc *prometheus.Desc
	Db            *sql.DB
}

func NewConnectCollector(db *sql.DB, mysqladdr string) *ConnectCollector {

	return &ConnectCollector{
		MaxConnDesc:   prometheus.NewDesc("Mysql_Variables_MaxConnect", "help mysql variables maxConnect", nil, prometheus.Labels{"MysqlAddr": mysqladdr}),
		ConnectedDesc: prometheus.NewDesc("Mysql_Status_Client_Connected", "help mysql status client connected", nil, prometheus.Labels{"MysqlAddr": mysqladdr}),
		Db:            db,
	}

}

func (c *ConnectCollector) Describe(desc chan<- *prometheus.Desc) {
	desc <- c.MaxConnDesc
	desc <- c.ConnectedDesc
}

func (c *ConnectCollector) Collect(metric chan<- prometheus.Metric) {

	var (
		varname string
		numbers float64
	)

	c.Db.QueryRow("show global variables where variable_name=?", "max_connections").Scan(&varname, &numbers)
	metric <- prometheus.MustNewConstMetric(c.MaxConnDesc, prometheus.CounterValue, numbers)

	c.Db.QueryRow("show global status where variable_name=?", "Threads_connected").Scan(&varname, &numbers)
	metric <- prometheus.MustNewConstMetric(c.ConnectedDesc, prometheus.CounterValue, numbers)
}
