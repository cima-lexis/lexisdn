package webdrops

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cima-lexis/lexisdn/config"
)

// Session ...
type Session struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    uint64 `json:"expires_in"`
	ClientID     string
	RefreshedAt  time.Time
	client       *http.Client
}

// Login ...
func (sess *Session) Login() error {

	data := url.Values{}
	data.Set("client_id", config.Config.ClientID)
	data.Set("grant_type", "password")
	data.Set("password", config.Config.Password)
	data.Set("username", config.Config.User)

	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Minute,
			KeepAlive: 30 * time.Minute,
		}).Dial,
		TLSHandshakeTimeout: 6 * time.Minute,
	}
	c := &http.Client{
		Transport: t,
	}

	sess.client = c
	req, err := http.NewRequest("POST", config.Config.AuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := sess.client.Do(req)
	if err != nil {
		return fmt.Errorf("error submitting HTTP request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error downloading HTTP response: %w", err)
	}
	sess.ClientID = config.Config.ClientID
	err = json.Unmarshal(body, sess)
	if err != nil {
		return fmt.Errorf("error parsing HTTP JSON response: %w", err)
	}
	sess.RefreshedAt = time.Now()
	sess.ExpiresIn = 30
	return err
}

func (sess *Session) refresh() error {
	secondsPassedFromToken := uint64(math.Floor(time.Since(sess.RefreshedAt).Seconds()))
	//fmt.Println("passed", secondsPassedFromToken, "of", sess.ExpiresIn)
	if secondsPassedFromToken < sess.ExpiresIn/2 {
		return nil
	}

	data := url.Values{}
	data.Set("client_id", config.Config.ClientID)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", sess.RefreshToken)

	req, err := http.NewRequest("POST", config.Config.AuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := sess.client.Do(req)
	if err != nil {
		return fmt.Errorf("error submitting HTTP request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error downloading HTTP response: %w", err)
	}
	//fmt.Printf("Before: \nTk:%s\nRefTk:%s\n", sess.Token, sess.RefreshToken)
	ret := json.Unmarshal(body, sess)
	//fmt.Printf("After: \nTk:%s\nRefTk:%s\n", sess.Token, sess.RefreshToken)
	sess.RefreshedAt = time.Now()
	sess.ExpiresIn = 30

	return ret

}
