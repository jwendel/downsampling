package core

import (
	"math"
)

func calculateTriangleArea[X, Y Number](pa Point[X, Y], pbX, pbY float64, pc Point[X, Y]) float64 {
	area := (float64(pa.X-pc.X)*(pbY-float64(pa.Y)) - (float64(pa.X)-pbX)*float64(pc.Y-pa.Y)) * 0.5
	return math.Abs(area)
}

func calculateAverageDataPoint[X, Y Number](points []Point[X, Y]) (float64, float64) {
	x, y := 0.0, 0.0
	for _, point := range points {
		x += float64(point.X)
		y += float64(point.Y)
	}
	l := float64(len(points))
	return x / l, y / l
}

func splitDataBucket[X, Y Number](data []Point[X, Y], threshold int) [][]Point[X, Y] {

	buckets := make([][]Point[X, Y], threshold)
	for i := range buckets {
		buckets[i] = make([]Point[X, Y], 0)
	}
	// First and last bucket are formed by the first and the last data points
	buckets[0] = append(buckets[0], data[0])
	buckets[threshold-1] = append(buckets[threshold-1], data[len(data)-1])

	// so we only have N - 2 buckets left to fill
	bucketSize := float64(len(data)-2) / float64(threshold-2)

	//slice remove the first and last point
	d := data[1 : len(data)-1]

	for i := 0; i < threshold-2; i++ {
		bucketStartIdx := int(math.Floor(float64(i) * bucketSize))
		bucketEndIdx := int(math.Floor(float64(i+1)*bucketSize)) + 1
		if i == threshold-3 {
			bucketEndIdx = len(d)
		}
		buckets[i+1] = append(buckets[i+1], d[bucketStartIdx:bucketEndIdx]...)
	}

	return buckets
}

func calculateAveragePoint[X, Y Number](points []Point[X, Y]) (x, y float64) {
	l := len(points)
	var p Point[X, Y]
	for i := 0; i < l; i++ {
		p.X += points[i].X
		p.Y += points[i].Y
	}

	return float64(p.X) / float64(l), float64(p.Y) / float64(l)
}
