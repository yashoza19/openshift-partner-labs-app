package grifts

import (
	"redhat_openshift_partner_labs_app/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
