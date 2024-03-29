// Copyright 2016 Ryan Boehning. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package q

import (
	"bytes"
	"regexp"

	//"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type color string

const (
	// ANSI color escape codes
	bold     color = "\033[1m"
	yellow   color = "\033[33m"
	cyan     color = "\033[36m"
	endColor color = "\033[0m" // "reset everything"

	maxLineWidth = 80
	bufSize      = 16384
)

// The q logger singleton
var std *logger

// logger writes pretty logs to the $TMPDIR/q file. It takes care of opening and
// closing the file. It is safe for concurrent use.
type logger struct {
	mu       sync.Mutex    // protects all the other fields
	buf      *bytes.Buffer // collects writes before they're flushed to the log file
	start    time.Time     // time of first write in the current log group
	timer    *time.Timer   // when it gets to 0, start a new log group
	lastFile string        // last file to call q.Q(). determines when to print header
	lastFunc string        // last function to call q.Q()
	logFile  string
}

// P (Regexp Pattern for functions) and O (Output destination) can be set to limit output
var P, O string

// init creates the standard logger.
func init() {
	// Starting with 0 time doesn't mean the timer is stopped, so we must
	// explicitly stop the timer.
	t := time.NewTimer(0)
	t.Stop()

	buf := &bytes.Buffer{}
	buf.Grow(bufSize)
	std = &logger{
		buf:   buf,
		timer: t,
	}
	//flag.StringVar(&P, "qp", "", "q package/function regexp pattern")
	//flag.StringVar(&O, "qo", "q", "q log output destination")
}

// header returns a formatted header string, e.g. [14:00:36 main.go main.main:122]
// if the 2s timer has expired, or the calling function or filename has changed.
// If none of those things are true, it returns an empty string.
func (l *logger) header(funcName, file string, line int) string {
	// Reset the 2s timer.
	timerExpired := l.resetTimer(2 * time.Second)

	if !timerExpired && funcName == l.lastFunc && file == l.lastFile {
		// Don't print a header line.
		return ""
	}

	l.lastFunc = funcName
	l.lastFile = file

	now := time.Now().UTC().Format("15:04:05")

	return fmt.Sprintf("[%s %s:%d %s]", now, shortFile(file), line, funcName)
}

// shortFile takes an absolute file path and returns just the <directory>/<file>,
// e.g. "foo/bar.go".
func shortFile(file string) string {
	dir := filepath.Base(filepath.Dir(file))
	file = filepath.Base(file)
	return filepath.Join(dir, file)
}

// resetTimer resets the logger's timer to the given time. It returns true if
// the timer had expired before it was reset.
func (l *logger) resetTimer(d time.Duration) (expired bool) {
	expired = !l.timer.Reset(d)
	if expired {
		l.start = time.Now()
	}
	return expired
}

func getPath() (string, error) {
	path := ""
	if strings.HasPrefix(O, "/") {
		path = O
	} else if strings.HasPrefix(O, "./") {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot get cwd", err)
			return "", fmt.Errorf("can't get cwd, %v\n", err)
		}
		path = filepath.Join(cwd, O)
	} else {
		path = filepath.Join(os.TempDir(), O)
	}
	return path, nil
}

// flush writes the logger's buffer to disk.
func (l *logger) flush() error {
	var f *os.File
	switch O {
	case "stderr":
		f = os.Stderr
	case "stdout":
		f = os.Stdout
	default:
		path, err := getPath()
		if err != nil {
			return err
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil && l.logFile != path {
			fmt.Fprintln(os.Stderr, "cannot open", path, err)
			l.logFile = path
			return fmt.Errorf("failed to open %q: %v", path, err)
		}
		if l.logFile != path {
			fmt.Fprintln(os.Stderr, "logging to", path)
		}
		l.logFile = path
		defer f.Close()
		return l.ioCopy(f)
	}
	return l.ioCopy(f)
}

func (l *logger) ioCopy(f *os.File) error {
	_, err := io.Copy(f, l.buf)
	l.buf.Reset()
	if err != nil {
		return fmt.Errorf("failed to flush q buffer: %v", err)
	}
	return nil
}

// output writes to the log buffer. Each log message is prepended with a
// timestamp. Long lines are broken at 80 characters.
func (l *logger) output(args ...string) {
	timestamp := fmt.Sprintf("%.3fs", time.Since(l.start).Seconds())
	timestampWidth := len(timestamp) + 1 // +1 for padding space after timestamp
	timestamp = colorize(timestamp, yellow)

	// preWidth is the length of everything before the log message.
	fmt.Fprint(l.buf, timestamp, " ")

	// Subsequent lines have to be indented by the width of the timestamp.
	indent := strings.Repeat(" ", timestampWidth)
	padding := "" // padding is the space between args.
	lineArgs := 0 // number of args printed on the current log line.
	lineWidth := timestampWidth
	for _, arg := range args {
		argWidth := argWidth(arg)
		lineWidth += argWidth + len(padding)

		// Some names in name=value strings contain newlines. Insert indentation
		// after each newline so they line up.
		arg = strings.Replace(arg, "\n", "\n"+indent, -1)

		// Break up long lines. If this is first arg printed on the line
		// (lineArgs == 0), it makes no sense to break up the line.
		if lineWidth > maxLineWidth && lineArgs != 0 {
			fmt.Fprint(l.buf, "\n", indent)
			lineArgs = 0
			lineWidth = timestampWidth + argWidth
			padding = ""
		}
		fmt.Fprint(l.buf, padding, arg)
		lineArgs++
		padding = " "
	}

	fmt.Fprint(l.buf, "\n")
}

// Q pretty-prints the given arguments to the $TMPDIR/q log file.
func Q(v ...interface{}) {
	if P == "" { //XXX hack
		return
	}

	std.mu.Lock()
	defer std.mu.Unlock()

	// Flush the buffered writes to disk.
	defer std.flush()

	args := formatArgs(v...)
	funcName, file, line, err := getCallerInfo()
	if err != nil {
		std.output(args...) // no name=value printing
		return
	}

	matched, err := regexp.MatchString(P, funcName)
	matched2, err := regexp.MatchString(P, file)
	if !matched && !matched2 {
		return
	}

	// Print a header line if this q.Q() call is in a different file or
	// function than the previous q.Q() call, or if the 2s timer expired.
	// A header line looks like this: [14:00:36 main.go main.main:122].
	header := std.header(funcName, file, line)
	if header != "" {
		fmt.Fprint(std.buf, "\n", header, "\n")
	}

	// q.Q(foo, bar, baz) -> []string{"foo", "bar", "baz"}
	names, err := argNames(file, line)
	if err != nil {
		std.output(args...) // no name=value printing
		return
	}

	// Convert the arguments to name=value strings.
	args = prependArgName(names, args)
	std.output(args...)
}
