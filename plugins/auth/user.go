package auth

import (
	"strconv"

	"github.com/hickeroar/enliven"
	"github.com/jinzhu/gorm"
)

// User describes the user database structure.
type User struct {
	gorm.Model

	DisplayName string
	Login       string `gorm:"type:varchar(100);unique_index"`
	Password    string
}

// NewUserPlugin generates and returns an instance of the UserPlugin
func NewUserPlugin(suppliedConfig enliven.Config) *UserPlugin {
	var config = enliven.Config{
		"user.login.route":  "/user/login",
		"user.logout.route": "/user/logout",
	}

	config = enliven.MergeConfig(config, suppliedConfig)

	return &UserPlugin{
		loginRoute:  config["user.login.route"],
		logoutRoute: config["user.logout.route"],
	}
}

// UserPlugin handles adding a route handler for static assets
type UserPlugin struct {
	loginRoute  string
	logoutRoute string
}

// Initialize sets up our plugin to handle embedded static asset requests
func (sap *UserPlugin) Initialize(ev *enliven.Enliven) {
	ev.GetDatabase().AutoMigrate(&User{})

	ev.AddMiddlewareFunc(UserSessionMiddleware)
}

// GetName returns the plugin's name
func (sap *UserPlugin) GetName() string {
	return "user"
}

// UserSessionMiddleware handles adding the elements to the context that carry the user's id and status
func UserSessionMiddleware(ctx *enliven.Context, next enliven.NextHandlerFunc) {
	if ctx.Session == nil {
		panic("The User plugin requires Session middleware to be registered.")
	}

	userID := ctx.Session.Get("UserPlugin_LoggedInUserID")

	// If there isn't a user id in the session, we set context items accordingly
	if userID == "" {
		ctx.Items["UserLoggedIn"] = "0"
		ctx.Items["UserID"] = "0"
	} else {
		ctx.Items["UserLoggedIn"] = "1"
		ctx.Items["UserID"] = userID
	}

	next(ctx)
}

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
	ctx.Enliven.GetDatabase().First(&user, dbUserID)

	// Caching the user lookup for later.
	ctx.Storage["User"] = &user

	return &user
}
