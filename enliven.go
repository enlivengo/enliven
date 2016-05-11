package enliven

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	// Adding DB requirements.
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// We'll use this to create and insert the initial enliven instance into the handlers
var enliven Enliven

// --------------------------------------------------

// MergeConfig takes a default config and merges a supplied one into it.
func MergeConfig(defaultConfig map[string]string, suppliedConfig map[string]string) map[string]string {
	for key, value := range suppliedConfig {
		defaultConfig[key] = value
	}

	return defaultConfig
}

// --------------------------------------------------

// IPlugin is an interface for writing Enliven plugins
// Plugins are basically packaged enliven setup code
type IPlugin interface {
	Initialize(ev *Enliven)
}

// ISession represents a session that session middleware must implement
type ISession interface {
	Set(key string, value string) error
	Get(key string) string
	Delete(key string) error
	Destroy() error
	SessionID() string
}

// IMiddlewareHandler is an interface to be used when writing Middleware
// Copied w/ alterations from github.com/codegangsta/negroni
type IMiddlewareHandler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context, next NextHandlerFunc)
}

// --------------------------------------------------

// CHandler Handles injecting the initial request context before passing handling on to the Middleware struct
type CHandler func(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context)

// ServeHTTP is the first handler that gets hit when a request comes in.
func (ch CHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Items: make(map[string]string),
	}
	ch(rw, r, enliven, ctx)
}

// ContextHandler sets up serving the first request, and the handing off of subsequent requests to the Middleware struct
func ContextHandler(h Middleware) CHandler {
	return CHandler(func(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context) {
		h.handler.ServeHTTP(rw, r, enliven, ctx, h.next.ServeHTTP)
	})
}

// --------------------------------------------------

// NextHandlerFunc allow use of ordinary functions middleware handlers
// Copied w/ alterations from github.com/codegangsta/negroni
type NextHandlerFunc func(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context)

// Copied w/ alterations from github.com/codegangsta/negroni
func (nh NextHandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context) {
	nh(rw, r, ev, ctx)
}

// --------------------------------------------------

// HandlerFunc allow use of ordinary functions middleware handlers
// Copied w/ alterations from github.com/codegangsta/negroni
type HandlerFunc func(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context, next NextHandlerFunc)

// Copied w/ alterations from github.com/codegangsta/negroni
func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context, next NextHandlerFunc) {
	h(rw, r, ev, ctx, next)
}

// --------------------------------------------------

// RouteHandlerFunc is an interface to be used when writing route handler functions
type RouteHandlerFunc func(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context)

// --------------------------------------------------

// Context stores context variables and the session that will be passed to requests
type Context struct {
	Session ISession
	Items   map[string]string
}

// --------------------------------------------------

// Middleware Represents a piece of middlewear
// Copied w/ alterations from github.com/codegangsta/negroni
type Middleware struct {
	handler IMiddlewareHandler
	next    *Middleware
}

// Copied w/ alterations from github.com/codegangsta/negroni
func (m Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context) {
	m.handler.ServeHTTP(rw, r, ev, ctx, m.next.ServeHTTP)
}

// --------------------------------------------------

// Enliven is....Enliven
type Enliven struct {
	services      map[string]interface{}
	routeHandlers map[string]RouteHandlerFunc
	middleware    Middleware
	handlers      []IMiddlewareHandler
}

// New (constructor) gets a new instance of enliven.
func New(config map[string]string) *Enliven {
	enliven = Enliven{
		services:      make(map[string]interface{}),
		routeHandlers: make(map[string]RouteHandlerFunc),
	}

	enliven.RegisterService("router", mux.NewRouter())
	enliven.registerConfig(config)
	enliven.registerDatabase()

	return &enliven
}

// addConfig created and registers the app config
func (ev *Enliven) registerConfig(suppliedConfig map[string]string) {
	var enlivenConfig = map[string]string{
		"db.driver":     "",
		"db.host":       "",
		"db.user":       "",
		"db.dbname":     "",
		"db.password":   "",
		"db.sslmode":    "disable",
		"db.port":       "",
		"db.connString": "",

		"server.port": "8000",
	}

	ev.RegisterService("config", MergeConfig(enlivenConfig, suppliedConfig))
}

// addDatabase Initializes a database given the values from the EnlivenConfig
func (ev *Enliven) registerDatabase() {
	config := ev.GetConfig()

	var driver string
	allowedDrivers := [4]string{"postgres", "mysql", "sqlite3", "mssql"}

	// Making sure the specified driver is in the list if allowed drivers
	for i := 0; i < 4; i++ {
		if allowedDrivers[i] == config["db.driver"] {
			driver = config["db.driver"]
			break
		}
	}

	// If we didn't set a driver, we return here.
	if driver == "" {
		return
	}

	var connString string

	// Someone can specify a whole connection string, or the parts of it
	if config["db.connString"] != "" {
		connString = config["db.connString"]
	} else {
		// driver specific connection string addons
		switch driver {

		case "sqlite3":
			// If the driver is sqlite3, but there wasn't a conn string, we return.
			if config["db.connString"] == "" {
				return
			}

		case "mysql", "mssql":
			connString = config["db.user"] + ":" + config["db.password"] + "@" + config["db.host"]

			// Adding a port if one was provided
			if len(config["db.port"]) > 0 {
				connString += ":" + config["db.port"]
			}

			connString += "/" + config["db.dbname"]

			if driver == "mysql" {
				connString += "?charset=utf8&parseTime=True&loc=Local"
			}

		case "postgres":
			var connStringParts []string
			connStringParts = append(connStringParts, "host="+config["db.host"])
			connStringParts = append(connStringParts, "user="+config["db.user"])
			connStringParts = append(connStringParts, "dbname="+config["db.dbname"])
			connStringParts = append(connStringParts, "sslmode="+config["db.sslmode"])
			connStringParts = append(connStringParts, "password="+config["db.password"])

			if len(config["db.port"]) > 0 {
				connStringParts = append(connStringParts, "port="+config["db.port"])
			}

			connString = strings.Join(connStringParts, " ")
		}
	}

	db, err := gorm.Open(driver, connString)

	// Making sure we got a database instance
	if err != nil {
		panic(err)
	}

	// Making sure we can ping the database
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	ev.RegisterService("database", db)
}

// RegisterService registers an enliven service or dependency
func (ev *Enliven) RegisterService(name string, service interface{}) {
	if _, ok := ev.services[name]; ok {
		panic("The service name you are attempting to register has already been registered.")
	}

	ev.services[name] = service
}

// GetService returns an enliven service or dependency
func (ev *Enliven) GetService(name string) interface{} {
	if _, ok := ev.services[name]; ok {
		return ev.services[name]
	}
	return nil
}

// InitPlugin initializes a provided plugin
func (ev *Enliven) InitPlugin(plugin IPlugin) {
	plugin.Initialize(ev)
}

// GetDatabase Gets an instance of the database
func (ev *Enliven) GetDatabase() *gorm.DB {
	if db, ok := ev.GetService("database").(*gorm.DB); ok {
		return db
	}
	return nil
}

// GetConfig Gets an instance of the config
func (ev *Enliven) GetConfig() map[string]string {
	config := ev.GetService("config").(map[string]string)
	return config
}

// GetRouter Gets the instance of the router
func (ev *Enliven) GetRouter() *mux.Router {
	router := ev.GetService("router").(*mux.Router)
	return router
}

// AddMiddleware adds a Handler onto the middleware stack.
// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) AddMiddleware(handler IMiddlewareHandler) {
	ev.handlers = append(ev.handlers, handler)
	ev.middleware = ev.buildMiddleware(ev.handlers)
}

// AddMiddlewareFunc adds a HandlerFunc onto the middleware stack.
// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) AddMiddlewareFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context, next NextHandlerFunc)) {
	ev.AddMiddleware(HandlerFunc(handlerFunc))
}

// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) buildMiddleware(handlers []IMiddlewareHandler) Middleware {
	var next Middleware

	voidMiddleware := Middleware{
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context, next NextHandlerFunc) {}), &Middleware{},
	}

	if len(handlers) == 0 {
		return voidMiddleware
	} else if len(handlers) > 1 {
		next = ev.buildMiddleware(handlers[1:])
	} else {
		next = voidMiddleware
	}

	return Middleware{handlers[0], &next}
}

// AddRoute Registers a handler for a given route.
// We register a dummy route with mux, and then store the provided handler
// which we'll use later in order to inject dependencies into the handler func.
func (ev *Enliven) AddRoute(path string, rhf func(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context)) *mux.Route {
	ev.routeHandlers[path] = RouteHandlerFunc(rhf)
	return ev.GetRouter().HandleFunc(path, func(http.ResponseWriter, *http.Request) {})
}

// Copied from github.com/gorilla/mux
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

// Copied w/ many alterations from github.com/gorilla/mux
// Hijacks the abilities of mux to add our DI handling to route handlers
func routeHandlerFunc(rw http.ResponseWriter, r *http.Request, ev Enliven, ctx *Context, next NextHandlerFunc) {

	// Clean path to canonical form and redirect.
	if p := cleanPath(r.URL.Path); p != r.URL.Path {
		url := *r.URL
		url.Path = p
		p = url.String()

		rw.Header().Set("Location", p)
		rw.WriteHeader(http.StatusMovedPermanently)
		return
	}

	if !ev.GetRouter().KeepContext {
		defer context.Clear(r)
	}

	var match mux.RouteMatch
	var handler http.Handler
	if enliven.GetRouter().Match(r, &match) {
		handler = match.Handler
	}
	if handler == nil {
		handler := http.NotFoundHandler()
		handler.ServeHTTP(rw, r)
	} else {
		// We use the request path to look up our stored route handler if it exists
		if routeHandler, ok := ev.routeHandlers[r.URL.Path]; ok {
			// Calling the route handler with Enliven passed in.
			routeHandler(rw, r, ev, ctx)
		} else {
			// Using handler request handling otherwise.
			handler.ServeHTTP(rw, r)
		}
	}

	next(rw, r, ev, ctx)
}

// Run executes the Enliven http server
func (ev *Enliven) Run(port string) {
	// Adding our route handler as the last piece of middleware
	ev.AddMiddlewareFunc(routeHandlerFunc)

	fmt.Println("Server is listening on port " + port + ".")
	http.ListenAndServe(":"+port, ContextHandler(ev.middleware))
	fmt.Println("Server has shut down.")
}
