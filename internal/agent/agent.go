package agent

import (
	"fmt"
	"time"

	"github.com/githubzjm/tuo/internal/pkg/metrics"
	_ "github.com/githubzjm/tuo/internal/pkg/metrics/all"
	"github.com/githubzjm/tuo/internal/pkg/rpc"
	log "github.com/sirupsen/logrus"
)

type Agent struct {
	Conf *AgentConf
}

type AgentConf struct {
	Interval string   `yaml:"interval"`
	Metrics  []string `yaml:"metrics"`
	Addr     string   `yaml:"addr"`
}

func NewAgent(conf *AgentConf) *Agent {
	if conf != nil {
		return &Agent{
			Conf: conf,
		}
	} else {
		return &Agent{
			Conf: &AgentConf{
				Interval: "10s",
				Metrics: []string{
					"cpu",
					"mem",
					"swap",
					"system",
					"kernel",
					"processes",
					"disk",
					"diskio",
				},
				Addr: "localhost:55555",
			},
		}
	}
}

func (a *Agent) Run(quit chan bool) error {
	var err error
	var interval time.Duration

	go rpc.StartGRPCServer(a.Conf.Addr)

	acc, err := metrics.NewAccumulator(100)
	if err != nil {
		return err
	}
	defer acc.Close() // will end report goruotines
	// Start report goruotines
	go acc.Report(3)

	interval, err = time.ParseDuration(a.Conf.Interval)
	if err != nil {
		return err
	}

	// Start collect goruotine
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for _, m := range a.Conf.Metrics {
				creator, ok := metrics.MetricCreators[m]
				if !ok {
					return fmt.Errorf("metric %s is illegal", m)
				}
				metric := creator()
				if err := metric.Gather(acc); err != nil {
					log.Infof("gather error: %v", err)
				}
			}
		case <-quit:
			log.Info("agent collect loop quit")
			return nil
		}
	}
}
