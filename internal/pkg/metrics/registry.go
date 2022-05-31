package metrics

type Metric interface {
	Gather(*Accumulator) error
}

type Creator func() Metric

//type Gather func(acc *Accumulator) error

var MetricCreators = map[string]Creator{}

func Add(name string, creator Creator) {
	MetricCreators[name] = creator
}
