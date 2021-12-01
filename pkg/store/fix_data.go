package store

import (
	"github.com/zxfishhack/cshell/pkg/model"
	"gorm.io/gorm"
)

func fixData() {
	hosts := GetSSHHostList("", true)
	check := make(map[string]bool)
	for _, h := range hosts {
		check[h] = true
	}

	fixHostAlias(check)
	fixHostVisible(check)
}

func fixHostVisible(check map[string]bool) {
	ids := make([]uint, 0)
	result := db.Find(&model.HostVisible{})
	if result.Error != nil {
		return
	}
	rows, err := result.Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		var ha model.HostVisible
		if db.ScanRows(rows, &ha) == nil {
			if _, ok := check[ha.Name]; !ok {
				ids = append(ids, ha.ID)
			}
		}
	}
	for _, id := range ids {
		db.Delete(&model.HostVisible{
			Model: gorm.Model{ID: id},
		})
	}
}

func fixHostAlias(check map[string]bool) {
	ids := make([]uint, 0)
	result := db.Find(&model.HostAlias{})
	if result.Error != nil {
		return
	}
	rows, err := result.Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		var ha model.HostAlias
		if db.ScanRows(rows, &ha) == nil {
			if _, ok := check[ha.Name]; !ok {
				ids = append(ids, ha.ID)
			}
		}
	}
	for _, id := range ids {
		db.Delete(&model.HostAlias{
			Model: gorm.Model{ID: id},
		})
	}
}
