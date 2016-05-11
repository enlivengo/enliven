package plugins

import (
	"net/http"

	"github.com/hickeroar/enliven"
	"github.com/rakyll/statik/fs"
)

// NewStatikAssetsPlugin Creates a new embedded statik asset plugin
func NewStatikAssetsPlugin(suppliedConfig enliven.Config) *StatikAssetsPlugin {
	var config = enliven.Config{
		"statik.assets.route": "/statik/",
	}

	config = enliven.MergeConfig(config, suppliedConfig)

	return &StatikAssetsPlugin{
		route: config["statik.assets.route"],
	}
}

// StatikAssetsPlugin handles adding a route handler for static assets
type StatikAssetsPlugin struct {
	route string
}

// Initialize sets up our plugin to handle embedded static asset requests
func (sap *StatikAssetsPlugin) Initialize(ev *enliven.Enliven) {
	router := ev.GetRouter()
	statikFS, _ := fs.New()
	handler := http.StripPrefix(sap.route, http.FileServer(statikFS))
	router.PathPrefix(sap.route).Handler(handler)
}
