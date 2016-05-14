package assets

import (
	"net/http"

	"github.com/hickeroar/enliven"
)

// NewStaticApp Creates a new static asset app instance
func NewStaticApp(suppliedConfig enliven.Config) *StaticApp {
	var config = enliven.Config{
		"assets.static.route": "/static/",
		"assets.static.path":  "./static/",
	}

	config = enliven.MergeConfig(config, suppliedConfig)

	// Making sure this route ends in a forward slash
	if config["assets.static.route"][len(config["assets.static.route"])-1:] != "/" {
		config["assets.static.route"] += "/"
	}

	return &StaticApp{
		route: config["assets.static.route"],
		path:  config["assets.static.path"],
	}
}

// StaticApp handles adding a route handler for static assets
type StaticApp struct {
	path  string
	route string
}

// Initialize sets up our app to ahndle static asset requests
func (sap *StaticApp) Initialize(ev *enliven.Enliven) {
	router := ev.GetRouter()
	handler := http.StripPrefix(sap.route, http.FileServer(http.Dir(sap.path)))
	router.PathPrefix(sap.route).Handler(handler)
}

// GetName returns the app's name
func (sap *StaticApp) GetName() string {
	return "static"
}
