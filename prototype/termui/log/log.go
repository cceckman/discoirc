// Package log adds some additional (flag-controlled) logging behavior, and log levels.

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
	// Error, warning, and info logs.
	elog *log.Logger
	wlog *log.Logger
	ilog *log.Logger

	elogtostderr = flag.Bool("elogtostderr", true, "Write error logs to stderr.")
	wlogtostderr = flag.Bool("wlogtostderr", true, "Write warning logs to stderr.")
	ilogtostderr = flag.Bool("llogtostderr", true, "Write info logs to stderr.")

	logpath = flag.String("logpath", "/tmp", "Path to create a log file in. If -logfile is not specified, generate a uniquely-named file in this directory.")
	logfile = flag.String("logfile", "", "File to log to; overrides -logpath.")
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
		filename := fmt.Sprintf(
			"%s.%s.%d.log",
			flag.Args()[0], // binary name
			time.Now().UTC().Format("2006-01-02.15:04:05"),
			os.Getpid(),
		)
		p = path.Join(*logpath, filename)
	}

	var f io.Writer = nil
	if p != "" {
		var err error
		f, err = os.OpenFile(p, openFileBits, 0)
		if err != nil {
			return fmt.Errorf("could not open file %s for writing logs: %v", p, err)
		}
	}

	// error log
	elogouts := []io.Writer{}
	if p != "" {
		elogouts = append(elogouts, f)
	}
	if *elogtostderr {
		elogouts = append(elogouts, os.Stderr)
	}
	elog = log.New(io.MultiWriter(elogouts...), "ERRO", logFlags)

	// warning
	wlogouts := []io.Writer{}
	if p != "" {
		wlogouts = append(wlogouts, f)
	}
	if *wlogtostderr {
		wlogouts = append(wlogouts, os.Stderr)
	}
	wlog = log.New(io.MultiWriter(wlogouts...), "WARN", logFlags)

	// info
	ilogouts := []io.Writer{}
	if p != "" {
		ilogouts = append(ilogouts, f)
	}
	if *ilogtostderr {
		ilogouts = append(ilogouts, os.Stderr)
	}
	ilog = log.New(io.MultiWriter(ilogouts...), "INFO", logFlags)

	return nil
}
