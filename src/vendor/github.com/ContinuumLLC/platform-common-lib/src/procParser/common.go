package procParser

import (
	"math"
	"strconv"
)

const (
	//Memory conversion from kB to mB to gB
	gBTokB = 1024 * 1024 * 1024
	mBTokB = 1024 * 1024
)

//GetBytes converts kB,mB,gB values to bytes
func GetBytes(size int64, measure string) int64 {
	switch measure {
	case "kB":
		size = size * 1024
	case "mB":
		size = size * mBTokB
	case "gB":
		size = size * gBTokB
	default:
		//TODO: log this errors.New(INVALIDMEMORYMEASURE + measure)
		return 0
	}
	return size
}

//GetInt64 returns an Int64 value of the input string
func GetInt64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}

//GetFormattedPercentUptoTwoDecimals returns an float64 with two decimals
func GetFormattedPercentUptoTwoDecimals(value float64) float64 {
	return (math.Floor(100*value) / 100)
}

//GetUint64 returns uint64 value of the input string
func GetUint64(val string) (uint64, error) {
	return strconv.ParseUint(val, 10, 64)
}
