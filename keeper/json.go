//
// Grafana API json utilities
//

package keeper

import (
    "os"
    "bytes"
    "encoding/json"
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
