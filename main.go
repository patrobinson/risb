package main

import (
	"time"
	log "github.com/sirupsen/logrus"
)

type run struct {
	values []*dataPoint
}

type dataPoint struct {
	tags map[string]string
	timestamp time.Time
	precision int // how precise the time series data is
	name string
	values map[string]float64
}

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Info("Starting...\n")
	setupInflux()
	lastRecordTime, err := lastRecord()
	if err != nil {
		log.Fatal("Failed to get a latest record from Influx: %v", err)
	}
	log.Info("Found the latest record retrieved:\n%v", lastRecordTime)
	pull(lastRecordTime)
}
