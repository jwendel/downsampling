package core_test

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/jwendel/downsampling/core"
)

func CheckError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func LoadPointsFromCSV(file string) []core.Point[float64, float64] {
	csvFile, err := os.Open(file)
	CheckError("Cannot Open the file.", err)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var data []core.Point[float64, float64]
	for {
		line, err2 := reader.Read()
		if err2 == io.EOF {
			break
		}
		CheckError("Read file error", err2)
		var d core.Point[float64, float64]
		d.X, _ = strconv.ParseFloat(line[0], 64)
		d.Y, _ = strconv.ParseFloat(line[1], 64)
		data = append(data, d)
	}
	return data
}

func BenchmarkLTD(b *testing.B) {
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
