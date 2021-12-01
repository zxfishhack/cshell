package tray

import (
	"context"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"github.com/zxfishhack/cshell/pkg/icon"
	"github.com/zxfishhack/cshell/pkg/store"
	"github.com/zxfishhack/cshell/pkg/utils"
	"golang.org/x/sync/errgroup"
	"log"
)

type Tray struct {
	g, g2           *errgroup.Group
	ctx, reloadCtx  context.Context
	cancel, cancel2 context.CancelFunc
	port            int

	termItems []*systray.MenuItem
}

var Inst *Tray

func New(port int) {
	Inst = &Tray{port: port}
	Inst.init()
	return
}

func (t *Tray) init() {
	var tCtx context.Context
	tCtx, t.cancel = context.WithCancel(context.Background())
	t.g, t.ctx = errgroup.WithContext(tCtx)
	tCtx, t.cancel2 = context.WithCancel(context.Background())
	t.g2, t.reloadCtx = errgroup.WithContext(tCtx)
}

func (t *Tray) Reload() {
	log.Print("reload....")
	systray.ClearAllMenuItem()
	t.cancel2()
	_ = t.g2.Wait()
	t.cancel()
	_ = t.g.Wait()
	t.init()
	_ = store.Reload()
	t.render()
}

func (t *Tray) render() {
	var m *systray.MenuItem
	m = systray.AddMenuItem("All", "")
	var hosts []string
	hosts = store.GetSSHHostList("", false)
	for _, host := range hosts {
		sm := m.AddSubMenuItem(host, fmt.Sprintf("ssh %s ...", host))
		t.handleSSH(host, sm)
	}
	tags := store.GetTagList()
	for _, tag := range tags {
		hosts = store.GetSSHHostList(tag, false)
		if len(hosts) == 0 {
			continue
		}
		m = systray.AddMenuItem(tag, "")
		for _, host := range hosts {
			sm := m.AddSubMenuItem(host, fmt.Sprintf("ssh %s ...", host))
			t.handleSSH(host, sm)
		}
	}
	systray.AddSeparator()
	typ := store.GetTerminalType()
	m = systray.AddMenuItemCheckbox("使用iTerm", "", typ == utils.ITerm2)
	t.handleTerm(m, utils.ITerm2)
	m = systray.AddMenuItemCheckbox("使用默认终端", "", typ == utils.DefaultTerminal)
	t.handleTerm(m, utils.DefaultTerminal)
	systray.AddSeparator()
	m = systray.AddMenuItem("重新加载", "")
	t.handleReload(m)
	systray.AddSeparator()
	m = systray.AddMenuItem("配置...", "")
	t.handleConfig(m)
	systray.AddSeparator()
	m = systray.AddMenuItem("退出", "")
	t.handleQuit(m)
}

func (t *Tray) Run() (err error) {
	systray.Run(func() {
		systray.SetTemplateIcon(icon.MustAsset("icon.png"), icon.MustAsset("icon.png"))
		systray.SetTooltip("CShell")
		t.render()
	}, nil)
	return
}

func (t *Tray) handleSSH(name string, m *systray.MenuItem) {
	t.g.Go(func() (err error) {
		run := true
		for run {
			select {
			case <-m.ClickedCh:
				_ = utils.OpenSSH(store.GetTerminalType(), name)
			case <-t.ctx.Done():
				run = false
			}
		}
		return nil
	})
}

func (t *Tray) handleReload(m *systray.MenuItem) {
	t.g.Go(func() error {
		run := true
		for run {
			select {
			case <-m.ClickedCh:
				run = false
				go t.Reload()
			case <-t.ctx.Done():
				run = false
			case <-t.reloadCtx.Done():
				run = false
			}
		}
		return nil
	})
}

func (t *Tray) handleConfig(m *systray.MenuItem) {
	t.g.Go(func() error {
		run := true
		for run {
			select {
			case <-m.ClickedCh:
				open.Run(fmt.Sprintf("http://localhost:%d", t.port))
			case <-t.ctx.Done():
				run = false
			}
		}
		return nil
	})
}

func (t *Tray) handleQuit(m *systray.MenuItem) {
	t.g2.Go(func() error {
		select {
		case <-m.ClickedCh:
			t.cancel()
			t.g.Wait()
			systray.Quit()
		case <-t.reloadCtx.Done():
		}
		return nil
	})
}

func (t *Tray) handleTerm(m *systray.MenuItem, termType utils.TerminalType) {
	t.g.Go(func() error {
		run := true
		for run {
			select {
			case <-m.ClickedCh:
				_ = store.SetTerminalType(termType)
				go t.Reload()
			case <-t.ctx.Done():
				run = false
			}
		}
		return nil
	})
}
