package format

import "time"

func Date(date time.Time) string {
	return date.Format("02.01.2006")
}
