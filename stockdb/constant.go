package stockdb

const (
	Second  int64 = 1
	Minute  int64 = 60 * Second
	Hour    int64 = 60 * Minute
	Day     int64 = 24 * Hour
	Week    int64 = 7 * Day
	Month   int64 = 30 * Day
	Quarter int64 = 3 * Month
	Year    int64 = 365 * Day
)
