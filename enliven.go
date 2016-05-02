package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// This is an accesible
var enliven Enliven

// EnlivenMiddleware is the interface making for any middleware added to enliven
type EnlivenMiddleware interface {
	GetName() string
	GetDependencies() []EnlivenMiddleware
	Handler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// EnlivenConfig holds config values
type EnlivenConfig struct {
	DatabaseDriver   string
	ConnectionString string
}

// New (constructor) gets a new instance of enliven.
func New(ec *EnlivenConfig) *Enliven {
	ng := negroni.New()

	enliven = Enliven{
		handlers:   make(map[string]EnlivenMiddleware),
		middleware: ng,
		config:     ec,
	}

	enliven.InitDatabase()

	return &enliven
}

// Enliven is....Enliven
type Enliven struct {
	handlers   map[string]EnlivenMiddleware
	middleware *negroni.Negroni
	database   *gorm.DB
	config     *EnlivenConfig
}

// InitDatabase Initializes a database given the values from the EnlivenConfig
func (e *Enliven) InitDatabase() {
	db, err := gorm.Open(e.config.DatabaseDriver, e.config.ConnectionString)

	if err != nil {
		panic(err)
	}

	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	e.database = db
}

// AddMiddleware Adds a piece of EnlivenMiddleware to the handler map and negroni
func (e *Enliven) AddMiddleware(em EnlivenMiddleware) {
	name := em.GetName()

	// If this middleware has been added already, we return
	if _, ok := e.handlers[name]; ok {
		return
	}

	// Adding all dependency middleware via recursion
	for _, dep := range em.GetDependencies() {
		e.AddMiddleware(dep)
	}

	// Storing the item in our map and then adding its middleware func to negroni
	e.handlers[name] = em
	e.middleware.UseFunc(e.handlers[name].Handler)
}

func main() {
	r := mux.NewRouter()

	port := flag.String("port", "8000", "The port the server should listen on.")
	flag.Parse()

	fmt.Println("Server is listening on port " + *port + ".")

	http.ListenAndServe(":"+*port, r)
}
