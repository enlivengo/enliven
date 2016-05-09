package plugins

import (
	"net/http"

	"github.com/hickeroar/enliven"
	"github.com/rakyll/statik/fs"
)

// NewStatikAssetPlugin Creates a new embedded statik asset plugin
func NewStatikAssetPlugin(route string) *StatikAssetPlugin {
	return &StatikAssetPlugin{
		route: route,
	}
}

// StatikAssetPlugin handles adding a route handler for static assets
type StatikAssetPlugin struct {
	route string
}

// Initialize sets up our plugin to handle static asset requests
func (sap *StatikAssetPlugin) Initialize(ev *enliven.Enliven) {
	router := ev.GetRouter()
	statikFS, _ := fs.New()
	handler := http.StripPrefix(sap.route, http.FileServer(statikFS))
	router.PathPrefix(sap.route).Handler(handler)
}
