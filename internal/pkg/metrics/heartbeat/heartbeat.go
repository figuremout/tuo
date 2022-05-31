package heartbeat

import (
	"time"

	"github.com/githubzjm/tuo/internal/pkg/metrics"
)

type HeartBeat struct {
}

func NewHeartBeat() *HeartBeat {
	return &HeartBeat{}
}

func (h *HeartBeat) Gather(acc *metrics.Accumulator) error {
	measurement := "heartbeat"
	now := time.Now()
	acc.Add(measurement, map[string]string{}, map[string]interface{}{
		"alive": true,
	}, now)
	return nil
}

func init() {
	metrics.Add("heartbeat", func() metrics.Metric {
		return NewHeartBeat()
	})
}
