package enliven

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

// Initialize initializes this function middleware by doing nothing
func (h HandlerFunc) Initialize(ev *Enliven) {}

// GetName returns an empty string for this function middleware
func (h HandlerFunc) GetName() string {
	return ""
}

// --------------------------------------------------

// RouteHandlerFunc is an interface to be used when writing route handler functions
type RouteHandlerFunc func(*Context)

// --------------------------------------------------

// DefaultAuth is a simple implementation of IAuthorizer to stand in for auth checking/adding
// This should be overridden by the user app or something else if permissions checking is needed.
type DefaultAuth struct{}

// HasPermission is the default permission checker and always returns true
func (da *DefaultAuth) HasPermission(permission string, ctx *Context) bool {
	return true
}

// AddPermission is the default permission adder, and does nothing.
func (da *DefaultAuth) AddPermission(permission string, ev *Enliven, groups ...string) {}
