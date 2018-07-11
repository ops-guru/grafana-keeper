//
// Grafana interface to access Grafana API
//
// Grafana API version 5.1 notes:
// field "readOnly" is returned different when get by Id and get by Name
// field "typeLogoUrl" is returned empty but filled by get datasources list
//

package keeper

import (
	"log"
	"path/filepath"
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
