# journeys-api

## Updating package versions
1. Update the version of the package in the go.mod file
2. Run `go get -u <package-name>` to update the package

## Running the server
Currently, Tampere GTFS config is available. To run the server with Tampere GTFS config, run the following command:
```bash
go run cmd/journeys/journeys.go
```

The command reads following environment variables:

| argument             | explanation                                                  |
|----------------------|--------------------------------------------------------------|
| JOURNEYS_GTFS_PATH   | path to directory where the GTFS files are located           |
| JOURNEYS_BASE_URL    | the base of the outputted URLs in responses                  |
| JOURNEYS_VA_BASE_URL | the base of the outputted vehicle activity URLs in responses |

There is also a Makefile command which uses :
```bash
make tre
```
