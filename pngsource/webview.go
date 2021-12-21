package main

import (
	"encoding/base64"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/fusion/pngsource/assets"
	"github.com/fusion/pngsource/lib"
	"github.com/ncruces/zenity"
	"github.com/rakyll/globalconf"
	"github.com/webview/webview"
	//"github.com/davecgh/go-spew/spew"
)

type StringValue struct {
	value string
}

func (s StringValue) String() string {
	return s.value
}

func (s StringValue) Set(newValue string) error {
	s.value = newValue
	return nil
}

func getActualString(sourceString string) string {
  comma := strings.Index(sourceString, ",")
  if comma == -1 {
    return sourceString
  }
  subcontent := content[comma:]
  debased, err := base64.StdEncoding.DecodeString(subcontent)
  if err != nil {
    l.Println("BASE64 Oopsy.")
    return "*err*"
  }
  return debased
}

func main() {
	debug := true
	l := log.New(os.Stderr, "", 0)

	f := flag.NewFlagSet("webview", flag.ContinueOnError)
	flagDestPath := f.String("dest", "", "destination path")

	globalconf.Register("webview", f)
	preferences, _ := globalconf.New("pngsource")
	preferences.ParseAll()

	css, _ := assets.Content.ReadFile("css/style.css")
	rawpage, _ := assets.Content.ReadFile("index.html")
	almostpage := strings.Replace(
		strings.Replace(
			string(rawpage), "{{STYLE}}", string(css), -1),
		"{{DESTPATH}}", *flagDestPath, -1)
	page := strings.Replace(almostpage, "%", "%25", -1)

	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("PngSource Thingamagig")
	w.SetSize(800, 720, webview.HintNone)

	// These bind statements refer to javascript functions... neat.
	w.Bind("wlog", func(msg string) {
		l.Println(msg)
	})

	w.Bind("wupdatedestfolderpref", func(value string) {
		preferences.Set("webview", &flag.Flag{Name: "dest", Value: StringValue{value}})
	})

	w.Bind("wrawimage", func(action, content string) string {
		return lib.Write_content_from_data(l, action, getActualString(content))
	})

	w.Bind("wsourcecode", func(action, content string) string {
		return lib.Write_content_from_data(l, action, getActualString(content))
	})

	w.Bind("wembedcode", func(
		writeFileName string,
		writepath string,
		sourcetype string,
		sourcepathorcode string,
		destfolder string) {

		// destFolder is not relevant at this time as we are using
		// a file dialog instead.
		selectedPath, err := zenity.SelectFileSave(
			zenity.Title("Save File with Embed"),
			zenity.Filename(writeFileName),
			zenity.FileFilters{{"PNG files", []string{"*.png"}}},
			zenity.ConfirmOverwrite(),
		)
		if err != nil { // e.g. dialog canceled
			w.Eval("resetEmbedButtonState(false)")
			return
		}
		if sourcetype == "string" {
			sourcepath := lib.Write_content_from_data(l, "put", []byte(sourcepathorcode))
			lib.Write_content_dynamic_config(l, writepath, sourcepath, selectedPath)
		} else {
			lib.Write_content_dynamic_config(l, writepath, sourcepathorcode, selectedPath)
		}

		w.Eval("resetEmbedButtonState(true)")
	})

	w.Bind("wtest", func() {
	})

	w.Bind("quit", func() {
		w.Terminate()
	})

	w.Navigate("data:text/html,<!doctype html>" + string(page))
	w.Run()
}
