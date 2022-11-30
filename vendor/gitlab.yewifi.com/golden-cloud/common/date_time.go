package common

import "time"

// GenStampSection return stamp section [begin, end] of specified month
func GenStampSection(year, month int32) (int64, int64) {
	secL := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.Local)
	secR := time.Date(int(year), time.Month(month+1), 1, 0, 0, 0, 0, time.Local)

	return secL.Unix(), secR.Unix()
}
