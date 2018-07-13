# grafana-keeper
Kubernetes application for persistence of Grafana dashboards to Kubernetes ConfigMaps

[![Go Report Card](https://goreportcard.com/badge/github.com/ops-guru/grafana-keeper)](https://goreportcard.com/report/github.com/ops-guru/grafana-keeper)

The Grafana-keeper was built to run Grafana in an easily replicable manner without the need to run a complicated database.

On start it deletes all datasources and dashboards in Grafana. Then it reads the set of objects from files matching *-datasource.json and
 *-dashboard.json in it's work directory and imports the datasources and dashboards to the serviced Grafana instance via Grafana's REST API.

While running the Grafana-keeper is checking Grafana's objects (datasources and dashboards) for changes each 30 seconds.
If any of this objects is changed or added a new one the Grafana-keeper saves changes to it's work directory.
On restart the set of objects will be automatically restored.
Please do not forget to delete corresponding files after delete or rename Grafana's objects.

The Grafana-keeper can be run in save-script mode to store the current state of Grafana's objects as files in work directory.
It may be useful before first time run the Grafana-keeper because it begin with delete all.
Then You can check for all objects are saved properly and run the Grafana-keeper in usual mode.

## How to use
### Parameters
| Parameter | Typical | Description | Required |
| --------- | ------- | ----------- | -------- |
| --grafana-url | http://localhost:3000 | URL to connect to Grafana API| Required |
| --work-dir | /var/grafana-objects | Directory to save datasources and dashboards | Required |
| --save-script | false | save-script mode (save and exit) | Optional, default=false |

### Environment variables
Grafaha-keeper must have admin access to Grafana's datasources and dashboards.
If Grafana is configured for authentication, username and password for admin user must be
exported in environment variables GRAFANA_USER and GRAFANA_PASSWORD correspondingly.
This could be done by append home/user/.profile file with lines:
```
export GRAFANA_USER=grafana-admin-user-name
export GRAFANA_PASSWORD=grafana-admin-user-password
```
Grafana's default values are admin : admin

### Building

**Prerequisites**

- golang environment
- docker (used for creating container images, etc.)
- kubernetes (optional)

**Build binary**

From project directory run:
```
make
```

**Build docker container**

From project directory run:
```
make image
```

### Running examples
Examples assume grafana-url and work-dir set to defaults.

**Run in save-script mode**

From the directory with built grafana-keeper binary run:
```sh
grafana-keeper/grafana-keeper --grafana-url=http://localhost:3000 --work-dir=/var/grafana-objects --save-script=true
```

**Run as standalone continuous service**

From the directory with built grafana-keeper binary run:
```sh
grafana-keeper/grafana-keeper --grafana-url=http://localhost:3000 --work-dir=/var/grafana-objects
```

**Run with docker**

Run in detached mode:
```
docker run --rm -d -v /var/grafana-objects:/var/grafana-objects:rw --net="host" -e GRAFANA_USER -e GRAFANA_PASSWORD opsguru.io/grafana-keeper:0.0.1-a91597e
```

Or run attached to terminal to see log messages:
```
docker run --rm -i -t -v /var/grafana-objects:/var/grafana-objects:rw --net="host" -e GRAFANA_USER -e GRAFANA_PASSWORD opsguru.io/grafana-keeper:0.0.1-a91597e
```

Run with different parameters:
```
docker run --rm -i -t -v /var/grafana-objects:/var/grafana-objects:rw --net="host" -e GRAFANA_USER -e GRAFANA_PASSWORD opsguru.io/grafana-keeper:0.0.1-a91597e  /grafana-keeper --grafana-url=http://localhost:3000 --work-dir=/var/grafana-objects
```

First '/var/grafana-objects' parameter is the path for storing Grafana objects on host computer.
This folder must be created before run the container. It's access rights must allow read and write for docker user.
The 'opsguru.io/grafana-keeper:0.0.1-a91597e' parameter is docker image name. It may vary depending on current grafana-keeper build.
