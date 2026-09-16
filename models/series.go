package models

// Series is a display-ready summary of one metric's rolling history.
type Series struct {
	Percent float64   `json:"percent"`
	Trend   []float64 `json:"trend"`
	Max     float64   `json:"max"`
	Avg     float64   `json:"avg"`
}

// NewSeries summarises h, downsampling the trend to at most points samples.
func NewSeries(h *MetricHistory, current float64, points int) Series {
	return Series{
		Percent: current,
		Trend:   h.Downsample(points),
		Max:     h.Max(),
		Avg:     h.Avg(),
	}
}
