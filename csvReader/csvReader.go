package csvReader

import (
	"encoding/csv"
	"log"
	"os"
	"io"
)

func GetDataTable(path string)  map[string][]string {
	var readings = make(map[string][]string)

	// Open the file, get a reader
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// Head
	reader.Read()

	// Body
	for {
		reading, err := reader.Read()
		if err == io.EOF {
			break
		}

		readings[reading[0]] = reading[1:5]
		// readings = append(readings, reading)
	}

	return readings
}


func GetDataArray(path string) [][]string {
	// Output
	var readings [][]string

	// Open the file, get a reader
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// Head
	reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		readings = append(readings, record)
	}

	return readings
}
