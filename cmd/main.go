package main

import (
	"fmt"
	"github.com/KhasanOrsaev/orse/internal/repository/config"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	v := viper.New()
	config.NewDefaultConfig(v)
}

func main() {
	c:= config.Config()
	defer c.CloseConnections()
}

