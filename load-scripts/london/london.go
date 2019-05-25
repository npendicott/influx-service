package main

import (
	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/npendicott/influx-service/csvReader"
	"github.com/npendicott/influx-service/influx"
	"github.com/npendicott/influx-service/schemas/london/daily"

	_ "influx-client-london/data/schemas/dailyReading"

	"fmt"
	"log"
	"strconv"
	_ "strings"
)

const (
	DAILY_TIMESTAMP_FORMAT = "2006-01-02"
	HH_TIMESTAMP_FORMAT    = "2006-01-02 15:04:05.0000000"
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
	londonPathRoot := "../../data/smart-meters-archive/"

	// ACORN
	// metaTable := csvReader.GetDataTable(londonPathRoot + "informations_households.csv")
	// _ = metaTable

	// blockMap := make(map[string][][]string)
	// _ = blockMap

	// Daily
	frequencyLevel := "daily_dataset/"

	readingBatch := make([]*client.Point, 4000)
	for blockIndex := 0; blockIndex < 3; blockIndex++ { // lol maybe they do start at 1 sometimes
		blockData := csvReader.GetDataArray(londonPathRoot + frequencyLevel + "block_" + strconv.Itoa(blockIndex) + ".csv")

		// batch := 0
		for i, reading := range blockData {
			fmt.Println("Block:", blockIndex)
			fmt.Println("Reading:", i)

			dailyReading := daily.ProcessCSVEnergyReading(reading)
			dailyPoint, err := influx.Marshall(dailyReading, "daily")

			if err != nil {
				panic(err)
			}
			readingBatch = append(readingBatch, dailyPoint)
			// control.Print(controlReading)

			if len(readingBatch)%4000 == 0 || i == len(readingBatch)-1 {
				fmt.Println(len(readingBatch))
				influx.WritePointBatch(influxClient, readingBatch)
				readingBatch = make([]*client.Point, 4000)
			}

			// readingBatch = append(readingBatch, processedReading)
			// batch++

			// // Block write logic
			// if batch == 4000 || i == len(blockData)-1 {
			// 	dailyReading.WriteEnergyReadingBatch(influxClient, readingBatch)
			// 	batch = 0
			// 	// TODO: Better way to clear this slice? Maybe preserve the space
			// 	readingBatch = nil
			// }

			fmt.Println()
		}
	}

	// // HalfHour
	// frequencyLevel := "halfhourly_dataset/"

	// for blockIndex := 0; blockIndex < 3; blockIndex++ {
	// 	blockData := csvReader.GetDataArray(londonPathRoot + frequencyLevel + "block_" + strconv.Itoa(blockIndex) + ".csv")

	// 	batch := 0
	// 	var readingBatch []halfhourlyReading.HalfhourlyReading
	// 	for i, reading := range blockData {
	// 		fmt.Println("Block:", blockIndex)
	// 		fmt.Println("Reading:", i)

	// 		processedReading := halfhourlyReading.ProcessCSVEnergyReading(reading, HH_TIMESTAMP_FORMAT)

	// 		readingBatch = append(readingBatch, processedReading)
	// 		batch++

	// 		// Block write logic
	// 		if batch == 4000 || i == len(blockData)-1 {
	// 			halfhourlyReading.WriteEnergyReadingBatch(influxClient, readingBatch)
	// 			batch = 0
	// 			// TODO: Better way to clear this slice? Maybe preserve the space
	// 			readingBatch = nil
	// 		}

	// 		fmt.Println()
	// 	}
	// }

}
