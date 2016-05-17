package assets

import (
	"net/http"

	"github.com/hickeroar/enliven"
	"github.com/rakyll/statik/fs"
)

// NewStatikApp Creates a new embedded statik asset app
func NewStatikApp() *StatikApp {
	return &StatikApp{}
}

// StatikApp handles adding a route handler for static assets
type StatikApp struct {
	route string
}

// Initialize sets up our app to handle embedded static asset requests
func (sap *StatikApp) Initialize(ev *enliven.Enliven) {
	var config = enliven.Config{
		"assets.statik.route": "/statik/",
	}

	config = enliven.MergeConfig(config, ev.GetConfig())

	// Making sure this route ends in a forward slash
	if config["assets.statik.route"][len(config["assets.statik.route"])-1:] != "/" {
		config["assets.statik.route"] += "/"
	}

	router := ev.GetRouter()
	statikFS, _ := fs.New()
	handler := http.StripPrefix(config["assets.statik.route"], http.FileServer(statikFS))
	router.PathPrefix(config["assets.statik.route"]).Handler(handler)
}

// GetName returns the apps's name
func (sap *StatikApp) GetName() string {
	return "statik"
}
