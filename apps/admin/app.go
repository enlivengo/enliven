package admin

//go:generate go-bindata -o views.go -pkg admin views/...

import (
	"github.com/hickeroar/enliven"
	"github.com/hickeroar/enliven/apps/database"
	"github.com/hickeroar/enliven/apps/user"
	"github.com/qor/qor"
)

// GetAdmin returns our instance of qor/admin
func GetAdmin(ev *enliven.Enliven) *Admin {
	if a, ok := ev.GetService("admin").(*Admin); ok {
		return a
	}
	return nil
}

// NewApp generates and returns an instance of the app
func NewApp() *App {
	return &App{}
}

// App is the admin application
type App struct {
}

// Initialize sets up the qor/admin module
func (aa *App) Initialize(ev *enliven.Enliven) {
	if !ev.AppInstalled("user") {
		panic("The Admin app requires that the User app is initialized first.")
	}

	db := database.GetDatabase(ev, "default")

	admin := New(&qor.Config{DB: db})
	admin.AddResource(&user.User{})
	admin.AddResource(&user.Group{})
	admin.AddResource(&user.Permission{})

	//admin.SetAuth(&AdminAuth{})
	admin.MountTo("/admin/", ev.GetRouter())

	ev.RegisterService("admin", admin)
}

// GetName returns the app's name
func (aa *App) GetName() string {
	return "admin"
}
