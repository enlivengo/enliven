package database

import (
	"strings"

	"github.com/enlivengo/enliven"
	"github.com/enlivengo/enliven/config"
	"github.com/jinzhu/gorm"

	// Adding DB requirements.
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// GetDatabase returns the requested database
func GetDatabase(ev *enliven.Enliven, namespace ...string) *gorm.DB {
	var name string
	if len(namespace) > 0 {
		name = namespace[0]
	} else {
		name = "default"
	}

	if db, ok := ev.GetService(name + "_database").(*gorm.DB); ok {
		return db
	}
	return nil
}

// NewApp creates and returns an instance of the database app
func NewApp(namespace ...string) *App {
	var name string
	if len(namespace) > 0 {
		name = namespace[0]
	} else {
		name = "default"
	}

	return &App{
		namespace: name,
	}
}

// App is the enliven application that sets up and manages the db
type App struct {
	namespace string
}

// Initialize sets up a database given the values from the supplied config
func (da *App) Initialize(ev *enliven.Enliven) {
	var namespace string
	if da.namespace != "default" {
		namespace = da.namespace + "_"
	}

	// Setting up the default config
	conf := config.Config{}
	conf[namespace+"database_driver"] = ""
	conf[namespace+"database_host"] = ""
	conf[namespace+"database_user"] = ""
	conf[namespace+"database_dbname"] = ""
	conf[namespace+"database_password"] = ""
	conf[namespace+"database_sslmode"] = "disable"
	conf[namespace+"database_port"] = ""
	conf[namespace+"database_connString"] = ""

	conf = config.UpdateConfig(config.MergeConfig(conf, config.GetConfig()))

	var driver string
	allowedDrivers := [4]string{"postgres", "mysql", "sqlite3", "mssql"}

	// Making sure the specified driver is in the list if allowed drivers
	for i := 0; i < 4; i++ {
		if allowedDrivers[i] == conf[namespace+"database_driver"] {
			driver = conf[namespace+"database_driver"]
			break
		}
	}

	// If we didn't set a driver, we return here.
	if driver == "" {
		return
	}

	var connString string

	// Someone can specify a whole connection string, or the parts of it
	if conf[namespace+"database_connString"] != "" {
		connString = conf[namespace+"database_connString"]
	} else {
		// driver specific connection string addons
		switch driver {

		case "sqlite3":
			// If the driver is sqlite3, but there wasn't a conn string, we return.
			if conf[namespace+"database_connString"] == "" {
				return
			}

		case "mysql", "mssql":
			connString = conf[namespace+"database_user"] + ":" + conf[namespace+"database_password"] + "@" + conf[namespace+"database_host"]

			// Adding a port if one was provided
			if len(conf[namespace+"database_port"]) > 0 {
				connString += ":" + conf[namespace+"database_port"]
			}

			connString += "/" + conf[namespace+"database_dbname"]

			if driver == "mysql" {
				connString += "?charset=utf8&parseTime=True&loc=Local"
			}

		case "postgres":
			var connStringParts []string
			connStringParts = append(connStringParts, "host="+conf[namespace+"database_host"])
			connStringParts = append(connStringParts, "user="+conf[namespace+"database_user"])
			connStringParts = append(connStringParts, "dbname="+conf[namespace+"database_dbname"])
			connStringParts = append(connStringParts, "sslmode="+conf[namespace+"database_sslmode"])
			connStringParts = append(connStringParts, "password="+conf[namespace+"database_password"])

			if len(conf[namespace+"database_port"]) > 0 {
				connStringParts = append(connStringParts, "port="+conf[namespace+"database_port"])
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

	ev.AddService(da.namespace+"_database", db)
}

// GetName gets the database app name
func (da *App) GetName() string {
	return da.namespace + "_database"
}
