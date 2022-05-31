package mem

import (
	"time"

	"github.com/githubzjm/tuo/internal/pkg/metrics"
	"github.com/githubzjm/tuo/internal/pkg/utils/json"
)

type MemStats struct {
}

func NewMemStats() *MemStats {
	return &MemStats{}
}

func (m *MemStats) Gather(acc *metrics.Accumulator) error {
	var err error

	measurement := "mem"
	now := time.Now()

	// Swap mem
	var swap *Swap
	swap, err = SwapMem()
	if err != nil {
		return err
	}
	acc.Add(measurement, map[string]string{
		"host": acc.Uid,
	}, json.StructToMap(swap), now)

	// Virtual mem
	var virtual *Virtual
	virtual, err = VirtualMem()
	if err != nil {
		return err
	}
	acc.Add(measurement, map[string]string{
		"host": acc.Uid,
	}, json.StructToMap(virtual), now)

	return nil
}

func init() {
	metrics.Add("mem", func() metrics.Metric {
		return NewMemStats()
	})
}
