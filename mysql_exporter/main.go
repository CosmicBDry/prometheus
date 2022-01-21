package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/CosmicBDry/prometheus/mysql_exporter/collectors/command"
	"github.com/CosmicBDry/prometheus/mysql_exporter/collectors/connect"
	"github.com/CosmicBDry/prometheus/mysql_exporter/collectors/queriesAll"
	"github.com/CosmicBDry/prometheus/mysql_exporter/collectors/slowquery"
	"github.com/CosmicBDry/prometheus/mysql_exporter/collectors/traffic"
	"github.com/CosmicBDry/prometheus/mysql_exporter/config"
	"github.com/CosmicBDry/prometheus/mysql_exporter/handlerAuth"
	"github.com/CosmicBDry/prometheus/mysql_exporter/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	var (
		configure, server, addr, init string
		help, h                       bool
	)
	flag.BoolVar(&help, "help", false, "./mysql_exporter -help useage")
	flag.BoolVar(&h, "h", false, "./mysql_exporter -h useage")
	flag.StringVar(&configure, "config", "./conf/config.yaml", "-config ./conf/config.yaml ")
	flag.StringVar(&server, "server", "", "-server localhost:9090  (default localhost:9090)")
	flag.StringVar(&init, "init", "", "-init password ( Generate an  encrypted Password ,According to the password you provided)")
	flag.Usage = func() {
		fmt.Println("usage: ./mysql_exporter [-config ./conf/config.yaml] [-server localhost:9090]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if help || h {
		flag.Usage()
		os.Exit(-1)
	}

	if init != "" {
		pass, _ := bcrypt.GenerateFromPassword([]byte(init), 6)
		fmt.Println("Token generated successfully: ", string(pass))
		os.Exit(0)
	}

	//自定义Getconfig函数，获取配置文件
	Conf, err := config.GetConfig(configure)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	//自定义一个开启日志记录的函数SetAppLog()
	Logger := logger.SetAppLog(Conf.Logger.FilePath, Conf.Logger.MaxSize, Conf.Logger.MaxAge, Conf.Logger.LocalTime, Conf.Logger.Compress)

	Logger.Info("成功读取配置：" + configure)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", Conf.Mysql.DbUser, Conf.Mysql.DbPassword, Conf.Mysql.Host, Conf.Mysql.Port)

	//代理指标启动监听的地址和端口
	if server == "" {
		if Conf.Web.Addr == "" {
			addr = "localhost:9090"
		} else {
			addr = Conf.Web.Addr
		}

	} else {
		addr = server
	}

	//指标采集的固定label值mysqladdr
	mysqladdr := fmt.Sprintf("%s:%d", Conf.Mysql.Host, Conf.Mysql.Port)

	DB, _ := sql.Open("mysql", dsn)
	err = DB.Ping()
	if err != nil {
		Logger.Error(err)
	} else {
		Logger.Info("MySQL数据库连接成功!")
	}

	user, password := Conf.Web.Basic_Auth.UserName, Conf.Web.Basic_Auth.PassWord //设置暴露访问认证
	//password, _ := bcrypt.GenerateFromPassword([]byte("admin@123"), 6)//bcrypt密码加密，6次hash加密
	//fmt.Println(string(password))

	Logger.Info("mysql_exporter服务启动成功，metrics访问地址为: " + "http://" + addr + "/metrics")
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
		Logger.Error(err)
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
	//将当前的Logger对象传入handlerAuth.Auth中，记录Auth中的错误日志
	http.Handle("/metrics/", handlerAuth.Auth(promhttp.Handler(), handlerAuth.AuthSecrets{user: password}, Logger))

	//启动http服务监听-------------------------------------------------------------------------------->
	http.ListenAndServe(addr, nil)

}
