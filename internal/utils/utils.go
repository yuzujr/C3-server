package utils

import "time"

func GetHourlyTime() string {
	return time.Now().Format("2006-01-02_15")
}
