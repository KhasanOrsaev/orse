package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Configuration struct {
	DB string
	DBConfig struct{
		Host string
		Port string
		Database string
		User string
		Password string
	}
	DBClient *gorm.DB
	Log struct{
		Format      string
		ServiceName string
		Level       int
	}
}

var (
	conf       Configuration
	err        error
)
func Config() *Configuration {
	return &conf
}

// NewConfig создание нового экземпляра конфига
func NewConfig(v *viper.Viper) *Configuration {
	var c = Configuration{}
	c.setConfig(v)
	return &c
}

// NewDefaultConfig обновить настройки экземпляра конфига
func NewDefaultConfig(v *viper.Viper) {
	conf.setConfig(v)
}

// setConfig initialize config from env
func (c *Configuration) setConfig(v *viper.Viper) {
	// если директории нет то создаем
	if _, err = os.Stat("./var/log"); os.IsNotExist(err) {
		err = os.MkdirAll("./var/log", os.ModePerm)
		if err != nil {
			fmt.Println("Error on create log directory:", err.Error())
			os.Exit(1)
		}
	}
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	v.BindEnv("log.level")
	v.SetDefault("log.level", 4)
	v.BindEnv("log.name")
	v.SetDefault("log.name", "dev")

	v.SetDefault("db.config.host","localhost")
	v.BindEnv("db.config.host")
	v.SetDefault("db.config.port","5432")
	v.BindEnv("db.config.port")
	v.SetDefault("db.config.user","test")
	v.BindEnv("db.config.user")
	v.SetDefault("db.config.password","test")
	v.BindEnv("db.config.password")
	v.SetDefault("db.config.db","test")
	v.BindEnv("db.config.db")
	v.SetDefault("db","postgres")
	v.BindEnv("db")

	c.DB = v.GetString("db")
	c.DBConfig.User = v.GetString("db.config.user")
	c.DBConfig.Password = v.GetString("db.config.password")
	c.DBConfig.Host = v.GetString("db.config.host")
	c.DBConfig.Port = v.GetString("db.config.port")
	c.DBConfig.Database = v.GetString("db.config.db")

	c.Log.Level = v.GetInt("log.level")
	c.Log.ServiceName = v.GetString("log.name")
	c.Log.Format = "[%s] %s.%s message: %s context: %s extra: %s"

	switch c.DB {
	case "postgres":
		c.DBClient,err = c.pqConnect()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		c.DBClient,err = c.pqConnect()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if err = c.DBClient.DB().Ping(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *Configuration) CloseConnections() error {
	if err := c.DBClient.Close(); err != nil {
		return err
	}
	return nil
}
