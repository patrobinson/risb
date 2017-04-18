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
	pull()
}