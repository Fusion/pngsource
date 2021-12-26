package web

type GUIHandler interface {
	Wlog(string)
	Wupdatedestfolderpref(string)
	Wrawimage(string, string) string
	Wsourcecode(string, string) string
	Wembedcode(string, string, string, string, string) bool
}

type AppHandler struct {
	States map[string]string
}
