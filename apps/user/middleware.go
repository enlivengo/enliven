package user

import "github.com/hickeroar/enliven"

// SessionMiddleware handles adding the elements to the context that carry the user's id and status
func SessionMiddleware(ctx *enliven.Context, next enliven.NextHandlerFunc) {
	if ctx.Session == nil {
		panic("The User app requires Session middleware to be registered.")
	}

	userID := ctx.Session.Get("UserApp_LoggedInUserID")

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
