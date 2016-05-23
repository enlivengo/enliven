package core

import (
	"html/template"

	"github.com/hickeroar/enliven/core/email"
	"github.com/hickeroar/enliven/core/templates"
	"github.com/hickeroar/enliven/core/util"
)

// Core holds core functionality for enliven that exists outside the enliven namespace
type Core struct {
	Email     email.Core
	Templates *template.Template
	Util      util.Core
}

// NewCore creates a new core struct instance for use in the enliven application
func NewCore() Core {
	return Core{
		Email:     email.Core{},
		Templates: templates.New(),
		Util:      util.Core{},
	}
}
