package main

import (
	"context"

	// init() will register its router into ruoters
	_ "github.com/githubzjm/tuo/api/v1/base"
	_ "github.com/githubzjm/tuo/api/v1/clusters"
	_ "github.com/githubzjm/tuo/api/v1/metrics"
	_ "github.com/githubzjm/tuo/api/v1/nodes"
	_ "github.com/githubzjm/tuo/api/v1/users"

	"github.com/githubzjm/tuo/internal/pkg/influxdb"
	"github.com/githubzjm/tuo/internal/pkg/logger"
	"github.com/githubzjm/tuo/internal/pkg/mysql"
	"github.com/githubzjm/tuo/routers"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	var isUp bool

	// init log
	logConf := &logger.LogConf{
		LogOutput: logger.LogOutputStdout,
	}
	logger.Setup(logConf)
	log.Info("logger init done")

	// init influxdb
	influxdb.InitDB("http://127.0.0.1:8086", "init-token", "init-org", "init-bucket")
	defer influxdb.Close()
	if isUp, err = influxdb.Client.Ping(context.Background()); err != nil {
		log.Fatalf("ping influxdb failed: %v", err)
	} else {
		log.Infof("ping influxdb: %v", isUp)
	}
	// influxdb.Query(`
	// from(bucket:"init-bucket")
	// 	|> range(start:-1d)
	// 	|> filter(fn: (r) => r.cpu == "cpu0" and r._field == "time_user")
	// 	|> keep(columns: ["_field", "_value", "_time"])
	// 	|> last()
	// `)

	// init mysql
	if err = mysql.InitDB("root", "tuo-mysql", "localhost", "3306", "tuo"); err != nil {
		log.Fatalf("connect mysql failed: %v", err)
	} else {
		log.Info("connect mysql success")
	}
	defer mysql.SqlDB.Close()
	if err = mysql.SqlDB.Ping(); err != nil {
		log.Fatalf("ping mysql failed: %v", err)
	} else {
		log.Info("ping mysql success")
	}

	router := routers.Init()
	if err := router.Run(":8080"); err != nil {
		log.Errorf("startup service failed, err: %v\n", err)
	}
	// Start gin server
	// Switch to "release" mode in production.
	// gin.SetMode(gin.ReleaseMode)
}
