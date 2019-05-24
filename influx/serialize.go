package influx

import (
	client "github.com/influxdata/influxdb1-client/v2"

	"fmt"
	"log"
	"reflect"
	"time"
)

// Marshall an object, returns a Point. Uses influx to tell what is what, json as labels
func Marshall(structure interface{}) (*client.Point, error) {
	// Schema
	timestamp := time.Time{}

	tags := map[string]string{
		//"mac_id": reading.MacId,
	}

	fields := map[string]interface{}{
		// "volt_al": reading.VoltAL,
	}

	// Reflection
	schema := reflect.TypeOf(structure)
	values := reflect.ValueOf(structure)

	fmt.Println(schema)

	// Iterate through fields
	for i := 0; i < schema.NumField(); i++ {
		field := schema.Field(i)
		value := values.Field(i)

		inftype, ok := field.Tag.Lookup("influx")
		if !ok {
			panic(ok)
		}

		label, ok := field.Tag.Lookup("json")
		if !ok {
			panic(ok)
		}

		switch inftype {
		case "timestamp":
			// https://stackoverflow.com/questions/17262238/how-to-cast-reflect-value-to-its-type
			timestamp = value.Interface().(time.Time)
		case "tag": // Just dump to string
			tags[label] = value.String()
		case "field": // Must be some number
			fields[label] = value
		}

		fmt.Println("Influx: ", inftype)
		fmt.Println("Label: ", label)
		fmt.Println("Val: ", value)
		fmt.Println()

	}

	// Create point
	pt, err := client.NewPoint(
		"Table",
		tags,
		fields,
		timestamp,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pt)
	fmt.Println()
	return pt, nil

}

// Takes a Point, and an object contents, and adds things to the treee
func Unmarshall() {

}
