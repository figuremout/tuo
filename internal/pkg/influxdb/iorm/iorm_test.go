package iorm

import (
	"fmt"
	"testing"
	"time"

	"github.com/githubzjm/tuo/internal/pkg/influxdb/iorm/schema"
)

type CPU struct {
	schema.Model

	Tags CPUTags

	Fields CPUFields
}

type CPUTags struct {
	CPU  string `iorm:"cpu"`
	Host string `iorm:"host"`
}

type CPUFields struct {
	// percent
	Percent int64 `iorm:"percent"`

	// info
	Model string

	// times
	User float64
}

func Test(t *testing.T) {
	cpu := CPU{
		Model:  schema.Model{Measurement: "cpum", Time: time.Now()},
		Tags:   CPUTags{CPU: "cpu-total", Host: "local"},
		Fields: CPUFields{Percent: 50},
	}
	p := schema.Parse(cpu)

	for _, f := range p.FieldList() {
		fmt.Printf("%#v\n", f)
	}

	for _, t := range p.TagList() {
		fmt.Printf("%#v\n", t)
	}

	t.Log("test register pass")
}
