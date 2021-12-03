package model

import "gorm.io/gorm"

type Config struct {
	gorm.Model
	Key   string `gorm:"config_key"`
	Value string
}
