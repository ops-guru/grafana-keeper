//
// Dashboards processing
//

package keeper

import (
	"os"
	"path/filepath"
	"encoding/json"
	"strings"
)

// Get all dashboards list
//
func GetAllDashboardsList(grafanaUrl string) ([]GrafanaDashboard, error) {

	grafanaRequestUrl := grafanaUrl + "/api/search"
	jsonData, err := apiGetRequest(grafanaRequestUrl)
	if err != nil {
		return nil, err
	}

	var dashboards []GrafanaDashboard
	err = json.Unmarshal(jsonData, &dashboards)
	if err != nil {
		return nil, err
	}

	return dashboards, nil
}

// Load dashboard from file
//
func loadDashboardFromFile(grafanaUrl string, filePath string) error {

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	grafanaRequestUrl := grafanaUrl + "/api/dashboards/db"
	err = apiPostRequest(grafanaRequestUrl, jsonFile)
	if err != nil {
		return err
	}

	return nil
}

// Save dashboard by uid
//
func SaveDashboardByUid(grafanaUrl string, workDir string, dashboard GrafanaDashboard) error {

	grafanaRequestUrl := grafanaUrl + "/api/dashboards/uid/" + dashboard.Uid
	jsonData, err := apiGetRequest(grafanaRequestUrl)
	if err != nil {
		return err
	}

	jsonResult, err := PrepareDashboardJson(jsonData)
	if err != nil {
		return err
	}

	dbName := strings.TrimPrefix(dashboard.Uri, "db/")
	fileName := dbName + "-dashboard.json"
	pathFileName := filepath.Join(workDir, fileName)
	err = writeJsonFile(pathFileName, jsonResult)
	if err != nil {
		return err
	}

	return nil
}

// Delete dashboard by uid
//
func DeleteDashboardByUid(grafanaUrl string, dashboardUid string) error {
	grafanaRequestUrl := grafanaUrl + "/api/dashboards/uid/" + dashboardUid
	return apiDeleteRequest(grafanaRequestUrl)
}

// Get dashboard crc32 checksum by id
//
func GetDashboardCrc32ByUid(grafanaUrl string, dashboard GrafanaDashboard) (uint32, error) {

	grafanaRequestUrl := grafanaUrl + "/api/dashboards/uid/" + dashboard.Uid
	jsonData, err := apiGetRequest(grafanaRequestUrl)
	if err != nil {
		return 0, err
	}

	crc32, err := Checksum32(jsonData)
	if err != nil {
		return 0, err
	}

	return crc32, nil
}
