package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/slack-go/slack"
	"io"
	"log"
	"net/http"
	"openshift-partner-labs-app/models"
	"strconv"
)

// Custom Flash Responses
var FlashResponse = `<div style="position: absolute; z-index: 1000; left: 0; right: 0; margin: auto; width: 30%; -moz-animation: cssAnimation 0s ease-in 15s forwards; -webkit-animation: cssAnimation 0s ease-in 15s forwards; -o-animation: cssAnimation 0s ease-in 15s forwards;
		animation: cssAnimation 0s ease-in 15s forwards; -webkit-animation-fill-mode: forwards; animation-fill-mode: forwards;">
		<div class="px-5 py-2.5 mb-5 flex items-center bg-red-50 border-l-4 border-red-500 rounded rounded-l-none alert-dismissible">
				<i class="mdi mdi-close-circle ltr:mr-2 rtl:ml-2 align-middle text-red-700 text-lg"></i>
				<p class="text-red-700">Lab Request Not Found</p>
				<button class="alert-close ltr:ml-auto rtl:mr-auto text-red-400 text-lg"><i class="mdi mdi-close"></i></button>
		</div>`

// Custom types for interaction with Google Chat
// only fields specific to our use case created

type GChatMessage struct {
	Cards []Cards `json:"cardsV2"`
}

type Header struct {
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	ImageURL     string `json:"imageUrl"`
	ImageType    string `json:"imageType"`
	ImageAltText string `json:"imageAltText,omitempty"`
}

type OpenLink struct {
	URL string `json:"url,omitempty"`
}

type OnClick struct {
	OpenLink OpenLink `json:"openLink,omitempty"`
}

type Buttons struct {
	Text    string  `json:"text,omitempty"`
	OnClick OnClick `json:"onClick,omitempty"`
}

type ButtonList struct {
	Buttons []Buttons `json:"buttons,omitempty"`
}

type DecoratedText struct {
	TopLabel string `json:"topLabel"`
	Text     string `json:"text"`
}

type Widgets struct {
	DecoratedText *DecoratedText `json:"decoratedText,omitempty"`
	ButtonList    *ButtonList    `json:"buttonList,omitempty"`
}

type Sections struct {
	Widgets []Widgets `json:"widgets,omitempty"`
}

type Card struct {
	Header   Header     `json:"header"`
	Sections []Sections `json:"sections"`
}

type Cards struct {
	CardID string `json:"cardId"`
	Card   Card   `json:"card"`
}

func sendGoogleNotification(c buffalo.Context, l *models.Lab) error {
	current_user := c.Value("current_user").(*models.User)
	// Build the message to send to Google Chat
	msg, err := json.Marshal(
		GChatMessage{
			Cards: []Cards{
				{
					CardID: l.ClusterID.String(),
					Card: Card{
						Header: Header{
							Title:     "New Lab Request - " + l.GeneratedName,
							Subtitle:  "Sponsor: " + l.Sponsor,
							ImageURL:  current_user.Picture,
							ImageType: "CIRCLE",
						},
						Sections: []Sections{
							{
								Widgets: []Widgets{
									{
										DecoratedText: &DecoratedText{
											TopLabel: "Company Name",
											Text:     l.CompanyName,
										},
									},
									{
										DecoratedText: &DecoratedText{
											TopLabel: "OpenShift Version",
											Text:     l.OpenShiftVersion,
										},
									},
									{
										DecoratedText: &DecoratedText{
											TopLabel: "Cluster Size",
											Text:     l.ClusterSize,
										},
									},
									{
										DecoratedText: &DecoratedText{
											TopLabel: "Request Type",
											Text:     l.RequestType,
										},
									},
									{
										DecoratedText: &DecoratedText{
											TopLabel: "Cloud Provider",
											Text:     l.CloudProvider,
										},
									},
									{
										DecoratedText: &DecoratedText{
											TopLabel: "Primary Contact",
											Text:     l.PrimaryFirst + " " + l.PrimaryLast + "\n" + l.PrimaryEmail,
										},
									},
									{
										ButtonList: &ButtonList{
											[]Buttons{
												{
													Text: "REVIEW",
													OnClick: OnClick{
														OpenLink: OpenLink{URL: fmt.Sprintf(envy.Get("HOST", "http://localhost:3000")+"/labs/%s", strconv.Itoa(l.ID))},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		})

	// Send the message to Google Chat
	req := bytes.NewBuffer(msg)

	// Leverage Go's HTTP Post function to make request
	resp, err := http.Post(envy.Get("GCHAT_WEBHOOK_URL", ""), "application/json", req)

	//Handle Error
	if err != nil {
		log.Printf("Erroring posting to Google Chat: %v", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Unable to properly close POST request for Google Chat message: %v", err)
		}
	}(resp.Body)

	return nil
}

func sendSlackNotification(l *models.Lab) error {
	var whm slack.WebhookMessage

	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "*New Lab Request - "+l.GeneratedName+":*\nSponsor: "+l.Sponsor, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Fields
	companyName := slack.NewTextBlockObject("mrkdwn", "*Company Name:*\n"+l.CompanyName, false, false)
	openshiftVersion := slack.NewTextBlockObject("mrkdwn", "*OpenShift Version:*\n"+l.OpenShiftVersion, false, false)
	clusterSize := slack.NewTextBlockObject("mrkdwn", "*Cluster Size:*\n"+l.ClusterSize, false, false)
	requestType := slack.NewTextBlockObject("mrkdwn", "*Request Type:*\n"+l.RequestType, false, false)
	cloudProvider := slack.NewTextBlockObject("mrkdwn", "*Cloud Provider:*\n"+l.CloudProvider, false, false)
	primaryContact := slack.NewTextBlockObject("mrkdwn", "*Primary Contact:*\n"+l.PrimaryFirst+" "+l.PrimaryLast+" <"+l.PrimaryEmail+">", false, false)

	fieldSlice := make([]*slack.TextBlockObject, 0)
	fieldSlice = append(fieldSlice, companyName)
	fieldSlice = append(fieldSlice, openshiftVersion)
	fieldSlice = append(fieldSlice, clusterSize)
	fieldSlice = append(fieldSlice, requestType)
	fieldSlice = append(fieldSlice, cloudProvider)
	fieldSlice = append(fieldSlice, primaryContact)

	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)

	viewBtnTxt := slack.NewTextBlockObject("plain_text", "View", false, false)
	viewBtn := slack.NewButtonBlockElement("", "view_request", viewBtnTxt)
	viewBtn.URL = envy.Get("HOST", "http://localhost:3000") + "/labs/" + strconv.Itoa(l.ID)

	// Approve and Deny Buttons
	actionBlock := slack.NewActionBlock("", viewBtn)

	msg := slack.Blocks{BlockSet: []slack.Block{
		headerSection,
		fieldsSection,
		actionBlock,
	}}

	whm.Blocks = &msg

	err = slack.PostWebhook(envy.Get("SLACK_WEBHOOK_URL", ""), &whm)
	if err != nil {
		return err
	}
	return nil
}

var WebsiteNotification = func(c buffalo.Context, icon string, color string, message string) string {
	prefix := `<div style="position: absolute; z-index: 1000; left: 0; right: 0; margin: auto; width: 30%; -moz-animation: cssAnimation 0s ease-in 15s forwards; -webkit-animation: cssAnimation 0s ease-in 15s forwards; -o-animation: cssAnimation 0s ease-in 15s forwards;
		animation: cssAnimation 0s ease-in 15s forwards; -webkit-animation-fill-mode: forwards; animation-fill-mode: forwards;">`

	infix := fmt.Sprintf("<div class=\"px-5 py-2.5 mb-5 flex items-center bg-%s-50 border-l-4 border-%s-500 rounded rounded-l-none alert-dismissible\">", color)
	if icon != "" {
		infix += fmt.Sprintf("<i class=\"mdi mdi-%s ltr:mr-2 rtl:ml-2 align-middle text-%s-700 text-lg\"></i>", icon, color)
	}
	infix += fmt.Sprintf("<p class=\"text-%s-700\">%s</p>", color, message)
	infix += fmt.Sprintf("<button class=\"alert-close ltr:ml-auto rtl:mr-auto text-%s-400 text-lg\"><i class=\"mdi mdi-close\"></i></button>", color)
	infix += "</div>"

	suffix := `</div>`

	return prefix + infix + suffix
}
