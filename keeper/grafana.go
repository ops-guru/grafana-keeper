//
// Access grafana API
//
// Grafana API version 5.1 notes:
// field "readOnly" is returned different when get by Id and get by Name
// field "typeLogoUrl" is returned empty but filled by get datasources list
//

package keeper

import (
	"log"
	"os"
	"path/filepath"
	"encoding/json"
	"strings"
	"strconv"
)

// Grafana datasource object key fields
//
type GrafanaDatasource struct {
	Id		int		`json:"id"`
	Name	string	`json:"name"`
}

// Grafana dashboard object key fields
//
type GrafanaDashboard struct {
	Id		int		`json:"id"`
	Uid		string	`json:"uid"`
	Title	string	`json:"title"`
	Uri		string	`json:"uri"`
}

////////////////////////////// Grafana interface //////////////////////////////
//
// Grafana interface
//
type grafanaInterface interface {
	IsSaveScriptMode() bool
	DeleteAllDatasources() error
	LoadAllDatasources() error
	SaveNewDatasources() error
	GetAllDatasourcesCrc32() error
	DeleteAllDashboards() error
	LoadAllDashboards() error
	SaveNewDashboards() error
	GetAllDashboardsCrc32() error
}

// Grafana interface internal data
//
type Grafana struct {
	BaseUrl		string
	WorkDir		string
	SaveFlag	bool
	DScrc32		map[int]uint32
	DBcrc32		map[string]uint32
}

// Create Grafana interface
//
func NewGrafana(baseUrl string, workDir string, saveFlag bool) grafanaInterface {
	return &Grafana {
		BaseUrl:	baseUrl,
		WorkDir:	workDir,
		SaveFlag:	saveFlag,
		DScrc32:	make(map[int]uint32),
		DBcrc32:	make(map[string]uint32),
	}
}

// Get save script mode status
//
func (grafana *Grafana) IsSaveScriptMode() bool {

	return grafana.SaveFlag
}

// Delete all datasources
//
func (grafana *Grafana) DeleteAllDatasources() error {

	dsList, err := GetAllDatasourcesList(grafana.BaseUrl)
	if err != nil {
		return err
	}

	for _, ds := range dsList {
		log.Printf("Delete datasource: '%s'\n", ds.Name)
		err = DeleteDatasourceById(grafana.BaseUrl, ds.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

// Load all datasources from work directory
//
func (grafana *Grafana) LoadAllDatasources() error {

	// Get datasource matching files list
	//
	fileList, err := filepath.Glob(filepath.Join(grafana.WorkDir, "*-datasource.json"))
	if err != nil {
		return err
	}

	for _, f := range fileList {
		log.Printf("Create datasource from: '%s'\n", f)
		err = loadDatasourceFromFile(grafana.BaseUrl, f)
		if err != nil {
			return err
		}
	}

	return nil
}

// Save all datasources
//
func (grafana *Grafana) SaveAllDatasources() error {

	dsList, err := GetAllDatasourcesList(grafana.BaseUrl)
	if err != nil {
		return err
	}

	for _, ds := range dsList {
		log.Printf("Save datasource: '%s'\n", ds.Name)
		err = SaveDatasourceById(grafana.BaseUrl, grafana.WorkDir, ds)
		if err != nil {
			return err
		}
	}

	return nil
}

// Save new datasources
//
func (grafana *Grafana) SaveNewDatasources() error {

	dsList, err := GetAllDatasourcesList(grafana.BaseUrl)
	if err != nil {
		return err
	}

	m := grafana.DScrc32
	grafana.DScrc32 = make(map[int]uint32)
	for _, ds := range dsList {
		crc32, err := GetDatasourceCrc32ById(grafana.BaseUrl, ds)
		if err != nil {
			return err
		}
		if crc32 == m[ds.Id] {
			grafana.DScrc32[ds.Id] = crc32
		} else {
			log.Printf("Save datasource: '%s'\n", ds.Name)
			err = SaveDatasourceById(grafana.BaseUrl, grafana.WorkDir, ds)
			if err != nil {
				return err
			}
			grafana.DScrc32[ds.Id] = crc32
		}
	}

	return nil
}

// Get all datasources crc32 checksum
//
func (grafana *Grafana) GetAllDatasourcesCrc32() error {

	dsList, err := GetAllDatasourcesList(grafana.BaseUrl)
	if err != nil {
		return err
	}

	grafana.DScrc32 = make(map[int]uint32)
	for _, ds := range dsList {
		crc32, err := GetDatasourceCrc32ById(grafana.BaseUrl, ds)
		if err != nil {
			return err
		}
		grafana.DScrc32[ds.Id] = crc32
	}

	return nil
}

// Delete all dashboards
//
func (grafana *Grafana) DeleteAllDashboards() error {

	dbList, err := GetAllDashboardsList(grafana.BaseUrl)
	if err != nil {
		return err
	}

	for _, db := range dbList {
		log.Printf("Delete dashboard: '%s'\n", db.Title)
		err = DeleteDashboardByUid(grafana.BaseUrl, db.Uid)
		if err != nil {
			return err
		}
	}

	return nil
}

// Load all dashboards from work directory
//
func (grafana *Grafana) LoadAllDashboards() error {

	// Get dashboard matching files list
	//
	fileList, err := filepath.Glob(filepath.Join(grafana.WorkDir, "*-dashboard.json"))
	if err != nil {
		return err
	}

	for _, f := range fileList {
		log.Printf("Create dashboard from: '%s'\n", f)
		err = loadDashboardFromFile(grafana.BaseUrl, f)
		if err != nil {
			return err
		}
	}

	return nil
}

// Save all dashboards
//
func (grafana *Grafana) SaveAllDashboards() error {

	dbList, err := GetAllDashboardsList(grafana.BaseUrl)
	if err != nil {
		return err
	}

	for _, db := range dbList {
		log.Printf("Save dashboard: '%s'\n", db.Title)
		err = SaveDashboardByUid(grafana.BaseUrl, grafana.WorkDir, db)
		if err != nil {
			return err
		}
	}

	return nil
}

// Save new dashboards
//
func (grafana *Grafana) SaveNewDashboards() error {

	dbList, err := GetAllDashboardsList(grafana.BaseUrl)
	if err != nil {
		return err
	}

	m := grafana.DBcrc32
	grafana.DBcrc32 = make(map[string]uint32)
	for _, db := range dbList {
		crc32, err := GetDashboardCrc32ByUid(grafana.BaseUrl, db)
		if err != nil {
			return err
		}
		if crc32 == m[db.Uid] {
			grafana.DBcrc32[db.Uid] = crc32
		} else {
			log.Printf("Save dashboard: '%s'\n", db.Title)
			err = SaveDashboardByUid(grafana.BaseUrl, grafana.WorkDir, db)
			if err != nil {
				return err
			}
			grafana.DBcrc32[db.Uid] = crc32
		}
	}

	return nil
}

// Get all dashboards crc32 checksum
//
func (grafana *Grafana) GetAllDashboardsCrc32() error {

	dbList, err := GetAllDashboardsList(grafana.BaseUrl)
	if err != nil {
		return err
	}

	grafana.DBcrc32 = make(map[string]uint32)
	for _, db := range dbList {
		crc32, err := GetDashboardCrc32ByUid(grafana.BaseUrl, db)
		if err != nil {
			return err
		}
		grafana.DBcrc32[db.Uid] = crc32
	}

	return nil
}

///////////////////////////////// Datasources /////////////////////////////////
//
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

///////////////////////////////// Dashboards //////////////////////////////////
//
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

// Save dashboard by uid ////// Check
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
