package actions

import (
	"errors"
	"fmt"
	"openshift-partner-labs-app/models"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo/worker"
)

var w worker.Worker

func init() {
	w = App().Worker

	w.Register("send_welcome_email", func(args worker.Args) error {
		fmt.Println(args)
		return nil
	})

	w.Register("send_login_info_email", func(args worker.Args) error {
		fmt.Println(args)
		return nil
	})

	w.Register("request_extension", func(args worker.Args) error {
		fmt.Println("Starting extension request job")

		tx := models.DB

		fmt.Println("Getting the lab request")
		lab := &models.Lab{}
		err = tx.Find(lab, args["lab_id"])
		if err != nil {
			return err
		}

		fmt.Println("Obtaining the new end date")
		endtime, err := setEndDate(lab.EndDate.Format("2006-01-02"), args["extension"].(string))
		if err != nil {
			return err
		}

		fmt.Println("Creating the request extension record")
		labid, _ := strconv.Atoi(args["lab_id"].(string))
		regext := &models.Regext{}
		regext.LabID = labid
		regext.CurrentUser = args["current_user"].(string)
		regext.Extension = args["extension"].(string)
		regext.Date = time.Now()

		_, err = tx.ValidateAndCreate(regext)
		if err != nil {
			return err
		}

		fmt.Println("Updating the lab request")
		lab = &models.Lab{}
		err = tx.Find(lab, labid)
		if err != nil {
			fmt.Println("Unable to find the lab request to update")
			return err
		}
		fmt.Println("Current End Date: ", lab.EndDate)
		lab.EndDate = endtime
		lab.State = "extended"
		fmt.Println("New End Date: ", lab.EndDate)
		err = tx.Update(lab)
		if err != nil {
			fmt.Println("Update the status of the request extension record: fail")
			regext.Status = "fail"
			err = tx.Update(regext)
			return err
		}

		fmt.Println("Update the status of the request extension record: pass")
		regext.Status = "pass"
		err = tx.Update(regext)
		if err != nil {
			return err
		}

		return nil
	})

	w.Register("email_extension_confirmation", func(args worker.Args) error {
		fmt.Println(args)
		return nil
	})
}

func setEndDate(startdate string, leasetime string) (endtime time.Time, err error) {
	layout := "2006-01-02"
	var enddate time.Time
	start, err := time.Parse(layout, startdate)
	if err != nil {
		return time.Time{}, errors.New("invalid start_date")
	}

	switch leasetime {
	case "1d":
		enddate = start.AddDate(0, 0, 1)
	case "1w":
		enddate = start.AddDate(0, 0, 7)
	case "2w":
		enddate = start.AddDate(0, 0, 14)
	case "1m":
		enddate = start.AddDate(0, 1, 0)
	case "2m":
		enddate = start.AddDate(0, 2, 0)
	case "3m":
		enddate = start.AddDate(0, 3, 0)
	case "6m":
		enddate = start.AddDate(0, 6, 0)
	default:
		return time.Time{}, errors.New("invalid lease_time")
	}

	endDate := enddate.Format("2006-01-02 15:04:05")
	endtime, err = time.Parse("2006-01-02 15:04:05", endDate)

	return endtime, nil
}
