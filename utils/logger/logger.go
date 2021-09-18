package logger

import (
	"log"
	"os"
	"time"

	"github.com/ItsClairton/Anny/utils"
	"github.com/TwinProduction/go-color"
)

var (
	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
)

func format(color string, tag string, content string) string {
	return utils.Fmt("%s[%s] [%s] %s", color, time.Now().Format("02/01/2006 - 15:04:05"), tag, content)
}

func Debug(s string, a ...interface{}) {
	stdout.Println(format(color.Cyan, "DEBUG", utils.Fmt(s, a...)))
}

func Info(s string, a ...interface{}) {
	stdout.Println(format(color.Green, "INFO", utils.Fmt(s, a...)))
}

func Warn(s string, a ...interface{}) {
	stdout.Println(format(color.Yellow, "WARN", utils.Fmt(s, a...)))
}

func Error(s string, a ...interface{}) {
	stderr.Println(format(color.Red, "ERROR", utils.Fmt(s, a...)))
}
