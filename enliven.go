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
)

// We'll use this to create and insert the initial enliven instance into the handlers
var enliven Enliven

// Enliven is....Enliven
type Enliven struct {
	services      map[string]interface{}
	routeHandlers map[string]map[string]RouteHandlerFunc
	middleware    Middleware
	handlers      []IMiddlewareHandler
	apps          []string
	permissions   IPermissionChecker
}

// New (constructor) gets a new instance of enliven.
func New(config Config) *Enliven {
	enliven = Enliven{
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

	enliven.SetPermissionChecker(&PermissionHandler{})
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

// This registers the default templates (header, footer, home, forbidden, and notfound)
// It's expected that developers will override at least the header, footer, and home templates.
func (ev *Enliven) registerTemplates() {
	headerTemplate, _ := templates.Asset("templates/header.html")
	footerTemplate, _ := templates.Asset("templates/footer.html")
	homeTemplate, _ := templates.Asset("templates/home.html")
	forbiddenTemplate, _ := templates.Asset("templates/forbidden.html")
	notfoundTemplate, _ := templates.Asset("templates/notfound.html")

	templates := template.New("enliven")
	templates.Parse(string(headerTemplate[:]))
	templates.Parse(string(footerTemplate[:]))
	templates.Parse(string(homeTemplate[:]))
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

// SetPermissionChecker sets the enliven permission checker.
func (ev *Enliven) SetPermissionChecker(checker IPermissionChecker) {
	ev.permissions = checker
}

// GetPermissionChecker returns the enliven permission checker
func (ev *Enliven) GetPermissionChecker() IPermissionChecker {
	return ev.permissions
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
			return ev.GetRouter().PathPrefix(prefix).HandlerFunc(func(http.ResponseWriter, *http.Request) {}).Methods(methods...)
		}
		return ev.GetRouter().HandleFunc(path, func(http.ResponseWriter, *http.Request) {}).Methods(methods...)
	}

	// We store a simple reference to their route handler without method expectations if none were provided
	ev.routeHandlers["ALL"][path] = RouteHandlerFunc(rhf)
	// Adding a dummy reference to a handler to mux which we'll override at execution-time, methods included
	if prefix != "" {
		return ev.GetRouter().PathPrefix(prefix).HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	}
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
		ctx.NotFound()
	} else {
		// We use the request path to look up our stored route handler if it exists
		if routeHandler, ok := ctx.Enliven.routeHandlers[strings.ToUpper(ctx.Request.Method)][ctx.Request.URL.Path]; ok {
			// Calling the route handle specific to a certain method if we stored one
			routeHandler(ctx)
		} else if routeHandler, ok := ctx.Enliven.routeHandlers["ALL"][ctx.Request.URL.Path]; ok {
			// Calling the route handler that handles all routes if we stored one
			routeHandler(ctx)
		} else if routeHandler, ok := ctx.Enliven.routePrefix(ctx, strings.ToUpper(ctx.Request.Method)); ok {
			// Per-method routing for path prefixes
			routeHandler(ctx)
		} else if routeHandler, ok := ctx.Enliven.routePrefix(ctx, "ALL"); ok {
			// All-method routing for path prefixes
			routeHandler(ctx)
		} else {
			// We didn't have a stored handler for this path/handler, so we execute the handler.
			handler.ServeHTTP(ctx.Response, ctx.Request)
		}
	}

	next(ctx)
}

func (ev *Enliven) routePrefix(ctx *Context, method string) (RouteHandlerFunc, bool) {
	var prefix string

	for path, hf := range ctx.Enliven.routeHandlers[method] {
		// Looking for paths that have ... on the end
		if len(path) < 3 || len(ctx.Request.URL.Path) < (len(path)-3) || string(path[(len(path)-3):]) != "..." {
			continue
		}
		prefix = string(path[:(len(path) - 3)])
		// If the request path prefix is the same as the path (minus the ...)
		if string(ctx.Request.URL.Path[:len(prefix)]) == prefix {
			return hf, true
		}
	}

	return nil, false
}

// Run executes the Enliven http server
func (ev *Enliven) Run(port string) {
	// Adding our route handler as the last piece of middleware
	ev.AddMiddlewareFunc(routeHandlerFunc)

	fmt.Println("Enliven server is listening on port " + port + ".")
	http.ListenAndServe(":"+port, ContextHandler(ev.middleware))
	fmt.Println("Enliven server has shut down.")
}
