package errorHandle

import "fmt"

const (
	pattern = "\npath: %s%s\nfunction: %s\nerror: %s"
)

func ErrorFormatString(path, file, function string, err error) string {
	return fmt.Sprintf(pattern, file, path, function, err)
}

func ErrorFormat(path, file, function string, err error) error {
	return fmt.Errorf(pattern, file, path, function, err)
}
