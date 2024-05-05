package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var DefaultConfig = initConfig()

type Config struct {
	viper *viper.Viper
	Grpc
}

type Grpc struct {
	Addr string
}

func initConfig() *Config {
	conf := &Config{viper: viper.New()}
	workdir, _ := os.Getwd()
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workdir + "/config")
	//conf.viper.AddConfigPath("/etc")
	if err := conf.viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}
	conf.LoadGrpcConfig()
	return conf
}

func (c *Config) LoadGrpcConfig() {
	grpc := Grpc{}
	grpc.Addr = c.viper.GetString("grpc.addr")
	c.Grpc = grpc
}
