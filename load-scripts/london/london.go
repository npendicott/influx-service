package main

import (
	// TODO: Go Path
	"influx-client-london/data/csvReader"
	_"influx-client-london/data/schemas/dailyReading"
	"influx-client-london/data/schemas/halfhourlyReading"
	
	"github.com/influxdata/influxdb1-client/v2"
	// "github.com/influxdata/influxdb/client/v2"

	"log"
	_"strings"
	"strconv"
	"fmt"
	
)

const (
	DAILY_TIMESTAMP_FORMAT = "2006-01-02"
	HH_TIMESTAMP_FORMAT = "2006-01-02 15:04:05.0000000"
)


func main() {
	// Influx
	// TODO: new client ugh
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://127.0.0.1:8086",
		// Addr: "http://influxdb:8086",
		// Addr: "http://192.168.0.24:8086",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer influxClient.Close()

	// Meta
	londonPathRoot := "../smart-meters-archive/"
	// ACORN
	metaTable := csvReader.GetDataTable(londonPathRoot + "informations_households.csv")
	_ = metaTable 

	// blockMap := make(map[string][][]string)
	// _ = blockMap


	// Daily
	// frequencyLevel := "daily_dataset/"

	// for blockIndex := 0; blockIndex < 3; blockIndex++ {  // lol maybe they do start at 1 sometimes
	// 	blockData := csvReader.GetDataArray(londonPathRoot + frequencyLevel + "block_" + strconv.Itoa(blockIndex) +".csv")

	// 	batch := 0
	// 	var readingBatch []dailyReading.DailyReading
	// 	for i, reading := range blockData {
	// 		fmt.Println("Block:", blockIndex)
	// 		fmt.Println("Reading:", i)

	// 		processedReading := dailyReading.ProcessCSVEnergyReading(reading, DAILY_TIMESTAMP_FORMAT)

	// 		readingBatch = append(readingBatch, processedReading)
	// 		batch++

	// 		// Block write logic
	// 		if batch == 4000 || i == len(blockData) - 1 {
	// 			dailyReading.WriteEnergyReadingBatch(influxClient, readingBatch)
	// 			batch = 0
	// 			// TODO: Better way to clear this slice? Maybe preserve the space
	// 			readingBatch = nil
	// 		}

	// 		fmt.Println()
	// 	}
	// }


	// HalfHour
	frequencyLevel := "halfhourly_dataset/"

	for blockIndex := 0; blockIndex < 3; blockIndex++ {
		blockData := csvReader.GetDataArray(londonPathRoot + frequencyLevel + "block_" + strconv.Itoa(blockIndex) +".csv")

		batch := 0
		var readingBatch []halfhourlyReading.HalfhourlyReading
		for i, reading := range blockData {
			fmt.Println("Block:", blockIndex)
			fmt.Println("Reading:", i)

			processedReading := halfhourlyReading.ProcessCSVEnergyReading(reading, HH_TIMESTAMP_FORMAT)

			readingBatch = append(readingBatch, processedReading)
			batch++

			// Block write logic
			if batch == 4000 || i == len(blockData) - 1 {
				halfhourlyReading.WriteEnergyReadingBatch(influxClient, readingBatch)
				batch = 0
				// TODO: Better way to clear this slice? Maybe preserve the space
				readingBatch = nil
			}

			fmt.Println()
		}
	}

}
