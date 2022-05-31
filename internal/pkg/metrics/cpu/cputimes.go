package cpu

import (
	"fmt"

	cpuPS "github.com/shirou/gopsutil/v3/cpu"
)

type Times struct {
	CPU       string  `json:"cpu"`
	User      float64 `json:"time_user"`
	System    float64 `json:"time_system"`
	Idle      float64 `json:"time_idle"`
	Nice      float64 `json:"time_nice"`
	Iowait    float64 `json:"time_iowait"`
	Irq       float64 `json:"time_irq"`
	Softirq   float64 `json:"time_softirq"`
	Steal     float64 `json:"time_steal"`
	Guest     float64 `json:"time_guest"`
	GuestNice float64 `json:"time_guestNice"`
}

func CPUTimes(percpu, totalcpu bool) ([]*Times, error) {
	// var points []*write.Point
	var result []*Times
	var err error
	var cpuTimesStats, totalTimesStats, perTimesStats []cpuPS.TimesStat

	// Collect total cpu times
	if totalcpu {
		totalTimesStats, err = cpuPS.Times(false)
		if err != nil {
			return nil, fmt.Errorf("collect CPUTimes error: %s", err)
		}
		cpuTimesStats = append(cpuTimesStats, totalTimesStats...)

	}

	// Collect per cpu times
	if percpu {
		perTimesStats, err = cpuPS.Times(true)
		if err != nil {
			return nil, err
		}
		cpuTimesStats = append(cpuTimesStats, perTimesStats...)
	}

	// measurement := "cpu"
	for _, cts := range cpuTimesStats {
		// tags := map[string]string{
		// 	"cpu": cts.CPU,
		// }

		cpuTimes := &Times{
			CPU:       cts.CPU,
			User:      cts.User,
			System:    cts.System,
			Idle:      cts.Idle,
			Nice:      cts.Nice,
			Iowait:    cts.Iowait,
			Irq:       cts.Irq,
			Softirq:   cts.Softirq,
			Steal:     cts.Steal,
			Guest:     cts.Guest,
			GuestNice: cts.GuestNice,
		}
		//m := json.StructToMap(cpuTimes)
		result = append(result, cpuTimes)

		// now := time.Now()
		// p := influxdb.NewPoint(measurement, tags, fields, now)
		// points = append(points, p)
	}
	return result, nil
}
