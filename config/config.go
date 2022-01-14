package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var path = pflag.StringP("config", "c", "", "项目配置文件")
var GobalConfig *Configuration

type Configuration struct {
	Name          string `yaml:"name"`
	IP            string `yaml:"ip"`
	Port          int    `yaml:"port"`
	MaxConnLen    int    `yaml:"maxConnLen"`
	WorkPoolSize  int    `yaml:"workPoolSize"`
	MaxPackerSize uint32 `yaml:"maxPackerSize"`
}

//初始化配置文件
func init() {
	pflag.Parse()
	v := viper.New()
	if *path != "" {
		fmt.Println(*path)
		v.SetConfigFile(*path) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		v.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		v.SetConfigName("config")
	}

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	var config Configuration
	if err := v.Unmarshal(&config); err != nil {
		panic(err)
	}

	GobalConfig = &config
}
