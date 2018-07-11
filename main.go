//
// The Grafana keeper was built to run Grafana in an easily replicable manner
// without the need to run a complicated database.
//
// On start it deletes all datasources and dashboards in Grafana
// then it reads the set of objects from files matching *-datasource.json and
// *-dashboard.json from it's work directory and imports the datasources and
// dashboards to a given Grafana instance via Grafana's REST API.
//
// While running the Grafana keeper is checking Grafana objects
// (datasources and dashboards) for changes each 30 seconds.
// If any of this objects is changed or added a new one
// the Grafana keeper saves changes to it's work directory.
// On restart the set of objects will be automatically restored.
//
// Before first time run Grafana keeper may be You'll want to save
// already created in Grafana datasources and dashboards.
// Please run the Grafana keeper in save-script mode
// to store present objects to work directory.
// Then You can check for all objects are saved and
// run the Grafana keeper in usual mode.
//

package main

import (
	"log"

	"grafana-keeper/keeper"
)

func main() {

	log.Println("### Grafana-keeper started...")

	Grafana := keeper.Init()

	if Grafana.IsSaveScriptMode() {
		// Save-script mode
		// Save all Grafana's objects and exit

		keeper.SaveAllObjects(Grafana)

	} else {
		// Normal keeping Grafana's objects mode
		// On error log and retry

		// Delete all current datasources and dashboards
		// Load datasources and dashboards from work directory
		// Do it ones on start service
		//
		keeper.LoadObjectsFromWorkDir(Grafana)

		// Save new datasources and dashboards periodically
		// Repeat each retryInterval while Grafana-keeper is active
		//
		keeper.SaveNewObjectsPeriodically(Grafana)

	}

	log.Println("### Grafana-keeper finished ok")
}
