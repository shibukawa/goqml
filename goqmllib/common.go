package goqmllib

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func Target() string {
	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	return goos
}

type sequentialRun struct {
	workDir string
	err     error
}

func SequentialRun(workDir string) *sequentialRun {
	return &sequentialRun{
		workDir: workDir,
	}
}

func (s *sequentialRun) Run(command string, args ...string) *sequentialRun {
	if s.err != nil {
		return s

	}
	cmd := Command(command, s.workDir, args...)
	cmd.Silent = true
	err := cmd.Run()
	if err != nil {
		s.err = fmt.Errorf("cmd: `%s %s` err: %s", command, strings.Join(args, " "), err.Error())
	}
	return s
}

func (s *sequentialRun) Finish() error {
	return s.err
}

func unixTime(paths ...string) int64 {
	var unixTime int64
	for _, path := range paths {
		stat, err := os.Stat(path)
		if os.IsNotExist(err) {
			return 0
		}
		if unixTime < stat.ModTime().Unix() {
			unixTime = stat.ModTime().Unix()
		}
	}
	return unixTime
}
