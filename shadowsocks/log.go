package shadowsocks

import (
	"log"
	"os"
)

type Logger bool

var logger = log.New(os.Stdout, "", log.Ltime)

func (l Logger) Printf(format string, args ...interface{}) {
	if l {
		logger.Printf(format, args...)
	}
}

func (l Logger) Println(args ...interface{}) {
	if l {
		logger.Println(args...)
	}
}
