package main

import (
	"context"
	"flag"
	"fmt"

	//    "gopkg.in/yaml.v2"
	log "github.com/sirupsen/logrus"

	"github.com/githubzjm/tuo/internal/agent"
	"github.com/githubzjm/tuo/internal/pkg/influxdb"
	"github.com/githubzjm/tuo/internal/pkg/logger"
	//"github.com/githubzjm/tuo/internal/agent/plugin"
	// "github.com/githubzjm/tuo/internal/agent/config"
)

// cmdline flags
var (
	// -h flag is set by default, which outputs supported flags
	fVersion = flag.Bool("v", false, "display version and exit")
	//fConfigFile = flag.String("f", "./configs/agent.yml", "specify config file")
)

var (
	version string // set by ldflags in Makefile
)

func init() {
	if version == "" {
		version = "unknown"
	}
}

func main() {
	var err error
	var isUp bool

	// parse flag
	flag.Parse()
	switch {
	case *fVersion:
		fmt.Printf("agent v%v\n", version)
		return
	}

	// load config
	// var conf *config.Config
	// conf, err = config.LoadFile(*fConfigFile)
	// if err != nil {
	// 	fmt.Printf("load config file failed: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(conf)

	// init log
	logConf := &logger.LogConf{
		LogOutput: logger.LogOutputStdout,
	}
	logger.Setup(logConf)
	log.Info("logger init done")

	reqLogger := logger.GetReqLogger("reqid", "userip")
	reqLogger.Info("test msg")
	//log.Info("log to file")

	// test metric
	//influxdbWriteAPI := influxdbUtil.NewWriteAPI(influxdbClient, conf.Influxdb.Org, conf.Influxdb.Bucket)

	// collect
	// cpuStats := cpu.NewCPUStats()
	// times, _ := cpu.CollectCPUTimes()
	// m, _ := cpu.TimesStatToMap(times)

	// init influxdb
	influxdb.InitDB("http://127.0.0.1:8086", "init-token", "init-org", "init-bucket")
	defer influxdb.Close()
	if isUp, err = influxdb.Client.Ping(context.Background()); err != nil {
		log.Fatalf("ping influxdb failed: %v", err)
	} else {
		log.Infof("ping influxdb: %v", isUp)
	}

	// init agent
	agentConf := &agent.AgentConf{
		Interval: "1s",
		Metrics: []string{
			"cpu",
			"mem",
			"host",
		},
		Addr: "localhost:55555",
	}
	agent := agent.NewAgent(agentConf)
	quit := make(chan bool)
	defer func() {
		quit <- true // TODO need waitgroup to make sure the quit procedure is done,
		// otherwise main goroutine will end before the goroutines
	}()

	err = agent.Run(quit)
	if err != nil {
		log.Info(err)
	}

	// log.Info("sleep start")
	// time.Sleep(time.Second * 12)
	// log.Info("sleep end")

}
