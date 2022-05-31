package mem

import (
	"github.com/shirou/gopsutil/v3/mem"
)

type Virtual struct {
	// Total amount of RAM on this system
	Total uint64 `json:"total"`

	// RAM available for programs to allocate
	//
	// This value is computed from the kernel specific values.
	Available uint64 `json:"available"`

	// RAM used by programs
	//
	// This value is computed from the kernel specific values.
	Used uint64 `json:"used"`

	// Percentage of RAM used by programs
	//
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`

	// This is the kernel's notion of free memory; RAM chips whose bits nobody
	// cares about the value of right now. For a human consumable number,
	// Available is what you really want.
	Free uint64 `json:"free"`

	// OS X / BSD specific numbers:
	// http://www.macyourself.com/2010/02/17/what-is-free-wired-active-and-inactive-system-memory-ram/
	// Active   uint64 `json:"active"`
	// Inactive uint64 `json:"inactive"`
	// Wired    uint64 `json:"wired"`

	// FreeBSD specific numbers:
	// https://reviews.freebsd.org/D8467
	// Laundry uint64 `json:"laundry"`

	// Linux specific numbers
	// https://www.centos.org/docs/5/html/5.1/Deployment_Guide/s2-proc-meminfo.html
	// https://www.kernel.org/doc/Documentation/filesystems/proc.txt
	// https://www.kernel.org/doc/Documentation/vm/overcommit-accounting
	// Buffers        uint64 `json:"buffers"`
	// Cached         uint64 `json:"cached"`
	// WriteBack      uint64 `json:"writeBack"`
	// Dirty          uint64 `json:"dirty"`
	// WriteBackTmp   uint64 `json:"writeBackTmp"`
	// Shared         uint64 `json:"shared"`
	// Slab           uint64 `json:"slab"`
	// Sreclaimable   uint64 `json:"sreclaimable"`
	// Sunreclaim     uint64 `json:"sunreclaim"`
	// PageTables     uint64 `json:"pageTables"`
	// SwapCached     uint64 `json:"swapCached"`
	// CommitLimit    uint64 `json:"commitLimit"`
	// CommittedAS    uint64 `json:"committedAS"`
	// HighTotal      uint64 `json:"highTotal"`
	// HighFree       uint64 `json:"highFree"`
	// LowTotal       uint64 `json:"lowTotal"`
	// LowFree        uint64 `json:"lowFree"`
	// SwapTotal      uint64 `json:"swapTotal"`
	// SwapFree       uint64 `json:"swapFree"`
	// Mapped         uint64 `json:"mapped"`
	// VmallocTotal   uint64 `json:"vmallocTotal"`
	// VmallocUsed    uint64 `json:"vmallocUsed"`
	// VmallocChunk   uint64 `json:"vmallocChunk"`
	// HugePagesTotal uint64 `json:"hugePagesTotal"`
	// HugePagesFree  uint64 `json:"hugePagesFree"`
	// HugePageSize   uint64 `json:"hugePageSize"`
}

func VirtualMem() (*Virtual, error) {
	var err error
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	virtual := &Virtual{
		Total:       vmem.Total,
		Available:   vmem.Available,
		Used:        vmem.Used,
		UsedPercent: vmem.UsedPercent,
		Free:        vmem.Free,
		// Active:         vmem.Active,
		// Inactive:       vmem.Inactive,
		// Wired:          vmem.Wired,
		// Laundry:        vmem.Laundry,
		// Buffers:        vmem.Buffers,
		// Cached:         vmem.Cached,
		// WriteBack:      vmem.WriteBack,
		// Dirty:          vmem.Dirty,
		// WriteBackTmp:   vmem.WriteBackTmp,
		// Shared:         vmem.Shared,
		// Slab:           vmem.Slab,
		// Sreclaimable:   vmem.Sreclaimable,
		// Sunreclaim:     vmem.Sunreclaim,
		// PageTables:     vmem.PageTables,
		// SwapCached:     vmem.SwapCached,
		// CommitLimit:    vmem.CommitLimit,
		// CommittedAS:    vmem.CommittedAS,
		// HighTotal:      vmem.HighTotal,
		// HighFree:       vmem.HighFree,
		// LowTotal:       vmem.LowTotal,
		// LowFree:        vmem.LowFree,
		// SwapTotal:      vmem.SwapTotal,
		// SwapFree:       vmem.SwapFree,
		// Mapped:         vmem.Mapped,
		// VmallocTotal:   vmem.VmallocTotal,
		// VmallocUsed:    vmem.VmallocUsed,
		// VmallocChunk:   vmem.VmallocChunk,
		// HugePagesTotal: vmem.HugePagesTotal,
		// HugePagesFree:  vmem.HugePagesFree,
		// HugePageSize:   vmem.HugePageSize,
	}
	return virtual, nil
}
