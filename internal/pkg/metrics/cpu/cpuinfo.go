package cpu

import (
	cpuPS "github.com/shirou/gopsutil/v3/cpu"
)

type Info struct {
	CPU        string   `json:"cpu"`
	VendorID   string   `json:"vendorID"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int32    `json:"stepping"`
	PhysicalID string   `json:"physicalId"`
	CoreID     string   `json:"coreId"`
	Cores      int32    `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int32    `json:"cacheSize"`
	Flags      []string `json:"flags"`
	Microcode  string   `json:"microcode"`
}

func CPUInfo() ([]*Info, error) {
	// var points []*write.Point
	var result []*Info

	cpuInfos, err := cpuPS.Info()
	if err != nil {
		return nil, err
	}

	//measurement := "cpu"
	for _, ci := range cpuInfos {
		// tags := map[string]string{
		// 	"cpu": strconv.FormatInt(int64(ci.CPU), 10),
		// }
		cpuInfo := &Info{
			CPU:        CPUName(int(ci.CPU)),
			VendorID:   ci.VendorID,
			Family:     ci.Family,
			Model:      ci.Model,
			Stepping:   ci.Stepping,
			PhysicalID: ci.PhysicalID,
			CoreID:     ci.CoreID,
			Cores:      ci.Cores,
			ModelName:  ci.ModelName,
			Mhz:        ci.Mhz,
			CacheSize:  ci.CacheSize,
			Flags:      ci.Flags,
			Microcode:  ci.Microcode,
		}
		//m := json.StructToMap(cpuInfo)
		// now := time.Now()
		// p := influxdb.NewPoint(measurement, tags, fields, now)
		// points = append(points, p)
		result = append(result, cpuInfo)
	}
	return result, nil
}
