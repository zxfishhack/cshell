package main

import (
	"github.com/zxfishhack/cshell/pkg/store"
	"github.com/zxfishhack/cshell/pkg/tray"
	"github.com/zxfishhack/cshell/pkg/web"
	"log"
)

//go:generate go-bindata -pkg icon -tags "linux darwin" -o pkg/icon/iconunix.go -prefix icon icon
//go:generate go-bindata -pkg resources -o pkg/resources/resouces.go -prefix console/dist console/dist/...
func main() {
	err := store.Init()
	if err != nil {
		log.Panicln(err)
	}
	port, err := web.InitWebService()
	if err != nil {
		log.Panicln(err)
	}
	tray.New(port)
	tray.Inst.Run()
}
