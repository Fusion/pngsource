//go:build !webview
// +build !webview

package web

import (
	_ "embed"
	"strings"

	"github.com/fusion/pngsource/assets"
	"github.com/wailsapp/wails"
)

type AppHandler struct {
	runtime *wails.Runtime
}

var localHandler *AppHandler

func (h *AppHandler) WailsInit(r *wails.Runtime) error {
	h.runtime = r
	return nil
}

func SelectSavePath(writeFileName string) string {
	return localHandler.runtime.Dialog.SelectSaveFile("Save PNG File", "*.png")
}

func MaybeExecute(payload string) {
	localHandler.runtime.Events.Emit("eval", payload)
}

func Instantiate(
	debug bool,
	destPath *string,
	h GUIHandler) {

	rawhtml, _ := assets.Content.ReadFile("index.html")
	html := strings.Replace(string(rawhtml), "{{STYLE}}", "", -1)
	css, _ := assets.Content.ReadFile("css/style.css")
	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "pngsource",
		HTML:   html,
		CSS:    string(css),
		Colour: "#131313",
	})

	var itsAnInterface interface{} = &h
	localHandler, _ := itsAnInterface.(AppHandler)

	// This is a bit gross.
	// I am registering twice:
	// - as a localhandler, so that I can get a ref' to runtime
	// - as a GUI handler, to wire my methods
	app.Bind(&localHandler)
	app.Bind(h)

	app.Run()
}
