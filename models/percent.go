package models

// Percent returns used as a percentage of total.
//
// A zero total means the metric has not reported a capacity yet, which reads as
// 0 rather than NaN so callers do not have to special-case it.
func Percent(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(used) / float64(total) * 100
}
