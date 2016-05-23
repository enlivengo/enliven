package user

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/enlivengo/admin"
	"github.com/enlivengo/enliven"
	"github.com/enlivengo/enliven/apps/database"
	"github.com/enlivengo/enliven/config"
	"github.com/jinzhu/gorm"
)

// GetUser returns an instance of the User model
func GetUser(ctx *enliven.Context) *User {
	// If they're not logged in or this app isn't installed, return 0
	if !ctx.Enliven.AppInstalled("user") || ctx.Booleans["UserLoggedIn"] == false {
		return nil
	}

	// The user may be cached from an earlier lookup.
	if user, ok := ctx.Storage["User"]; ok {
		return user.(*User)
	}

	var user User
	dbUserID, _ := ctx.Integers["UserID"]
	database.GetDatabase(ctx.Enliven).Preload("Groups").Preload("Groups.Permissions").First(&user, dbUserID)

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
func VerifyPasswordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return (err == nil)
}

// User describes the user database structure.
type User struct {
	gorm.Model

	DisplayName      string
	Username         string  `gorm:"type:varchar(100);unique_index;"`
	Email            string  `gorm:"type:varchar(100);unique_index;"`
	Password         string  `gorm:"type:varchar(100);"`
	VerificationCode string  `gorm:"type:varchar(64);index;"`
	Status           int     `gorm:"default:0;"`
	Groups           []Group `gorm:"many2many:user_group;"`
	Superuser        bool    `gorm:"index;"`
}

// UserHasPermission checks if a user has a specific permission
func (u *User) UserHasPermission(name string) bool {
	if u.Superuser {
		return true
	}

	for _, group := range u.Groups {
		for _, perm := range group.Permissions {
			if name == perm.Name {
				return true
			}
		}
	}

	return false
}

// GetDisplayName gets the proper display name for a user.
func (u *User) GetDisplayName() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}
	return u.Username
}

// Group describes the user group database structure
type Group struct {
	gorm.Model

	Name        string       `gorm:"not null;unique;"`
	Permissions []Permission `gorm:"many2many:group_permission;"`
}

// Permission describes a permission that can be linked to many groups
type Permission struct {
	gorm.Model

	Name string `gorm:"not_null;unique_index;"`
}

// NewApp generates and returns an instance of the app
func NewApp() *App {
	return &App{}
}

// App handles adding a route handler for static assets
type App struct{}

// Initialize sets up our app to handle embedded static asset requests
func (ua *App) Initialize(ev *enliven.Enliven) {
	if !ev.MiddlewareInstalled("session") {
		panic("The User app requires Session middleware to be registered.")
	}
	if !ev.AppInstalled("default_database") {
		panic("The User app requires that the Database app is initialized with a default connection.")
	}

	var conf = config.Config{
		"user_login_route":    "/user/login/",
		"user_logout_route":   "/user/logout/",
		"user_register_route": "/user/register/",
		"user_verify_route":   "/user/verify/{code}/",
		"user_password_route": "/user/password/",
		"user_profile_route":  "/user/profile/",

		// Where the user will be redirected after these successful actions
		"user_login_redirect":    "/",
		"user_logout_redirect":   "/",
		"user_register_redirect": "/user/login/",
		"user_password_redirect": "/",
		"user_profile_redirect":  "/user/profile/",

		"user_default_group":        "Member",
		"user_require_verification": "1",
	}

	conf = config.UpdateConfig(config.MergeConfig(conf, config.GetConfig()))

	db := database.GetDatabase(ev)

	// Migrating the user tables
	db.AutoMigrate(&Permission{}, &Group{}, &User{})
	ua.initDefaultUserModels(db)

	// Routing setup
	ev.AddRoute(conf["user_login_route"], LoginGetHandler, "GET")
	ev.AddRoute(conf["user_login_route"], LoginPostHandler, "POST")
	ev.AddRoute(conf["user_register_route"], RegisterGetHandler, "GET")
	ev.AddRoute(conf["user_register_route"], RegisterPostHandler, "POST")
	ev.AddRoute(conf["user_profile_route"], ProfileGetHandler, "GET")
	ev.AddRoute(conf["user_profile_route"], ProfilePostHandler, "POST")
	ev.AddRoute(conf["user_verify_route"], VerifyHandler)
	ev.AddRoute(conf["user_logout_route"], LogoutHandler)

	// Handles the setup of context variables to support user session management
	ev.AddMiddlewareFunc(SessionMiddleware)

	for _, templateType := range []string{"login", "password", "register", "verify", "profile", "verify_email"} {
		ev.Core.Templates.Parse(getTemplate(templateType))
	}

	admin.AddResources(&User{}, &Group{}, &Permission{})

	// Setting this app as the authorizer
	ev.Auth = ua
}

// initDefaultUser will set up the default admin user if the user database is empty.
func (ua *App) initDefaultUserModels(db *gorm.DB) {
	user := User{}
	var count int
	db.Find(&user).Count(&count)

	if count > 0 {
		return
	}

	member := Group{Name: "Member"}
	db.Create(&member)
	moderator := Group{Name: "Moderator"}
	db.Create(&moderator)
	admin := Group{Name: "Administrator"}
	db.Create(&admin)

	user = User{
		DisplayName:      "Administrator",
		Username:         "admin",
		Email:            "admin@admin.admin",
		Password:         GeneratePasswordHash("admin"),
		VerificationCode: "",
		Status:           1,
		Groups:           []Group{member, moderator, admin},
		Superuser:        true,
	}
	db.Create(&user)
}

// GetName returns the app's name
func (ua *App) GetName() string {
	return "user"
}

// HasPermission checks if a user has a permission matching the permission string passed in
func (ua *App) HasPermission(permission string, ctx *enliven.Context) bool {
	u := GetUser(ctx)
	if u != nil && u.UserHasPermission(permission) {
		return true
	}
	return false
}

// AddPermission adds a new permission to the user table if it doesn't exist
func (ua *App) AddPermission(permission string, ev *enliven.Enliven, groups ...string) {
	db := database.GetDatabase(ev)

	perm := Permission{}
	db.Where(&Permission{Name: permission}).First(&perm)

	if perm.ID == 0 {
		newPerm := Permission{
			Name: permission,
		}
		db.Create(&newPerm)

		// Adding this new permission to any specified groups.
		for _, groupName := range groups {
			group := Group{}
			db.Where("Name = ?", groupName).First(&group)

			if group.ID != 0 {
				group.Permissions = append(group.Permissions, newPerm)
			}

			db.Save(group)
		}
	}
}

// getTemplate looks up a template in config or embedded assets and returns its contents
func getTemplate(templateType string) string {
	conf := config.GetConfig()

	requestedTemplate := conf["user_"+templateType+"_template"]

	if requestedTemplate == "" {
		temp, _ := Asset("templates/" + templateType + ".html")

		if len(temp) > 0 {
			requestedTemplate = string(temp[:])
		}
	}

	return requestedTemplate
}
