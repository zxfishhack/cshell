package dto

import "github.com/zxfishhack/ssh_config"

type HostConfig struct {
	Host    string           `json:"host"`
	Visible bool             `json:"visible"`
	Tags    []string         `json:"tags"`
	Items   []*ssh_config.KV `json:"items"`
}
