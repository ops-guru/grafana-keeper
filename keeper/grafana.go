//
// Grafana interface to access Grafana API
//
// Grafana API version 5.1 notes:
// field "readOnly" is returned different when get datasource by ID and get by Name
// field "typeLogoUrl" is returned empty but filled by get datasources list
//

package keeper

import (
	"log"
	"path/filepath"
)

type grafanaDatasource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type grafanaDashboard struct {
	ID    int    `json:"id"`
	UID   string `json:"uid"`
	Title string `json:"title"`
	URI   string `json:"uri"`
}

// GrafanaInterface to access Grafana API
//
type GrafanaInterface interface {
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

// Grafana is internal data of GrafanaInterface
//
type Grafana struct {
	BaseURL  string
	WorkDir  string
	SaveFlag bool
	DScrc32  map[int]uint32
	DBcrc32  map[string]uint32
}

// NewGrafana creates GrafanaInterface
//
func NewGrafana(baseURL string, workDir string, saveFlag bool) GrafanaInterface {
	return &Grafana{
		BaseURL:  baseURL,
		WorkDir:  workDir,
		SaveFlag: saveFlag,
		DScrc32:  make(map[int]uint32),
		DBcrc32:  make(map[string]uint32),
	}
}

// IsSaveScriptMode returns save-script mode status
//
func (grafana *Grafana) IsSaveScriptMode() bool {

	return grafana.SaveFlag
}

// DeleteAllDatasources deletes all Grafana's datasources
//
func (grafana *Grafana) DeleteAllDatasources() error {

	dsList, err := getAllDatasourcesList(grafana.BaseURL)
	if err != nil {
		return err
	}

	for _, ds := range dsList {
		log.Printf("Delete datasource: '%s'\n", ds.Name)
		err = deleteDatasourceByID(grafana.BaseURL, ds.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadAllDatasources loads datasources from work directory files
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
		err = loadDatasourceFromFile(grafana.BaseURL, f)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveAllDatasources saves all Grafana's datasources to work directory files
//
func (grafana *Grafana) SaveAllDatasources() error {

	dsList, err := getAllDatasourcesList(grafana.BaseURL)
	if err != nil {
		return err
	}

	for _, ds := range dsList {
		log.Printf("Save datasource: '%s'\n", ds.Name)
		err = saveDatasourceByID(grafana.BaseURL, grafana.WorkDir, ds)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveNewDatasources saves all new and changed
// Grafana's datasources to files in work directory
//
func (grafana *Grafana) SaveNewDatasources() error {

	dsList, err := getAllDatasourcesList(grafana.BaseURL)
	if err != nil {
		return err
	}

	m := grafana.DScrc32
	grafana.DScrc32 = make(map[int]uint32)
	for _, ds := range dsList {
		crc32, err := getDatasourceCrc32ByID(grafana.BaseURL, ds)
		if err != nil {
			return err
		}
		if crc32 == m[ds.ID] {
			grafana.DScrc32[ds.ID] = crc32
		} else {
			log.Printf("Save datasource: '%s'\n", ds.Name)
			err = saveDatasourceByID(grafana.BaseURL, grafana.WorkDir, ds)
			if err != nil {
				return err
			}
			grafana.DScrc32[ds.ID] = crc32
		}
	}

	return nil
}

// GetAllDatasourcesCrc32 get list of all datasources,
// request json data of each and calculate crc32 checksum
//
func (grafana *Grafana) GetAllDatasourcesCrc32() error {

	dsList, err := getAllDatasourcesList(grafana.BaseURL)
	if err != nil {
		return err
	}

	grafana.DScrc32 = make(map[int]uint32)
	for _, ds := range dsList {
		crc32, err := getDatasourceCrc32ByID(grafana.BaseURL, ds)
		if err != nil {
			return err
		}
		grafana.DScrc32[ds.ID] = crc32
	}

	return nil
}

// DeleteAllDashboards deletes all Grafana's dashboards
//
func (grafana *Grafana) DeleteAllDashboards() error {

	dbList, err := getAllDashboardsList(grafana.BaseURL)
	if err != nil {
		return err
	}

	for _, db := range dbList {
		log.Printf("Delete dashboard: '%s'\n", db.Title)
		err = deleteDashboardByUID(grafana.BaseURL, db.UID)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadAllDashboards loads dashboards from work directory files
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
		err = loadDashboardFromFile(grafana.BaseURL, f)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveAllDashboards saves all Grafana's dashboards to work directory files
//
func (grafana *Grafana) SaveAllDashboards() error {

	dbList, err := getAllDashboardsList(grafana.BaseURL)
	if err != nil {
		return err
	}

	for _, db := range dbList {
		log.Printf("Save dashboard: '%s'\n", db.Title)
		err = saveDashboardByUID(grafana.BaseURL, grafana.WorkDir, db)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveNewDashboards saves all new and changed
// Grafana's dashboards to files in work directory
//
func (grafana *Grafana) SaveNewDashboards() error {

	dbList, err := getAllDashboardsList(grafana.BaseURL)
	if err != nil {
		return err
	}

	m := grafana.DBcrc32
	grafana.DBcrc32 = make(map[string]uint32)
	for _, db := range dbList {
		crc32, err := getDashboardCrc32ByUID(grafana.BaseURL, db)
		if err != nil {
			return err
		}
		if crc32 == m[db.UID] {
			grafana.DBcrc32[db.UID] = crc32
		} else {
			log.Printf("Save dashboard: '%s'\n", db.Title)
			err = saveDashboardByUID(grafana.BaseURL, grafana.WorkDir, db)
			if err != nil {
				return err
			}
			grafana.DBcrc32[db.UID] = crc32
		}
	}

	return nil
}

// GetAllDashboardsCrc32 get list of all dashboards,
// request json data of each and calculate crc32 checksum
//
func (grafana *Grafana) GetAllDashboardsCrc32() error {

	dbList, err := getAllDashboardsList(grafana.BaseURL)
	if err != nil {
		return err
	}

	grafana.DBcrc32 = make(map[string]uint32)
	for _, db := range dbList {
		crc32, err := getDashboardCrc32ByUID(grafana.BaseURL, db)
		if err != nil {
			return err
		}
		grafana.DBcrc32[db.UID] = crc32
	}

	return nil
}
