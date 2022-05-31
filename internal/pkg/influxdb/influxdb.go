package influxdb

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	log "github.com/sirupsen/logrus"
)

var Client influxdb2.Client
var WriteAPI api.WriteAPI
var QueryAPI api.QueryAPI

func InitDB(url, token, org, bucket string) {
	// Init client
	// default options:
	// &Options{batchSize: 5_000, flushInterval: 1_000, precision: time.Nanosecond, useGZip: false,
	// retryBufferLimit: 50_000, defaultTags: make(map[string]string),
	// maxRetries: 5, retryInterval: 5_000, maxRetryInterval: 125_000, maxRetryTime: 180_000, exponentialBase: 2}
	options := influxdb2.DefaultOptions().
		SetBatchSize(5000).
		SetFlushInterval(1000).
		SetUseGZip(true)
	Client = influxdb2.NewClientWithOptions(url, token, options)

	// Init WriteAPI
	// WriteAPI returns the asynchronous, non-blocking, Write client.
	// Ensures using a single WriteAPI instance for each org/bucket pair.
	WriteAPI = Client.WriteAPI(org, bucket)
	// Get errors channel
	errorsCh := WriteAPI.Errors()
	// Create go proc for reading and logging errors
	go func() {
		for err := range errorsCh {
			log.Errorf("influxdb write error: %s\n", err.Error())
		}
	}()

	// Init QueryAPI
	QueryAPI = Client.QueryAPI(org)
}

func Close() {
	// Force all unwritten data to be sent
	WriteAPI.Flush()
	Client.Close()
}

func Query(flux string) {
	// get QueryTableResult
	result, err := QueryAPI.Query(context.Background(), flux)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			// if result.TableChanged() {
			// 	log.Infof("table: %s\n", result.TableMetadata().String())
			// }
			// Access data
			log.Infof("value: %v\n", result.Record().Values())
			// return result.Record() // data of one row
		}
		// check for an error
		if result.Err() != nil {
			log.Infof("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
}

// Check detailed info about the server status, along with version string
func Health() (*domain.HealthCheck, error) {
	return Client.Health(context.Background())
}

// Check Server uptime info
func Ready() (*domain.Ready, error) {
	return Client.Ready(context.Background())
}

func NewPoint(measurement string, tags map[string]string,
	fields map[string]interface{}, ts time.Time) *write.Point {
	return influxdb2.NewPoint(measurement, tags, fields, ts)
}
