package logger

import (
	"io"
	"os"
	"path/filepath"
	glog "github.com/labstack/gommon/log"
)

func NewLogger(dir, filename string) (*glog.Logger, error) {
	logger := glog.New("")
	logger.SetHeader("{\"time\":\"${time_rfc3339}\", \"level\":\"${level}\"}")
	fpath := filepath.Join(dir, filename)
	if fpath == "" {
		logger.SetOutput(os.Stdout)
	} else {
		f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(io.MultiWriter(os.Stdout, f))
	}
	return logger, nil
}
