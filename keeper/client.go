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

// httpCodeMessage returns http status detailed message
// for use by http client functions
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

// apiGetRequest send get request to Grafana API
//
func apiGetRequest(requestURL string) ([]byte, error) {

	resp, err := http.Get(requestURL)
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

// apiPostRequest send post request to Grafana API
//
func apiPostRequest(requestURL string, jsonData io.Reader) error {

	req, err := http.NewRequest("POST", requestURL, jsonData)
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

// apiDeleteRequest send delete request to Grafana API
//
func apiDeleteRequest(requestURL string) error {

	req, err := http.NewRequest("DELETE", requestURL, nil)
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
