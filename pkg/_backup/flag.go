package command

import (
	"fmt"
)

type Flag struct {
	Name        string
	Description string
	HasValue    bool
	Owner       *Node

	IsSet bool
	// Value string // assigned after FlagSet.Parse
	Value interface{}
}

func (f *Flag) Refresh() {
	f.Value = ""
	f.IsSet = false
}

func (f *Flag) Verify() error {
	if f.IsSet && !f.HasValue && f.Value != "" {
		return fmt.Errorf("flag with extra arguments: %s", f.Name)
	}
	return nil
}
