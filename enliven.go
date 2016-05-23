package enliven

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/enlivengo/enliven/config"
	"github.com/enlivengo/enliven/core"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// Our instance of enliven that will be set up in request contexts
var enliven Enliven

// Enliven is....Enliven
type Enliven struct {
	Auth   IAuthorizer
	Core   core.Core
	Router *mux.Router

	services      map[string]interface{}
	routeHandlers map[string]map[string]RouteHandlerFunc
	middleware    Middleware
	handlers      []IMiddlewareHandler

	// Supports the AppInstalled and MiddlewareInstalled boolean methods
	installedApps       []string
	installedMiddleware []string
}

// New gets a new instance of enliven.
func New(conf config.Config) *Enliven {
	config.CreateConfig(config.MergeConfig(DefaultEnlivenConfig, conf))

	enliven = Enliven{
		Auth:   &DefaultAuth{},
		Core:   core.NewCore(),
		Router: mux.NewRouter(),

		services: make(map[string]interface{}),
		routeHandlers: map[string]map[string]RouteHandlerFunc{
			"ALL":    make(map[string]RouteHandlerFunc),
			"GET":    make(map[string]RouteHandlerFunc),
			"DELETE": make(map[string]RouteHandlerFunc),
			"PATCH":  make(map[string]RouteHandlerFunc),
			"POST":   make(map[string]RouteHandlerFunc),
			"PUT":    make(map[string]RouteHandlerFunc),
		},
	}

	return &enliven
}

// AddService registers an enliven service or dependency
func (ev *Enliven) AddService(name string, service interface{}) {
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
	ev.installedApps = append(ev.installedApps, app.GetName())
}

// AppInstalled returns true if a given app has already been installed
func (ev *Enliven) AppInstalled(name string) bool {
	for _, value := range ev.installedApps {
		if name == value {
			return true
		}
	}
	return false
}

// AddMiddleware adds a Handler onto the middleware stack.
// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) AddMiddleware(handler IMiddlewareHandler) {
	// We track middleware names that are not empty
	if handler.GetName() != "" && ev.MiddlewareInstalled(handler.GetName()) {
		panic("The '" + handler.GetName() + "' middleware has already been added.")
	} else if handler.GetName() != "" {
		ev.installedMiddleware = append(ev.installedApps, handler.GetName())
	}

	handler.Initialize(ev)
	ev.handlers = append(ev.handlers, handler)
	ev.middleware = ev.buildMiddleware(ev.handlers)
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

// AddMiddlewareFunc adds a HandlerFunc onto the middleware stack.
// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) AddMiddlewareFunc(handlerFunc func(*Context, NextHandlerFunc)) {
	ev.AddMiddleware(HandlerFunc(handlerFunc))
}

// MiddlewareInstalled returns true if a given middleware has already been installed
func (ev *Enliven) MiddlewareInstalled(name string) bool {
	for _, value := range ev.installedMiddleware {
		if name == value {
			return true
		}
	}
	return false
}

// AddRoute Registers a handler for a given route.
// We register a dummy route with mux, and then store the provided handler
// which we'll use later in order to inject dependencies into the handler func.
func (ev *Enliven) AddRoute(path string, rhf func(*Context), methods ...string) *mux.Route {
	var prefix string
	if len(path) > 3 {
		if string(path[(len(path)-3):]) == "..." {
			prefix = string(path[:(len(path) - 3)])
		}
	}

	if len(methods) > 0 {
		// We store a reference to their handler for each of the methods if they passed some in
		for _, method := range methods {
			// If they provided a legit method, we silo this handler into that method
			if _, ok := ev.routeHandlers[strings.ToUpper(method)]; ok {
				ev.routeHandlers[strings.ToUpper(method)][path] = RouteHandlerFunc(rhf)
			}
		}
		// Adding a dummy reference to a handler to mux which we'll override at execution-time, methods included
		if prefix != "" {
			return ev.Router.PathPrefix(prefix).HandlerFunc(func(http.ResponseWriter, *http.Request) {}).Methods(methods...)
		}
		return ev.Router.HandleFunc(path, func(http.ResponseWriter, *http.Request) {}).Methods(methods...)
	}

	// We store a simple reference to their route handler without method expectations if none were provided
	ev.routeHandlers["ALL"][path] = RouteHandlerFunc(rhf)
	// Adding a dummy reference to a handler to mux which we'll override at execution-time, methods included
	if prefix != "" {
		return ev.Router.PathPrefix(prefix).HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	}
	return ev.Router.HandleFunc(path, func(http.ResponseWriter, *http.Request) {})
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

	if !ctx.Enliven.Router.KeepContext {
		defer context.Clear(ctx.Request)
	}

	var match mux.RouteMatch
	var handler http.Handler
	if enliven.Router.Match(ctx.Request, &match) {
		handler = match.Handler
		ctx.Vars = match.Vars
	}

	if handler == nil {
		ctx.NotFound()
	} else {
		// The routing path match.
		urlPath, _ := match.Route.GetPathTemplate()

		// We use the request path to look up our stored route handler if it exists
		if routeHandler, ok := ctx.Enliven.routeHandlers[strings.ToUpper(ctx.Request.Method)][urlPath]; ok {
			// Calling the route handle specific to a certain method if we stored one
			routeHandler(ctx)
		} else if routeHandler, ok := ctx.Enliven.routeHandlers["ALL"][urlPath]; ok {
			// Calling the route handler that handles all routes if we stored one
			routeHandler(ctx)
		} else if routeHandler, ok := ctx.Enliven.routePrefix(ctx, urlPath, strings.ToUpper(ctx.Request.Method)); ok {
			// Per-method routing for path prefixes
			routeHandler(ctx)
		} else if routeHandler, ok := ctx.Enliven.routePrefix(ctx, urlPath, "ALL"); ok {
			// All-method routing for path prefixes
			routeHandler(ctx)
		} else {
			// We didn't have a stored handler for this path/handler, so we execute the handler.
			handler.ServeHTTP(ctx.Response, ctx.Request)
		}
	}

	next(ctx)
}

func (ev *Enliven) routePrefix(ctx *Context, urlPath string, method string) (RouteHandlerFunc, bool) {
	var prefix string

	for path, hf := range ctx.Enliven.routeHandlers[method] {
		// Looking for paths that have ... on the end
		// Example: Admin app, route.go, MountTo method
		if len(path) < 3 || len(urlPath) < (len(path)-3) || string(path[(len(path)-3):]) != "..." {
			continue
		}
		prefix = string(path[:(len(path) - 3)])
		// If the request path prefix is the same as the path (minus the ...)
		if string(urlPath[:len(prefix)]) == prefix {
			return hf, true
		}
	}

	return nil, false
}

// Run executes the Enliven http server
func (ev *Enliven) Run() {
	// Adding our route handler as the last piece of middleware
	ev.AddMiddlewareFunc(routeHandlerFunc)

	address := config.GetConfig()["server_address"]

	fmt.Println("Enliven server is listening on " + address + ".")
	http.ListenAndServe(address, ContextHandler(ev.middleware))
	fmt.Println("Enliven server has shut down.")
}
