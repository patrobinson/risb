package main

import (
	"github.com/strava/go.strava"
	"os"
	"strconv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func pull() error {
	streamTypes := []strava.StreamType{
		"time",
		"distances",
		"altitude",
		"heartrate",
		"moving",
		"grade_smooth",
		"velocity_smooth",
	}
	fmt.Printf("Pulling...\n")
	accessToken := os.Getenv("STRAVA_ACCESS_TOKEN")
	client := strava.NewClient(accessToken)
	cAthlete := strava.NewCurrentAthleteService(client)
	summaries, err := cAthlete.ListActivities().Do()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to get activities %v", err))
	}
	streams := strava.NewActivityStreamsService(client)
	for _, s := range summaries {
		streamSets, err := streams.Get(s.Id, streamTypes).Do()
		if err != nil {
			fmt.Printf("Failed to get Activity %v: %v", s.Id, err)
			continue
		}
		var previousTime int
		previousTime = 0
		for i, t := range streamSets.Time.RawData {
			elapsed := (*t - previousTime)
			data := dataPoint{
				name: "run",
				tags: map[string]string{
					"Id": strconv.FormatInt(s.Id, 10),
					"Athlete": strconv.FormatInt(s.Athlete.AthleteMeta.Id, 10),
					"MovingTime": strconv.Itoa(s.MovingTime),
					"ElapsedTime": strconv.Itoa(s.ElapsedTime),
				},
				timestamp: s.StartDate.Add(time.Second * time.Duration(*t)),
				precision: elapsed,
				values: map[string]float64{
					"Time": float64(*t),
				},
			}
			if distance := streamSets.Distance.RawData[i]; distance != nil {
				data.values["Distance"] = *distance
				data.values["MinutesPerKM"] = (((1000 / *distance) * float64(elapsed)) / 60)
			}
			if speed := streamSets.Speed.RawData[i]; speed != nil {
				data.values["MetersPerSecond"] = *speed
			}
			if altitude := streamSets.Elevation.RawData[i]; altitude != nil {
				data.values["Altitude"] = *altitude
			}
			if heartrate := streamSets.HeartRate.RawData[i]; heartrate != nil {
				data.values["HeartRate"] = float64(*heartrate)
			}
			if grade := streamSets.Grade.RawData[i]; grade != nil {
				data.values["Grade"] = *grade
			}
			data.tags["Moving"] = strconv.FormatBool(streamSets.Moving.Data[i])
			previousTime = *t
			sink(data)
		}
	}
	return nil
}
