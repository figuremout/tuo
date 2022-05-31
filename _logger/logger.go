package logger

import (
    "log"
    "os"
    "io"
    "time"
    "github.com/dustin/go-humanize"
    "githubzjm/tuo-agent/config"
    "githubzjm/tuo-agent/internal/logger/rotate"
)

type logLevelNo int

const (
    _ logLevelNo = iota
    DEBUG
    INFO
    WARN
    ERROR
)

var (
    logLevels = map[logLevelNo]string {
        DEBUG: "[DEBUG]",
        INFO: "[INFO]",
        WARN: "[WARN]",
        ERROR: "[ERROR]",
    }
    Debug, Info, Warn, Error *log.Logger
)

func Setup(logConf *config.LogConf) {
    //log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
    //log.SetPrefix("[tuo-agent]")
    //logWriter, _ := createWriter(logConf)
    //log.SetOutput(logWriter)
    logFlags := log.Ldate | log.Ltime | log.Llongfile
    logWriter, _ := createWriter(logConf)
    Debug = log.New(logWriter, logLevels[DEBUG], logFlags)
    Info = log.New(logWriter, logLevels[INFO], logFlags)
    Warn = log.New(logWriter, logLevels[WARN], logFlags)
    Error = log.New(logWriter, logLevels[ERROR], logFlags)
}

func createWriter(logConf *config.LogConf) (io.Writer, error) {
    var writer, defaultWriter io.Writer
    defaultWriter = os.Stderr


    switch logConf.LogTarget {
    case config.LogTargetFile:
        if logConf.LogFile != "" {
            var err error

            // parse RotationMaxSize
            var rotationMaxSizeBytes uint64
            rotationMaxSizeBytes, err = humanize.ParseBytes(logConf.Rotation.RotationMaxSize)
            if err != nil {
                log.Fatalf("Parse rotation_max_size error: %v\n", err)
            }

            // parse RotationInterval
            var rotationIntervalDuration time.Duration
            rotationIntervalDuration, err = time.ParseDuration(logConf.Rotation.RotationInterval)
            if err != nil {
                log.Fatalf("Parse rotation_interval error: %v\n", err)
            }

            if writer, err = rotate.NewFileWriter(logConf.LogFile, rotationIntervalDuration,
                int64(rotationMaxSizeBytes), logConf.Rotation.RotationMaxArchives); err != nil {
                log.Printf("Unable to open %s (%s), using stderr\n", logConf.LogFile, err)
                writer = defaultWriter
            }
        } else {
            writer = defaultWriter
        }
    case config.LogTargetStderr, "":
        writer = defaultWriter
    default:
        log.Printf("Unsupport logtarget: %s, using stderr", logConf.LogTarget)
        writer = defaultWriter
    }
    return writer, nil
}
