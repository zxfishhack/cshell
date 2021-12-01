package main

import (
	"context"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/zxfishhack/cshell/pkg/icon"
	"github.com/zxfishhack/cshell/pkg/store"
	"github.com/zxfishhack/cshell/pkg/utils"
	"golang.org/x/sync/errgroup"
	"log"
)

func main() {
	err := store.Init()
	if err != nil {
		log.Panicln(err)
	}

	systray.Run(onReady, nil)
}

func onReady() {
	hosts := store.GetSSHHostList()

	ctx, cancel := context.WithCancel(context.Background())
	g, _ := errgroup.WithContext(ctx)

	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTooltip("CShell")
	for _, host := range hosts {
		item := systray.AddMenuItem(host, fmt.Sprintf("ssh %s ...", host))
		func(h string) {
			g.Go(func() error {
				for {
					select {
					case <-ctx.Done():
						return nil
					case <-item.ClickedCh:
						log.Printf("%s clicked.", h)
						err := utils.OpenSSH(h)
						log.Printf("ssh result: %v", err)
					}
				}
			})
		}(host)
	}

	mQuit := systray.AddMenuItem("退出", "")
	go func() {
		<-mQuit.ClickedCh
		cancel()
		g.Wait()
		systray.Quit()
	}()
}
