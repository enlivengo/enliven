package admin

//go:generate go-bindata -o views.go -pkg admin views/...

import (
	"github.com/hickeroar/enliven"
	"github.com/hickeroar/enliven/apps/database"
	"github.com/qor/qor"
)

var adminResources []interface{}

// AddResources adds models to qor/admin
func AddResources(resources ...interface{}) {
	for _, res := range resources {
		adminResources = append(adminResources, res)
	}
}

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
	if !ev.AppInstalled("default_database") {
		panic("The Admin app requires that the Database app is initialized with a default connection.")
	}

	db := database.GetDatabase(ev, "default")

	admin := New(&qor.Config{DB: db})

	for _, resource := range adminResources {
		admin.AddResource(resource)
	}

	//admin.SetAuth(&AdminAuth{})
	admin.MountTo("/admin", ev.GetRouter())

	ev.RestrictRoute("/admin")

	ev.RegisterService("admin", admin)
}

// GetName returns the app's name
func (aa *App) GetName() string {
	return "admin"
}
