package assets

import (
	"net/http"

	"github.com/hickeroar/enliven"
)

// NewStaticPlugin Creates a new static asset plugin instance
func NewStaticPlugin(suppliedConfig enliven.Config) *StaticPlugin {
	var config = enliven.Config{
		"assets.static.route": "/static/",
		"assets.static.path":  "./static/",
	}

	config = enliven.MergeConfig(config, suppliedConfig)

	return &StaticPlugin{
		route: config["assets.static.route"],
		path:  config["assets.static.path"],
	}
}

// StaticPlugin handles adding a route handler for static assets
type StaticPlugin struct {
	path  string
	route string
}

// Initialize sets up our plugin to ahndle static asset requests
func (sap *StaticPlugin) Initialize(ev *enliven.Enliven) {
	router := ev.GetRouter()
	handler := http.StripPrefix(sap.route, http.FileServer(http.Dir(sap.path)))
	router.PathPrefix(sap.route).Handler(handler)
}

// GetName returns the plugin's name
func (sap *StaticPlugin) GetName() string {
	return "static"
}
