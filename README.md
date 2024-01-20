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

### Install MySQL
[MySQL](https://dev.mysql.com/doc/mysql-getting-started/en/)
MySQL Workbench
Execute the SQL within tech_test_device_interactions.sql to initialize the schema

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