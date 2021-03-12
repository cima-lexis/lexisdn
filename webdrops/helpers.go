package webdrops

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const maxRetry = 5

/*
func (sess *Session) login() error {

	data := url.Values{}
	data.Set("client_id", config.Config.ClientID)
	data.Set("grant_type", "password")
	data.Set("password", config.Config.Password)
	data.Set("username", config.Config.User)

	client := &http.Client{
		Timeout: time.Second,
	}
	req, err := http.NewRequest("POST", config.Config.AuthURL, strings.NewReader(data.Encode()))
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
	//fmt.Println(string(body))
	if err != nil {
		return fmt.Errorf("Error downloading HTTP response: %w", err)
	}
	sess.ClientID = config.Config.ClientID
	err = json.Unmarshal(body, sess)
	sess.RefreshedAt = time.Now()
	sess.ExpiresIn = 30
	return err
}
*/

// DoGet ...
func (sess *Session) DoGet(url string) (res []byte, err error) {

	for i := time.Duration(0); i < maxRetry; i++ {
		err = sess.refresh()
		if err != nil {
			return nil, err
		}

		res, err = sess.get(url)
		if err == nil {
			return
		}
		fmt.Fprintf(os.Stderr, "An error occurred while getting from %s:%s\n", url, err.Error())
		time.Sleep(i * 10 * time.Second)

	}

	return
}

// DoPost ...
func (sess *Session) DoPost(url string, body interface{}) (res []byte, err error) {

	for i := time.Duration(0); i < maxRetry; i++ {
		err = sess.refresh()
		if err != nil {
			return nil, err
		}

		res, err = sess.post(url, body)
		if err == nil {
			return
		}

		fmt.Fprintf(os.Stderr, "An error occurred while posting to %s:%s\n", url, err.Error())
		time.Sleep(i * 10 * time.Second)

	}

	return
}

func (sess *Session) get(url string) ([]byte, error) {
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

	bodybuf := bufio.NewReaderSize(res.Body, 1024*1024)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(bodybuf)
	if err != nil {
		return nil, fmt.Errorf("Error downloading HTTP response: %w", err)
	}
	return body, nil
}

func (sess *Session) post(url string, body interface{}) ([]byte, error) {
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
