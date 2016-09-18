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
	badrequestTemplate, _ := files.Asset("files/badrequest.html")

	baseTemplate := template.New("enliven")
	baseTemplate.Parse(string(headerTemplate[:]))
	baseTemplate.Parse(string(footerTemplate[:]))
	baseTemplate.Parse(string(homeTemplate[:]))
	baseTemplate.Parse(string(forbiddenTemplate[:]))
	baseTemplate.Parse(string(notfoundTemplate[:]))
	baseTemplate.Parse(string(badrequestTemplate[:]))

	return baseTemplate
}
