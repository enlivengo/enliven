package user

//go:generate go-bindata -o templates.go -pkg user templates/...

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/enlivengo/enliven"
	"github.com/enlivengo/enliven/apps/database"
	"github.com/enlivengo/enliven/config"
	"github.com/jmcvetta/randutil"
)

// FormError represents an error in form validation
type FormError struct {
	Message string
	Field   string
}

// LoginGetHandler handles get requests to the login route
func LoginGetHandler(ctx *enliven.Context) {
	ctx.Template("user_login")
}

// LoginPostHandler handles the form submission for logging a user in.
func LoginPostHandler(ctx *enliven.Context) {
	ctx.Request.ParseForm()
	login := ctx.Request.Form.Get("username")
	password := ctx.Request.Form.Get("password")

	conf := config.GetConfig()
	db := database.GetDatabase(ctx.Enliven)

	user := User{}
	var where string
	if strings.Contains(login, "@") {
		where = "Email = ?"
	} else {
		where = "Username = ?"
		login = strings.ToLower(login)
	}
	db.Where(where+" AND Status = ?", login, 1).First(&user)

	if user.ID == 0 || !VerifyPasswordHash(user.Password, password) {
		ctx.Strings["LoginError"] = "Invalid Login or Password."
		LoginGetHandler(ctx)
		return
	}

	ctx.Session.Set("user_id", strconv.FormatUint(uint64(user.ID), 10))
	ctx.Redirect(conf["user_login_redirect"])
}

// RegisterGetHandler handles get requests to the register route
func RegisterGetHandler(ctx *enliven.Context) {
	ctx.Strings["FormErrors"] = "[]"
	ctx.Template("user_register")
}

// RegisterPostHandler handles get requests to the register route
func RegisterPostHandler(ctx *enliven.Context) {
	ctx.Request.ParseForm()
	var errors []FormError
	db := database.GetDatabase(ctx.Enliven)
	conf := config.GetConfig()
	verifyAccount := (conf["user_require_verification"] == "1" && ctx.Enliven.Core.Email.Enabled())

	// Making sure none of the required fields are empty
	for _, field := range []string{"username", "email", "password", "verifyPassword"} {
		if len(strings.TrimSpace(ctx.Request.Form.Get(field))) == 0 {
			errors = append(errors, FormError{
				Message: "Field '" + strings.Title(field) + "' is required.",
				Field:   field,
			})
		}
	}

	username := strings.TrimSpace(ctx.Request.Form.Get("username"))
	email := strings.TrimSpace(ctx.Request.Form.Get("email"))
	password := ctx.Request.Form.Get("password")
	verifyPassword := ctx.Request.Form.Get("verifyPassword")

	if len(username) < 3 {
		errors = append(errors, FormError{
			Message: "Your username must be three characters in length or longer.",
			Field:   "username",
		})
	}

	if !govalidator.IsAlphanumeric(username) {
		errors = append(errors, FormError{
			Message: "Username must contain only letters and numbers.",
			Field:   "username",
		})
	}

	userNameCheck := User{}
	db.Where("Username = ?", username).First(&userNameCheck)
	if userNameCheck.ID != 0 {
		errors = append(errors, FormError{
			Message: "The Username you have entered is already registered.",
			Field:   "username",
		})
	}

	if !govalidator.IsEmail(email) {
		errors = append(errors, FormError{
			Message: "The provided email address is invalid.",
			Field:   "email",
		})
	}

	userEmailCheck := User{}
	db.Where("Email = ?", email).First(&userEmailCheck)
	if userEmailCheck.ID != 0 {
		errors = append(errors, FormError{
			Message: "The Email you have entered is already registered.",
			Field:   "email",
		})
	}

	if password != verifyPassword {
		errors = append(errors, FormError{
			Message: "The passwords do not match.",
			Field:   "verifyPassword",
		})
	}

	if len(errors) > 0 {
		jsonResponse, _ := json.Marshal(errors)
		ctx.Strings["FormErrors"] = string(jsonResponse[:])
		ctx.Strings["RegisterUsername"] = username
		ctx.Strings["RegisterEmail"] = email
		ctx.Template("user_register")
		return
	}

	newUser := User{
		DisplayName:      username,
		Username:         strings.ToLower(username),
		Email:            email,
		Password:         GeneratePasswordHash(password),
		VerificationCode: "",
		Status:           1,
		Superuser:        false,
	}

	userGroup := Group{}
	db.Where("Name = ?", conf["user_default_group"]).First(&userGroup)

	if userGroup.ID != 0 {
		newUser.Groups = []Group{userGroup}
	}

	if verifyAccount {
		verificationCode, _ := randutil.AlphaString(32)
		newUser.VerificationCode = verificationCode
		newUser.Status = 0
	}

	db.Create(&newUser)

	if verifyAccount {
		templateData := map[string]interface{}{
			"User":    newUser,
			"Context": ctx,
			"Config":  conf,
		}
		var bMessage bytes.Buffer
		eerr := ctx.Enliven.Core.Templates.ExecuteTemplate(&bMessage, "verify_email", templateData)
		if eerr != nil {
			fmt.Println(eerr.Error())
		}

		email := ctx.Enliven.Core.Email.New()
		email.AddRecipient(newUser.Email)
		email.Subject = conf["site_name"] + ": Please Verify Your Account"
		email.Message = bMessage.String()
		err := email.Send()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	ctx.Redirect(conf["user_register_redirect"])
}

// ProfileGetHandler displays the profile editor
func ProfileGetHandler(ctx *enliven.Context) {
	u := GetUser(ctx)
	if u.ID == 0 {
		ctx.Forbidden()
		return
	}
	ctx.Strings["FormErrors"] = "[]"
	ctx.Template("user_profile")
}

// ProfilePostHandler handles the updating of a user's profile
func ProfilePostHandler(ctx *enliven.Context) {
	ctx.Request.ParseForm()
	var errors []FormError
	db := database.GetDatabase(ctx.Enliven)
	conf := config.GetConfig()
	u := GetUser(ctx)

	// Making sure none of the required fields are empty
	for _, field := range []string{"username", "email"} {
		if len(strings.TrimSpace(ctx.Request.Form.Get(field))) == 0 {
			errors = append(errors, FormError{
				Message: "Field '" + strings.Title(field) + "' is required.",
				Field:   field,
			})
		}
	}

	username := strings.TrimSpace(ctx.Request.Form.Get("username"))
	email := strings.TrimSpace(ctx.Request.Form.Get("email"))
	password := ctx.Request.Form.Get("password")
	verifyPassword := ctx.Request.Form.Get("verifyPassword")

	u.DisplayName = strings.TrimSpace(ctx.Request.Form.Get("displayName"))

	if len(username) < 3 {
		errors = append(errors, FormError{
			Message: "Username must be three characters in length or longer.",
			Field:   "username",
		})
	}

	if !govalidator.IsAlphanumeric(username) {
		errors = append(errors, FormError{
			Message: "Username must contain only letters and numbers.",
			Field:   "username",
		})
	}

	userNameCheck := User{}
	db.Where("Username = ? AND ID <> ?", username, u.ID).First(&userNameCheck)
	if userNameCheck.ID != 0 {
		errors = append(errors, FormError{
			Message: "The Username you have entered is already registered.",
			Field:   "username",
		})
	}

	u.Username = username

	if !govalidator.IsEmail(email) {
		errors = append(errors, FormError{
			Message: "The provided email address is invalid.",
			Field:   "email",
		})
	}

	userEmailCheck := User{}
	db.Where("Email = ? AND ID <> ?", email, u.ID).First(&userEmailCheck)
	if userEmailCheck.ID != 0 {
		errors = append(errors, FormError{
			Message: "The Email you have entered is already registered.",
			Field:   "email",
		})
	}

	u.Email = email

	if len(strings.TrimSpace(password)) > 0 {
		if password != verifyPassword {
			errors = append(errors, FormError{
				Message: "The passwords do not match.",
				Field:   "verifyPassword",
			})
		} else {
			u.Password = GeneratePasswordHash(password)
		}
	}

	if len(errors) > 0 {
		jsonResponse, _ := json.Marshal(errors)
		ctx.Strings["FormErrors"] = string(jsonResponse[:])
		ctx.Storage["User"] = u
		ctx.Template("user_profile")
		return
	}

	db.Save(u)
	ctx.Redirect(conf["user_profile_redirect"])
}

// VerifyHandler checks a user's verification code to see if it matches.
func VerifyHandler(ctx *enliven.Context) {
	ctx.Booleans["Verified"] = false

	code, ok := ctx.Vars["code"]
	if !ok {
		ctx.Template("user_verify")
		return
	}

	db := database.GetDatabase(ctx.Enliven)
	u := User{}
	db.Where("Verification_Code = ? AND Status = ?", code, 0).First(&u)

	if u.ID == 0 {
		ctx.Template("user_verify")
		return
	}

	u.Status = 1
	db.Save(&u)

	ctx.Booleans["Verified"] = true
	ctx.Strings["LoginURL"] = config.GetConfig()["user_login_route"]
	ctx.Template("user_verify")
}

// LogoutHandler logs a user out and redirects them to the configured page.
func LogoutHandler(ctx *enliven.Context) {
	ctx.Session.Destroy()
	ctx.Redirect(config.GetConfig()["user_logout_redirect"])
}
