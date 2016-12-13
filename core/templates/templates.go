package templates

import (
	"html/template"

	"github.com/enlivengo/enliven/core/templates/files"
)

//go:generate go-bindata -ignore \.go -o files/files.go -pkg files files/...

// TemplateManager manages our templates
type TemplateManager struct {
	BaseTemplate *template.Template
	Templates    map[string]*template.Template
}

// CreateTemplate duplicates the base template, parses templates text into it, and stores it as a new template
func (tm TemplateManager) CreateTemplate(name string, text string) {
	newTemplate, _ := tm.BaseTemplate.Clone()
	newTemplate.Parse(text)
	tm.Templates[name] = newTemplate
}

// NewTemplateManager returns an instance of our temlate manager
func NewTemplateManager() TemplateManager {
	headerTemplate, _ := files.Asset("files/header.html")
	footerTemplate, _ := files.Asset("files/footer.html")
	homeTemplate, _ := files.Asset("files/home.html")
	forbiddenTemplate, _ := files.Asset("files/forbidden.html")
	notfoundTemplate, _ := files.Asset("files/notfound.html")
	badrequestTemplate, _ := files.Asset("files/badrequest.html")

	baseTemplate := template.New("enliven")
	baseTemplate.Parse(string(headerTemplate[:]))
	baseTemplate.Parse(string(footerTemplate[:]))

	// These are the core full templates which can be called directly with ctx.ExecuteBaseTemplate
	baseTemplate.Parse(string(homeTemplate[:]))
	baseTemplate.Parse(string(forbiddenTemplate[:]))
	baseTemplate.Parse(string(notfoundTemplate[:]))
	baseTemplate.Parse(string(badrequestTemplate[:]))

	tm := TemplateManager{
		BaseTemplate: baseTemplate,
		Templates:    make(map[string]*template.Template),
	}

	return tm
}
