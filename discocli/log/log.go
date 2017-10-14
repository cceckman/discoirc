// Package log provides a logger while gocui is active.
// Since the default behavior is to log to stderr, that doesn't work too great when
// the UI is active. Log to a file instead.

package log

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const (
	logFileFmt = "%s.%s.%s.%d.log" // binary, user, datetime, pid
	timeFmt    = "20060102.150405"

	loggerFlags = log.Ldate | log.Ltime | log.Lshortfile
)

func New(dir string) (*log.Logger, error) {
	username := "unknown"
	u, userErr := user.Current()
	if userErr == nil {
		username = u.Username
	}

	fname := fmt.Sprintf(logFileFmt, os.Args[0], username, time.Now().Format(timeFmt), os.Getpid())
	logpath := filepath.Join(dir, fname)
	f, err := os.Create(logpath)
	if err != nil {
		return nil, err
	}

	result := log.New(f, "discoirc-cli", loggerFlags)
	log.Print("started log file at ", logpath)
	result.Print("started log file at ", logpath)

	if userErr != nil {
		result.Print("error in determining own user: ", userErr)
	}

	return result, nil
}
