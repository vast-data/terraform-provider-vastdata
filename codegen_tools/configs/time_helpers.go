package configs

import "time"

func MinToDuration(m int64) *time.Duration {
	min := time.Duration(m) * time.Minute / time.Microsecond
	return &min
}

func SecondToDuration(s int64) *time.Duration {
	sec := time.Duration(s) * time.Second / time.Microsecond
	return &sec
}

func HoursToDuration(h int64) *time.Duration {
	hour := (time.Duration(h) * time.Second) / time.Microsecond
	return &hour
}
