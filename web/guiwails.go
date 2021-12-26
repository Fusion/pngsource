//go:build !webview
// +build !webview

package web

import (
	_ "embed"
	"strings"

	"github.com/fusion/pngsource/assets"
	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails"
)

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
	app.Bind(h)
	app.Run()
}
