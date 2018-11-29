package timeutil

import (
	"fmt"
	"strconv"
	"time"
)

//ToLongYYYYMMDDHH returns the time in int format as YYYYMMDDHH
func ToLongYYYYMMDDHH(tm *time.Time) (int, error) {
	return strconv.Atoi(fmt.Sprintf("%d%02d%02d%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour()))
}

//ToLongYYYYMMDD returns the time in int format as YYYYMMDD
func ToLongYYYYMMDD(tm *time.Time) (int, error) {
	return strconv.Atoi(fmt.Sprintf("%d%02d%02d", tm.Year(), tm.Month(), tm.Day()))
}

//ToHourLong returns the time array in long format for the time range given in the input
func ToHourLong(fromTime, toTime time.Time) []int {
	var tmLong []int
	toTm, _ := ToLongYYYYMMDD(&toTime)
	for tm := fromTime; ; {
		tmInt, _ := ToLongYYYYMMDD(&tm)
		if tmInt > toTm {
			break
		}
		tmLong = append(tmLong, tmInt)
		tm = tm.Add(24 * time.Hour)

	}
	return tmLong
}
