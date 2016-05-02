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

// EnlivenConfig holds config values
type EnlivenConfig struct {
	DatabaseDriver   string
	ConnectionString string
}

// New (constructor) gets a new instance of enliven.
func New(ec *EnlivenConfig) *Enliven {
	ng := negroni.New()
	r := mux.NewRouter()

	enliven = Enliven{
		config:     ec,
		middleware: ng,
		router:     r,
	}

	if len(ec.DatabaseDriver) > 0 {
		enliven.InitDatabase()
	}

	return &enliven
}

// Enliven is....Enliven
type Enliven struct {
	config     *EnlivenConfig
	database   *gorm.DB
	middleware *negroni.Negroni
	router     *mux.Router
}

// InitDatabase Initializes a database given the values from the EnlivenConfig
func (e *Enliven) InitDatabase() {
	db, err := gorm.Open(e.config.DatabaseDriver, e.config.ConnectionString)

	// Making sure we got a database instance
	if err != nil {
		panic(err)
	}

	// Making sure we can ping the database
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	e.database = db
}

// AddMiddleware Adds a piece of EnlivenMiddleware to the handler map and negroni
func (e *Enliven) AddMiddleware(middlewareFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	// Storing the item in our map and then adding its middleware func to negroni
	e.middleware.UseFunc(middlewareFunc)
}

// GetRouter returns our mux instance
func (e *Enliven) GetRouter() *mux.Router {
	return e.router
}

// Run executes the Enliven http server
func (e *Enliven) Run(addr string) {
	e.middleware.UseHandler(e.router)
	fmt.Println("Server is listening on " + addr + ".")
	http.ListenAndServe(addr, e.middleware)
}

// GetEnlivenDatabase returns the enliven database instance.
func GetEnlivenDatabase() *gorm.DB {
	return enliven.database
}

func main() {
	en := New(&EnlivenConfig{
		DatabaseDriver:   "",
		ConnectionString: "",
	})

	en.GetRouter().HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("It's working!"))
	}).Methods("POST")

	port := flag.String("port", "8000", "The port the server should listen on.")
	flag.Parse()

	en.Run(":" + *port)
}
