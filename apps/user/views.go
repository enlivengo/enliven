package user

//go:generate go-bindata -o templates.go -pkg user templates/...

import (
	"strconv"

	"github.com/hickeroar/enliven"
	"github.com/hickeroar/enliven/apps/database"
)

// LoginGetHandler handles get requests to the login route
func LoginGetHandler(ctx *enliven.Context) {
	templates := ctx.Enliven.GetTemplates()
	ctx.NamedTemplate(templates, "login")
}

// LoginPostHandler handles the form submission for logging a user in.
func LoginPostHandler(ctx *enliven.Context) {
	ctx.Request.ParseForm()
	username := ctx.Request.Form.Get("username")
	password := ctx.Request.Form.Get("password")

	config := ctx.Enliven.GetConfig()
	db := database.GetDatabase(ctx.Enliven)

	user := User{}
	db.Where("Login = ?", username).First(&user)

	if user.ID == 0 || !VerifyPasswordHash(user.Password, password) {
		LoginGetHandler(ctx)
		return
	}

	ctx.Session.Set("user_id", strconv.FormatUint(uint64(user.ID), 10))
	ctx.Redirect(config["user_login_redirect"])
}

// LogoutHandler logs a user out and redirects them to the configured page.
func LogoutHandler(ctx *enliven.Context) {
	ctx.Session.Destroy()
	ctx.Redirect(ctx.Enliven.GetConfig()["user_logout_redirect"])
}
