package kylogr

import (
	"os"
	"path"
	"strconv"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_LOG_DIR            = "./logs"
	DEFAULT_LOG_NAME_PREFIX    = "log"
	DEFAULT_LOG_LEVEL          = "info"
	DEFAULT_LOG_NAME_SUFFIX    = "_%Y_%m_%d_%H_%M_%S.log"
	DEFAULT_LOG_ROTATION_COUNT = "7"
	DEFAULT_LOG_ROTATION_TIME  = "24"
	DEFAULT_LOG_MAX_AGE_HOUR   = "null"
	DEFAULT_LOG_FORMATTER      = "text"
)

func init() {
	InitLog()
}

func InitLog() {
	logDir := GetEnvWithDefault("LOG_DIR", DEFAULT_LOG_DIR)
	logNamePrefix := GetEnvWithDefault("LOG_NAME_PREFIX", DEFAULT_LOG_NAME_PREFIX)
	logNameSuffix := GetEnvWithDefault("LOG_NAME_SUFFIX", DEFAULT_LOG_NAME_SUFFIX)
	logRotationTimeStr := GetEnvWithDefault("LOG_ROTATION_TIME", DEFAULT_LOG_ROTATION_TIME)
	logRotationCountStr := GetEnvWithDefault("LOG_ROTATION_COUNT", DEFAULT_LOG_ROTATION_COUNT)
	logFormatterStr := GetEnvWithDefault("LOG_FORMATTER", DEFAULT_LOG_FORMATTER)
	logLevelStr := GetEnvWithDefault("LOG_LEVEL", DEFAULT_LOG_LEVEL)

	switch logLevelStr {
	case "debug":
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stderr)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logRotationTimeInt64, _ := strconv.ParseInt(logRotationTimeStr, 10, 64)
	logRotationTime := time.Hour * time.Duration(logRotationTimeInt64)

	logRotationCountUint64, _ := strconv.ParseUint(logRotationCountStr, 10, 64)
	logRotationCount := uint(logRotationCountUint64)

	logFileName := logNamePrefix + logNameSuffix

	var logFormatter logrus.Formatter
	switch logFormatterStr {
	case "json":
		logFormatter = &logrus.JSONFormatter{}
	default:
		logFormatter = &logrus.TextFormatter{}
	}

	var lfsHook *lfshook.LfsHook
	logMaxAgeHourStr := GetEnvWithDefault("LOG_MAX_AGE_HOUR", DEFAULT_LOG_MAX_AGE_HOUR)
	if "null" != logMaxAgeHourStr {
		logMaxAgeHourInt64, _ := strconv.ParseInt(logMaxAgeHourStr, 10, 64)
		logMaxAge := time.Hour * time.Duration(logMaxAgeHourInt64)
		lfsHook = lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: logWriterWithMaxAge(logDir, "debug", logFileName, logMaxAge, logRotationTime), // 为不同级别设置不同的输出目的
			logrus.InfoLevel:  logWriterWithMaxAge(logDir, "info", logFileName, logMaxAge, logRotationTime),
			logrus.WarnLevel:  logWriterWithMaxAge(logDir, "warn", logFileName, logMaxAge, logRotationTime),
			logrus.ErrorLevel: logWriterWithMaxAge(logDir, "error", logFileName, logMaxAge, logRotationTime),
			logrus.FatalLevel: logWriterWithMaxAge(logDir, "fatal", logFileName, logMaxAge, logRotationTime),
			logrus.PanicLevel: logWriterWithMaxAge(logDir, "panic", logFileName, logMaxAge, logRotationTime),
		}, logFormatter)
	}

	lfsHook = lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: logWriterWithRotationCount(logDir, "debug", logFileName, logRotationCount, logRotationTime),
		logrus.InfoLevel:  logWriterWithRotationCount(logDir, "info", logFileName, logRotationCount, logRotationTime),
		logrus.WarnLevel:  logWriterWithRotationCount(logDir, "warn", logFileName, logRotationCount, logRotationTime),
		logrus.ErrorLevel: logWriterWithRotationCount(logDir, "error", logFileName, logRotationCount, logRotationTime),
		logrus.FatalLevel: logWriterWithRotationCount(logDir, "fatal", logFileName, logRotationCount, logRotationTime),
		logrus.PanicLevel: logWriterWithRotationCount(logDir, "panic", logFileName, logRotationCount, logRotationTime),
	}, logFormatter)

	logrus.AddHook(lfsHook)
}

func GetEnvWithDefault(key, defaultValue string) string {
	value, had := os.LookupEnv(key)
	if !had {
		return defaultValue
	}
	return value
}

func logWriterWithMaxAge(logDir, level, logFileName string, logMaxAge, logRotationTime time.Duration) *rotatelogs.RotateLogs {
	linkFilePath := path.Join(logDir, level)
	logFilePath := path.Join(logDir, level, logFileName)

	logier, err := rotatelogs.New(
		logFilePath,
		rotatelogs.WithLinkName(linkFilePath),
		rotatelogs.WithMaxAge(logMaxAge),
		rotatelogs.WithRotationTime(logRotationTime),
	)

	if err != nil {
		panic(err)
	}
	return logier
}

func logWriterWithRotationCount(logDir, level, logFileName string, logRotationCount uint, logRotationTime time.Duration) *rotatelogs.RotateLogs {
	linkFilePath := path.Join(logDir, level, level)
	logFilePath := path.Join(logDir, level, logFileName)

	logier, err := rotatelogs.New(
		logFilePath,
		rotatelogs.WithLinkName(linkFilePath),
		rotatelogs.WithRotationCount(logRotationCount),
		rotatelogs.WithRotationTime(logRotationTime),
	)

	if err != nil {
		panic(err)
	}
	return logier
}
