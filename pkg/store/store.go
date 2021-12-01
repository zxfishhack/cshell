package store

import (
	"errors"
	"github.com/zxfishhack/cshell/pkg/model"
	"github.com/zxfishhack/cshell/pkg/sshc"
	"github.com/zxfishhack/cshell/pkg/utils"
	"github.com/zxfishhack/ssh_config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var sshConfig *ssh_config.Config
var db *gorm.DB

func Init() (err error) {
	db, err = gorm.Open(sqlite.Open(filepath.Join(os.Getenv("HOME"), ".ssh", ".cshell")),
		&gorm.Config{Logger: logger.Discard})
	if err != nil {
		return
	}

	_ = db.AutoMigrate(&model.Tag{})
	_ = db.AutoMigrate(&model.HostVisible{})
	_ = db.AutoMigrate(&model.HostAlias{})
	_ = db.AutoMigrate(&model.Config{})

	var cfg model.Config
	if errors.Is(db.First(&cfg, model.Config{Key: model.TerminalType}).Error, gorm.ErrRecordNotFound) {
		db.Create(&model.Config{Key: "terminal-type", Value: string(utils.DefaultTerminal)})
	}

	return Reload()
}

func Reload() (err error) {
	sshConfig, err = sshc.LoadSSHConfig(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	if err != nil {
		return
	}

	fixData()

	return
}

func SaveSSHConfig() (err error) {
	return sshc.SaveSSHConfig(filepath.Join(os.Getenv("HOME"), ".ssh", "config"), sshConfig)
}

func GetTerminalType() utils.TerminalType {
	var cfg model.Config
	if err := db.First(&cfg, model.Config{Key: model.TerminalType}).Error; err != nil {
		return utils.DefaultTerminal
	}
	return utils.TerminalType(cfg.Value)
}

func SetTerminalType(typ utils.TerminalType) (err error) {
	cfg := &model.Config{Key: model.TerminalType, Value: string(typ)}
	return db.Model(&model.Config{}).Where(model.Config{Key: model.TerminalType}).UpdateColumns(cfg).Error
}

func GetTagList() (tags []string) {
	result := db.Find(&model.Tag{})
	if result.Error != nil {
		return
	}
	rows, err := result.Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		var tag model.Tag
		if db.ScanRows(rows, &tag) == nil {
			tags = append(tags, tag.Name)
		}
	}
	return
}

func DeleteHost(hostId string) bool {
	idx := -1
	for i, h := range sshConfig.Hosts {
		for _, p := range h.Patterns {
			if p.String() == hostId {
				idx = i
				break
			}
			break
		}
	}
	if idx != -1 {
		sshConfig.Hosts = append(sshConfig.Hosts[:idx], sshConfig.Hosts[idx+1:]...)
		_ = SaveSSHConfig()
	}
	return idx != -1
}

func IsHostVisible(name string, all bool) bool {
	if name == "*" {
		return false
	}
	if all {
		return true
	}
	var hv model.HostVisible
	err := db.First(&hv, model.HostVisible{Name: name}).Error
	if err != nil {
		return true
	}
	return hv.Visible
}

func SetHostVisible(name string, visible bool) (update bool) {
	var hv model.HostVisible
	update = IsHostVisible(name, false) == visible
	err := db.First(&hv, model.HostVisible{Name: name}).Error
	if err != nil {
		db.Create(&model.HostVisible{Name: name, Visible: visible})
	} else {
		db.Where(model.HostVisible{Name: name}).Updates(&model.HostVisible{Name: name, Visible: visible})
	}
	return
}

func GetSSHHostList(tag string, all bool) (hosts []string) {
	if tag == "" {
		if sshConfig != nil {
			for _, v := range sshConfig.Hosts {
				for _, p := range v.Patterns {
					if !IsHostVisible(p.String(), all) {
						continue
					}
					hosts = append(hosts, p.String())
				}
			}
		}
		return
	}
	result := db.Find(&model.HostAlias{}, model.HostAlias{Tag: tag})
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
			if !IsHostVisible(ha.Name, all) {
				continue
			}
			hosts = append(hosts, ha.Name)
		}
	}
	return
}

func GetHostKV(hostId string) (res []*ssh_config.KV) {
	for _, host := range sshConfig.Hosts {
		for _, p := range host.Patterns {
			if hostId == p.String() {
				for _, n := range host.Nodes {
					if kv, ok := n.(*ssh_config.KV); ok {
						res = append(res, kv)
					}
				}
				return
			}
		}
	}
	return
}

func SaveHostKV(hostId string, newHostId string, kv []*ssh_config.KV) (update bool) {
	found := false
	for _, host := range sshConfig.Hosts {
		for _, p := range host.Patterns {
			if hostId == p.String() {
				found = true
				break
			}
		}
		if found {
			host.Nodes = make([]ssh_config.Node, 0)
			for _, i := range kv {
				i.SetLeadingSpace(4)
				host.Nodes = append(host.Nodes, i)
			}
			if newHostId != hostId {
				if p, err := ssh_config.NewPattern(newHostId); err == nil {
					host.Patterns = nil
					host.Patterns = append(host.Patterns, p)
				}
			}
			_ = SaveSSHConfig()
			break
		}
	}
	if !found {
		p, err := ssh_config.NewPattern(newHostId)
		if err != nil {
			return
		}
		update = true
		host := &ssh_config.Host{
			Patterns:   nil,
			Nodes:      nil,
			EOLComment: "",
		}
		host.Patterns = append(host.Patterns, p)
		for _, i := range kv {
			i.SetLeadingSpace(4)
			host.Nodes = append(host.Nodes, i)
		}
		sshConfig.Hosts = append(sshConfig.Hosts, host)
		_ = SaveSSHConfig()
	}
	return
}

func GetHostTags(hostId string) (res []string) {
	result := db.Find(&model.HostAlias{Name: hostId})
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
			res = append(res, ha.Tag)
		}
	}
	return
}

func SaveHostTags(hostId string, tags []string) (updated bool) {
	oldTags := GetHostTags(hostId)
	addTags := diff(tags, oldTags)
	removeTags := diff(oldTags, tags)
	if len(addTags) == 0 && len(removeTags) == 0 {
		return
	}
	updated = true
	for _, tag := range addTags {
		db.Create(&model.HostAlias{Name: hostId, Tag: tag})
	}
	for _, tag := range removeTags {
		db.Delete(&model.HostAlias{Name: hostId, Tag: tag})
	}
	addTags = diff(tags, GetTagList())
	for _, tag := range addTags {
		db.Create(&model.Tag{Name: tag})
	}
	return
}

func ChangeName(hostId string, newHostId string) {
	db.Where(model.HostAlias{Name: hostId}).Update("name", newHostId)
	db.Where(model.HostVisible{Name: hostId}).Update("name", newHostId)
}

// Set Difference: A - B
func diff(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func GetKeys() (res []string) {
	keys := make(map[string]bool)
	for _, h := range sshConfig.Hosts {
		for _, item := range h.Nodes {
			if v, ok := item.(*ssh_config.KV); ok {
				if v.Key == "IdentityFile" {
					keys[v.Value] = true
				}
			}
		}
	}
	filepath.Walk(filepath.Join(os.Getenv("HOME"), ".ssh", "keys"), func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		// 默认pub结尾的为公钥
		if strings.HasSuffix(path, ".pub") {
			return nil
		}
		keys[strings.Replace(path, os.Getenv("HOME"), "~", 1)] = true
		return nil
	})
	for k := range keys {
		res = append(res, k)
	}
	sort.Strings(res)
	return
}