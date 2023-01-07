package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func (c *Configuration)pqConnect()(*gorm.DB, error) {
	return gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.DBConfig.User,
		c.DBConfig.Password,
		c.DBConfig.Host,
		c.DBConfig.Database))
}
