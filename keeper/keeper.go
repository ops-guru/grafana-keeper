//
// Grafana-keeper main logic
//

package keeper

import (
	"flag"
	"log"
	"net/url"
	"os"
	"time"
)

const retryInterval = 30 * time.Second

// Init should be called first to arrange the Grafana-keeper environment
// It parse and check command line arguments, set running mode,
// prepare Grafana's base url with authentication
//
func Init() GrafanaInterface {

	// Parse command-line arguments
	// grafana-url, usually "http://localhost:3000"
	// work-dir, usually "/var/grafana-dashboards"
	//
	grafanaURLPtr := flag.String("grafana-url", "", "Grafana server url")
	workDirPtr := flag.String("work-dir", "", "Directory to save grafana objects")
	saveFlagPtr := flag.String("save-script", "false", "Save-script mode")
	flag.Parse()
	if *grafanaURLPtr == "" {
		log.Fatalln("Missing parameter grafana-url")
	}
	if *workDirPtr == "" {
		log.Fatalln("Missing parameter work-dir")
	}
	saveFlag := *saveFlagPtr != "false"
	log.Printf("grafana-url: %s\n", *grafanaURLPtr)
	log.Printf("work-dir: %s\n", *workDirPtr)
	if saveFlag {
		log.Println("save-script mode on")
	}

	// Prepare Grafana's base url with authentication
	// If Grafana is configured for authentication, username and password
	// must be set in environment variables GRAFANA_USER and GRAFANA_PASSWORD.
	// Grafana's defaults is admin:admin
	//
	grafanaURLObj, err := url.Parse(*grafanaURLPtr)
	if err != nil {
		log.Fatalf("Grafana URL could not be parsed: %s\n", *grafanaURLPtr)
	}
	grafanaUser := os.Getenv("GRAFANA_USER")
	grafanaPass := os.Getenv("GRAFANA_PASSWORD")
	if grafanaUser != "" {
		if grafanaPass == "" {
			grafanaURLObj.User = url.User(grafanaUser)
		} else {
			grafanaURLObj.User = url.UserPassword(grafanaUser, grafanaPass)
		}
	}
	grafanaURL := grafanaURLObj.String()

	// Init Grafana interface
	//
	return NewGrafana(grafanaURL, *workDirPtr, saveFlag)
}

// SaveAllObjects is what Grafana-keeper do in save-script mode
// It saves all datasources and dashboards to work directory
// Call on start when checksum lists are empty
// for all current objects to be saved
// Function terminates main process on error
//
func SaveAllObjects(Grafana GrafanaInterface) {

	err := Grafana.SaveNewDatasources()
	if err != nil {
		log.Fatalln("Save datasources error:", err, "Grafana-keeper terminated")
	}
	err = Grafana.SaveNewDashboards()
	if err != nil {
		log.Fatalln("Save dashboards error:", err, "Grafana-keeper terminated")
	}
}

// LoadObjectsFromWorkDir is first stage when Grafana-keeper
// is in Normal keeping Grafana's objects mode
// Function deletes all datasources and dashboards in Grafana
// Then it loads objects from work directory
// Repeat on error with retryInterval until load all
// Finally save crc32 checksum of all objects
// Return after all operations will be finished
//
func LoadObjectsFromWorkDir(Grafana GrafanaInterface) {

	isStarting := true
	for {
		if isStarting {
			isStarting = false
		} else {
			time.Sleep(retryInterval)
		}

		// Delete all datasources and dashboards
		//
		err := Grafana.DeleteAllDatasources()
		if err != nil {
			log.Println("Delete datasources error:", err)
			continue
		}
		err = Grafana.DeleteAllDashboards()
		if err != nil {
			log.Println("Delete dashboards error:", err)
			continue
		}

		// Load datasources and dashboards from work directory
		//
		err = Grafana.LoadAllDatasources()
		if err != nil {
			log.Println("Load datasources error:", err)
			continue
		}
		err = Grafana.LoadAllDashboards()
		if err != nil {
			log.Println("Load dashboards error:", err)
			continue
		}

		// Get all datasources and dashboards crc32 checksum
		//
		err = Grafana.GetAllDatasourcesCrc32()
		if err != nil {
			log.Println("Get datasources crc32 error:", err)
			continue
		}
		err = Grafana.GetAllDashboardsCrc32()
		if err != nil {
			log.Println("Get dashboards crc32 error:", err)
			continue
		}
		break
	}
}

// SaveNewObjectsPeriodically repeat each retryInterval:
// compare current Grafana objects's checksum with saved
// on previous step to check if the object has been changed,
// save all new and changed datasources and dashboards,
// renew checksum each time while checking objects,
// continue the loop while Grafana-keeper is active
//
func SaveNewObjectsPeriodically(Grafana GrafanaInterface) {

	isStarting := true
	for {
		if isStarting {
			isStarting = false
		} else {
			time.Sleep(retryInterval)
		}

		// Save new datasources and dashboards
		//
		err := Grafana.SaveNewDatasources()
		if err != nil {
			log.Println("Save datasources error:", err)
		}
		err = Grafana.SaveNewDashboards()
		if err != nil {
			log.Println("Save dashboards error:", err)
		}
	}
}
