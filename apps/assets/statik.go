package assets

import (
	"net/http"

	"github.com/hickeroar/enliven"
	"github.com/hickeroar/enliven/config"
	"github.com/rakyll/statik/fs"
)

// NewStatikApp Creates a new embedded statik asset app
func NewStatikApp() *StatikApp {
	return &StatikApp{}
}

// StatikApp handles adding a route handler for static assets
type StatikApp struct{}

// Initialize sets up our app to handle embedded static asset requests
func (sa *StatikApp) Initialize(ev *enliven.Enliven) {
	var conf = config.Config{
		"assets_statik_route": "/statik/",
	}

	conf = config.UpdateConfig(config.MergeConfig(conf, config.GetConfig()))

	// Making sure this route ends in a forward slash
	if conf["assets_statik_route"][len(conf["assets_statik_route"])-1:] != "/" {
		conf["assets_statik_route"] += "/"
	}

	router := ev.GetRouter()
	statikFS, _ := fs.New()
	handler := http.StripPrefix(conf["assets_statik_route"], http.FileServer(statikFS))
	router.PathPrefix(conf["assets_statik_route"]).Handler(handler)
}

// GetName returns the apps's name
func (sa *StatikApp) GetName() string {
	return "statik"
}
