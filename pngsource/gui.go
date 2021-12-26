package main

import (
	"encoding/base64"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/fusion/pngsource/lib"
	"github.com/fusion/pngsource/web"
	"github.com/rakyll/globalconf"
	//"github.com/davecgh/go-spew/spew"
)

var l *log.Logger
var preferences *globalconf.GlobalConf

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

func getActualBytes(sourceString string) []byte {
	comma := strings.Index(sourceString, ",")
	if comma == -1 {
		return []byte(sourceString)
	}
	subcontent := sourceString[comma+1:]
	debased, err := base64.StdEncoding.DecodeString(subcontent)
	if err != nil {
		l.Println("BASE64 Oopsy.")
		return []byte("*err*")
	}
	return debased
}

// "Localize" type so that we can use it as a receiver
type AppHandler web.AppHandler

func (h *AppHandler) Wlog(msg string) {
	l.Println(msg)
}

func (h *AppHandler) Wupdatedestfolderpref(value string) {
	preferences.Set("webview", &flag.Flag{Name: "dest", Value: StringValue{value}})
}

func (h *AppHandler) Wrawimage(action string, content string) string {
	return lib.Write_content_from_data(l, action, getActualBytes(content))
}

func (h *AppHandler) Wsourcecode(action string, content string) string {
	return lib.Write_content_from_data(l, action, getActualBytes(content))
}

func (h *AppHandler) Wembedcode(
	writeFileName string,
	writepath string,
	sourcetype string,
	sourcepathorcode string,
	destfolder string) bool {

	// destFolder is not relevant at this time as we are using
	// a file dialog instead.
	selectedPath := web.SelectSavePath(writeFileName)
	if selectedPath == "" {
		return false
	}
	if sourcetype == "string" {
		sourcepath := lib.Write_content_from_data(l, "put", []byte(sourcepathorcode))
		lib.Write_content_dynamic_config(l, writepath, sourcepath, selectedPath)
	} else {
		lib.Write_content_dynamic_config(l, writepath, sourcepathorcode, selectedPath)
	}

	return true
}

func main() {
	debug := true
	l = log.New(os.Stderr, "", 0)

	f := flag.NewFlagSet("gui", flag.ContinueOnError)
	flagDestPath := f.String("dest", "", "destination path")

	globalconf.Register("gui", f)
	preferences, _ := globalconf.New("pngsource")
	preferences.ParseAll()

	web.Instantiate(
		debug,
		flagDestPath,
		&AppHandler{States: make(map[string]string)})
}
