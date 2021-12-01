package model

import (
	"gorm.io/gorm"
)

type HostAlias struct {
	gorm.Model
	Name string
	Tag  string
}

type HostVisible struct {
	gorm.Model
	Name    string
	Visible bool
}
