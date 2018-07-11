# grafana-keeper
Kubernetes application for persistence of Grafana dashboards to Kubernetes ConfigMaps

[![Go Report Card](https://goreportcard.com/badge/github.com/ops-guru/grafana-keeper)](https://goreportcard.com/report/github.com/ops-guru/grafana-keeper)

The Grafana-keeper was built to run Grafana in an easily replicable manner without the need to run a complicated database.

On start it deletes all datasources and dashboards in Grafana. Then it reads the set of objects from files matching *-datasource.json and
 *-dashboard.json in it's work directory and imports the datasources and dashboards to the serviced Grafana instance via Grafana's REST API.

While running the Grafana-keeper is checking Grafana's objects (datasources and dashboards) for changes each 30 seconds.
If any of this objects is changed or added a new one the Grafana-keeper saves changes to it's work directory.
On restart the set of objects will be automatically restored.

The Grafana-keeper can be run in save-script mode to store the current state of Grafana's objects as files in work directory.
It may be useful before first time run the Grafana-keeper because it begin with delete all.
Then You can check for all objects are saved properly and run the Grafana-keeper in usual mode.

## How to use
### Parameters
| Parameter | Typical | Description | Required |
| --------- | ------- | ----------- | -------- |
| --grafana-url | http://localhost:3000 | URL to connect to Grafana API| Required |
| --work-dir | /var/grafana-dashboards | Directory to save datasources and dashboards | Required |
| --save-script | false | save-script mode (save and exit) | Optional, default=false |

### Examples
Run as continuous service
```sh
grafana-keeper --grafana-url=http://localhost:3000 --work-dir=/var/grafana-dashboards
```
Run as save-script
```sh
grafana-keeper --grafana-url=http://localhost:3000 --work-dir=/var/grafana-dashboards --save-script=true
```
