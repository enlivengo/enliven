package user

import (
	"strconv"

	"github.com/enlivengo/enliven"
)

// SessionMiddleware handles adding the elements to the context that carry the user's id and status
func SessionMiddleware(ctx *enliven.Context, next enliven.NextHandlerFunc) {
	userID := ctx.Session.Get("user_id")

	// If there isn't a user id in the session, we set context items accordingly
	if userID == "" {
		ctx.Booleans["UserLoggedIn"] = false
		ctx.Integers["UserID"] = 0
		ctx.Booleans["UserSuperUser"] = false
	} else {
		ctx.Booleans["UserLoggedIn"] = true
		ctx.Integers["UserID"], _ = strconv.Atoi(userID)
		u := GetUser(ctx)
		ctx.Booleans["UserSuperUser"] = u.Superuser
	}

	next(ctx)
}
