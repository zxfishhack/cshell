package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/zxfishhack/cshell/pkg/dto"
	"github.com/zxfishhack/cshell/pkg/store"
	"github.com/zxfishhack/cshell/pkg/tray"
	"net/http"
)

type CShellController struct {
	Ctx iris.Context
}

func (*CShellController) GetHostList() []string {
	return store.GetSSHHostList("", true)
}

func (*CShellController) GetHostListTagBy(tag string) []string {
	return store.GetSSHHostList(tag, true)
}

func (*CShellController) GetKeys() []string {
	if v := store.GetKeys(); v == nil {
		return make([]string, 0)
	} else {
		return v
	}
}

func (*CShellController) GetTags() []string {
	if v := store.GetTagList(); v == nil {
		return make([]string, 0)
	} else {
		return v
	}
}

func (*CShellController) GetHostByConfig(hostId string) (res dto.HostConfig) {
	res.Host = hostId
	res.Items = store.GetHostKV(hostId)
	res.Visible = store.IsHostVisible(hostId, false)
	res.Tags = store.GetHostTags(hostId)
	return
}

func (c *CShellController) PostHostByConfig(hostId string) {
	var hc dto.HostConfig
	err := c.Ctx.ReadJSON(&hc)
	if err != nil {
		c.Ctx.StatusCode(http.StatusBadRequest)
		c.Ctx.WriteString(err.Error())
		return
	}
	update := false
	update = store.SaveHostKV(hostId, hc.Host, hc.Items) || update
	update = store.SaveHostTags(hostId, hc.Tags) || update
	update = store.SetHostVisible(hostId, hc.Visible) || update
	update = update || (hostId != hc.Host)
	if hostId != hc.Host {
		store.ChangeName(hostId, hc.Host)
	}
	if update {
		tray.Inst.Reload()
	}
}

func (*CShellController) DeleteHostByConfig(hostId string) {
	if store.DeleteHost(hostId) {
		tray.Inst.Reload()
	}
}
