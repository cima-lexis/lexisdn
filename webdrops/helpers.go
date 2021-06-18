package webdrops

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const maxRetry = 5

// DoGet ...
func (sess *Session) DoGet(url string, expectedContentType string) (res []byte, err error) {
	//fmt.Println("GET", url)
	for i := time.Duration(0); i < maxRetry; i++ {
		err = sess.refresh()
		if err != nil {
			time.Sleep(i * 1 * time.Second)
			continue
		}

		res, err = sess.get(url, expectedContentType)
		if err == nil {
			return
		}
		fmt.Fprintf(os.Stderr, "An error occurred while getting from %s:%s\n", url, err.Error())
		time.Sleep(i * 1 * time.Second)

	}

	return
}

/*
// DoPost ...
func (sess *Session) DoPost(url string, body interface{}) (res []byte, err error) {
	//fmt.Println("POST", url)
	for i := time.Duration(0); i < maxRetry; i++ {
		err = sess.refresh()
		if err != nil {
			time.Sleep(i * 1 * time.Second)
			continue
		}

		res, err = sess.post(url, body)
		if err == nil {
			return
		}

		fmt.Fprintf(os.Stderr, "an error occurred while posting to %s:%s\n", url, err.Error())
		time.Sleep(i * 1 * time.Second)

	}

	return
}
*/
//"application/json"

func (sess *Session) get(url string, expectedContentType string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+sess.Token)
	//req.Header.Set("Accept-Encoding", "gzip, deflate")

	res, err := sess.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error submitting HTTP request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		//fmt.Println(string(body))

		return nil, fmt.Errorf("error in response: HTTP status: %s\nResponse Body:\n%s", res.Status, string(body))
	}
	encoding := res.Header.Get("Content-Type")
	if encoding != expectedContentType {
		return nil, fmt.Errorf("Response with status 200, but content type different than expected.\n expecting `%s`, got `%s`", expectedContentType, encoding)
	}

	bodybuf := bufio.NewReaderSize(res.Body, 10*1024)

	defer res.Body.Close()
	/*
		respWriter := bytes.NewBuffer([]byte{})
		bodyResp := bufio.NewWriterSize(respWriter, 10*1024)
		_, err = io.Copy(bodyResp, bodybuf)
		body := respWriter.Bytes()
	*/
	body, err := io.ReadAll(bodybuf)

	if err != nil {
		return nil, fmt.Errorf("error downloading HTTP response: %w", err)
	}
	return body, nil
}

/*
func (sess *Session) post(url string, body interface{}) ([]byte, error) {
	bodyJ, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error converting body to JSON: %w", err)
	}
	bodyR := bytes.NewBuffer(bodyJ)

	req, err := http.NewRequest("POST", url, bodyR)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+sess.Token)
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")

	res, err := sess.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error submitting HTTP request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error submitting request: HTTP status: %s", res.Status)
	}
	//fmt.Println("Content-Encoding: ", res.Header.Get("Content-Encoding"))

	bodybuf := bufio.NewReaderSize(res.Body, 10*1024)

	defer res.Body.Close()
	// respWriter := bytes.NewBuffer([]byte{})
	// bodyRespW := bufio.NewWriterSize(respWriter, 10*1024)
	// _, err = io.Copy(bodyRespW, bodybuf)
	// bodyResp := respWriter.Bytes()
	bodyResp, err := io.ReadAll(bodybuf)

	if err != nil {
		return nil, fmt.Errorf("error downloading: HTTP response: %w", err)
	}
	return bodyResp , nil
}
*/

// Domain is
type Domain struct {
	MinLat, MinLon, MaxLat, MaxLon float64
}

/*// ItalyDomain ...
var ItalyDomain = Domain{
	MaxLat: 66,
	MinLat: 23,
	MinLon: -18,
	MaxLon: 48,
}
*/
