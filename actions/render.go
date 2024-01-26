package actions

import (
	"math"
	"openshift-partner-labs-app/public"
	"openshift-partner-labs-app/templates"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/render"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "app.plush.html",

		// fs.FS containing templates
		TemplatesFS: templates.FS(),

		// fs.FS containing assets
		AssetsFS: public.FS(),

		// Add template helpers here:
		Helpers: render.Helpers{
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,

			"progress": func(start, end time.Time) int {
				totalDuration := end.Sub(start)
				elapsedDuration := time.Now().Sub(start)
				timepassed := int(math.Round(elapsedDuration.Seconds() / totalDuration.Seconds() * 100))
				if timepassed > 100 {
					timepassed = 100
				}

				if timepassed < 0 {
					timepassed = 0
				}

				return timepassed
			},
			"setstartdate": func(start time.Time) string {
				return start.Format("2006-01-02")
			},
			"setenddate": func(end time.Time) string {
				return end.Format("2006-01-02")
			},
			"stripocp": func(version string) string {
				return strings.Join(strings.Split(version, "-")[1:], "-")
			},
			"extdate": func(ext time.Time) string {
				return ext.Format("01/02/2006")
			},
			"auditdate": func(audit time.Time) string {
				return audit.Format("01/02/2006 15:04:05")
			},
			"toupper": func(s string) string {
				return strings.ToUpper(s)
			},
			"color": func(pct int) string {
				switch {
				case pct > 0 && pct <= 50:
					return "green"
				case pct > 50 && pct <= 90:
					return "yellow"
				case pct > 90 && pct <= 100:
					return "red"
				default:
					return "grey"
				}
			},
		},
	})
}
