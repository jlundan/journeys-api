# journeys-api

Journeys API serves GTFS data as JSON. This project was originally created for City of Tampere, but is now open sourced.
You can find more information about the project in the [ITS Factory Wiki](https://wiki.itsfactory.fi/index.php/Journeys_API), where the current closed source API is discussed. This open source version of the API tries to match the proprietary version as much as possible.

<i>You should consider this application as alpha-state for now. It is not yet ready for production use. You should be able to test it locally though.</i>

## Running the server
First, download the dependencies:
```bash
go mod download
```
Start the server:
```bash
go run cmd/journeys/journeys.go
```

The command reads following environment variables:

| argument             | explanation                                                  |
|----------------------|--------------------------------------------------------------|
| JOURNEYS_GTFS_PATH   | path to directory where the GTFS files are located           |
| JOURNEYS_BASE_URL    | the base of the outputted URLs in responses                  |
| JOURNEYS_VA_BASE_URL | the base of the outputted vehicle activity URLs in responses |

The GTFS files for the Tampere region can be downloaded from [ITS Factory](https://data.itsfactory.fi/journeys/files/gtfs/). We currently use the Tampere GTFS files for development, but you should be able to use other cities' GTFS files as well. 

There is also a Makefile command which uses defaults for the environment variables:
```bash
# Run the server with the default development environment variables (localhost for internal URL links)
make tre-dev
# Run the server with the default development environment variables (data.itsfactory.fi for internal URL links)
make tre
```
Using the Makefile is not required, you can run the server with `go` command and set the environment variables manually.

## Updating module versions
1. Update the version of the package in the `go.mod` file
2. Run `go get -u <package-name>` to update the package
