package plugins

import (
	"net/http"

	"github.com/hickeroar/enliven"
)

// NewStaticAssetsPlugin Creates a new static asset plugin instance
func NewStaticAssetsPlugin(suppliedConfig map[string]string) *StaticAssetsPlugin {
	var config = map[string]string{
		"static.assets.route": "/static/",
		"static.assets.path":  "./static/",
	}

	config = enliven.MergeConfig(config, suppliedConfig)

	return &StaticAssetsPlugin{
		route: config["static.assets.route"],
		path:  config["static.assets.path"],
	}
}

// StaticAssetsPlugin handles adding a route handler for static assets
type StaticAssetsPlugin struct {
	path  string
	route string
}

// Initialize sets up our plugin to ahndle static asset requests
func (sap *StaticAssetsPlugin) Initialize(ev *enliven.Enliven) {
	router := ev.GetRouter()
	handler := http.StripPrefix(sap.route, http.FileServer(http.Dir(sap.path)))
	router.PathPrefix(sap.route).Handler(handler)
}
