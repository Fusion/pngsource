//go:build webview
// +build webview

package web

import (
	"strings"

	"github.com/fusion/pngsource/assets"
	"github.com/ncruces/zenity"
	"github.com/webview/webview"
)

type AppHandler struct {
}

func SelectSavePath(writeFileName string) string {
	selectedPath, err := zenity.SelectFileSave(
		zenity.Title("Save File with Embed"),
		zenity.Filename(writeFileName),
		zenity.FileFilters{{
			Name:     "PNG files",
			Patterns: []string{"*.png"}}},
		zenity.ConfirmOverwrite(),
	)
	if err != nil { // e.g. dialog canceled
		return ""
	}
	return selectedPath
}

var w webview.WebView

func MaybeExecute(payload string) {
	w.Eval(payload)
}

func Instantiate(
	debug bool,
	destPath *string,
	h GUIHandler) {
	css, _ := assets.Content.ReadFile("css/style.css")
	rawpage, _ := assets.Content.ReadFile("index.html")
	almostpage := strings.Replace(
		strings.Replace(
			string(rawpage), "{{STYLE}}", string(css), -1),
		"{{DESTPATH}}", *destPath, -1)
	page := strings.Replace(almostpage, "%", "%25", -1)

	w = webview.New(debug)
	defer w.Destroy()

	w.Bind("wlog", h.Wlog)
	w.Bind("wupdatedestfolderpref", h.Wupdatedestfolderpref)
	w.Bind("wrawimage", h.Wrawimage)
	w.Bind("wsourcecode", h.Wsourcecode)
	w.Bind("wembedcode", h.Wembedcode)

	w.SetTitle("PngSource Thingamagig")
	w.SetSize(800, 720, webview.HintNone)

	w.Navigate("data:text/html,<!doctype html>" + string(page))
	w.Run()
}
