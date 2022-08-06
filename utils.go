package intranet

import "time"

const dateFormat = "2006-01-02"

func Date(d time.Time) string {
	return d.Format(dateFormat)
}
