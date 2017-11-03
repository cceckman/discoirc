// Package log provides a logger while gocui is active.
// Since the default behavior is to log to stderr, that doesn't work too great when
// the UI is active. Log to a file instead.

package log

import (
	"fmt"
	"log"
	"io/ioutil"
	"os"
	"io"
	"os/user"
	"path/filepath"
	"time"
)

const (
	logFileFmt = "%s.%s.%s.%d.log" // binary, user, datetime, pid
	timeFmt    = "20060102.150405"

	loggerFlags = log.Ldate | log.Ltime | log.Lshortfile
)

func New() (*log.Logger, error) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	username := "unknown"
	u, userErr := user.Current()
	if userErr == nil {
		username = u.Username
	}

	binname := filepath.Base(os.Args[0])

	fname := fmt.Sprintf(logFileFmt, binname, username, time.Now().Format(timeFmt), os.Getpid())
	logpath := filepath.Join(dir, fname)
	log.Print("started log file at ", logpath)
	f, err := os.Create(logpath)
	if err != nil {
		return nil, err
	}
	return NewFor(f)
}

func NewFor(f io.Writer) (*log.Logger, error) {
	result := log.New(f, "discoirc-cli", loggerFlags)
	result.Print("starting new debug log")

	return result, nil
}
