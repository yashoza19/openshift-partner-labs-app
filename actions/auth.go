package actions

import (
	"encoding/gob"
	"fmt"
	"github.com/gobuffalo/pop/v6"
	"net/http"
	"openshift-partner-labs-app/models"
	"os"
	"strings"

	"github.com/gobuffalo/envy"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func init() {
	gob.Register(&models.User{})

	gothic.Store = App().SessionStore

	scopes := []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"}

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/google/callback"), scopes...),
	)
}

func AuthCallback(c buffalo.Context) error {
	userauth, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	userdata := &models.User{}

	userdata.UserID = userauth.UserID
	userdata.FirstName = userauth.FirstName
	userdata.LastName = userauth.LastName
	userdata.FullName = userauth.Name
	userdata.Email = userauth.Email
	userdata.Picture = userauth.AvatarURL
	userdata.Admin = false

	admins := strings.Split(envy.Get("ADMIN_USERS", ""), ",")
	for _, admin := range admins {
		if admin == userauth.Email {
			userdata.Admin = true
		}
	}

	c.Session().Set("current_user", userdata)
	_ = c.Session().Save()

	tx := c.Value("tx").(*pop.Connection)
	count, _ := tx.Count(userdata)
	if count > 0 {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	err = tx.Create(userdata)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func Authorized(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if c.Session().Get("current_user") == nil {
			//c.Flash().Add("danger", "You must be logged in to access this page")
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
		return next(c)
	}
}

func SetUserData(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if c.Session().Get("current_user") != nil {
			c.Set("current_user", c.Session().Get("current_user"))
		}
		return next(c)
	}
}

func AuthDestroy(c buffalo.Context) error {
	gothic.Logout(c.Response(), c.Request())
	c.Session().Clear()
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
