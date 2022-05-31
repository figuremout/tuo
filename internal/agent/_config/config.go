package config

import (
	"github.com/githubzjm/tuo/internal/agent"
	"github.com/githubzjm/tuo/internal/pkg/influxdb"
	"github.com/githubzjm/tuo/internal/pkg/logger"
	jsonUtil "github.com/githubzjm/tuo/internal/pkg/utils/json"
	yamlUtil "github.com/githubzjm/tuo/internal/pkg/utils/yaml"
)

// const (
// 	RotationPrecisionNs = "ns"
// 	RotationPrecisionUs = "us"
// 	RotationPrecisionMs = "ms"
// 	RotationPrecisionS  = "s"
// )

var (
	// Default total config
	config_default = &Config{
		Influxdb: nil,
		Log:      logger.LogConf_default,
		Agent:    agent.AgentConf_default,
	}
)

/*
 * Total config
 */
type Config struct {
	Influxdb *influxdb.InfluxDBConf `yaml:"influxdb"`
	Log      *logger.LogConf        `yaml:"log"`
	Agent    *agent.AgentConf       `yaml:"agent"`
}

func NewConfig() *Config {
	return &Config{
		Agent: agent.NewAgent(nil),
	}
}

func (c *Config) String() string {
	return jsonUtil.Stringfy(c, "")
}

//func (c *Config) SetDefaults() {
//    // fill in default values
//    if c.Log == nil {
//        c.Log = logger.Log_default
//    }
//    if c.Collector == nil {
//        c.Collector = collector_default
//    }
//}

func LoadFile(path string) (*Config, error) {
	var conf = config_default
	err := yamlUtil.LoadYaml(path, conf)
	return conf, err
}
