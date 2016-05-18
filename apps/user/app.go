package user

//go:generate go-bindata -o templates.go -pkg user templates/...

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/hickeroar/enliven"
	"github.com/hickeroar/enliven/apps/database"
	"github.com/jinzhu/gorm"
)

// GetUser returns an instance of the User model
func GetUser(ctx *enliven.Context) *User {
	// If they're not logged in, return 0
	if ctx.Items["UserLoggedIn"] == "0" {
		return nil
	}

	// The user may be cached from an earlier lookup.
	if user, ok := ctx.Storage["User"]; ok {
		return user.(*User)
	}

	var user User
	dbUserID, _ := strconv.Atoi(ctx.Items["UserID"])
	database.GetDatabase(ctx, ctx.Enliven.GetConfig()["user.database.namespace"]).First(&user, dbUserID)

	// Caching the user lookup for later.
	ctx.Storage["User"] = &user

	return &user
}

// GeneratePasswordHash produces a bcrypt hash and returns it
func GeneratePasswordHash(password string, cost ...int) string {
	var bcryptCost int
	if len(cost) > 0 {
		bcryptCost = cost[0]
	} else {
		bcryptCost = 12
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(hash[:])
}

// VerifyPasswordHash checks a password for validity
func VerifyPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return (err == nil)
}

// User describes the user database structure.
type User struct {
	gorm.Model

	DisplayName      string
	Age              int
	Login            string `gorm:"type:varchar(100);unique_index;"`
	Email            string `gorm:"type:varchar(100);unique_index;"`
	Password         string `gorm:"type:varchar(100);"`
	VerificationCode string `gorm:"type:varchar(100);unique_index;"`
	Group            Group  `gorm:"not null"`
	Superuser        bool
}

// HasPermission checks if a user has a specific permission
func (u *User) HasPermission(name string) bool {
	if u.Superuser {
		return true
	}

	var groupStack []string
	return u.hasPermission(name, &u.Group, groupStack)
}

// hasPermission recursively looks through a group's inheritance chain to look for a permission
func (u *User) hasPermission(name string, group *Group, groupStack []string) bool {
	// Checking if this group has a permission matching the one we're looking for
	for _, perm := range group.Permisions {
		if name == perm.Name {
			return true
		}
	}

	// If this group inherits from another groups, we check that group for a permission
	if group.Inherits != nil && group.Name != group.Inherits.Name {

		// We're avoiding infinite group loops here by making sure we don't check a group more than once.
		for _, stackGroup := range groupStack {
			if stackGroup == group.Inherits.Name {
				return false
			}
		}
		groupStack = append(groupStack, group.Inherits.Name)

		// Recursively checking the group's inheritance chain.
		return u.hasPermission(name, group.Inherits, groupStack)
	}

	return false
}

// Group describes the user group database structure
type Group struct {
	gorm.Model

	Name       string `gorm:"not null;unique;"`
	Inherits   *Group
	Permisions []Permission
}

// Permission describes a permission that can be linked to many groups
type Permission struct {
	gorm.Model

	Name string `gorm:"not_null;unique;"`
}

// NewApp generates and returns an instance of the app
func NewApp() *App {
	return &App{}
}

// App handles adding a route handler for static assets
type App struct{}

// Initialize sets up our app to handle embedded static asset requests
func (ua *App) Initialize(ev *enliven.Enliven) {
	var config = enliven.Config{
		"user.login.route":    "/user/login/",
		"user.logout.route":   "/user/logout/",
		"user.register.route": "/user/register/",
		"user.verify.route":   "/user/verify/",
		"user.password.route": "/user/password/",

		// Full text template text
		"user.login.template":    "",
		"user.logout.template":   "",
		"user.register.template": "",
		"user.verify.template":   "",
		"user.password.template": "",

		// Where the user will be redirected after these successful actions.
		"user.login.redirect":    "/",
		"user.logout.redirect":   "/",
		"user.register.redirect": "/",
		"user.password.redirect": "/",
		"user.verify.redirect":   "/",

		"user.database.namespace": "default",
	}

	config = enliven.MergeConfig(config, ev.GetConfig())
	ev.AppendConfig(config)

	db := database.GetDatabase(&enliven.Context{Enliven: ev}, config["user.database.namespace"])

	if db == nil {
		panic("The User app is unable to locate the '" + config["user.database.namespace"] + "' database. A valid database is required.")
	}

	// Migrating the user tables
	db.AutoMigrate(&User{}, &Group{}, &Permission{})

	// Routing setup
	ev.AddRoute(config["user.login.route"], LoginGetHandler, "GET")
	ev.AddRoute(config["user.login.route"], LoginPostHandler, "POST")

	// Handles the setup of context variables to support user session management
	ev.AddMiddlewareFunc(SessionMiddleware)
}

// GetName returns the app's name
func (ua *App) GetName() string {
	return "user"
}
