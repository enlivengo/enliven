package statik

import (
	"net/http"

	"github.com/enlivengo/enliven"
	"github.com/enlivengo/enliven/config"
	"github.com/rakyll/statik/fs"
)

// NewApp Creates a new embedded statik asset app
func NewApp() *App {
	return &App{}
}

// App handles adding a route handler for static assets
type App struct{}

// Initialize sets up our app to handle embedded static asset requests
func (sa *App) Initialize(ev *enliven.Enliven) {
	var conf = config.Config{
		"assets_statik_route": "/statik/",
	}

	conf = config.UpdateConfig(config.MergeConfig(conf, config.GetConfig()))

	// Making sure this route ends in a forward slash
	if conf["assets_statik_route"][len(conf["assets_statik_route"])-1:] != "/" {
		conf["assets_statik_route"] += "/"
	}

	statikFS, _ := fs.New()
	handler := http.StripPrefix(conf["assets_statik_route"], http.FileServer(statikFS))
	ev.Router.PathPrefix(conf["assets_statik_route"]).Handler(handler)
}

// GetName returns the apps's name
func (sa *App) GetName() string {
	return "statik"
}
