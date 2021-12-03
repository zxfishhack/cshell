package model

import (
	"gorm.io/gorm"
)

type HostAlias struct {
	gorm.Model
	Name string `gorm:"index:name_tag"`
	Tag  string `gorm:"index:name_tag"`
}

type HostVisible struct {
	gorm.Model
	Name    string `gorm:"index:name"`
	Visible bool
}
