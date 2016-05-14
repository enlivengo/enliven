package assets

import (
	"net/http"

	"github.com/hickeroar/enliven"
	"github.com/rakyll/statik/fs"
)

// NewStatikPlugin Creates a new embedded statik asset plugin
func NewStatikPlugin(suppliedConfig enliven.Config) *StatikPlugin {
	var config = enliven.Config{
		"assets.statik.route": "/statik/",
	}

	config = enliven.MergeConfig(config, suppliedConfig)

	return &StatikPlugin{
		route: config["assets.statik.route"],
	}
}

// StatikPlugin handles adding a route handler for static assets
type StatikPlugin struct {
	route string
}

// Initialize sets up our plugin to handle embedded static asset requests
func (sap *StatikPlugin) Initialize(ev *enliven.Enliven) {
	router := ev.GetRouter()
	statikFS, _ := fs.New()
	handler := http.StripPrefix(sap.route, http.FileServer(statikFS))
	router.PathPrefix(sap.route).Handler(handler)
}

// GetName returns the plugin's name
func (sap *StatikPlugin) GetName() string {
	return "statik"
}
