package core

import (
	"math"
)

func calculateLinearRegressionCoefficients[X, Y Number](points []Point[X, Y]) (float64, float64) {

	avgX, avgY := calculateAveragePoint(points)

	aNumerator := 0.0
	aDenominator := 0.0
	for i := 0; i < len(points); i++ {
		aNumerator += (float64(points[i].X) - avgX) * (float64(points[i].Y) - avgY)
		aDenominator += (float64(points[i].X) - avgX) * (float64(points[i].X) - avgX)
	}

	a := aNumerator / aDenominator
	b := avgY - a*avgX

	return a, b
}

func calculateSSEForBucket[X, Y Number](points []Point[X, Y]) float64 {
	a, b := calculateLinearRegressionCoefficients(points)
	sumStandardErrorsSquared := 0.0
	for _, p := range points {
		standardError := float64(p.Y) - (a*float64(p.X) + b)
		sumStandardErrorsSquared += standardError * standardError
	}
	return sumStandardErrorsSquared
}

func calculateSSEForBuckets[X, Y Number](buckets [][]Point[X, Y]) []float64 {
	sse := make([]float64, len(buckets)-2)

	// We skip the first and last buckets since they only contain one data point
	for i := 1; i < len(buckets)-1; i++ {
		prevBucket := buckets[i-1]
		currBucket := buckets[i]
		nextBucket := buckets[i+1]
		// var bucketWithAdjacentPoints []Point
		// bucketWithAdjacentPoints = append(bucketWithAdjacentPoints, prevBucket[len(prevBucket)-1])
		// bucketWithAdjacentPoints = append(bucketWithAdjacentPoints, currBucket...)
		// bucketWithAdjacentPoints = append(bucketWithAdjacentPoints, nextBucket[0])
		bucketWithAdjacentPoints := make([]Point[X, Y], len(currBucket)+2)
		bucketWithAdjacentPoints[0] = prevBucket[len(prevBucket)-1]
		bucketWithAdjacentPoints[len(bucketWithAdjacentPoints)-1] = nextBucket[0]
		for i := 1; i < len(currBucket); i++ {
			bucketWithAdjacentPoints[i] = currBucket[i-1]
		}

		sse[i-1] = calculateSSEForBucket(bucketWithAdjacentPoints)
	}

	sse = append(sse, 0)
	return sse
}

func findLowestSSEAdjacentBucketsIndex(sse []float64, ignoreIndex int) int {
	minSSE := float64(math.MaxInt64)
	minSSEIndex := -1
	for i := 1; i < len(sse)-2; i++ {
		if i == ignoreIndex || i+1 == ignoreIndex {
			continue
		}

		if sse[i]+sse[i+1] < minSSE {
			minSSE = sse[i] + sse[i+1]
			minSSEIndex = i
		}
	}
	return minSSEIndex
}

func findHighestSSEBucketIndex[X, Y Number](buckets [][]Point[X, Y], sse []float64) int {
	maxSSE := 0.0
	maxSSEIdx := -1
	for i := 1; i < len(sse)-1; i++ {
		if len(buckets[i]) > 1 && sse[i] > maxSSE {
			maxSSE = sse[i]
			maxSSEIdx = i
		}
	}
	return maxSSEIdx
}

func splitBucketAt[X, Y Number](buckets [][]Point[X, Y], index int) [][]Point[X, Y] {
	if index < 0 || index >= len(buckets) {
		return buckets
	}
	bucket := buckets[index]
	bucketSize := len(bucket)
	if bucketSize < 2 {
		return buckets
	}

	bucketALength := int(math.Ceil(float64(bucketSize / 2)))
	bucketA := bucket[0 : bucketALength+1]
	bucketB := bucket[bucketALength:]

	var newBuckets [][]Point[X, Y]
	newBuckets = append(newBuckets, buckets[0:index]...)
	newBuckets = append(newBuckets, bucketA, bucketB)
	newBuckets = append(newBuckets, buckets[index+1:]...)

	return newBuckets
}

func mergeBucketAt[X, Y Number](buckets [][]Point[X, Y], index int) [][]Point[X, Y] {
	if index < 0 || index >= len(buckets)-1 {
		return buckets
	}
	mergedBucket := buckets[index]
	mergedBucket = append(mergedBucket, buckets[index+1]...)

	var newBuckets [][]Point[X, Y]
	newBuckets = append(newBuckets, buckets[0:index]...)
	newBuckets = append(newBuckets, mergedBucket)
	newBuckets = append(newBuckets, buckets[index+2:]...)

	return newBuckets
}

// LTD - Largest triangle dynamic(LTD) data downsampling algorithm implementation
//   - Require: data . The original data
//   - Require: threshold . Number of data points to be returned
func LTD[X, Y Number](data []Point[X, Y], threshold int) []Point[X, Y] {

	if threshold >= len(data) || threshold == 0 {
		return data // Nothing to do
	}

	// 1: Split the data into equal number of buckets as the threshold but have the first
	// bucket only containing the first data point and the last bucket containing only
	// the last data point . First and last buckets are then excluded in the bucket
	// resizing
	// 2: Calculate the SSE for the buckets accordingly . With one point in adjacent
	// buckets overlapping
	// 3: while halting condition is not met do . For example, using formula 4.2
	// 4: Find the bucket F with the highest SSE
	// 5: Find the pair of adjacent buckets A and B with the lowest SSE sum . The
	// pair should not contain F
	// 6: Split bucket F into roughly two equal buckets . If bucket F contains an odd
	// number of points then one bucket will contain one more point than the other
	// 7: Merge the buckets A and B
	// 8: Calculate the SSE of the newly split up and merged buckets
	// 9: end while.
	// 10: Use the Largest-Triangle-Three-Buckets algorithm on the resulting bucket configuration
	// to select one point per buckets

	//1: Split the data into equal number of buckets as the threshold.
	buckets := splitDataBucket(data, threshold)
	numIterations := len(data) * 10 / threshold
	for iter := 0; iter < numIterations; iter++ {

		// 2: Calculate the SSE for the buckets accordingly.
		sseForBuckets := calculateSSEForBuckets(buckets)

		// 4: Find the bucket F with the highest SSE
		highestSSEBucketIndex := findHighestSSEBucketIndex(buckets, sseForBuckets)
		if highestSSEBucketIndex < 0 {
			break
		}

		// 5: Find the pair of adjacent buckets A and B with the lowest SSE sum .
		lowestSSEAdjacentBucketIndex := findLowestSSEAdjacentBucketsIndex(sseForBuckets, highestSSEBucketIndex)
		if lowestSSEAdjacentBucketIndex < 0 {
			break
		}

		// 6: Split bucket F into roughly two equal buckets . If bucket F contains an odd
		// number of points then one bucket will contain one more point than the other
		buckets = splitBucketAt(buckets, highestSSEBucketIndex)

		// 7: Merge the buckets A and B
		if lowestSSEAdjacentBucketIndex > highestSSEBucketIndex {
			lowestSSEAdjacentBucketIndex++
		}
		buckets = mergeBucketAt(buckets, lowestSSEAdjacentBucketIndex)

	}
	// 10: Use the Largest-Triangle-Three-Buckets algorithm on the resulting bucket
	return LTTBForBuckets(buckets)
}
