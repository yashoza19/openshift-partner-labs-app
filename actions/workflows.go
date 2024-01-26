package actions

import (
	"bufio"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
	"net/http"
	"openshift-partner-labs-app/models"
	"os"
	"strings"
)

func AuditWorkflow(c buffalo.Context) error {
	if c.Request().Header.Get("X_WORKFLOW_HEADER_AUTH_TOKEN") != envy.Get("X_WORKFLOW_HEADER_AUTH_TOKEN", "") {
		return c.Render(http.StatusUnauthorized, r.JSON(map[string]string{"status": "unauthorized"}))
	}

	audit := &models.Audit{}

	err = c.Bind(audit)
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	if audit.LoginType != "console" {
		fmt.Println("Login type is not console")
		if audit.LoginType == "system:masters" {
			fmt.Println("Login type is system:masters; ignoring audit")
			return c.Render(http.StatusOK, r.JSON(map[string]string{"status": "ignored"}))
		}

		if checkIfIgnoredProject(audit.LoginName) {
			return c.Render(http.StatusOK, r.JSON(map[string]string{"status": "ignored"}))
		}

		fmt.Println("Removing system service account")
		audit.LoginName = removeSystemServiceAccount(audit.LoginName)
	}

	tx := c.Value("tx").(*pop.Connection)
	err = tx.Create(audit)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"status": "ok"}))
}

func removeSystemServiceAccount(loginName string) string {
	loginNameSplit := strings.Split(loginName, ":")
	if len(loginNameSplit) < 3 {
		return loginName
	}
	newLoginName := strings.Join(loginNameSplit[2:], ":")
	return newLoginName
}

func checkIfIgnoredProject(loginName string) bool {
	loginNameSplit := strings.Split(loginName, ":")
	if len(loginNameSplit) < 2 {
		return false
	}

	osproject := loginNameSplit[2]

	file, err := os.Open("osprojectignores")
	if err != nil {
		fmt.Printf("Unable to open osprojectignores file: %s", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return contains(lines, osproject)
}

func contains(lines []string, serviceaccount string) bool {
	for _, line := range lines {
		if line == serviceaccount {
			return true
		}
	}
	return false
}
