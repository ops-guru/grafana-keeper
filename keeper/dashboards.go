//
// Dashboards processing
//

package keeper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// getAllDashboardsList requests from Grafana json containing
// list of all dashboards with a limited set of parameters
//
func getAllDashboardsList(grafanaURL string) ([]grafanaDashboard, error) {

	grafanaRequestURL := grafanaURL + "/api/search"
	jsonData, err := apiGetRequest(grafanaRequestURL)
	if err != nil {
		return nil, err
	}

	var dashboards []grafanaDashboard
	err = json.Unmarshal(jsonData, &dashboards)
	if err != nil {
		return nil, err
	}

	return dashboards, nil
}

func loadDashboardFromFile(grafanaURL string, filePath string) error {

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	grafanaRequestURL := grafanaURL + "/api/dashboards/db"
	err = apiPostRequest(grafanaRequestURL, jsonFile)
	if err != nil {
		return err
	}

	return nil
}

func saveDashboardByUID(grafanaURL string, workDir string, dashboard grafanaDashboard) error {

	grafanaRequestURL := grafanaURL + "/api/dashboards/uid/" + dashboard.UID
	jsonData, err := apiGetRequest(grafanaRequestURL)
	if err != nil {
		return err
	}

	jsonResult, err := prepareDashboardJSON(jsonData)
	if err != nil {
		return err
	}

	dbName := strings.TrimPrefix(dashboard.URI, "db/")
	fileName := dbName + "-dashboard.json"
	pathFileName := filepath.Join(workDir, fileName)
	err = writeJSONFile(pathFileName, jsonResult)
	if err != nil {
		return err
	}

	return nil
}

func deleteDashboardByUID(grafanaURL string, dashboardUID string) error {
	grafanaRequestURL := grafanaURL + "/api/dashboards/uid/" + dashboardUID
	return apiDeleteRequest(grafanaRequestURL)
}

func getDashboardCrc32ByUID(grafanaURL string, dashboard grafanaDashboard) (uint32, error) {

	grafanaRequestURL := grafanaURL + "/api/dashboards/uid/" + dashboard.UID
	jsonData, err := apiGetRequest(grafanaRequestURL)
	if err != nil {
		return 0, err
	}

	crc32, err := checksum32(jsonData)
	if err != nil {
		return 0, err
	}

	return crc32, nil
}
