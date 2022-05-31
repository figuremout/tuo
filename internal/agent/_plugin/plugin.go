package plugin

type Plugin interface {
    PluginDescriber
    // Gather takes in an accumulator and adds the metrics that the Input
	// gathers. This is called every agent.interval
	Gather(Accumulator) error
}

type PluginDescriber interface {
    // SampleConfig returns the default configuration of the Processor
	SampleConfig() string
	// Description returns a one-sentence description on the Processor
	Description() string
}

// Initializer is an interface that all plugin types: Inputs, Outputs,
// Processors, and Aggregators can optionally implement to initialize the
// plugin.
type Initializer interface {
	// Init performs one time setup of the plugin and returns an error if the
	// configuration is invalid.
	Init() error
}
