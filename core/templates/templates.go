package templates

import (
	"html/template"

	"github.com/enlivengo/enliven/core/templates/files"
)

//go:generate go-bindata -o files/files.go -pkg files files/...

// New returns an instance of core templates
func New() *template.Template {
	headerTemplate, _ := files.Asset("files/header.html")
	footerTemplate, _ := files.Asset("files/footer.html")
	homeTemplate, _ := files.Asset("files/home.html")
	forbiddenTemplate, _ := files.Asset("files/forbidden.html")
	notfoundTemplate, _ := files.Asset("files/notfound.html")

	templates := template.New("enliven")
	templates.Parse(string(headerTemplate[:]))
	templates.Parse(string(footerTemplate[:]))
	templates.Parse(string(homeTemplate[:]))
	templates.Parse(string(forbiddenTemplate[:]))
	templates.Parse(string(notfoundTemplate[:]))

	return templates
}
