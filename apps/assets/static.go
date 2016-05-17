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
type StaticApp struct {
	path  string
	route string
}

// Initialize sets up our app to ahndle static asset requests
func (sap *StaticApp) Initialize(ev *enliven.Enliven) {
	var config = enliven.Config{
		"assets.static.route": "/static/",
		"assets.static.path":  "./static/",
	}

	config = enliven.MergeConfig(config, ev.GetConfig())

	// Making sure this route ends in a forward slash
	if config["assets.static.route"][len(config["assets.static.route"])-1:] != "/" {
		config["assets.static.route"] += "/"
	}

	router := ev.GetRouter()
	handler := http.StripPrefix(config["assets.static.route"], http.FileServer(http.Dir(config["assets.static.path"])))
	router.PathPrefix(config["assets.static.route"]).Handler(handler)
}

// GetName returns the app's name
func (sap *StaticApp) GetName() string {
	return "static"
}
