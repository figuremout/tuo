package logger

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	// FilePerm defines the permissions that Writer will use for all the files it creates
	FilePerm = os.FileMode(0644)

	LogOutputStderr = "stderr"
	LogOutputStdout = "stdout"
)

// var (
// 	// Default log config
// 	LogConf_default = &LogConf{
// 		LogOutput: LogOutputStderr,
// 		Rotation: &RotationConf{
// 			RotationMaxSize:     "50MB",
// 			RotationInterval:    "0",
// 			RotationMaxArchives: 0,
// 		},
// 	}
// )

type LogConf struct {
	LogOutput string `yaml:"logoutput"`
	// Rotation  *RotationConf `yaml:"rotation"`
}

// type RotationConf struct {
// 	RotationMaxSize     string `yaml:"rotation_max_size"`
// 	RotationInterval    string `yaml:"rotation_interval"`
// 	RotationMaxArchives int    `yaml:"rotation_max_archives"`
// }

// Setup standard logger
func Setup(logConf *LogConf) {
	var writer, defaultWriter io.Writer
	defaultWriter = os.Stderr

	// Setup formatter
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Setup report caller
	log.SetReportCaller(true)

	// Setup level
	// log level: Panic > Fatal > Error > Warn > Info > Debug > Trace
	log.SetLevel(log.TraceLevel)

	// Setup output
	switch logConf.LogOutput {
	case LogOutputStderr, "":
		writer = defaultWriter
	case LogOutputStdout:
		writer = os.Stdout
	default:
		fw, err := os.OpenFile(logConf.LogOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, FilePerm)
		if err != nil {
			writer = defaultWriter
			log.Infof("Failed to log to file: %v, using default stderr\n", logConf.LogOutput)
		} else {
			writer = fw
		}
	}
	log.SetOutput(writer)
}

// Also need setup standard logger first, new loggeres will inherit the features
func GetReqLogger(request_id, user_ip string) *log.Entry {
	var fields = map[string]interface{}{
		"request_id": request_id,
		"user_ip":    user_ip,
	}
	return getLoggerWithFields(fields)
}

func getLoggerWithFields(fields map[string]interface{}) *log.Entry {
	return log.WithFields(fields)
}
