# journeys-api

Journeys API serves a subset of GTFS static data as JSON. This project was originally created for City of Tampere, but is now open sourced.
You can find more information about the original project in
the [ITS Factory Wiki](https://wiki.itsfactory.fi/index.php/Journeys_API), where the current closed source API is
discussed. This open source version of the API tries to match the proprietary version as much as possible.

You can participate freely in [discussions](https://github.com/jlundan/journeys-api/discussions) if you have questions about the API, create [issues](https://github.com/jlundan/journeys-api/issues), or join our [Discord Server](https://discord.gg/AvZvJxq8BM) if you want to discuss the API in less formal setting (we speak English and Finnish). 

## API
### Introduction

Journeys API follows a REST-styled design pattern. Each entity has an `url` field which allows the clients to access it. 
You can see this in action in Tampere GTFS data, where you would access stop point with short name 0001 using the URL:

```
https://data.itsfactory.fi/journeys/api/1/stop-points/0001
```
(In the Tampere environment, the `v1` part of the URL is replaced with `1` for historical reasons and backwards compatibility). 

The server responds:

```json
{
  "status": "success",
  "data": {
    "headers": {
      "paging": {
        "startIndex": 0,
        "pageSize": 1,
        "moreData": false
      }
    }
  },
  "body": [
    {
      "location": "61.49754,23.76152",
      "municipality": {
        "name": "Tampere",
        "shortName": "837",
        "url": "https://data.itsfactory.fi/journeys/api/1/municipalities/837"
      },
      "name": "Keskustori H",
      "shortName": "0001",
      "tariffZone": "A",
      "url": "https://data.itsfactory.fi/journeys/api/1/stop-points/0001"
    }
  ]
}
```
A collection of all stop points (and other entities, respectively) are accessed by omitting the entity id from the URL:

```
https://data.itsfactory.fi/journeys/api/1/stop-points
``` 

### Response Structure

Journeys API response is structured as follows:

```
{
    "status" : "success",
    "data" : {
        "headers" : {
            "paging" : {
                "startIndex" : 0,
                "pageSize" : 1,
                "moreData" : false
            }
        },
        "body": [
            ...
        ]
    }
}
```

The response has headers and body elements. Body elements contains the entities returned, and its content varies
depending on the request made by the client. Headers contain metadata-like information about the response. startIndex
tells the index of the first returned element, pageSize tells how many items were returned.

`moreData` property was used in the original implementation of the API, but is now kept for backwards compatibility. It will currently always return `false`.

### Entity Queries

Journeys API allows the client use optional URL parameters in queries. For example a client could query:

```
<base url>/v1/stop-points
```
which returns all stop points in the GTFS data.
```
<base url>/v1/stop-points?municpalityShortName=837&tariffZone=C
```

returns all stop points in Tampere within tariff zone C.

Journeys API allows (comma-separated) exclusion of returned fields. For example, a client could query:

```
<base url>/v1/stop-points?exclude-fields=municipality,url
```

Which would return all stop points without the municipality and url fields.

```
...
body": [
    {
        "location": "61.49754,23.76152",
        "name": "Keskustori H",
        "shortName": "0001",
        "tariffZone": "A"
    }
]
```

#### Query reference

The reference format is

```
[endpoint]
    - [parameter] : [description]
```

The endpoints and the query parameters are listed below. Tariff Zones depend on the GTFS data itself, the reference
lists the zones used in the Tampere GTFS data.

```
<base url>/v1/lines
	- description : string

<base url>/v1/routes
	- lineId : string
	- name : string

<base url>/v1/journey-patterns
	- lineId : string
	- name : string
	- firstStopPointId : string
	- lastStopPointId : string
	- stopPointId : string

<base url>/v1/journeys
	- lineId : string
	- routeId : string
	- journeyPatternId : string
	- dayTypes : comma separated list of: monday, tuesday, wednesday, friday, saturday, sunday
	- departureTime : hh:mm
	- arrivalTime : hh:mm
	- firstStopPointId : string
	- lastStopPointId : string
	- stopPointId : string
    - gtfsTripId: string

<base url>/v1/stop-points
	- name: string 
	- location: lat,lon or lat1,lon1:lat2,lon2 (upper left corner of a box : lower right corner of a box)
	- tariffZone : one of: A,B or C (https://www.nysse.fi/en/tickets-and-fares/zones.html)
	- municipalityName: string
	- municipalityShortName: string

<base url>/v1/municipalities
	- name: string
	- shortName: string

```

### Entities

Please note that the entity contents is based on the GTFS data. The field values might change with the GTFS data
used in the API, however the response structure and field names are accurate.

#### Lines
```json
{
    "status" : "success",
    "data" : {
        "headers" : {
            "paging" : {
                "startIndex" : 0,
                "pageSize" : 1,
                "moreData" : false
            }
        },
        "body": [
            {
                "href" : "<base url>/v1/lines/1",
                "name" : "1",
                "description" : "Kaupin kampus - Keskustori - Lentävänniemi"
            }
        ]
    }
}
```

#### Routes
```json
{
  "status": "success",
  "data": {
    "headers": {
      "paging": {
        "startIndex": 0,
        "pageSize": 1,
        "moreData": false
      }
    },
    "body": [
      {
        "url": "<base url>/v1/routes/288",
        "lineUrl": "<base url>/v1/lines/50",
        "name": "Lempäälä - Koskipuisto H",
        "journeyPatterns": [
          {
            "url": "<base url>/v1/journey-patterns/507",
            "originStop": "<base url>/v1/stop-points/7559",
            "destinationStop": "<base url>/v1/stop-points/0517",
            "name": "Lempäälä - Koskipuisto H"
          }
        ],
        "journeys": [
          {
            "url": "<base url>/v1/journeys/8726",
            "journeyPatternUrl": "<base url>/v1/journey-patterns/507",
            "departureTime": "07:10:00",
            "arrivalTime": "07:52:30",
            "dayTypes": [
              "friday"
            ],
            "dayTypeExceptions": [
              {
                "from": "2015-04-30",
                "to": "2015-04-30",
                "runs": "yes"
              }
            ]
          },
          {
            "url": "<base url>/v1/journeys/8728",
            "journeyPatternUrl": "<base url>/v1/journey-patterns/507",
            "departureTime": "07:10:00",
            "arrivalTime": "07:52:30",
            "dayTypes": [
              "monday",
              "tuesday",
              "wednesday",
              "thursday"
            ],
            "dayTypeExceptions": [
              {
                "from": "2015-05-14",
                "to": "2015-05-14",
                "runs": "no"
              }
            ]
          }
        ],
        "geographicCoordinateProjection": "6131429,2375268:-349,183:-259,177:-334,61:-225,-67 ..."
      }
    ]
  }
}
```
geographicCoordinateProjection contains information on how to draw the route on a map. This field is encoded to save bandwidth and the client must decode the fields value. The field takes form:
```
lat1,lon1:delta_lat2,delta_lon2:delta_lat3,delta_lon3 ...
```
A client would decode the field by subtracting `delta_lon2` from `lon1` and `delta_lat2` from `lat1` and dividing result with `10000`. This results in a coordinate pair from which `delta_lat3` and `delta_lon3` could be subtracted again and so on. First `lat1` and `lon1` should be just divided with `10000`.

Assuming the projection would start with `6131429,2375268:-349,183:-259,177:-334,61:-225,-67 ...`, the first coordinate pair can be obtained by dividing comma separated values with `10000`. Therefore, first coordinate pair would be:

```
lat1: 6131429 / 100000 = 61.31429
lon1: 2375268 / 100000 = 23.75268

==> lat1,lon1 = 61.31429,23.75268
```
Second coordinate pair would be obtained like this:
```
lat2: (6131429 - (-349)) = 6131778 => 6131778 / 100000 = 61.31778
lon2: (2375268 - 183) = 2375085 =>  2375085 / 100000 = 23.75085

==> lat2,lon2 = 61.31778,23.75085
```
Respectively, third coordinate pair would be obtained like this:
``` 
lat3: (6131778 - (-259)) = 6132037 => 6132037 / 100000 = 61.32037
lon3: (2375085 - 177) = 2374908 =>  2374908 / 100000 = 23.74908

==> lat3,lon3 = 61.32037,23.74908
```
And so on.

#### Journeys
```json
{
  "status": "success",
  "data": {
    "headers": {
      "paging": {
        "startIndex": 0,
        "pageSize": 1,
        "moreData": false
      }
    }
  },
  "body": [
    {
      "activityUrl": "<vehicle activity base url>/vehicle-activity/4_1340_1005_3596",
      "arrivalTime": "14:51:00",
      "calls": [
        {
          "arrivalTime": "14:51:00",
          "departureTime": "14:51:00",
          "stopPoint": {
            "location": "61.51211,23.68481",
            "municipality": {
              "name": "Tampere",
              "shortName": "837",
              "url": "<base url>/v1/municipalities/837"
            },
            "name": "Hiedanranta D",
            "shortName": "1005",
            "tariffZone": "B",
            "url": "<base url>/v1/stop-points/1005"
          }
        }
      ],
      "dayTypeExceptions": [],
      "dayTypes": [
        "sunday"
      ],
      "departureTime": "13:40:00",
      "directionId": "0",
      "gtfs": {
        "tripId": "77_15831_9189616"
      },
      "headSign": "Hiedanranta",
      "journeyPatternUrl": "<base url>/v1/journey-patterns/2212da15031a5cbf3a3c8ffecc59a00f",
      "lineUrl": "<base url>/v1/lines/4",
      "routeUrl": "<base url>/v1/routes/2318969642",
      "url": "<base url>/v1/journeys/77_15831_9189616",
      "wheelchairAccessible": true
    }
  ]
}
```
The activityUrl points to a service which hosts vehicle activity data. Currently, for Tampere, the vehicle activity is available at
```
https://data.itsfactory.fi/journeys/api/1/vehicle-activity
```
which means the activityUrl would be
```
https://data.itsfactory.fi/journeys/api/1
```
You can set this via the `JOURNEYS_VA_BASE_URL` environment variable. For example, in a Linux shell:
```
export JOURNEYS_VA_BASE_URL=https://data.itsfactory.fi/journeys/api/1
```
#### Journey Patterns
```json
{
  "status": "success",
  "data": {
    "headers": {
      "paging": {
        "startIndex": 0,
        "pageSize": 1,
        "moreData": false
      }
    }
  },
  "body": [
    {
      "destinationStop": "<base url>/v1/stop-points/3521",
      "direction": "1",
      "journeys": [
        {
          "arrivalTime": "04:48:00",
          "dayTypeExceptions": [
            {
              "from": "2025-02-06",
              "runs": "no",
              "to": "2025-02-06"
            }
          ],
          "dayTypes": [
            "monday",
            "tuesday",
            "wednesday",
            "thursday"
          ],
          "departureTime": "03:55:00",
          "headSign": "Hervanta",
          "journeyPatternUrl": "<base url>/v1/journey-patterns/b299a71359332e86a3b4c0c0cbfefa4b",
          "url": "<base url>/v1/journeys/77_15838_9308598"
        }
      ],
      "lineUrl": "<base url>/v1/lines/8B",
      "name": "Haukiluoma - Hervanta",
      "originStop": "<base url>/v1/stop-points/1668",
      "routeUrl": "<base url>/v1/routes/691625816",
      "stopPoints": [
        {
          "location": "61.51861,23.60668",
          "municipality": {
            "name": "Tampere",
            "shortName": "837",
            "url": "<base url>/v1/municipalities/837"
          },
          "name": "Haukiluoma",
          "shortName": "1668",
          "tariffZone": "B",
          "url": "<base url>/v1/stop-points/1668"
        },
        {
          "location": "61.4552,23.849",
          "municipality": {
            "name": "Tampere",
            "shortName": "837",
            "url": "<base url>/v1/municipalities/837"
          },
          "name": "Hervanta",
          "shortName": "3521",
          "tariffZone": "B",
          "url": "<base url>/v1/stop-points/3521"
        }
      ],
      "url": "<base url>/v1/journey-patterns/b299a71359332e86a3b4c0c0cbfefa4b"
    }
  ]
}
```
#### Stop Points
```json
{
  "status": "success",
  "data": {
    "headers": {
      "paging": {
        "startIndex": 0,
        "pageSize": 3446,
        "moreData": false
      }
    }
  },
  "body": [
    {
      "location": "61.49754,23.76152",
      "municipality": {
        "name": "Tampere",
        "shortName": "837",
        "url": "<base url>/v1/municipalities/837"
      },
      "name": "Keskustori H",
      "shortName": "0001",
      "tariffZone": "A",
      "url": "<base url>/v1/stop-points/0001"
    }
  ]
}
```
#### Municipalities
```json
{
  "status": "success",
  "data": {
    "headers": {
      "paging": {
        "startIndex": 0,
        "pageSize": 1,
        "moreData": false
      }
    }
  },
  "body": [
    {
      "name": "Tampere",
      "shortName": "837",
      "url": "<base url>/v1/municipalities/837"
    }
  ]
}
```
## About caching
Journeys API supports caching of responses. The cache has two modes: short cache and long cache. Short cache can be enabled for a specific time period (defined in hours), for example from 0 to 5. This can be used to fine tune the cache expiration close to the time when the current service day ends and a new service day begins. The service day change should ideally fall between the short cache period, and the short cache duration should be set so that it is shorter than the time between the stop of the last journey of the current service day and the start of the first journey on the next service day. This allows the cache to evict the previous day's journeys before the new journeys on the next service day start (with possibly on a different schedule), so that stale items from the previous day are no longer in the cache. Please see the "Environment variables" section for the environment variables that control the cache. 

## Running the server binary
After downloading the binary, run 
```bash
chmod +x ./journeys.api-linux-amd64
./journeys.api-linux-amd64 start
```

View help (and supported command line options):

```bash
./journeys.api-linux-amd64 help
```

## Environment variables

| argument                         | explanation                                                  |
|----------------------------------|--------------------------------------------------------------|
| JOURNEYS_GTFS_PATH               | path to directory where the GTFS files are located           |
| JOURNEYS_BASE_URL                | the base of the outputted URLs in responses                  |
| JOURNEYS_VA_BASE_URL             | the base of the outputted vehicle activity URLs in responses |
| JOURNEYS_PORT                    | the port where the service will run. defaults to 8080        |
| JOURNEYS_SHORT_CACHE_LOWER_BOUND | the lower hour of short cache period. defaults to 0          |
| JOURNEYS_SHORT_CACHE_UPPER_BOUND | the upper hour of short cache period. defaults to 5          |
| JOURNEYS_SHORT_CACHE_DURATION    | the short cache duration. defaults to 30 minutes             |
| JOURNEYS_LONG_CACHE_DURATION     | the long cache duration. defaults to 2 hours                 |


## Development Environment

After cloning the repository, download the dependencies:

```bash
go mod download
```

Start the server:

```bash
go run cmd/journeys/journeys.go start
```

View help (and supported command line options):

```bash
go run cmd/journeys/journeys.go help
```

You can update dependency versions (if needed) in standard Go way:

1. Update the version of the (dependency) package in the `go.mod` file
2. Run `go get -u <package-name>` to update the package

## The Tampere environment

The GTFS files for the Tampere region can be downloaded from [ITS Factory](https://data.itsfactory.fi/journeys/files/gtfs/). We currently use the Tampere GTFS files for
development, but you should be able to use other cities' GTFS files as well, assuming the GTFS data is similar.

There are a couple of Makefile targets to run the server with the Tampere environment:

```bash
# Run the server with the default development environment variables (localhost for internal URL links)
make tre-dev
# Run the server with the default development environment variables (data.itsfactory.fi for internal URL links)
make tre
```

Using the Makefile is not required, you can run the server with `go` command and set the environment variables manually.

## Endpoint compatibility with the proprietary Journeys API for the City of Tampere 

This repository provides a subset of functionality that is available in the proprietary Journeys API for City of Tampere. You should expect endpoint compatibility with the proprietary Journeys API according to the following table:

| Endpoint                            | Status           | Notes                   |
|-------------------------------------|------------------|-------------------------|
| /lines                              | Fully compatible | -                       |
| /routes                             | Fully compatible | -                       |
| /journey-patterns                   | Fully compatible | -                       |
| /journeys                           | Fully compatible | -                       |
| /stop-points                        | Fully compatible | -                       |
| /municipalities                     | Fully compatible | -                       |
| /lines/:lineId                      | Fully compatible | -                       |
| /routes/:routeId                    | Fully compatible | -                       |
| /journey-patterns/:journeyPatternId | Fully compatible | -                       |
| /journeys/:journeyId                | Fully compatible | -                       |
| /stop-points/:stopPointId           | Fully compatible | -                       |
| /municipalities/:municipalityId     | Fully compatible | -                       |
| /vehicle-activity                   | n/a              | Not supported currently |
| /stop-monitoring                    | n/a              | Not supported currently |
| /files/gtfs                         | n/a              | Not supported currently |

Please note:

* This repository might add the missing endpoints in the future
* You should not expect the response object properties to be in the same order as in the proprietary API
* You should not expect the response array items to be in the same order as in the proprietary API
* You should not expect the response page sizes to match the proprietary API. We aim to return all items in one page,
  however the paging information is maintained for backwards compatibility


