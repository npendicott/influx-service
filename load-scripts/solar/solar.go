package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2" //"github.com/influxdata/influxdb/client/v2"

	"github.com/npendicott/influx-service/csvReader"
	"github.com/npendicott/influx-service/influx"
	"github.com/npendicott/influx-service/schemas/solar/control"
	"github.com/npendicott/influx-service/schemas/solar/pando"
	_ "github.com/npendicott/influx-service/schemas/solar/pando"

	_ "strconv"
	_ "strings"
)

func main() {
	// // ENVs
	// err := godotenv.Load("../..")  // TODO: ENV in the folder dangus
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// Influx
	influxClient := influx.GetConnection()
	defer influxClient.Close()
	pointBatch := make([]*client.Point, 4000) // Batch

	// Control
	controlData := csvReader.GetDataArray("../../data/solar/control.csv")

	for i, reading := range controlData {
		fmt.Println("Reading:", i)

		reading := control.ProcessCSVReading(reading)
		fmt.Println("Time", reading.Timestamp)

		point, err := influx.Marshall(reading, "control")
		if err != nil {
			panic(err)
		}

		pointBatch = append(pointBatch, point)
		if len(pointBatch)%4000 == 0 || i == len(pointBatch)-1 {
			fmt.Println(len(pointBatch))
			influx.WritePointBatch(influxClient, pointBatch)
			pointBatch = make([]*client.Point, 4000)
		}

		fmt.Println()
	}

	// PandO
	pandoData := csvReader.GetDataArray("../../data/solar/pando_2k.csv")

	for i, reading := range pandoData {
		fmt.Println("Reading:", i)

		reading := pando.ProcessCSVReading(reading)
		fmt.Println("Time", reading.Timestamp)

		point, err := influx.Marshall(reading, "pando_2k")
		if err != nil {
			panic(err)
		}

		pointBatch = append(pointBatch, point)
		if len(pointBatch)%4000 == 0 || i == len(pointBatch)-1 {
			fmt.Println(len(pointBatch))
			influx.WritePointBatch(influxClient, pointBatch)
			pointBatch = make([]*client.Point, 4000)
		}

		fmt.Println()
	}
}
