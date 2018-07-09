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
    "flag"
    "log"
    "net/url"
    "os"
    "time"

    "grafana-keeper/keeper"
)

const startRetries = 10
const retryInterval = 3 * time.Second

func main() {

    // Parse command-line arguments
    // grafana-url, usually "http://localhost:3000"
    // work-dir, usually "/var/grafana-dashboards"
    //
    grafanaUrlPtr := flag.String("grafana-url", "", "Grafana server url")
    workDirPtr := flag.String("work-dir", "", "Directory to save grafana objects")
    saveFlagPtr := flag.String("save-script", "false", "Save-script mode")
    flag.Parse()
    if *grafanaUrlPtr == "" {
        log.Fatal("Missing parameter grafana-url")
    }
    if *workDirPtr == "" {
        log.Fatal("Missing parameter work-dir")
    }
    

    // Prepare Grafana's base url with authentication
    // If Grafana is configured for authentication, username and password
    // must be set in environment variables GRAFANA_USER and GRAFANA_PASSWORD.
    // Grafana's defaults is admin:admin
    //
    grafanaUrlObj, err := url.Parse(*grafanaUrlPtr)
    if err != nil {
        log.Fatalf("Grafana URL could not be parsed: %s", *grafanaUrlPtr)
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

    // Start the Grafana-keeper
    //
    log.Println("### Grafana-keeper started...")
    log.Printf("grafana-url: %s\n", *grafanaUrlPtr)
    log.Printf("work-dir: %s\n", *workDirPtr)

    // Init Grafana interface
    //
    Grafana := keeper.NewGrafana(grafanaUrl, *workDirPtr)

    // Save all datasources and dashboards to work directory
    //
    if *saveFlagPtr != "false" {
        log.Println("save-script mode on")
        err = Grafana.SaveAllDatasources()
        if err != nil {
            panic(err)
        }
        err = Grafana.SaveAllDashboards()
        if err != nil {
            panic(err)
        }

    } else {

      // Delete all datasources and dashboards
      //
      err = Grafana.DeleteAllDatasources()
      if err != nil {
          panic(err)
      }
      err = Grafana.DeleteAllDashboards()
      if err != nil {
          panic(err)
      }

      // Load datasources and dashboards from work directory
      //
      err = Grafana.LoadAllDatasources()
      if err != nil {
          panic(err)
      }
      err = Grafana.LoadAllDashboards()
      if err != nil {
          panic(err)
      }

    }


//    jsonFileName := *workDirPtr + "/prometheus-datasource.json"

//    log.Println("Creating datasource from:", jsonFileName)
//    for i := 0; i < startRetries; i++ {
//        err = keeper.CreateDatasourceFromFile(grafanaUrl, jsonFileName)
//        if err == nil { break }
//        log.Println("Retry on error:", err)
//        time.Sleep(retryInterval)
//    }

//    if err != nil {
//        panic(err)
//    }







    log.Println("### Grafana-keeper finished ok")

}



