package common

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/jwendel/downsampling/core"
)

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

func SavePointsToCSV(file string, points []core.Point[float64, float64]) {
	fp, err := os.Create(file)
	CheckError("Cannot create file", err)
	defer CheckError("Failed to close file", fp.Close())

	writer := csv.NewWriter(fp)
	defer writer.Flush()

	for _, point := range points {
		x := fmt.Sprintf("%f", point.X)
		y := fmt.Sprintf("%f", point.Y)
		err := writer.Write([]string{x, y})
		CheckError("Cannot write to file", err)
	}
}
