package main

import (
	"flag"
	"fmt"
	"net/http"
	"path"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// This is an accesible
var enliven Enliven

// MiddlewareHandler is an interface to be used when writing Middleware
// Copied w/ alterations from github.com/codegangsta/negroni
type MiddlewareHandler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ev Enliven)
}

// HandlerFunc allow use of ordinary functions middleware handlers
// Copied w/ alterations from github.com/codegangsta/negroni
type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ev Enliven)

// Copied w/ alterations from github.com/codegangsta/negroni
func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ev Enliven) {
	h(rw, r, next, enliven)
}

// RouteHandlerFunc is an interface to be used when writing route handler functions
type RouteHandlerFunc func(rw http.ResponseWriter, r *http.Request, ev Enliven)

// Middleware Represents a piece of middlewear
// Copied w/ alterations from github.com/codegangsta/negroni
type Middleware struct {
	handler MiddlewareHandler
	next    *Middleware
}

// Copied w/ alterations from github.com/codegangsta/negroni
func (m Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(rw, r, m.next.ServeHTTP, enliven)
}

// New (constructor) gets a new instance of enliven.
func New(config map[string]string) *Enliven {
	r := mux.NewRouter()

	enliven = Enliven{
		services:      make(map[string]interface{}),
		router:        r,
		routeHandlers: make(map[string]RouteHandlerFunc),
	}

	enliven.addConfig(config)
	enliven.addDatabase()

	return &enliven
}

// addConfig created and registers the app config
func (ev *Enliven) addConfig(suppliedConfig map[string]string) {
	var enlivenConfig = map[string]string{
		"db.driver":           "",
		"db.connectionString": "",
		"server.port":         "8000",
	}

	for key, value := range suppliedConfig {
		enlivenConfig[key] = value
	}

	ev.Register("config", enlivenConfig)
}

// addDatabase Initializes a database given the values from the EnlivenConfig
func (ev *Enliven) addDatabase() {
	config := ev.GetConfig()

	if len(config["db.driver"]) == 0 || len(config["db.connectionString"]) == 0 {
		return
	}

	db, err := gorm.Open(config["db.driver"], config["db.connectionString"])

	// Making sure we got a database instance
	if err != nil {
		panic(err)
	}

	// Making sure we can ping the database
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	ev.Register("database", db)
}

// Enliven is....Enliven
type Enliven struct {
	services      map[string]interface{}
	router        *mux.Router
	routeHandlers map[string]RouteHandlerFunc
	middleware    Middleware
	handlers      []MiddlewareHandler
}

// Register registers an enliven service or dependency
func (ev *Enliven) Register(name string, service interface{}) {
	if _, ok := ev.services[name]; ok {
		panic("The service name you are attempting to register has already been registered.")
	}

	ev.services[name] = service
}

// Get returns an enliven service or dependency
func (ev *Enliven) Get(name string) interface{} {
	if _, ok := ev.services[name]; ok {
		return ev.services[name]
	}
	return nil
}

// GetDatabase Gets an instance of the database
func (ev *Enliven) GetDatabase() *gorm.DB {
	if db, ok := ev.Get("database").(*gorm.DB); ok {
		return db
	}
	return nil
}

// GetConfig Gets an instance of the database
func (ev *Enliven) GetConfig() map[string]string {
	config := ev.Get("config").(map[string]string)
	return config
}

// AddMiddlewareHandler adds a Handler onto the middleware stack.
// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) AddMiddlewareHandler(handler MiddlewareHandler) {
	ev.handlers = append(ev.handlers, handler)
	ev.middleware = ev.buildMiddleware(ev.handlers)
}

// AddMiddleware adds a HandlerFunc onto the middleware stack.
// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) AddMiddleware(handlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ev Enliven)) {
	ev.AddMiddlewareHandler(HandlerFunc(handlerFunc))
}

// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) buildMiddleware(handlers []MiddlewareHandler) Middleware {
	var next Middleware

	voidMiddleware := Middleware{
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ev Enliven) {}), &Middleware{},
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

// GetRouter returns our mux instance
func (ev *Enliven) GetRouter() *mux.Router {
	return ev.router
}

// AddRoute Registers a handler for a given route.
// We register a dummy route with mux, and then store the provided handler
// which we'll use later in order to inject dependencies into the handler func.
func (ev *Enliven) AddRoute(path string, rhf RouteHandlerFunc) *mux.Route {
	ev.routeHandlers[path] = rhf
	return ev.router.HandleFunc(path, func(http.ResponseWriter, *http.Request) {})
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
func routeHandlerFunc(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ev Enliven) {

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
	if enliven.router.Match(r, &match) {
		handler = match.Handler
	}
	if handler == nil {
		handler := http.NotFoundHandler()
		handler.ServeHTTP(rw, r)
	} else {
		// We use the request path to look up our stored route handler
		routeHandler := ev.routeHandlers[r.URL.Path]
		// Calling the route handler with Enliven passed in.
		routeHandler(rw, r, ev)
	}

	next(rw, r)
}

// Run executes the Enliven http server
func (ev *Enliven) Run(port string) {
	// Adding our route handler as the last piece of middleware
	ev.AddMiddlewareHandler(HandlerFunc(routeHandlerFunc))

	fmt.Println("Server is listening on port " + port + ".")
	http.ListenAndServe(":"+port, ev.middleware)
	fmt.Println("Server has shut down.")
}

// Example/Test usage
func main() {
	ev := New(make(map[string]string))

	ev.AddRoute("/", RouteHandlerFunc(func(rw http.ResponseWriter, r *http.Request, ev Enliven) {
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte("It's working!!"))
	}))

	port := flag.String("port", "8000", "The port the server should listen on.")
	flag.Parse()

	ev.Run(*port)
}
