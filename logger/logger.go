package logger

import (
	"fmt"
	"log"
)

var globalLogLevel int

func init() {
	globalLogLevel = INFO
}

const (
	FATAL = iota
	ERROR
	WARN
	INFO
	DEBUG
)

func SetGlobalLogLevel(level int) {
	globalLogLevel = level
}

func Debug(format string, v ...any) {
	if globalLogLevel >= DEBUG {
		if len(v) == 0 {
			log.Println("[DEBUG] " + format)
		} else {
			log.Printf("[DEBUG] "+format+"\n", v...)
		}
	}
}

func Info(format string, v ...any) {
	if globalLogLevel >= INFO {
		if len(v) == 0 {
			log.Println("[INFO] " + format)
		} else {
			log.Printf("[INFO] "+format+"\n", v...)
		}
	}
}

func Warn(format string, v ...any) {
	if globalLogLevel >= WARN {
		if len(v) == 0 {
			log.Println("[WARN] " + format)
		} else {
			log.Printf("[WARN] "+format+"\n", v...)
		}
	}
}

func Error(err error, format string, v ...any) {
	if globalLogLevel >= ERROR {
		log.Printf("[ERROR] "+format+"\n%v\n", append(v, fmt.Errorf("cause:%w", err))...)
	}
}

func Fatal(err error, format string, v ...any) {
	log.Fatalf("[ERROR] "+format+"\nCause:%w\n", append(v, err)...)
}

func Fatalf(format string, v ...any) {
	log.Fatalf("[ERROR] "+format+"\n", v...)
}
