# go_api_demo
This is a simple REST Api using Go Lang for a demo

There is only a single post endpoint for this demo "/submit-iot-data"

It accepts a request of the following format:
```
type Submission struct {
	Timestamp time.Time `json:"timestamp"`
	Location  struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Devices []Device `json:"devices"`
}

type Device struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
```

### How to run the application
Navigate to the api directory of the project 
```cd api```
and run the following command:
```go run main.go```

### Prerequisites

### Install MySQL or Postgres
[MySQL](https://dev.mysql.com/doc/mysql-getting-started/en/)
[Postgres](https://www.postgresql.org/docs/current/tutorial-start.html)
MySQL Workbench, DBeaver, or PGAdmin are recommended for managing the database
Execute the SQL within *_tech_test_device_interactions.sql to initialize the schema
Configure BigQuery to connect to the database

### Install Go
[Go](https://go.dev/doc/install)
[Getting Started](https://go.dev/doc/tutorial/getting-started)

### Initialize Go Modules
go mod init go-demo 

### Install Go Dep
go get github.com/go-sql-driver/mysql

### Testing the API
[POST](http://localhost:8080/submit-iot-data)

Sample Request:

```{
  "timestamp": "2023-01-02T12:34:56Z",
  "location": {
    "latitude": 37.7749,
    "longitude": -122.4194
  },
  "devices": [
    {
      "id": 1,
      "name": "DeviceA"
    },
    {
      "id": 2,
      "name": "DeviceB"
    }
  ]
}
```

### Executing Unit Tests
go test

### Executing Individual Unit Tests
 go test -run ^UnitTestName$ 

 ### Anomaly Detetion

 There are a few things that can be done to detect anomalies in the data.  The following are some examples of queries that can be run to detect anomalies. Currently validation only takes into account timestamp anomalies and request structure. An argument can be made for not allowing updates to the Device Name as well for DeviceIds that have been seen before.  This would be a business decision. The expected behavior could be to update the device name to the most recent value, insert the new value as a new record, or reject the request.  Rejecting the request would be the most conservative approach.  
 
 
 The following queries can be used to detect anomalies:

 These queries can be run on a schedule to detect anomalies in the data.  The queries can be run on a schedule and the results can be sent to a monitoring system.

Unusually High Interaction Frequency (fixed threshold):
```
SELECT
  TIMESTAMP_TRUNC(timestamp, 'custom_time_interval') AS time_interval,
  COUNT(*) AS interaction_count
FROM
  device_interactions
WHERE
  timestamp BETWEEN 'start_timestamp' AND 'end_timestamp'
GROUP BY
  time_interval
HAVING
  interaction_count > 'custom_threshold';
```

Unusually Higher/Lower Interaction Frequency than Average for a Device:

```
SELECT
  device_id,
  AVG(COUNT(*)) AS average_interaction_count
FROM
  device_interactions
WHERE
  timestamp BETWEEN 'start_timestamp' AND 'end_timestamp'
GROUP BY
  device_id
HAVING
  average_interaction_count < 'custom_low_threshold' OR average_interaction_count > 'custom_high_threshold';
```

Geographical Outliers:

```
SELECT
  latitude,
  longitude,
  COUNT(*) AS interaction_count
FROM
  device_interactions
WHERE
  timestamp BETWEEN 'start_timestamp' AND 'end_timestamp'
GROUP BY
  latitude, longitude
HAVING
  distance(latitude, longitude, 'expected_latitude', 'expected_longitude') > 'custom_distance_threshold';

```

Recently Updated Device Name

```
SELECT
  device_id,
  device_name,
  COUNT(*) AS interaction_count
FROM
  device_interactions
WHERE
  timestamp BETWEEN 'start_timestamp' AND 'end_timestamp'
GROUP BY
  device_id, device_name
HAVING
  device_id NOT IN ('custom_device_id_list') OR device_name NOT IN ('custom_device_name_list');
```

Inactive Periods

```
SELECT
  TIMESTAMP_DIFF(MAX(timestamp), MIN(timestamp), SECOND) AS inactive_duration
FROM
  device_interactions
WHERE
  timestamp BETWEEN 'start_timestamp' AND 'end_timestamp'
HAVING
  inactive_duration > 'custom_inactive_threshold';
```