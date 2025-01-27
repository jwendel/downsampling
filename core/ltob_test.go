package core_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jwendel/downsampling/core"
)

func BenchmarkLTOB(b *testing.B) {
	dir, err := os.Getwd()
	if err != nil {
		b.Fatal(err)
	}
	source := filepath.Join(dir, "..", "demo", "data", "source.csv")

	const sampledCount = 500
	rawdata := LoadPointsFromCSV(source)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.LTOB(rawdata, sampledCount)
	}
}
