package enliven

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
	Initialize(*Enliven)
	GetName() string
	ServeHTTP(*Context, NextHandlerFunc)
}

// IAuthorizer is an interface to be used when writing a struct for checking a permission
type IAuthorizer interface {
	HasPermission(string, *Context) bool
	// Name of the new permission, the enliven instance, groups that we want to add permission to
	AddPermission(string, *Enliven, ...string)
}
