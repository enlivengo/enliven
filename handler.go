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

// --------------------------------------------------

// RouteHandlerFunc is an interface to be used when writing route handler functions
type RouteHandlerFunc func(*Context)

// --------------------------------------------------

// PermissionHandler is a stub to stand in as the default permission checker.
// This should be overridden by the user app or something else if permissions checking is needed.
type PermissionHandler struct{}

// HasPermission is the default permission checker and always returns true
func (ph *PermissionHandler) HasPermission(permission string, ctx *Context) bool {
	return true
}

// AddPermission is the default permission adder, and does nothing.
func (ph *PermissionHandler) AddPermission(permission string, ev *Enliven) {}
