package cpu

import (
	"strconv"
	"time"

	"github.com/githubzjm/tuo/internal/pkg/metrics"

	"github.com/githubzjm/tuo/internal/pkg/utils/json"
)

type CPUStats struct {
	PerCPU         bool
	TotalCPU       bool
	CollectCPUTime bool
	ReportActive   bool
}

// type Collector interface {
// 	CPUInfo() ([]map[string]interface{}, error)
// 	CPUPercent(interval time.Duration, percpu, totalcpu bool) ([]map[string]interface{}, error)
// 	CPUTimes(percpu, totalcpu bool) ([]map[string]interface{}, error)
// }

func NewCPUStats() *CPUStats {
	return &CPUStats{
		PerCPU:         true,
		TotalCPU:       true,
		CollectCPUTime: true,
		ReportActive:   true,
	}
}

func (c *CPUStats) Gather(acc *metrics.Accumulator) error {
	var err error

	measurement := "cpu"
	now := time.Now()

	// CPU Info
	var infos []*Info
	infos, err = CPUInfo()
	if err != nil {
		return err
	}
	for _, info := range infos {
		acc.Add(measurement, map[string]string{
			"cpu":  info.CPU,
			"host": acc.Uid,
		}, json.StructToMap(info), now)
	}

	// CPU Percent
	var percents []*Percent
	percents, err = CPUPercent(0, c.PerCPU, c.TotalCPU) // TODO interval
	if err != nil {
		return err
	}
	for _, percent := range percents {
		acc.Add(measurement, map[string]string{
			"cpu":  percent.CPU,
			"host": acc.Uid,
		}, json.StructToMap(percent), now)
	}

	// CPU Times
	var times []*Times
	times, err = CPUTimes(c.PerCPU, c.TotalCPU)
	if err != nil {
		return err
	}
	for _, time := range times {
		acc.Add(measurement, map[string]string{
			"cpu":  time.CPU,
			"host": acc.Uid,
		}, json.StructToMap(time), now)
	}

	return nil
}

func CPUName(n int) string {
	if n < 0 {
		return "cpu-total"
	} else {
		return "cpu" + strconv.FormatInt(int64(n), 10)
	}
}

func init() {
	metrics.Add("cpu", func() metrics.Metric {
		return NewCPUStats()
	})
}
