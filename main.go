package main

import (
	"time"
	"fmt"
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

func main() {
	fmt.Printf("Starting...\n")
	pull()
}