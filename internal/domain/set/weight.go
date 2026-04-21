package set

import "fmt"

type Weight float64

func (w Weight) Float64() float64 {
	return float64(w)
}

func NewWeight(w float64) (Weight, error) {
	if w < 0 {
		return Weight(0), fmt.Errorf("invalid weight: %f", w)
	}
	return Weight(w), nil
}
