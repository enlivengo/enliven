package static

import (
	"net/http"

	"github.com/enlivengo/enliven"
	"github.com/enlivengo/enliven/config"
)

// NewApp Creates a new static asset app instance
func NewApp() *App {
	return &App{}
}

// App handles adding a route handler for static assets
type App struct{}

// Initialize sets up our app to ahndle static asset requests
func (sa *App) Initialize(ev *enliven.Enliven) {
	var conf = config.Config{
		"assets_static_route": "/static/",
		"assets_static_path":  "./static/", // Path relative to the executable
	}

	conf = config.UpdateConfig(config.MergeConfig(conf, config.GetConfig()))

	// Making sure this route ends in a forward slash
	if conf["assets_static_route"][len(conf["assets_static_route"])-1:] != "/" {
		conf["assets_static_route"] += "/"
	}

	handler := http.StripPrefix(conf["assets_static_route"], http.FileServer(http.Dir(conf["assets_static_path"])))
	ev.Router.PathPrefix(conf["assets_static_route"]).Handler(handler)
}

// GetName returns the app's name
func (sa *App) GetName() string {
	return "static"
}
