package main

import (
	// TODO: Remove
	"errors"

	"github.com/npendicott/influx-service/influx"
	_ "github.com/npendicott/influx-service/schemas/london/daily"
	_ "github.com/npendicott/influx-service/schemas/london/haflhourlyReading"

	client "github.com/influxdata/influxdb1-client/v2" //"github.com/influxdata/influxdb/client/v2"
	"github.com/joho/godotenv"

	// Net stuff
	"encoding/json"
	"net/http"

	// OS stuff
	"fmt"
	"log"
	"os"
	"regexp"

	// TODO: RM
	_ "github.com/influxdata/influxdb1-client/v2"
	//_"github.com/influxdata/influxdb/client/v2"
	_ "strconv"
	_ "strings"
)

const (
	PORT = ":9090" // TODO: Remove
)

var validPath = regexp.MustCompile("^/(readings|reading)/([a-zA-Z0-9]+)$")
var influxClient client.Client

// Utility
func addCrossSiteOriginHeader(resWrt http.ResponseWriter) http.ResponseWriter {
	resWrt.Header().Set("Access-Control-Allow-Origin", "*")
	// These all have to be lowercase, this might be bad HTTP
	// TODO: maybe an "expectedTokens" slice?
	resWrt.Header().Set("Access-Control-Allow-Headers", "mac-id")

	return resWrt
}

func getSeries(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Series")
	}
	return m[2], nil // The title is the second subexpression.
}

// So, not usre if this is the best way, but seems better than copying for each.
// I guess these should just be headers?? I don't think so though
func getRequiredParam(req *http.Request, key string) (string, []byte) {
	result, ok := req.URL.Query()[key]
	if !ok || len(result[0]) < 1 {
		errString := fmt.Sprintf("One %s param is required.\n", key)
		errOut := []byte(errString)
		log.Println(errString)

		return "", errOut
	}

	return result[0], nil
}

// Routs
func root(resWrt http.ResponseWriter, req *http.Request) {
	out := []byte("Root!")
	resWrt.Write(out)

}

func readings(resWrt http.ResponseWriter, req *http.Request) {
	series := req.URL.Path[len("/readings/"):] // TODO: Is this really the only way to do this?

	// Headers
	macID := req.Header.Get("mac-id")

	// Get date range
	start, errOut := getRequiredParam(req, "start")
	if errOut != nil {
		resWrt.Write(errOut)
		return
	}
	fmt.Println(start)

	end, errOut := getRequiredParam(req, "end")
	if errOut != nil {
		resWrt.Write(errOut)
		return
	}
	fmt.Println(end)

	// Headers
	var resp *client.Response
	if macID == "" {
		resp = influx.QueryDateRange(influxClient, series, start, end)
	} else {
		fmt.Println(macID)
		resp = influx.QueryDateRangeWithMAC(influxClient, series, macID, start, end)
	}

	// Get response

	// Writing Response Body
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Fatal(err)
	}

	// Header
	resWrt = addCrossSiteOriginHeader(resWrt)

	resWrt.Write(out)
	return
}

// Server
func main() {
	// ENVs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to InfluxDB
	InfluxAddress := os.Getenv("INFLUX_ADDRESS")
	fmt.Println(InfluxAddress)

	influxClient = influx.GetConnection()
	defer influxClient.Close()

	// Routes
	http.HandleFunc("/", root)
	http.HandleFunc("/readings/", readings)
	// http.HandleFunc("/mac_readings/", mac_readings)
	// http.HandleFunc("/readings/daily", readingsDaily)
	// http.HandleFunc("/readings/hhourly", readingsHHourly)

	// Start server
	fmt.Println("Listeninggggggg")
	serveErr := http.ListenAndServe(PORT, nil)
	if serveErr != nil {
		log.Fatal("ListenAndServe: ", serveErr)
	}
}
