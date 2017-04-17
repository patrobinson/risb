package main

import (
	"github.com/influxdata/influxdb/client/v2"
	"fmt"
	"os"
	log "github.com/sirupsen/logrus"
)

func sink(rawData dataPoint) error {
	iAddr := "http://" + os.Getenv("INFLUX_HOSTNAME") + ":" + os.Getenv("INFLUX_PORT")
	conf := client.HTTPConfig{
    	Addr: iAddr,
	}

	c, err := client.NewHTTPClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	precision, err := precisionToString(rawData.precision)
	if err != nil {
		return err
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  os.Getenv("DB_NAME"),
		Precision: precision,
	})
	if err != nil {
		log.Fatal(err)
	}

	fields := make(map[string]interface{})

	for key, value := range rawData.values {        
	    fields[key] = value
	}

	pt, err := client.NewPoint(
			rawData.name,
			rawData.tags,
			fields,
			rawData.timestamp,
		)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
	return nil
}

func precisionToString(precision int) (string, error) {
	switch precision {
	case 1:
		return "s", nil
	case 60:
		return "m", nil
	case 3600:
		return "h", nil
	case 86400:
		return "d", nil
	default:
		return "", fmt.Errorf("Invalid precision for InfluxDB %d", precision)
	}
}