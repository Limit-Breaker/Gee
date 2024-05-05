package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var DefaultConfig = initConfig()

type Config struct {
	viper *viper.Viper
	App
	Grpc
}

type App struct {
	Name string
	Addr string
}

type Grpc struct {
	Addr string
}

// Log TODO
type Log struct {
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
	conf.LoadAppConfig()
	conf.LoadGrpcConfig()
	return conf
}

func (c *Config) LoadAppConfig() {
	app := App{}
	app.Name = c.viper.GetString("app.name")
	app.Addr = c.viper.GetString("app.addr")
	c.App = app
}

func (c *Config) LoadGrpcConfig() {
	grpc := Grpc{}
	grpc.Addr = c.viper.GetString("grpc.addr")
	c.Grpc = grpc
}

func (c *Config) LoadLogConfig() {
	// TODO
}
