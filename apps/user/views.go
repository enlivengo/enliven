package user

import (
	"html/template"

	"github.com/hickeroar/enliven"
	"github.com/hickeroar/enliven/apps/database"
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
	tmpl, _ := template.New("User_LoginGetHandler").Parse(getTemplate(ctx, "login"))
	ctx.Template(tmpl)
}

// LoginPostHandler handles the form submission for logging a user in.
func LoginPostHandler(ctx *enliven.Context) {
	ctx.Request.ParseForm()
	username := ctx.Request.Form.Get("username")
	password := ctx.Request.Form.Get("password")

	config := ctx.Enliven.GetConfig()
	db := database.GetDatabase(ctx, config["user.database.namespace"])

	user := User{}
	db.Where("Login = ?", username).First(&user)

	if user.ID == 0 || !VerifyPasswordHash(password, user.Password) {
		tmpl, _ := template.New("User_LoginGetHandler").Parse(getTemplate(ctx, "login"))
		ctx.Template(tmpl)
		return
	}

	ctx.Redirect(config["user.login.redirect"])
}
