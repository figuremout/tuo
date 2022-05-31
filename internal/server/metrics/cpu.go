package metrics

import (
	"context"

	"github.com/githubzjm/tuo/internal/pkg/influxdb"
)

func CPUPercent() (map[string]interface{}, error) {
	// get QueryTableResult
	result, err := influxdb.QueryAPI.Query(context.Background(), `
	from(bucket:"init-bucket")
		|> range(start:-5s, stop: now())
		|> filter(fn: (r) => r._measurement == "cpu" and r.cpu == "cpu-total" and r._field == "percent")
		|> keep(columns: ["cpu", "_value", "_time"])
		|> last()`)
	if err == nil {
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
	return nil, err
}
