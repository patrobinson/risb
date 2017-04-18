package main

import (
	"github.com/strava/go.strava"
	"os"
	"strconv"
	log "github.com/sirupsen/logrus"
	"time"
)

func pull() {
	log.Info("Pulling...\n")
	accessToken := os.Getenv("STRAVA_ACCESS_TOKEN")
	client := strava.NewClient(accessToken)
	cAthlete := strava.NewCurrentAthleteService(client)
	page := 1
	for {
		summaries, err := cAthlete.ListActivities().Page(page).PerPage(200).Do()
		if err != nil {
			log.Fatal("Failed to get activities %v", err)
		}
		if len(summaries) == 0 {
			log.Info("End of activities reached at page %v", page)
			break
		}
		page++
		streams := strava.NewActivityStreamsService(client)
		for _, s := range summaries {
			err = getActivityStream(*s, streams)
			if err != nil {
				log.Info("Failed to get Activity %v: %v", s.Id, err)
			}
		}
	}
}

func getActivityStream(activity strava.ActivitySummary, streams *strava.ActivityStreamsService) error {
	streamTypes := []strava.StreamType{
		"time",
		"distances",
		"altitude",
		"heartrate",
		"moving",
		"grade_smooth",
		"velocity_smooth",
	}
	streamSets, err := streams.Get(activity.Id, streamTypes).Do()
	if err != nil {
		return err
	}
	var previousTime int
	previousTime = 0
	for i, t := range streamSets.Time.RawData {
		elapsed := (*t - previousTime)
		data := dataPoint{
			name: "run",
			tags: map[string]string{
				"Id": strconv.FormatInt(activity.Id, 10),
				"Athlete": strconv.FormatInt(activity.Athlete.AthleteMeta.Id, 10),
				"MovingTime": strconv.Itoa(activity.MovingTime),
				"ElapsedTime": strconv.Itoa(activity.ElapsedTime),
			},
			timestamp: activity.StartDate.Add(time.Second * time.Duration(*t)),
			precision: elapsed,
			values: map[string]float64{
				"Time": float64(*t),
			},
		}
		if distance := streamSets.Distance.RawData[i]; distance != nil {
			data.values["Distance"] = *distance
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
	return nil
}
