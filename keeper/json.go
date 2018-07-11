//
// Grafana API json utilities
//

package keeper

import (
	"bytes"
	"encoding/json"
	"hash/crc32"
	"os"
)

// Prepare json for create datasource
// field 'id' must be deleted for create
// datasource by Grafana API properly
//
func PrepareDatasourceJson(jsonData []byte) ([]byte, error) {

	var jsonInterface interface{}
	err := json.Unmarshal(jsonData, &jsonInterface)
	if err != nil {
		return nil, err
	}

	mapData := jsonInterface.(map[string]interface{})
	delete(mapData, "id")
	jsonResult, err := json.Marshal(mapData)
	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}

// Prepare json for create dashboard
// in "dashboard" section top level fields
// 'id' and 'uid' must be set to null for
// create dashboard by Grafana API properly

//
func PrepareDashboardJson(jsonData []byte) ([]byte, error) {

	var jsonInterface interface{}
	err := json.Unmarshal(jsonData, &jsonInterface)
	if err != nil {
		return nil, err
	}

	mapData := jsonInterface.(map[string]interface{})
	mapData["dashboard"].(map[string]interface{})["id"] = nil
	mapData["dashboard"].(map[string]interface{})["uid"] = nil

	jsonResult, err := json.Marshal(mapData)
	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}

// (Re)Write json file
//
func writeJsonFile(jsonFileName string, jsonData []byte) error {

	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var buf bytes.Buffer
	err = json.Indent(&buf, jsonData, "", "\t")
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Crc32 checksum of json
// Grafana time to time return json data in different order
// To fix records order for stable 32-bit checksum
// data must be sorted equally by Unmarshal/Marshal
//
func Checksum32(jsonData []byte) (uint32, error) {

	var jsonInterface interface{}
	err := json.Unmarshal(jsonData, &jsonInterface)
	if err != nil {
		return 0, err
	}

	jsonResult, err := json.Marshal(jsonInterface)
	if err != nil {
		return 0, err
	}

	crc32q := crc32.MakeTable(0xD5828281)
	return crc32.Checksum(jsonResult, crc32q), nil
}
