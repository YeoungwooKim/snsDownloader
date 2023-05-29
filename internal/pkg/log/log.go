package log

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var levels = [...]string{
	"TRACE",
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

var logger = newLogger(TRACE, false, true)

type Logger struct {
	LogLevel int
	SaveLog  bool
	PostLog  bool
}

type Caller struct {
	file     string
	funcName string
	line     int
}

func SetLogLevel(level int) {
	logger.LogLevel = level
}

// Options to save conditions according to level, save log to local, send logfile to remote
func newLogger(level int, saveLog bool, postLog bool) *Logger {
	catchShutdown()
	return &Logger{
		LogLevel: level,
		SaveLog:  saveLog,
		PostLog:  postLog,
	}
}

// Catch shutdown signal
func catchShutdown(gracefulShutdownFunc ...func()) {
	// create channel and asign signal (1 receive)
	var sigs = make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGTERM, // 15
		syscall.SIGHUP,  // 1
		syscall.SIGINT,  // 2
		syscall.SIGQUIT, // 3
		os.Interrupt,    // == SIGINT
	)

	// close
	go func() {
		sig := <-sigs
		log.SetPrefix("")
		log.Println("::: Terminating... :::\ncaught signal : ", sig)

		for i := 0; i < len(gracefulShutdownFunc); i++ {
			gracefulShutdownFunc[i]()
		}
		os.Exit(0)
	}()
}

// returns specific location info that log declared.
func getCallerInfo() *Caller {
	pc, file, line, _ := runtime.Caller(3)
	funcName := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return &Caller{
		file:     file,
		funcName: funcName[len(funcName)-1],
		line:     line,
	}
}

func (logger *Logger) print(logLevel int, message string, a ...interface{}) {
	var logType string
	logType, message = logger.highlightMode(logLevel, fmt.Sprintf(message, a...))
	callerInfo := getCallerInfo()
	caller := fmt.Sprintf("%v::%v(%v)", callerInfo.file, callerInfo.funcName, callerInfo.line)
	if len(caller) > 25 {
		caller = caller[len(caller)-25:]
	}

	if logLevel >= logger.LogLevel {
		// log.SetOutput(multiWriter)
		log.SetPrefix(fmt.Sprintf("[%v][%v][%v]| ", getDate(), logType, caller))
		log.SetFlags(0)

		log.Println(message)
	}
}

func (logger *Logger) highlightMode(logLevel int, message string) (string, string) {
	// TRACE, DEBUG, INFO, WARN, ERROR, FATAL.
	var textPreset, messagePreset string
	switch logLevel {
	case TRACE:
		// lowest log level
	case DEBUG:
		textPreset = fmt.Sprintf("%v%v", LightWhite, BgCyan)
	case INFO:
		textPreset = fmt.Sprintf("%v%v", LightWhite, BgLightBlue)
	case WARN:
		textPreset = fmt.Sprintf("%v%v", LightWhite, BgMagenta)
	case ERROR:
		textPreset = fmt.Sprintf("%v%v", LightWhite, BgYellow)
		messagePreset = fmt.Sprintf("%v%v%v%v",
			BoldOn, "", Yellow, BgLightGray)
	case FATAL:
		textPreset = fmt.Sprintf("%v%v", LightWhite, BgRed)
		messagePreset = fmt.Sprintf("%v%v%v%v",
			BoldOn, UnderLineOn, Red, BgWhite)
	}
	return fmt.Sprintf("%v%v%v", textPreset, levels[logLevel], Reset),
		fmt.Sprintf("%v%v%v", messagePreset, message, Reset)
}

func getDate() interface{} {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

// TRACE, DEBUG, INFO, WARN, ERROR, FATAL.
func Trace(message string, a ...interface{}) {
	logger.print(TRACE, message, a...)
}

func Debug(message string, a ...interface{}) {
	logger.print(DEBUG, message, a...)
}

func Info(message string, a ...interface{}) {
	logger.print(INFO, message, a...)
}

func Warn(message string, a ...interface{}) {
	logger.print(WARN, message, a...)
}

func Error(message string, a ...interface{}) {
	logger.print(ERROR, message, a...)
}

func Fatal(message string, a ...interface{}) {
	logger.print(FATAL, message, a...)
}

////////////////////////////////////////////////////////////////////////////
// var defaultLogger *zap.SugaredLogger

// func init() {
// 	defaultLogger = initLogger()
// }

// func initLogger() *zap.SugaredLogger {
// 	config := zap.NewProductionConfig()
// 	encoderConfig := zap.NewProductionEncoderConfig()
// 	encoderConfig.TimeKey = "timestamp"
// 	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
// 	encoderConfig.StacktraceKey = ""
// 	config.EncoderConfig = encoderConfig

// 	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
// 	level := zap.InfoLevel

// 	core := zapcore.NewTee(
// 		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
// 	)
// 	logger := zap.New(core)
// 	// logger, err := config.Build(zap.AddCallerSkip(int(zapcore.ErrorLevel)))
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	defer logger.Sync()

// 	return logger.Sugar()
// }

// func Info(msg string, fields ...interface{}) {
// 	defaultLogger.Infof(msg, fields)
// }

// func Debug(msg string, fields ...interface{}) {
// 	defaultLogger.Debugf(msg, fields)
// }

// func Error(msg string, fields ...interface{}) {
// 	defaultLogger.Errorf(msg, fields)
// }

// func Fatal(msg string, fields ...interface{}) {
// 	defaultLogger.Fatalf(msg, fields)
// }
