package nflogger

import (
	"fmt"
	"log"
)

// Make makes a new custom logger
func Make(prefix string) *log.Logger {
	return log.New(log.Writer(), fmt.Sprintf("[%s]  ", prefix), 0)
}
