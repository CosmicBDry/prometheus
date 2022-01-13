package command

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

//1.创建一个结构体CmdCollector，且定义Describe和Collect两个方法，使得满足prometheus.Collector接口中所有方法------------------------------------------>
type CmdCollector struct {
	Desc *prometheus.Desc
	Db   *sql.DB
}

//2.初始化CmdCollector结构体，方便main.go中注册此结构体对象---------------------------------------------------------------------------------------->
func NewCmdCollector(db *sql.DB, mysqladdr string) *CmdCollector {

	return &CmdCollector{
		Desc: prometheus.NewDesc("Mysql_Status_Command", "help mysql status command", []string{"cmd"}, prometheus.Labels{"MysqlAddr": mysqladdr}),
		Db:   db,
	}
}

//3.常见Describe方法，用于定义指标类型：如定义指标名称、help信息、可变标签、固定标签等----------------------------------------------------------------->
func (c *CmdCollector) Describe(desc chan<- *prometheus.Desc) {
	desc <- c.Desc
}

//4.创建Collect方法，用于指定将被采集的数据--------------------------------------------------------------------------------------------------------->
func (c *CmdCollector) Collect(metric chan<- prometheus.Metric) {
	var (
		labelvalue string
		count      float64
	)

	cmdName := []string{"Insert", "Delete", "Update", "Select"}

	for _, cmdname := range cmdName {
		c.Db.QueryRow("show global status where variable_name=?", "Com_"+cmdname).Scan(&labelvalue, &count)
		metric <- prometheus.MustNewConstMetric(c.Desc, prometheus.CounterValue, count, cmdname) //（）中参数分别是指标名称、指标类型、返回的样本值、可变标签值
	}

}
