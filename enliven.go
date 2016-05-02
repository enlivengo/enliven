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

// EnlivenConfig holds config values
type EnlivenConfig struct {
	DatabaseDriver   string
	ConnectionString string
}

// New (constructor) gets a new instance of enliven.
func New(ec *EnlivenConfig) *Enliven {
	r := mux.NewRouter()

	enliven = Enliven{
		config:        ec,
		router:        r,
		routeHandlers: make(map[string]RouteHandlerFunc),
	}

	if len(ec.DatabaseDriver) > 0 && len(ec.ConnectionString) > 0 {
		enliven.InitDatabase()
	}

	return &enliven
}

// Enliven is....Enliven
type Enliven struct {
	config        *EnlivenConfig
	database      *gorm.DB
	router        *mux.Router
	routeHandlers map[string]RouteHandlerFunc
	middleware    Middleware
	handlers      []MiddlewareHandler
}

// InitDatabase Initializes a database given the values from the EnlivenConfig
func (ev *Enliven) InitDatabase() {
	db, err := gorm.Open(ev.config.DatabaseDriver, ev.config.ConnectionString)

	// Making sure we got a database instance
	if err != nil {
		panic(err)
	}

	// Making sure we can ping the database
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	ev.database = db
}

// Use adds a Handler onto the middleware stack.
// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) Use(handler MiddlewareHandler) {
	ev.handlers = append(ev.handlers, handler)
	ev.middleware = ev.buildMiddleware(ev.handlers)
}

// UseFunc adds a HandlerFunc onto the middleware stack.
// Copied w/ alterations from github.com/codegangsta/negroni
func (ev *Enliven) UseFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ev Enliven)) {
	ev.Use(HandlerFunc(handlerFunc))
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
	ev.Use(HandlerFunc(routeHandlerFunc))
	fmt.Println("Server is listening on port " + port + ".")
	http.ListenAndServe(":"+port, ev.middleware)
	fmt.Println("Server has shut down.")
}

func main() {
	ev := New(&EnlivenConfig{
		DatabaseDriver:   "",
		ConnectionString: "",
	})

	ev.AddRoute("/", RouteHandlerFunc(func(rw http.ResponseWriter, r *http.Request, ev Enliven) {
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte("It's working!!"))
	}))

	port := flag.String("port", "8000", "The port the server should listen on.")
	flag.Parse()

	ev.Run(*port)
}
