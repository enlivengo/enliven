package user

import (
	"strconv"

	"github.com/hickeroar/enliven"
)

// SessionMiddleware handles adding the elements to the context that carry the user's id and status
func SessionMiddleware(ctx *enliven.Context, next enliven.NextHandlerFunc) {
	if ctx.Session == nil {
		panic("The User app requires Session middleware to be registered.")
	}

	userID := ctx.Session.Get("UserApp_LoggedInUserID")

	// If there isn't a user id in the session, we set context items accordingly
	if userID == "" {
		ctx.Booleans["UserLoggedIn"] = false
		ctx.Integers["UserID"] = 0
	} else {
		ctx.Booleans["UserLoggedIn"] = true
		ctx.Integers["UserID"], _ = strconv.Atoi(userID)

		// Caching the user so we can use it via storage.
		GetUser(ctx)
	}

	next(ctx)
}
