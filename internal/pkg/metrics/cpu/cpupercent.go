package cpu

import (
	"time"

	cpuPS "github.com/shirou/gopsutil/v3/cpu"
)

type Percent struct {
	CPU     string  `json:"cpu"`
	Percent float64 `json:"percent"`
}

func CPUPercent(interval time.Duration, percpu, totalcpu bool) ([]*Percent, error) {
	// var points []*write.Point
	var err error
	//measurement := "cpu"

	var result []*Percent

	// Collect total
	if totalcpu {
		var total []float64
		total, err = cpuPS.Percent(interval, false)
		if err != nil {
			return nil, err
		}
		p := &Percent{
			CPU:     CPUName(-1),
			Percent: total[0],
		}
		// m := json.StructToMap(s)
		result = append(result, p)
	}

	// COllect per
	if percpu {
		var pers []float64
		pers, err = cpuPS.Percent(interval, true)
		if err != nil {
			return nil, err
		}
		for i, p := range pers {
			p := &Percent{
				CPU:     CPUName(i),
				Percent: p,
			}
			// m := json.StructToMap(s)
			result = append(result, p)
		}

	}

	// for i, p := range percents {
	// 	// var id string
	// 	// if i == 0 {
	// 	// 	id = "cpu-total"
	// 	// } else {
	// 	// 	id = "cpu-" + strconv.Itoa(i)
	// 	// }
	// 	// tags := map[string]string{
	// 	// 	"cpu": id,
	// 	// }
	// 	cpuPercent := &CPUPercent{
	// 		Percent: p,
	// 	}
	// 	fields := json.StructToMap(cpuPercent)
	// 	now := time.Now()
	// 	p := influxdb.NewPoint(measurement, tags, fields, now)
	// 	points = append(points, p)
	// }
	return result, nil

}
