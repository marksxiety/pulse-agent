package models

const HistoryCapacity = 2160

type MetricHistory struct {
	buf   []float64
	head  int
	count int
}

func NewMetricHistory() *MetricHistory {
	return &MetricHistory{
		buf: make([]float64, HistoryCapacity),
	}
}

func (h *MetricHistory) Push(value float64) {
	h.buf[h.head] = value
	h.head = (h.head + 1) % len(h.buf)
	if h.count < len(h.buf) {
		h.count++
	}
}

func (h *MetricHistory) Len() int {
	return h.count
}

func (h *MetricHistory) get(i int) float64 {
	if h.count < len(h.buf) {
		return h.buf[i]
	}
	return h.buf[(h.head+i)%len(h.buf)]
}

func (h *MetricHistory) Min() float64 {
	if h.count == 0 {
		return 0
	}
	m := h.get(0)
	for i := 1; i < h.count; i++ {
		if v := h.get(i); v < m {
			m = v
		}
	}
	return m
}

func (h *MetricHistory) Max() float64 {
	if h.count == 0 {
		return 0
	}
	m := h.get(0)
	for i := 1; i < h.count; i++ {
		if v := h.get(i); v > m {
			m = v
		}
	}
	return m
}

func (h *MetricHistory) Avg() float64 {
	if h.count == 0 {
		return 0
	}
	sum := 0.0
	for i := 0; i < h.count; i++ {
		sum += h.get(i)
	}
	return sum / float64(h.count)
}

func (h *MetricHistory) Downsample(n int) []float64 {
	if h.count == 0 {
		return nil
	}

	data := make([]float64, h.count)
	for i := 0; i < h.count; i++ {
		data[i] = h.get(i)
	}

	if h.count <= n {
		return data
	}

	result := make([]float64, n)
	bucketSize := float64(h.count) / float64(n)

	for i := 0; i < n; i++ {
		start := int(float64(i) * bucketSize)
		end := int(float64(i+1) * bucketSize)
		if end > h.count {
			end = h.count
		}

		sum := 0.0
		for j := start; j < end; j++ {
			sum += data[j]
		}
		result[i] = sum / float64(end-start)
	}

	return result
}
