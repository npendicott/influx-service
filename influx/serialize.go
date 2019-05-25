package influx

import (
	client "github.com/influxdata/influxdb1-client/v2"

	"log"
	"reflect"
	"time"
)

// Marshall an object, returns a Point. Uses influx to tell what is what, json as labels
func Marshall(structure interface{}, table string) (*client.Point, error) {
	// Point structures
	timestamp := time.Time{}

	tags := map[string]string{
		//"mac_id": reading.MacId,
	}

	fields := map[string]interface{}{
		// "volt_al": reading.VoltAL,
	}

	// Reflect
	schema := reflect.TypeOf(structure)
	values := reflect.ValueOf(structure)

	// Iterate through fields
	for i := 0; i < schema.NumField(); i++ {
		field := schema.Field(i)
		value := values.Field(i)

		// Get Influx type
		inftype, ok := field.Tag.Lookup("influx")
		if !ok {
			panic(ok)
		}

		// Get label
		label, ok := field.Tag.Lookup("json")
		if !ok {
			panic(ok)
		}

		// Add to correct slice
		switch inftype {
		case "timestamp":
			timestamp = value.Interface().(time.Time) // https://stackoverflow.com/questions/17262238/how-to-cast-reflect-value-to-its-type
		case "tag": // Just dump to string
			tags[label] = value.String()
		case "field": // Must be some number
			fields[label] = value
		default:
			break
		}
	}

	// Create point
	pt, err := client.NewPoint(
		table,
		tags,
		fields,
		timestamp,
	)
	if err != nil {
		log.Fatal(err)
	}

	return pt, nil
}

// Takes a Point, and an object contents, and adds things to the treee
func Unmarshall() {

}
