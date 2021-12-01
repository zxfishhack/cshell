package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/zxfishhack/cshell/pkg/resources"
	"mime"
	"net/http"
	"path/filepath"
)

type ConsoleController struct {
}

func (c *ConsoleController) GetByWildcard(ctx iris.Context, path string) error {
	if path == "" || path == "/" {
		path = "index.html"
	}
	if _, err := resources.Asset(path); err != nil {
		path = "index.html"
	}
	if b, err := resources.Asset(path); err != nil {
		ctx.StatusCode(http.StatusNotFound)
	} else {
		if path != "index.html" {
			ctx.Header("Cache-Control", "public, max-age=604800, immutable")
		}
		ctx.ServeFile(path, false)
		ctx.ContentType(mime.TypeByExtension(filepath.Ext(path)))
		ctx.Write(b)
	}
	return nil
}
