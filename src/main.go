package main

import "C"
import (
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile("./config.toml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("读取配置文件出错", "err", err)
	}

	var cfg BuilderConfig

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("解析配置参数出错", "err", err)
	}
	builder := NewBuilder(cfg)

	builder.Build()
}
