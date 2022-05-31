package iorm

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

type Engine struct {
	Client influxdb2.Client

	OrgMap map[string]*Org // currently using orgs, not all orgs in influxdb
}

func NewEngine(url, token string) *Engine {
	e := &Engine{
		OrgMap: make(map[string]*Org),
	}

	// Init client
	// default options:
	// &Options{batchSize: 5_000, flushInterval: 1_000, precision: time.Nanosecond, useGZip: false,
	// retryBufferLimit: 50_000, defaultTags: make(map[string]string),
	// maxRetries: 5, retryInterval: 5_000, maxRetryInterval: 125_000, maxRetryTime: 180_000, exponentialBase: 2}
	options := influxdb2.DefaultOptions().
		SetBatchSize(5000).
		SetFlushInterval(1000).
		SetUseGZip(true)
	e.Client = influxdb2.NewClientWithOptions(url, token, options)

	return e
}

func (e *Engine) AddOrg(name string) {
	org := &Org{
		Engine:    e,
		Name:      name,
		QueryAPI:  e.Client.QueryAPI(name),
		BucketMap: make(map[string]*Bucket),
	}
	e.OrgMap[name] = org
}

func (e *Engine) Close() {
	e.Client.Close()
	// Force all unwritten data to be sent
	for _, org := range e.OrgMap {
		for _, bucket := range org.BucketMap {
			bucket.WriteAPI.Flush()
		}
	}
}

type Org struct {
	Engine    *Engine
	Name      string
	QueryAPI  api.QueryAPI // Ensures using a single QueryAPI instance each org
	BucketMap map[string]*Bucket
}

func (o *Org) Query(flux string, dest interface{}) (map[string]interface{}, error) {
	result, err := o.QueryAPI.Query(context.Background(), flux)
	if err != nil {
		return nil, err
	}
	// Iterate over query response
	for result.Next() {
		// Notice when group key has changed
		// if result.TableChanged() {
		// 	log.Infof("table: %s\n", result.TableMetadata().String())
		// }
		// Access data
		// log.Infof("value: %v\n", result.Record().Values())
		// return result.Record() // data of one row
		return result.Record().Values(), nil
	}
	// check for an error
	if result.Err() != nil {
		// log.Infof("query parsing error: %s\n", result.Err().Error())
		return nil, err
	}

}

func (o *Org) AddBucket(name string) {
	bucket := &Bucket{
		Org:      o,
		Name:     name,
		WriteAPI: o.Engine.Client.WriteAPI(o.Name, name),
	}
	o.BucketMap[name] = bucket

	// Get errors channel
	errorsCh := bucket.WriteAPI.Errors()
	// Create go proc for reading and logging errors
	go func() {
		for err := range errorsCh {
			log.Errorf("influxdb write error: %s\n", err.Error())
		}
	}()
}

type Bucket struct {
	Org      *Org
	Name     string
	WriteAPI api.WriteAPI // Ensures using a single WriteAPI instance for each org/bucket pair.
}
