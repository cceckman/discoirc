// Package log configures the behavior of Go's standard logging library.

package log

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"io"
	"time"
)

var (
	logtostderr = flag.Bool("logtostderr", true, "Write logs to stderr.")

	logpath = flag.String("logpath", "/tmp", "Path to create a log file in. If -logfile is not specified, generate a uniquely-named file in this directory.")
	logfile = flag.String("logfile", "", "File to log to; overrides -logpath.")

	logperms = flag.Int("logperms", 0755, "Permissions for the log file (Unix mask, RWXRWXRWX only.)")
)

const (
	openFileBits = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	logFlags     = log.LstdFlags | log.LUTC | log.Lshortfile
)

func Init() error {
	// Attach loggers to the relevant outputs.

	p := ""
	if *logfile != "" {
		p = *logfile
	} else if *logpath != "" {
		binname := path.Base(os.Args[0])
		filename := fmt.Sprintf(
			"%s.%s.%d.log",
			binname,
			time.Now().UTC().Format("2006-01-02.15:04:05"),
			os.Getpid(),
		)
		p = path.Join(*logpath, filename)
	}

	var f io.Writer = nil
	if p != "" {
		var err error
		perms := os.FileMode(*logperms & 0777)

		f, err = os.OpenFile(p, openFileBits, perms)
		if err != nil {
			return fmt.Errorf("could not open file %s for writing logs: %v", p, err)
		}
	}

	// Configure the library.
	logouts := []io.Writer{}
	if p != "" {
		logouts = append(logouts, f)
	}
	if *logtostderr {
		logouts = append(logouts, os.Stderr)
	}
	log.SetOutput(io.MultiWriter(logouts...))
	log.SetFlags(logFlags)

	log.Printf("Log settings applied.")

	return nil
}

func LogArgs() {
	log.Println("Arguments:")
	for _, arg := range os.Args {
		log.Printf("\t%s\n", arg)
	}
	log.Println("(end argument list)")
}
