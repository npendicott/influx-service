package influx

import (
	client "github.com/influxdata/influxdb1-client/v2" //"github.com/influxdata/influxdb/client/v2"

	"fmt"
	"log"
)

const (
	// DatabaseName is the influxdb being targeted
	DatabaseName = "energydb"
	// Verbose allows for busy logs
	Verbose = false
)

var (
	// InfluxClient is the client lib for talking to Influx.
	// InfluxClient client.Client

	// InfluxAddress is the IP address of the Influx DB.
	InfluxAddress string

	bpConfig = client.BatchPointsConfig{
		Database: DatabaseName,
		// Precision: "s",
	}
)

// GetConnection creates an InfluxClient based on env INFLUX_ADDRESS. The returned connection must be closed
// TODO: store connection externally?
func GetConnection() client.Client {
	// InfluxAddress = os.Getenv("INFLUX_ADDRESS")
	InfluxAddress = "http://localhost:8086"

	// Influx
	//https://stackoverflow.com/questions/31218008/sharing-a-globally-defined-db-conn-with-multiple-packages-in-golang
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: InfluxAddress,
	})
	if err != nil {
		log.Fatal(err)
	}

	return influxClient
}

// WritePoint takes a Point and writes it to the provided connection.
func WritePoint(clnt client.Client, pt *client.Point) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(bpConfig)
	if err != nil {
		log.Fatal(err)
	}
	if Verbose {
		fmt.Println("Created Batchpoint")
	}

	// Add to batch
	bp.AddPoint(pt)
	if Verbose {
		fmt.Println("Added point to table")
		fmt.Println()
	}

	// Write
	if Verbose {
		fmt.Println("Write Batch")
	}
	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}
}

// WritePoints takes a slice of points and
func WritePoints(clnt client.Client, pts []*client.Point) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(bpConfig)
	if err != nil {
		log.Fatal(err)
	}

	for i, pt := range pts {
		// Add to batch
		bp.AddPoint(pt)

		if i%4000 == 0 || i == len(pts)-1 { // TODO: pull out batch size
			if err := clnt.Write(bp); err != nil {
				log.Fatal(err)
			}

			// CLEAR BP?
			bp, err = client.NewBatchPoints(bpConfig)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// WritePointBatch writes a single batch based on an array of points.
// It was initially used in place of WritePoints, but requres iteration over the overall input slice twice.
// NOTE: the WrotePoints assumes all points are read into memory tho, so maybe this is the best way to go about?
func WritePointBatch(clnt client.Client, pts []*client.Point) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(bpConfig)
	if err != nil {
		log.Fatal(err)
	}
	if Verbose {
		fmt.Println("Created Batchpoint")
	}

	for _, pt := range pts {
		// Add to batch
		bp.AddPoint(pt)

		if Verbose {
			fmt.Println("Added point to table")
			fmt.Println()
		}
	}

	// Write
	if Verbose {
		fmt.Println("Write Batch")
	}
	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}
}

// QueryDateRange gets a response (Influx client official lib) from the current connection, based on a series (table) arugment, ID string, and date range.
// Two examples:
// SELECT * FROM daily_energy_readings WHERE time >= '2012-04-12 10:30:00.0000000' and time < '2012-05-12 10:30:00.0000000'
// SELECT * FROM daily_energy_readings WHERE time >= 1392249600000000000 and time < 1393545600000000000
// TODO: return pointer or contents?
// TODO: These start/ends need to be time types
func QueryDateRange(influxClient client.Client, series string, macID string, startDate string, endDate string) *client.Response {
	queryString := fmt.Sprintf("SELECT * FROM %s WHERE time >= '%s' AND time < '%s' AND mac_id = '%s'", series, startDate, endDate, macID)

	query := client.Query{
		Command:  queryString,
		Database: DatabaseName,
	}

	resp, err := influxClient.Query(query)
	if err != nil {
		log.Println(err)
		fmt.Println("Probably: could not connect to", InfluxAddress)
		// log.Fatal(err)
	}

	return resp
}
