package user

import (
	"html/template"

	"github.com/hickeroar/enliven"
)

// getTemplate looks up a template in config or embedded assets and returns its contents
func getTemplate(ctx *enliven.Context, templateType string) string {
	config := ctx.Enliven.GetConfig()

	requestedTemplate := config["user."+templateType+".template"]

	if requestedTemplate == "" {
		temp, _ := Asset("templates/" + templateType + ".html")

		if len(temp) > 0 {
			requestedTemplate = string(temp[:])
		}
	}

	return requestedTemplate
}

// LoginGetHandler handles get requests to the login route
func LoginGetHandler(ctx *enliven.Context) {
	tmpl, _ := template.New("LoginGetHandler").Parse(getTemplate(ctx, "login"))
	err := tmpl.Execute(ctx.Response, nil)
	if err != nil {
		ctx.String(err.Error())
	}
}
