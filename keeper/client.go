//
// Grafana API http client
//

package keeper

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Http status code message
//
func httpCodeMessage(resp *http.Response) error {

	strInfo := ""
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		strInfo = "[Invalid credentials] "
	case http.StatusConflict:
		strInfo = "[Creating duplicate object] "
	}

	return fmt.Errorf("%sStatus code returned from Grafana API (got: %d, expected: 200, msg:%s)\n", strInfo, resp.StatusCode, resp.Status)

}

// Grafana API http get request
//
func apiGetRequest(requestUrl string) ([]byte, error) {

	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, httpCodeMessage(resp)
	}

	// Get json string from http response
	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// Grafana API http post request
//
func apiPostRequest(requestUrl string, jsonData io.Reader) error {

	req, err := http.NewRequest("POST", requestUrl, jsonData)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return httpCodeMessage(resp)
	}

	return nil
}

// Grafana API http delete request
//
func apiDeleteRequest(requestUrl string) error {

	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return httpCodeMessage(resp)
	}

	return nil
}
