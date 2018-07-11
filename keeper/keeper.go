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

// Initialize the Grafana keeper
// Parse command line parameters
// Prepare Grafana's base url with authentication
//
func Init() grafanaInterface {

	// Parse command-line arguments
	// grafana-url, usually "http://localhost:3000"
	// work-dir, usually "/var/grafana-dashboards"
	//
	grafanaUrlPtr := flag.String("grafana-url", "", "Grafana server url")
	workDirPtr := flag.String("work-dir", "", "Directory to save grafana objects")
	saveFlagPtr := flag.String("save-script", "false", "Save-script mode")
	flag.Parse()
	if *grafanaUrlPtr == "" {
		log.Fatalln("Missing parameter grafana-url")
	}
	if *workDirPtr == "" {
		log.Fatalln("Missing parameter work-dir")
	}
	saveFlag := *saveFlagPtr != "false"
	log.Printf("grafana-url: %s\n", *grafanaUrlPtr)
	log.Printf("work-dir: %s\n", *workDirPtr)
	if saveFlag {
		log.Println("save-script mode on")
	}

	// Prepare Grafana's base url with authentication
	// If Grafana is configured for authentication, username and password
	// must be set in environment variables GRAFANA_USER and GRAFANA_PASSWORD.
	// Grafana's defaults is admin:admin
	//
	grafanaUrlObj, err := url.Parse(*grafanaUrlPtr)
	if err != nil {
		log.Fatalf("Grafana URL could not be parsed: %s\n", *grafanaUrlPtr)
	}
	grafanaUser := os.Getenv("GRAFANA_USER")
	grafanaPass := os.Getenv("GRAFANA_PASSWORD")
	if grafanaUser != "" {
		if grafanaPass == "" {
			grafanaUrlObj.User = url.User(grafanaUser)
		} else {
			grafanaUrlObj.User = url.UserPassword(grafanaUser, grafanaPass)
		}
	}
	grafanaUrl := grafanaUrlObj.String()

	// Init Grafana interface
	//
	return NewGrafana(grafanaUrl, *workDirPtr, saveFlag)
}

// Save all datasources and dashboards to work directory
// Called on start when checksum lists are empty
// so all current objects will be saved
// Exit on error
//
func SaveAllObjects(Grafana grafanaInterface) {

	err := Grafana.SaveNewDatasources()
	if err != nil {
		log.Fatalln("Save datasources error:", err, "Grafana-keeper terminated")
	}
	err = Grafana.SaveNewDashboards()
	if err != nil {
		log.Fatalln("Save dashboards error:", err, "Grafana-keeper terminated")
	}
}

// Delete all datasources and dashboards in Grafana
// Then load objects from work directory
// Repeat on error with retryInterval until load all
// Finally save crc32 checksum of all objects
// Return after all operations will be finished
//
func LoadObjectsFromWorkDir(Grafana grafanaInterface) {

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

// Save new or changed datasources and dashboards
// Check each retryInterval while Grafana-keeper is active
// Use current objects's checksum saved on previous step
// to check if the object has been changed
// Renew checksum each time while checking objects
//
func SaveNewObjectsPeriodically(Grafana grafanaInterface) {

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
