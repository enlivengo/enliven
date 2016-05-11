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

// Config represents string kvps of application configuration
type Config map[string]string

// MergeConfig takes a default config and merges a supplied one into it.
func MergeConfig(defaultConfig Config, suppliedConfig Config) Config {
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
	ServeHTTP(*Context, NextHandlerFunc)
}

// --------------------------------------------------

// Context stores context variables and the session that will be passed to requests
type Context struct {
	Session  ISession
	Items    map[string]string
	Enliven  Enliven
	Response http.ResponseWriter
	Request  *http.Request
}

// String sets up string headers and outputs a string response
func (ctx *Context) String(output string) {
	ctx.Response.Header().Set("Content-Type", "text/plain")
	ctx.Response.Write([]byte(output))
}

// HTML sets up HTML headers and outputs an HTML response
func (ctx *Context) HTML(output string) {
	ctx.Response.Header().Set("Content-Type", "text/html")
	ctx.Response.Write([]byte(output))
}

// JSON sets up JSON headers and outputs a JSON response
// Expects to recieve the result of json marshalling ([]byte)
func (ctx *Context) JSON(output []byte) {
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.Write(output)
}

// --------------------------------------------------

// CHandler Handles injecting the initial request context before passing handling on to the Middleware struct
type CHandler func(*Context)

// ServeHTTP is the first handler that gets hit when a request comes in.
func (ch CHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Items:    make(map[string]string),
		Enliven:  enliven,
		Response: rw,
		Request:  r,
	}
	ch(ctx)
}

// ContextHandler sets up serving the first request, and the handing off of subsequent requests to the Middleware struct
func ContextHandler(h Middleware) CHandler {
	return CHandler(func(ctx *Context) {
		h.handler.ServeHTTP(ctx, h.next.ServeHTTP)
	})
}

// --------------------------------------------------

// NextHandlerFunc allow use of ordinary functions middleware handlers
// Copied w/ alterations from github.com/codegangsta/negroni
type NextHandlerFunc func(*Context)

// Copied w/ alterations from github.com/codegangsta/negroni
func (nh NextHandlerFunc) ServeHTTP(ctx *Context) {
	nh(ctx)
}

// --------------------------------------------------

// HandlerFunc allow use of ordinary functions middleware handlers
// Copied w/ alterations from github.com/codegangsta/negroni
type HandlerFunc func(*Context, NextHandlerFunc)

// Copied w/ alterations from github.com/codegangsta/negroni
func (h HandlerFunc) ServeHTTP(ctx *Context, next NextHandlerFunc) {
	h(ctx, next)
}

// --------------------------------------------------

// RouteHandlerFunc is an interface to be used when writing route handler functions
type RouteHandlerFunc func(*Context)

// --------------------------------------------------

// Middleware Represents a piece of middlewear
// Copied w/ alterations from github.com/codegangsta/negroni
type Middleware struct {
	handler IMiddlewareHandler
	next    *Middleware
}

// Copied w/ alterations from github.com/codegangsta/negroni
func (m Middleware) ServeHTTP(ctx *Context) {
	m.handler.ServeHTTP(ctx, m.next.ServeHTTP)
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
func New(config Config) *Enliven {
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
func (ev *Enliven) registerConfig(suppliedConfig Config) {
	var enlivenConfig = Config{
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
func (ev *Enliven) GetConfig() Config {
	config := ev.GetService("config").(Config)
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
func (ev *Enliven) AddMiddlewareFunc(handlerFunc func(*Context, NextHandlerFunc)) {
	ev.AddMiddleware(HandlerFunc(handlerFunc))
}

// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) buildMiddleware(handlers []IMiddlewareHandler) Middleware {
	var next Middleware

	voidMiddleware := Middleware{
		HandlerFunc(func(*Context, NextHandlerFunc) {}), &Middleware{},
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
func (ev *Enliven) AddRoute(path string, rhf func(*Context)) *mux.Route {
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
func routeHandlerFunc(ctx *Context, next NextHandlerFunc) {

	// Clean path to canonical form and redirect.
	if p := cleanPath(ctx.Request.URL.Path); p != ctx.Request.URL.Path {
		url := *ctx.Request.URL
		url.Path = p
		p = url.String()

		ctx.Response.Header().Set("Location", p)
		ctx.Response.WriteHeader(http.StatusMovedPermanently)
		return
	}

	if !ctx.Enliven.GetRouter().KeepContext {
		defer context.Clear(ctx.Request)
	}

	var match mux.RouteMatch
	var handler http.Handler
	if enliven.GetRouter().Match(ctx.Request, &match) {
		handler = match.Handler
	}
	if handler == nil {
		handler := http.NotFoundHandler()
		handler.ServeHTTP(ctx.Response, ctx.Request)
	} else {
		// We use the request path to look up our stored route handler if it exists
		if routeHandler, ok := ctx.Enliven.routeHandlers[ctx.Request.URL.Path]; ok {
			// Calling the route handler with Enliven passed in.
			routeHandler(ctx)
		} else {
			// Using handler request handling otherwise.
			handler.ServeHTTP(ctx.Response, ctx.Request)
		}
	}

	next(ctx)
}

// Run executes the Enliven http server
func (ev *Enliven) Run(port string) {
	// Adding our route handler as the last piece of middleware
	ev.AddMiddlewareFunc(routeHandlerFunc)

	fmt.Println("Server is listening on port " + port + ".")
	http.ListenAndServe(":"+port, ContextHandler(ev.middleware))
	fmt.Println("Server has shut down.")
}
