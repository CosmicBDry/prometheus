package config

import "github.com/spf13/viper"

//定义一个接收配置文件信息的结构体类型----------------------------------------->
type Config struct {
	Mysql struct {
		Host       string `mapstructure:"host"`
		Port       int64  `mapstructure:"port"`
		DbUser     string `mapstructure:"db_user"`
		DbPassword string `mapstructure:"db_password"`
	} `mapstructure: "Mysql"`

	Web struct {
		Addr       string `mapstructure:"addr"`
		Basic_Auth struct {
			UserName string `mapstructure:"username"`
			PassWord string `mapstructure:"user_password"`
		} `mapstructure:"basic_auth"`
	} `mapstructure:"Web"`
	Logger struct {
		FilePath  string `mapstructure:"file_path"`
		MaxSize   int    `mapstructure:"log_file_max_size"`
		MaxAge    int    `mapstructure:"log_retain_recent_days"`
		LocalTime bool   `mapstructure:"log_backup_local_time"`
		Compress  bool   `mapstructure:"compress"`
	} `mapstructure:"Logger"`
}

//定义一个Getconfig函数来获取配置信息，通过viper工具来实现----------------------------------------->
func GetConfig(path string) (*Config, error) {
	config := &Config{}
	//viper设置将要获取的配置文件路径
	viper.SetConfigFile(path)
	//若配置文件没指定，可设置如下默认配置
	viper.SetDefault("Mysql.port", 3306)
	//将配置信息读入到viper中
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	//通过viper.Unmashal反序列化将配置文件内容解析到config结构体中
	// 也可以直接通过viper.Get 、viper.GetString、 viper.GetSInt、来获取配置key对应的value: 如viper.GetInt("Mysql.port")、viper.Getstring("Mysql.dbname")
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	//解析成功后返回一个config和nil
	return config, nil
}
