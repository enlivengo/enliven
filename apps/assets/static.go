package assets

import (
	"net/http"

	"github.com/hickeroar/enliven"
)

// NewStaticApp Creates a new static asset app instance
func NewStaticApp() *StaticApp {
	return &StaticApp{}
}

// StaticApp handles adding a route handler for static assets
type StaticApp struct{}

// Initialize sets up our app to ahndle static asset requests
func (sa *StaticApp) Initialize(ev *enliven.Enliven) {
	var config = enliven.Config{
		"assets_static_route": "/static/",
		"assets_static_path":  "./static/",
	}

	config = enliven.MergeConfig(config, ev.GetConfig())

	// Making sure this route ends in a forward slash
	if config["assets_static_route"][len(config["assets_static_route"])-1:] != "/" {
		config["assets_static_route"] += "/"
	}

	router := ev.GetRouter()
	handler := http.StripPrefix(config["assets_static_route"], http.FileServer(http.Dir(config["assets_static_path"])))
	router.PathPrefix(config["assets_static_route"]).Handler(handler)
}

// GetName returns the app's name
func (sa *StaticApp) GetName() string {
	return "static"
}
