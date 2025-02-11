package utils

import "time"

const (
	timeFormatString = "2006-01-02 15:04:05"
)
func SerializeTime(t time.Time) string {
	timeFormat := "2006-01-02 15:04:05"
	return t.Format(timeFormat)
}

func DeserializeTime(s string) (time.Time,error) {
	t,err := time.Parse(timeFormatString, s)
	if err != nil {
		return time.Time{},nil
	}
	return t,nil
}