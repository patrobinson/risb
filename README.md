# Running is Beautiful

Strava gives a nice amount of data and while [Stravistix](https://chrome.google.com/webstore/detail/stravistix-for-strava/dhiaggccakkgdfcadnklkbljcgicpckn?hl=en) expands on that, it doesn't quite show me what I need. I'm too much of a tight arse to pay for Premium and I actually really enjoy building data analytics tools.

## How to use this repo

You need Docker installed to run this code. To import data into InfluxDB create a [Strava API Token](https://www.strava.com/settings/api) and export the Access Token to your environment. Then bring up the containers:

```
 export STRAVA_ACCESS_TOKEN=abcdefghijklmnopqrstuvwxyz1234567890
docker-compose up
```

This will take a while to fill in all the historical data. Unfortunately, at the moment, we trawl through every activity no matter if we have pulled them down before.

Once downloaded you can browse data via the grafana interface at http://localhost:3000
