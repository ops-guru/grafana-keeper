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

// Get all datasources list
//
func GetAllDatasourcesList(grafanaUrl string) ([]GrafanaDatasource, error) {

	grafanaRequestUrl := grafanaUrl + "/api/datasources"
	jsonData, err := apiGetRequest(grafanaRequestUrl)
	if err != nil {
		return nil, err
	}

	var datasources []GrafanaDatasource
	err = json.Unmarshal(jsonData, &datasources)
	if err != nil {
		return nil, err
	}

	return datasources, nil
}

// Load datasource from file
//
func loadDatasourceFromFile(grafanaUrl string, filePath string) error {

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	grafanaRequestUrl := grafanaUrl + "/api/datasources"
	err = apiPostRequest(grafanaRequestUrl, jsonFile)
	if err != nil {
		return err
	}

	return nil
}

// Save datasource by id
//
func SaveDatasourceById(grafanaUrl string, workDir string, datasource GrafanaDatasource) error {

	grafanaRequestUrl := grafanaUrl + "/api/datasources/" + strconv.Itoa(datasource.Id)
	jsonData, err := apiGetRequest(grafanaRequestUrl)
	if err != nil {
		return err
	}

	jsonResult, err := PrepareDatasourceJson(jsonData)
	if err != nil {
		return err
	}

	fileName := datasource.Name + "-datasource.json"
	pathFileName := filepath.Join(workDir, fileName)
	err = writeJsonFile(pathFileName, jsonResult)
	if err != nil {
		return err
	}

	return nil
}

// Save datasource by name
//
func SaveDatasourceByName(grafanaUrl string, workDir string, datasource GrafanaDatasource) error {

	grafanaRequestUrl := grafanaUrl + "/api/datasources/name/" + datasource.Name
	jsonData, err := apiGetRequest(grafanaRequestUrl)
	if err != nil {
		return err
	}

	jsonResult, err := PrepareDatasourceJson(jsonData)
	if err != nil {
		return err
	}

	fileName := datasource.Name + "-datasource.json"
	pathFileName := filepath.Join(workDir, fileName)
	err = writeJsonFile(pathFileName, jsonResult)
	if err != nil {
		return err
	}

	return nil
}

// Delete datasource by id
//
func DeleteDatasourceById(grafanaUrl string, datasourceId int) error {

	grafanaRequestUrl := grafanaUrl + "/api/datasources/" + strconv.Itoa(datasourceId)
	return apiDeleteRequest(grafanaRequestUrl)
}

// Get datasource crc32 checksum by id
//
func GetDatasourceCrc32ById(grafanaUrl string, datasource GrafanaDatasource) (uint32, error) {

	grafanaRequestUrl := grafanaUrl + "/api/datasources/" + strconv.Itoa(datasource.Id)
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
