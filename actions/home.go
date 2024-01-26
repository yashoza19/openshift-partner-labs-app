package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	c.Set("title", "Dashboard")
	c.Set("page_title", "Dashboard")
	c.Set("sub_title", "Menu")

	return c.Render(http.StatusOK, r.HTML("pages/home/index.plush.html"))
}

func LoginHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("pages/home/login.plush.html", "layouts/layout-login.plush.html"))
}
