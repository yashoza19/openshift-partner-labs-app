package actions

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"net/http"
	"sync"

	"openshift-partner-labs-app/locales"
	"openshift-partner-labs-app/models"
	"openshift-partner-labs-app/public"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/csrf"
	"github.com/gobuffalo/middleware/forcessl"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app     *buffalo.App
	appOnce sync.Once
	T       *i18n.Translator
	err     error
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	appOnce.Do(func() {
		appSession := sessions.NewFilesystemStore("/tmp", []byte(envy.Get("SESSION_SECRET", "secret")))
		appSession.Options.SameSite = http.SameSiteLaxMode

		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionName:  "_openshift_partner_labs_app_session",
			SessionStore: appSession,
			/*Worker: gwa.New(gwa.Options{
				Pool: &redis.Pool{
					MaxIdle:   5,
					MaxActive: 5,
					Wait:      true,
					Dial: func() (redis.Conn, error) {
						return redis.Dial("tcp", envy.Get("REDIS_URL", ":6379"),
							redis.DialUsername(envy.Get("REDIS_USERNAME", "")),
							redis.DialPassword(envy.Get("REDIS_PASSWORD", "")))
					},
				},
				Name:           "openshift-partner-labs-app",
				MaxConcurrency: 25,
			}),*/
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		_csrf := csrf.New
		app.Use(_csrf)

		// Wraps each request in a transaction.
		//   c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))
		// Setup and use translations:
		app.Use(translations())

		app.Use(Authorized)
		app.Use(SetUserData)

		app.GET("/", HomeHandler)

		auth := app.Group("/auth")
		AuthLogin := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)
		auth.GET("/{provider}", AuthLogin)
		auth.GET("/{provider}/callback", AuthCallback)
		auth.GET("/{provider}/logout", AuthDestroy)
		auth.Middleware.Skip(Authorized, AuthCallback, AuthLogin)

		wrkflw := app.Group("/workflows")
		wrkflw.Middleware.Skip(_csrf, AuditWorkflow)
		wrkflw.Middleware.Skip(Authorized, AuditWorkflow)
		wrkflw.POST("/audit", AuditWorkflow)

		app.GET("/labs/archive", LabsArchive)
		app.Resource("/labs", LabsResource{})
		app.POST("/labs/approve/{lab_id}", LabApprove).Name("labApprove")
		app.POST("/labs/deny/{lab_id}", LabDeny).Name("labDeny")
		app.POST("/labs/extension", LabExtension)
		app.POST("/labs/{lab_id}/notes", CreateLabNote)
		//app.GET("/labs/extensions/{lab_id}", GetLabExtensions)

		app.GET("/login", LoginHandler)
		app.Middleware.Skip(Authorized, LoginHandler)

		app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory
	})

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
