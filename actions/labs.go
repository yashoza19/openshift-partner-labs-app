package actions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/buffalo/worker"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/x/responder"

	"github.com/coreos/go-semver/semver"
	"github.com/itchyny/gojq"

	"openshift-partner-labs-app/models"
)

// LabsResource is the resource for the Lab model
type LabsResource struct {
	buffalo.Resource
}

func CreateLabNote(c buffalo.Context) error {
	note := &models.Note{}
	err = c.Bind(note)
	if err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	err = tx.Create(note)
	if err != nil {
		return c.Render(http.StatusInternalServerError, r.Func("text/html", func(w io.Writer, d render.Data) error {
			wsn := WebsiteNotification(c, "close-circle", "red", "Unable to create note")
			_, err = w.Write([]byte(wsn))
			if err != nil {
				return err
			}
			return nil
		}))
	}

	return c.Render(http.StatusOK, r.Func("text/html", func(w io.Writer, d render.Data) error {
		wsn := WebsiteNotification(c, "check-all", "green", "Note Created")
		_, err = w.Write([]byte(wsn))
		if err != nil {
			return err
		}
		return nil
	}))
}

func LabApprove(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.Func("text/html", func(w io.Writer, d render.Data) error {
		wsn := WebsiteNotification(c, "check-all", "green", "Lab Request Approved")
		_, err = w.Write([]byte(wsn))
		if err != nil {
			return err
		}
		return nil
	}))
}

func LabDeny(c buffalo.Context) error {
	lab := &models.Lab{}

	tx := c.Value("tx").(*pop.Connection)

	if err = tx.Find(lab, c.Param("lab_id")); err != nil {
		return c.Render(http.StatusNotFound, r.Func("text/html", func(w io.Writer, d render.Data) error {
			wsn := WebsiteNotification(c, "close-circle", "red", "Lab Request Not Found")
			_, err = w.Write([]byte(wsn))
			if err != nil {
				return err
			}
			return nil
		}))
	}

	return c.Render(http.StatusOK, r.Func("text/html", func(w io.Writer, d render.Data) error {
		wsn := WebsiteNotification(c, "check-all", "green", "Lab Request Denied")
		_, err = w.Write([]byte(wsn))
		if err != nil {
			return err
		}
		return nil
	}))
}

func LabExtension(c buffalo.Context) error {
	current_user := c.Value("current_user").(*models.User)

	err = w.Perform(worker.Job{
		Queue: "default",
		Args: worker.Args{
			"lab_id":       c.Param("lab_id"),
			"current_user": current_user.FullName,
			"extension":    c.Param("extension"),
		},
		Handler: "request_extension",
	})

	if err != nil {
		return c.Render(http.StatusInternalServerError, r.Func("text/html", func(w io.Writer, d render.Data) error {
			reqexttr := fmt.Sprintf("<tr><td class=\"px-6 py-3.5 dark:text-zinc-100\">%s</td><td class=\"px-6 py-3.5 dark:text-zinc-100\">%s</td><td class=\"px-6 py-3.5 dark:text-zinc-100\">%s</td></tr>", current_user.FullName, time.Now().Format("01/02/2006"), c.Param("extension"))
			_, err = w.Write([]byte(reqexttr))
			if err != nil {
				return err
			}
			return nil
		}))
	}

	c.Response().Header().Set("HX-Trigger", "log-extension")

	return c.Render(http.StatusOK, r.Func("text/html", func(w io.Writer, d render.Data) error {
		reqexttr := fmt.Sprintf("<tr><td class=\"px-6 py-3.5 dark:text-zinc-100\">%s</td><td class=\"px-6 py-3.5 dark:text-zinc-100\">%s</td><td class=\"px-6 py-3.5 dark:text-zinc-100\">%s</td></tr>", current_user.FullName, time.Now().Format("01/02/2006"), c.Param("extension"))
		_, err = w.Write([]byte(reqexttr))
		if err != nil {
			return err
		}
		return nil
	}))
}

// List gets all Labs. This function is mapped to the path
// GET /labs
func (v LabsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	labs := &models.Labs{}

	// Retrieve all Labs from the DB
	if err = tx.Where("state != 'complete' and state != 'denied'").All(labs); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context, so it can be used in the template.
		c.Set("labs", labs)
		return c.Render(http.StatusOK, r.HTML("pages/labs/index.plush.html", "layouts/layout-labs.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(labs))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(labs))
	}).Respond(c)
}

// Show gets the data for one Lab. This function is mapped to
// the path GET /labs/{lab_id}
func (v LabsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Lab
	lab := &models.Lab{}
	exts := &models.Regexts{}
	audits := &models.Audits{}
	notes := &models.Notes{}

	// To find the Lab the parameter lab_id is used.
	if err = tx.Find(lab, c.Param("lab_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	if err = tx.Where("lab_id = ?", c.Param("lab_id")).All(exts); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	if err = tx.Where("generated_name = ?", lab.GeneratedName).All(audits); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	if err = tx.Where("lab_id = ?", c.Param("lab_id")).All(notes); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("lab", lab)
		c.Set("exts", exts)
		c.Set("audits", audits)
		c.Set("notes", notes)

		return c.Render(http.StatusOK, r.HTML("pages/labs/show.plush.html", "layouts/layout-labs.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(lab))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(lab))
	}).Respond(c)
}

// New renders the form for creating a new Lab.
// This function is mapped to the path GET /labs/new
func (v LabsResource) New(c buffalo.Context) error {
	c.Set("lab", &models.Lab{})

	openShiftReleaseImages := getOpenShiftReleaseImages()
	c.Set("releaseimages", openShiftReleaseImages)

	return c.Render(http.StatusOK, r.HTML("pages/labs/new.plush.html", "layouts/layout-labs.plush.html"))
}

// Create adds a Lab to the DB. This function is mapped to the
// path POST /labs
func (v LabsResource) Create(c buffalo.Context) error {
	// Allocate an empty Lab
	lab := &models.Lab{}

	// Bind lab to the html form elements
	if err = c.Bind(lab); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(lab)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template

			c.Set("errors", verrs.Errors)
			c.Set("errorscount", verrs.Count())

			// Render again the new.html template that the user can
			// correct the input.
			c.Set("lab", lab)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("pages/labs/new.plush.html", "layouts/layout-labs.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	/*
		err = sendGoogleNotification(c, lab)
		if err != nil {
			c.Flash().Add("danger", "Unable to send notification to Google Chat")
		}
	*/

	err = sendSlackNotification(lab)
	if err != nil {
		c.Flash().Add("danger", "Unable to send notification to Slack")
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", "Lab request was created successfully!")

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/labs/%v", lab.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(lab))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(lab))
	}).Respond(c)
}

// Edit renders an edit form for a Lab. This function is
// mapped to the path GET /labs/{lab_id}/edit
func (v LabsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Lab
	lab := &models.Lab{}

	openShiftReleaseImages := getOpenShiftReleaseImages()
	c.Set("releaseimages", openShiftReleaseImages)

	if err = tx.Find(lab, c.Param("lab_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("lab", lab)
	return c.Render(http.StatusOK, r.HTML("pages/labs/edit.plush.html", "layouts/layout-labs.plush.html"))
}

// Update changes a Lab in the DB. This function is mapped to
// the path PUT /labs/{lab_id}
func (v LabsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Lab
	lab := &models.Lab{}

	if err = tx.Find(lab, c.Param("lab_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Lab to the html form elements
	if err = c.Bind(lab); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(lab)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)
			fmt.Println(verrs)
			c.Set("errorscount", verrs.Count())

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("lab", lab)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("pages/labs/edit.plush.html", "layouts/layout-labs.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", "Lab request was updated successfully!")

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/labs/%v", lab.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(lab))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(lab))
	}).Respond(c)
}

// Destroy deletes a Lab from the DB. This function is mapped
// to the path DELETE /labs/{lab_id}
func (v LabsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Lab
	lab := &models.Lab{}

	// To find the Lab the parameter lab_id is used.
	if err = tx.Find(lab, c.Param("lab_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err = tx.Destroy(lab); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", "Lab request was destroyed successfully!")

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/labs")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(lab))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(lab))
	}).Respond(c)
}

// this needs to become a background job that runs every 30 or 60 minutes
// and updates the database with the latest release images
// this will allow pulling the latest release images from the database
// instead of querying the API every time the page is loaded
func getOpenShiftReleaseImages() (tags []*semver.Version) {
	resp, err := http.Get("https://quay.io/api/v1/repository/openshift-release-dev/ocp-release?includeTags=true")
	if err != nil {
		fmt.Println("Unable to query image source API: ", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read response body: ", err)
	}
	defer resp.Body.Close()

	repository := make(map[string]interface{})

	if json.Unmarshal(body, &repository) != nil {
		log.Fatalln(err)
	}

	query, err := gojq.Parse(".tags|with_entries(select(.key|match(\"[4|5].[0-9]+.[0-9]+(-(rc|ec).[0-9]+)?-x86_64\")))")
	if err != nil {
		fmt.Println("Unable to compile gojq query string: ", err)
	}

	iquery := query.Run(repository)
	images, _ := iquery.Next()

	for _, image := range images.(map[string]interface{}) {
		name := strings.Split(image.(map[string]interface{})["name"].(string), "-x86_64")[0]
		tag := semver.New(name)
		tags = append(tags, tag)
	}

	semver.Sort(tags)
	for i, j := 0, len(tags)-1; i < j; i, j = i+1, j-1 {
		tags[i], tags[j] = tags[j], tags[i] // reverse the tags
	}
	return tags
}
