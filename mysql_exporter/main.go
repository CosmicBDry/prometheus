package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/CosmicBDry/prometheus/mysql_exporter/collectors/command"
	"github.com/CosmicBDry/prometheus/mysql_exporter/collectors/connect"
	"github.com/CosmicBDry/prometheus/mysql_exporter/collectors/queriesAll"
	"github.com/CosmicBDry/prometheus/mysql_exporter/collectors/slowquery"
	"github.com/CosmicBDry/prometheus/mysql_exporter/collectors/traffic"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	dsn := "root:123456@tcp(192.168.30.32:3306)/gocmdb"
	addr := "localhost:9090"
	mysqladdr := "192.168.30.32:3306/gocmdb"
	DB, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	//定义指标类型--------------------------------------------------------------------------------->
	//mysql存活探测指标
	mysql_up := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "Mysql_Status_Up",
		Help:        "help mysql up",
		ConstLabels: map[string]string{"addr": mysqladdr},
	}, func() float64 {
		if err := DB.Ping(); err == nil {

			return 1
		}
		return 0
	})

	//指标类型注册--------------------------------------------------------------------------------->
	//mysql存活探测指标注册
	prometheus.MustRegister(mysql_up)

	//mysql的增删改查统计指标注册
	prometheus.MustRegister(command.NewCmdCollector(DB, mysqladdr))

	//慢查询指标注册
	prometheus.MustRegister(slowquery.NewSlowQueryCollector(DB, mysqladdr))

	//最大连接和已连接指标注册
	prometheus.MustRegister(connect.NewConnectCollector(DB, mysqladdr))

	//所有查询包括show、select...等Query统计指标注册
	prometheus.MustRegister(queriesAll.NewQeuryClollector(DB, mysqladdr))

	//数据库已接收和已发送的流量统计注册
	prometheus.MustRegister(traffic.NewTrafficCollector(DB, mysqladdr))

	//暴露指标，url与处理器关系绑定-------------------------------------------------------------------->
	http.Handle("/metrics/", promhttp.Handler())
	
	//启动http服务监听-------------------------------------------------------------------------------->
	http.ListenAndServe(addr, nil)

}
