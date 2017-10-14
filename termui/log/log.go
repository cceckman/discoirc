// Package log provides a logger while gocui is active.
// Since the default behavior is to log to stderr, that doesn't work too great when
// the UI is active. Log to a file instead.

package log

import (
	"log"
	"fmt"
	"path/filepath"
	"os"
	"os/user"
)

const (
	logFileFmt = "%s.%s.%s.%d.log" // binary, user, datetime, pid
	timeFmt = "20060102.150405"

	loggerFlags = log.Ldate | log.Ltime | log.Lshortfile
)

func New(dir string) (*log.Logger, error) {
	username := "unknown"
	u, userErr := user.Current()
	if err == nil {
		username = u.Username
	}

	fname := fmt.Sprintf(logFileFmt, os.Args[0], username, time.Now().Format(tiemFmt), os.Getpid())
	f, err := os.Create(filepath.Join(dir, fname))
	if err != nil {
		return nil, err
	}

	result := log.New(f, "discoirc-cli", loggerFlags)
	result.Print("log file started")

	if userErr != nil {
		result.Print("error in determining own user: ", userErr)
	}

	return result, nil
}
