package enliven

//go:generate go-bindata -o templates/templates.go -pkg templates templates/...

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/hickeroar/enliven/templates"

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

// IApp is an interface for writing Enliven apps
// Apps are basically packaged code to extend Enliven's functionality
type IApp interface {
	Initialize(*Enliven)
	GetName() string
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
	Strings  map[string]string
	Integers map[string]int
	Booleans map[string]bool
	Storage  map[string]interface{}
	Enliven  *Enliven
	Response http.ResponseWriter
	Request  *http.Request
}

// routeRestricted checks if a given route is restricted to superusers
func (ctx *Context) routeRestricted() bool {
	for path, length := range ctx.Enliven.restrictedRoutes {
		// Checking if the current route is in the list of restricted one
		if len(ctx.Request.URL.Path) >= length && string(ctx.Request.URL.Path[:length]) == path {
			// If the user app isnt installed, or they're not a super user, this route is restricted
			if !ctx.Enliven.AppInstalled("user") || !ctx.Booleans["UserSuperUser"] {
				return true
			}
			return false
		}
	}
	return false
}

// String sets up string headers and outputs a string response
func (ctx *Context) String(output string) {
	ctx.Response.Header().Set("Content-Type", "text/plain")
	ctx.Response.Write([]byte(output))
}

// HTML sets up HTML headers and outputs a string response
func (ctx *Context) HTML(output string) {
	ctx.Response.Header().Set("Content-Type", "text/html")
	ctx.Response.Write([]byte(output))
}

// Template sets up HTML headers and outputs an html/template response
func (ctx *Context) Template(tmpl *template.Template) {
	ctx.Response.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(ctx.Response, ctx)
	if err != nil {
		ctx.String(err.Error())
	}
}

// NamedTemplate sets up HTML headers and outputs an html/template response for a specific template definition
func (ctx *Context) NamedTemplate(tmpl *template.Template, templateName string) {
	ctx.Response.Header().Set("Content-Type", "text/html")
	err := tmpl.ExecuteTemplate(ctx.Response, templateName, ctx)
	if err != nil {
		ctx.String(err.Error())
	}
}

// JSON sets up JSON headers and outputs a JSON response
// Expects to recieve the result of json marshalling ([]byte)
func (ctx *Context) JSON(output []byte) {
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.Write(output)
}

// Redirect is a shortcut for redirecting a browser to a new URL
func (ctx *Context) Redirect(location string, status ...int) {
	var statusCode int
	if len(status) > 0 {
		statusCode = status[0]
	} else {
		statusCode = 302
	}

	http.Redirect(ctx.Response, ctx.Request, location, statusCode)
}

// Forbidden returns a 403 status and the forbidden page.
func (ctx *Context) Forbidden() {
	ctx.Response.WriteHeader(http.StatusForbidden)
	tmpl := ctx.Enliven.GetTemplates()
	ctx.NamedTemplate(tmpl, "forbidden")
}

// NotFound returns a 404 status and the not-found page
func (ctx *Context) NotFound() {
	ctx.Response.WriteHeader(http.StatusNotFound)
	tmpl := ctx.Enliven.GetTemplates()
	ctx.NamedTemplate(tmpl, "notfound")
}

// --------------------------------------------------

// CHandler Handles injecting the initial request context before passing handling on to the Middleware struct
type CHandler func(*Context)

// ServeHTTP is the first handler that gets hit when a request comes in.
func (ch CHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Strings:  make(map[string]string),
		Integers: make(map[string]int),
		Booleans: make(map[string]bool),
		Storage:  make(map[string]interface{}),
		Enliven:  &enliven,
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
	services         map[string]interface{}
	restrictedRoutes map[string]int
	routeHandlers    map[string]map[string]RouteHandlerFunc
	middleware       Middleware
	handlers         []IMiddlewareHandler
	apps             []string
}

// New (constructor) gets a new instance of enliven.
func New(config Config) *Enliven {
	enliven = Enliven{
		services:         make(map[string]interface{}),
		restrictedRoutes: make(map[string]int),
		routeHandlers: map[string]map[string]RouteHandlerFunc{
			"ALL":    make(map[string]RouteHandlerFunc),
			"GET":    make(map[string]RouteHandlerFunc),
			"DELETE": make(map[string]RouteHandlerFunc),
			"PATCH":  make(map[string]RouteHandlerFunc),
			"POST":   make(map[string]RouteHandlerFunc),
			"PUT":    make(map[string]RouteHandlerFunc),
		},
	}

	enliven.RegisterService("router", mux.NewRouter())
	enliven.registerConfig(config)
	enliven.registerTemplates()

	return &enliven
}

// addConfig created and registers the app config
func (ev *Enliven) registerConfig(suppliedConfig Config) {
	var enlivenConfig = Config{
		"server_port": "8000",
	}

	ev.RegisterService("config", MergeConfig(enlivenConfig, suppliedConfig))
}

func (ev *Enliven) registerTemplates() {
	headerTemplate, _ := templates.Asset("templates/header.html")
	footerTemplate, _ := templates.Asset("templates/footer.html")
	forbiddenTemplate, _ := templates.Asset("templates/forbidden.html")
	notfoundTemplate, _ := templates.Asset("templates/notfound.html")

	templates, _ := template.New("enliven").Parse(string(headerTemplate[:]))
	templates.Parse(string(footerTemplate[:]))
	templates.Parse(string(forbiddenTemplate[:]))
	templates.Parse(string(notfoundTemplate[:]))

	ev.RegisterService("templates", templates)
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

// AddApp initializes a provided enliven application
func (ev *Enliven) AddApp(app IApp) {
	if ev.AppInstalled(app.GetName()) {
		panic("The '" + app.GetName() + "' app has already been added.")
	}

	app.Initialize(ev)
	ev.apps = append(ev.apps, app.GetName())
}

// AppInstalled returns true if a given app has already been installed
func (ev *Enliven) AppInstalled(name string) bool {
	for _, value := range ev.apps {
		if name == value {
			return true
		}
	}

	return false
}

// AppendConfig merges and adds config to the enliven config
func (ev *Enliven) AppendConfig(suppliedConfig Config) {
	ev.services["config"] = MergeConfig(ev.GetConfig(), suppliedConfig)
}

// GetConfig Gets an instance of the config
func (ev *Enliven) GetConfig() Config {
	config := ev.GetService("config").(Config)
	return config
}

// GetTemplates gets the template instance so devs can add to and use it
func (ev *Enliven) GetTemplates() *template.Template {
	templates := ev.GetService("templates").(*template.Template)
	return templates
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
func (ev *Enliven) AddRoute(path string, rhf func(*Context), methods ...string) *mux.Route {
	if len(methods) > 0 {
		// We store a reference to their handler for each of the methods if they passed some in
		for _, method := range methods {
			// If they provided a legit method, we silo this handler into that method
			if _, ok := ev.routeHandlers[strings.ToUpper(method)]; ok {
				ev.routeHandlers[strings.ToUpper(method)][path] = RouteHandlerFunc(rhf)
			}
		}
		// Adding a dummy reference to a handler to mux which we'll override at execution-time, methods included
		return ev.GetRouter().HandleFunc(path, func(http.ResponseWriter, *http.Request) {}).Methods(methods...)
	}

	// We store a simple reference to their route handler without method expectations if none were provided
	ev.routeHandlers["ALL"][path] = RouteHandlerFunc(rhf)
	// Adding a dummy reference to a handler to mux which we'll override at execution-time, methods included
	return ev.GetRouter().HandleFunc(path, func(http.ResponseWriter, *http.Request) {})
}

// RestrictRoute restricts a specified route to superusers only
func (ev *Enliven) RestrictRoute(path string) {
	ev.restrictedRoutes[path] = len(path)
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
		ctx.NotFound()
	} else if ctx.routeRestricted() {
		ctx.Forbidden()
	} else {
		// We use the request path to look up our stored route handler if it exists
		if routeHandler, ok := ctx.Enliven.routeHandlers[strings.ToUpper(ctx.Request.Method)][ctx.Request.URL.Path]; ok {
			// Calling the route handle specific to a certain method if we stored one
			routeHandler(ctx)
		} else if routeHandler, ok := ctx.Enliven.routeHandlers["ALL"][ctx.Request.URL.Path]; ok {
			// Calling the route handler that handles all routes if we stored one
			routeHandler(ctx)
		} else {
			// We didn't have a stored handler for this path/handler, so we execute the handler.
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
