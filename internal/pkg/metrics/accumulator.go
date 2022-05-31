package metrics

import (
	"sync"
	"time"

	"github.com/githubzjm/tuo/internal/pkg/influxdb"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/shirou/gopsutil/v3/host"
	log "github.com/sirupsen/logrus"
)

// Use of channel example:
// https://pkg.go.dev/github.com/influxdata/influxdb-client-go/v2#readme-concurrency
type Accumulator struct {
	PointCh chan *write.Point
	Uid     string
}

func NewAccumulator(pointChLen int) (*Accumulator, error) {
	pointCh := make(chan *write.Point, pointChLen) // buffered channel
	hostid, err := host.HostID()
	if err != nil {
		return nil, err
	}
	return &Accumulator{
		PointCh: pointCh,
		Uid:     hostid,
	}, nil
}

func (acc *Accumulator) Add(measurement string, tags map[string]string, fields map[string]interface{}, now time.Time) {
	p := influxdb.NewPoint(measurement, tags, fields, now)
	acc.PointCh <- p
}

func (acc *Accumulator) Close() {
	close(acc.PointCh)
}

func (acc *Accumulator) Report(threads int) {
	var wg sync.WaitGroup
	// Launch write routines
	for t := 0; t < threads; t++ {
		wg.Add(1)
		go func() {
			for p := range acc.PointCh { // will block until channel is closed
				//log.Info("write point")
				influxdb.WriteAPI.WritePoint(p)
			}
			wg.Done()
		}()
	}
	// Wait for writes complete
	wg.Wait()
	log.Info("report routines quit")
}
