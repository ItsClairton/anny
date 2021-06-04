package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
)

func getPretty(color string, tag string, content string) string {
	return fmt.Sprintf("[%s] %s[%s]%s %s", time.Now().Format("02/01/2006 - 15:04:05"), color, tag, "\u001B[0m", content)
}

func Info(s string, args ...interface{}) {
	stdout.Println(getPretty("\u001B[32m", "INFO", fmt.Sprintf(s, args...)))
}

func Debug(s string, args ...interface{}) {
	stdout.Println(getPretty("\u001B[35m", "DEBUG", fmt.Sprintf(s, args...)))
}

func Error(s string, args ...interface{}) {
	stderr.Println(getPretty("\u001B[31m", "ERROR", fmt.Sprintf(s, args...)))
}

func ErrorAndExit(s string, args ...interface{}) {
	stderr.Println(getPretty("\u001B[31m", "ERROR", fmt.Sprintf(s, args...)))
	os.Exit(0)
}

func Warn(s string, args ...interface{}) {
	stdout.Println(getPretty("\u001B[33m", "WARN", fmt.Sprintf(s, args...)))
}
