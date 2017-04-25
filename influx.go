package main

import (
	"github.com/influxdata/influxdb/client/v2"
	"fmt"
	"os"
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"strconv"
)

var iClient client.Client

func setupInflux() {
	iAddr := "http://" + os.Getenv("INFLUX_HOSTNAME") + ":" + os.Getenv("INFLUX_PORT")
	conf := client.HTTPConfig{
    	Addr: iAddr,
	}

	var err error
	iClient, err = client.NewHTTPClient(conf)
	if err != nil {
		log.Fatal(err)
	}
}

func lastRecord() (int64, error) {
	q := client.NewQuery("SELECT LAST(Time) FROM Run", "risb", "s")

	// Set the last record to The Epoch by default
    emptyResponse := int64(0)

	if response, err := iClient.Query(q); err == nil {
		if response.Error() != nil {
			return emptyResponse, response.Error()
		}
		if len(response.Results) > 0 && len(response.Results[0].Series) > 0 && len(response.Results[0].Series[0].Values) > 0 {
			lastTimestamp := response.Results[0].Series[0].Values[0][0].(json.Number)
			i64, _ := strconv.ParseInt(string(lastTimestamp), 10, 64)
			return i64, nil
		}
		// Fall through to the default
	} else {
		return emptyResponse, err
	}

	return emptyResponse, nil
}

func extractUniqueIdTags(results []client.Result) []string {
	var uniqueIds []string
	for _, result := range results {
		for _, row := range result.Series {
			uniqueIds = append(uniqueIds, row.Tags["Id"])
		}
	}
	return uniqueIds
}

func sink(rawData dataPoint) error {
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
	if err := iClient.Write(bp); err != nil {
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