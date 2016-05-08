package plugins

import (
	"net/http"

	"github.com/hickeroar/enliven"
)

// NewStaticAssetPlugin Creates a new static asset middleware item and
func NewStaticAssetPlugin(route string, path string) *StaticAssetPlugin {
	return &StaticAssetPlugin{
		route: route,
		path:  path,
	}
}

// StaticAssetPlugin handles adding a route handler for static assets
type StaticAssetPlugin struct {
	path  string
	route string
}

// Initialize sets up our plugin to ahndle static asset requests
func (sap *StaticAssetPlugin) Initialize(ev *enliven.Enliven) {
	router := ev.GetRouter()
	handler := http.StripPrefix(sap.route, http.FileServer(http.Dir(sap.path)))
	router.PathPrefix(sap.route).Handler(handler)
}
