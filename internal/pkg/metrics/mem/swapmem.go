package mem

import (
	"github.com/shirou/gopsutil/v3/mem"
)

type Swap struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
	// Sin         uint64  `json:"sin"`
	// Sout        uint64  `json:"sout"`
	// PgIn        uint64  `json:"pgIn"`
	// PgOut       uint64  `json:"pgOut"`
	// PgFault     uint64  `json:"pgFault"`

	// Linux specific numbers
	// https://www.kernel.org/doc/Documentation/cgroup-v2.txt
	// PgMajFault uint64 `json:"pgMajFault"`
}

func SwapMem() (*Swap, error) {
	var err error
	swapmem, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}
	swap := &Swap{
		Total:       swapmem.Total,
		Used:        swapmem.Used,
		Free:        swapmem.Free,
		UsedPercent: swapmem.UsedPercent,
		// Sin:         swapmem.Sin,
		// Sout:        swapmem.Sout,
		// PgIn:        swapmem.PgIn,
		// PgOut:       swapmem.PgOut,
		// PgFault:     swapmem.PgFault,
		// PgMajFault:  swapmem.PgMajFault,
	}
	return swap, nil
}
