package colorLog

import (
	"fmt"
	"headless/internal/pkg/colorPreset"
	"io"
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

var logger = newLogger(TRACE, true, true, false)
var SENDER, RECEIVER = "this_is_sender_email_address", "this_is_receiver_email_address"

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
func newLogger(level int, saveLog bool, postLog bool, mailingService bool) *Logger {
	catchShutdown(
		func(isEnableMailingService bool) func() {
			if isEnableMailingService {
				// return calling function which sends email...
				return nil
			} else {
				// fmt.Printf("no mailing service \n")
				return nil
			}
		}(mailingService),
	)
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
		// file, multiWriter := initFileIo("output.log")
		// defer file.Close()
		// log.SetOutput(multiWriter)
		// log.SetPrefix("")
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

// "must" close file properly
func initFileIo(filename string) (*os.File, io.Writer) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("error occured while opening file %v", err)
		return nil, nil
	}
	return file, io.MultiWriter(file, os.Stdout)
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
		// file, multiWriter := initFileIo("output.log")
		// defer file.Close()
		// log.SetOutput(multiWriter)
		// log.SetPrefix(fmt.Sprintf("[%v][%v][%v]| ", getDate(), logType, caller))
		// log.SetFlags(0)
		logType += ""

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
		textPreset = fmt.Sprintf("%v%v", colorPreset.LightWhite, colorPreset.BgGreen)
	case INFO:
		textPreset = fmt.Sprintf("%v%v", colorPreset.LightWhite, colorPreset.BgLightGray)
	case WARN:
		textPreset = fmt.Sprintf("%v%v", colorPreset.LightWhite, colorPreset.BgMagenta)
	case ERROR:
		textPreset = fmt.Sprintf("%v%v", colorPreset.LightWhite, colorPreset.BgYellow)
		messagePreset = fmt.Sprintf("%v%v%v%v",
			colorPreset.BoldOn, "", colorPreset.Yellow, colorPreset.BgLightGray)
	case FATAL:
		textPreset = fmt.Sprintf("%v%v", colorPreset.LightWhite, colorPreset.BgRed)
		messagePreset = fmt.Sprintf("%v%v%v%v",
			colorPreset.BoldOn, colorPreset.UnderLineOn, colorPreset.Red, colorPreset.BgWhite)
	}
	return fmt.Sprintf("%v%v%v", textPreset, levels[logLevel], colorPreset.Reset),
		fmt.Sprintf("%v%v%v", messagePreset, message, colorPreset.Reset)
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
