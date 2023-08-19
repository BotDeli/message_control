package errorHandle

import (
	"github.com/labstack/gommon/log"
)

const (
	pattern = "\npath: %s%s\nfunction: %s\nerror: %s"
)

func Commit(path, file, function string, err error) {
	log.Errorf(pattern, file, path, function, err)
}
