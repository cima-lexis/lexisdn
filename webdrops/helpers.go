package webdrops

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (sess *Session) DoGet(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating HTTP request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+sess.Token)

	res, err := client.Do(req)
	if err != nil {

		return nil, fmt.Errorf("Error submitting HTTP request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		//body, _ := ioutil.ReadAll(res.Body)
		//fmt.Println(string(body))

		return nil, fmt.Errorf("Error submitting request: HTTP status: %s", res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error downloading HTTP response: %w", err)
	}
	return body, nil
}

func (sess *Session) DoPost(url string, body interface{}) ([]byte, error) {
	bodyJ, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("Error converting body to JSON: %w", err)
	}
	bodyR := bytes.NewBuffer(bodyJ)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bodyR)
	if err != nil {
		return nil, fmt.Errorf("Error creating HTTP request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+sess.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error submitting HTTP request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error submitting request: HTTP status: %s", res.Status)
	}

	defer res.Body.Close()
	bodyResp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error downloading HTTP response: %w", err)
	}
	return bodyResp, nil
}

// Domain is
type Domain struct {
	MinLat, MinLon, MaxLat, MaxLon float64
}

var ItalyDomain = Domain{
	MaxLat: 66,
	MinLat: 23,
	MinLon: -18,
	MaxLon: 48,
}
