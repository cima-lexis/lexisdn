package webdrops

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cima-lexis/lexisdn/config"
)

type Session struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string
}

func parseTemplate(tmplS string, args interface{}) (io.Reader, error) {
	tmpl := template.New("query")
	_, err := tmpl.Parse("client_id={{.ClientID}}&grant_type=password&password={{.Password}}&username={{.User}}")
	if err != nil {
		return nil, fmt.Errorf("Error parsing querystring template: %w", err)
	}

	var tmplResult bytes.Buffer

	err = tmpl.Execute(&tmplResult, config.Config)
	if err != nil {
		return nil, fmt.Errorf("Error building querystring: %w", err)
	}

	return strings.NewReader(tmplResult.String()), nil
}

func (sess *Session) Login() error {
	payload, err := parseTemplate(
		"client_id={{.ClientID}}&grant_type=password&password={{.Password}}&username={{.User}}",
		config.Config,
	)

	client := &http.Client{}
	req, err := http.NewRequest("POST", config.Config.AuthURL, payload)
	if err != nil {
		return fmt.Errorf("Error creating HTTP request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error submitting HTTP request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error downloading HTTP response: %w", err)
	}
	sess.ClientID = config.Config.ClientID
	return json.Unmarshal(body, sess)
}

func (sess *Session) Refresh() error {
	payload, err := parseTemplate(
		"client_id={{.ClientID}}&grant_type=refresh_token&refresh_token={{.RefreshToken}}",
		config.Config,
	)

	client := &http.Client{}
	req, err := http.NewRequest("POST", config.Config.AuthURL, payload)
	if err != nil {
		return fmt.Errorf("Error creating HTTP request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error submitting HTTP request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error downloading HTTP response: %w", err)
	}

	return json.Unmarshal(body, sess)

}
