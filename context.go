package enliven

import (
	"html/template"
	"net/http"
)

// Context stores context variables and the session that will be passed to requests
type Context struct {
	Session  ISession
	Vars     map[string]string // This map specifically to hold route key value pairs.
	Strings  map[string]string
	Integers map[string]int
	Booleans map[string]bool
	Storage  map[string]interface{}
	Enliven  *Enliven
	Response http.ResponseWriter
	Request  *http.Request
}

// String sets up string headers and outputs a string response
func (ctx *Context) String(output string) {
	ctx.Response.Header().Set("Content-Type", "text/plain")
	ctx.Response.Write([]byte(output))
}

// HTML sets up HTML headers and outputs a string response
func (ctx *Context) HTML(output string) {
	ctx.Response.Header().Set("Content-Type", "text/html")
	ctx.Response.Write([]byte(output))
}

// AnonymousTemplate sets up HTML headers and outputs an html/template response
func (ctx *Context) AnonymousTemplate(tmpl *template.Template) {
	ctx.Response.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(ctx.Response, ctx)
	if err != nil {
		ctx.String(err.Error())
	}
}

// Template sets up HTML headers and outputs an html/template response for a specific template definition
func (ctx *Context) Template(templateName string) {
	ctx.Response.Header().Set("Content-Type", "text/html")
	err := ctx.Enliven.Core.Templates.ExecuteTemplate(ctx.Response, templateName, ctx)
	if err != nil {
		ctx.String(err.Error())
	}
}

// JSON sets up JSON headers and outputs a JSON response
// Expects to recieve the result of json marshalling ([]byte)
func (ctx *Context) JSON(output []byte) {
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.Write(output)
}

// Redirect is a shortcut for redirecting a browser to a new URL
func (ctx *Context) Redirect(location string, status ...int) {
	var statusCode int
	if len(status) > 0 {
		statusCode = status[0]
	} else {
		statusCode = 302
	}

	http.Redirect(ctx.Response, ctx.Request, location, statusCode)
}

// Forbidden returns a 403 status and the forbidden page.
func (ctx *Context) Forbidden() {
	ctx.Response.WriteHeader(http.StatusForbidden)
	ctx.Template("forbidden")
}

// NotFound returns a 404 status and the not-found page
func (ctx *Context) NotFound() {
	ctx.Response.WriteHeader(http.StatusNotFound)
	ctx.Template("notfound")
}

// BadRequest returns a 400 status and the bad-request page
func (ctx *Context) BadRequest() {
	ctx.Response.WriteHeader(http.StatusBadRequest)
	ctx.Template("badrequest")
}

// EmptyOK outputs a 200 status with nothing else
func (ctx *Context) EmptyOK() {
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.String("")
}

// --------------------------------------------------

// CHandler Handles injecting the initial request context before passing handling on to the Middleware struct
type CHandler func(*Context)

// ServeHTTP is the first handler that gets hit when a request comes in.
func (ch CHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Vars:     make(map[string]string),
		Strings:  make(map[string]string),
		Integers: make(map[string]int),
		Booleans: make(map[string]bool),
		Storage:  make(map[string]interface{}),
		Enliven:  &enliven,
		Response: rw,
		Request:  r,
	}
	ch(ctx)
}

// ContextHandler sets up serving the first request, and the handing off of subsequent requests to the Middleware struct
func ContextHandler(h Middleware) CHandler {
	return CHandler(func(ctx *Context) {
		h.handler.ServeHTTP(ctx, h.next.ServeHTTP)
	})
}
