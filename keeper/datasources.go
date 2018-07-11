//
// Datasources processing
//

package keeper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

// getAllDatasourcesList requests from Grafana json containing
// list of all atasources with a limited set of parameters
//
func getAllDatasourcesList(grafanaURL string) ([]grafanaDatasource, error) {

	grafanaRequestURL := grafanaURL + "/api/datasources"
	jsonData, err := apiGetRequest(grafanaRequestURL)
	if err != nil {
		return nil, err
	}

	var datasources []grafanaDatasource
	err = json.Unmarshal(jsonData, &datasources)
	if err != nil {
		return nil, err
	}

	return datasources, nil
}

func loadDatasourceFromFile(grafanaURL string, filePath string) error {

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	grafanaRequestURL := grafanaURL + "/api/datasources"
	err = apiPostRequest(grafanaRequestURL, jsonFile)
	if err != nil {
		return err
	}

	return nil
}

func saveDatasourceByID(grafanaURL string, workDir string, datasource grafanaDatasource) error {

	grafanaRequestURL := grafanaURL + "/api/datasources/" + strconv.Itoa(datasource.ID)
	jsonData, err := apiGetRequest(grafanaRequestURL)
	if err != nil {
		return err
	}

	jsonResult, err := prepareDatasourceJSON(jsonData)
	if err != nil {
		return err
	}

	fileName := datasource.Name + "-datasource.json"
	pathFileName := filepath.Join(workDir, fileName)
	err = writeJSONFile(pathFileName, jsonResult)
	if err != nil {
		return err
	}

	return nil
}

// Grafana API version 5.1 notes:
// field "readOnly" is returned different when get datasource by ID and get by Name
// field "typeLogoUrl" is returned empty but filled by get datasources list
//
func saveDatasourceByName(grafanaURL string, workDir string, datasource grafanaDatasource) error {

	grafanaRequestURL := grafanaURL + "/api/datasources/name/" + datasource.Name
	jsonData, err := apiGetRequest(grafanaRequestURL)
	if err != nil {
		return err
	}

	jsonResult, err := prepareDatasourceJSON(jsonData)
	if err != nil {
		return err
	}

	fileName := datasource.Name + "-datasource.json"
	pathFileName := filepath.Join(workDir, fileName)
	err = writeJSONFile(pathFileName, jsonResult)
	if err != nil {
		return err
	}

	return nil
}

func deleteDatasourceByID(grafanaURL string, datasourceID int) error {

	grafanaRequestURL := grafanaURL + "/api/datasources/" + strconv.Itoa(datasourceID)
	return apiDeleteRequest(grafanaRequestURL)
}

func getDatasourceCrc32ByID(grafanaURL string, datasource grafanaDatasource) (uint32, error) {

	grafanaRequestURL := grafanaURL + "/api/datasources/" + strconv.Itoa(datasource.ID)
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
