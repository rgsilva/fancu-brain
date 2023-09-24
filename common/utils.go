package common

import (
	"image"
	"math"
	"sort"
)

func PointDistance(p1 image.Point, p2 image.Point) float64 {
	a := float64(p1.X - p2.X)
	b := float64(p1.Y - p2.Y)
	return math.Sqrt((a * a) + (b * b))
}

func SortFoundObjects[T any](fms []T, sorters []func(left, right T) bool) {
	for _, sorter := range sorters {
		sort.SliceStable(fms, func(i, j int) bool {
			return sorter(fms[i], fms[j])
		})
	}
}

func IntAbs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
