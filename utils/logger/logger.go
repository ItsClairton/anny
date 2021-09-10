package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TwinProduction/go-color"
)

var (
	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
)

func getPretty(color string, tag string, content string) string {
	return fmt.Sprintf("[%s] %s[%s]%s %s", time.Now().Format("02/01/2006 - 15:04:05"), color, tag, "\u001B[0m", content)
}

func Info(s string, args ...interface{}) {
	stdout.Println(getPretty(color.Green, "INFO", fmt.Sprintf(s, args...)))
}

func Debug(s string, args ...interface{}) {
	stdout.Println(getPretty(color.Cyan, "DEBUG", fmt.Sprintf(s, args...)))
}

func Error(s string, args ...interface{}) {
	stderr.Println(getPretty(color.Red, "ERROR", fmt.Sprintf(s, args...)))
}

func ErrorAndExit(s string, args ...interface{}) {
	stderr.Println(getPretty(color.Red, "ERROR", fmt.Sprintf(s, args...)))
	os.Exit(0)
}

func Warn(s string, args ...interface{}) {
	stdout.Println(getPretty(color.Yellow, "WARN", fmt.Sprintf(s, args...)))
}
