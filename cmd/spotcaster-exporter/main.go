package main

import (
	"flag"
	"log"
	"os"
	"time"

	stdhttp "net/http"

	"github.com/codeandship/spotcaster/http"
	"github.com/codeandship/spotcaster/metrics"
	"github.com/codeandship/spotcaster/storage"
	"github.com/peterbourgon/ff"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)

	fs := flag.NewFlagSet("spotcaster-exporter", flag.ExitOnError)
	var (
		clientID = fs.String("client-id", "", "spotify client id for authentication (does not need to be set for now)")
		// debug           = fs.Bool("debug", false, "enable debug output")
		showID          = fs.String("show-id", "", "podcast show id")
		spdcCookieValue = fs.String("spdc-cookie-value", "", "sp_dc cookie value from spotify's webapp when you are logged in")
		timeSchedule    = fs.String("time-sched", "24h", "set schedule")
		tokenFileName   = fs.String("token-file", "token.json", "access token file in json format")
		listenAddr      = fs.String("listen-addr", ":8080", "set listener address")
	)
	ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithEnvVarPrefix("SPOTCASTER"),
	)

	if *clientID == "" {
		*clientID = http.ClientID
	}
	if *showID == "" {
		log.Fatal("show-id not set")
	}
	if *spdcCookieValue == "" {
		log.Fatal("spdc-cookie-val not set")
	}

	var (
		dur time.Duration
		err error
	)
	if dur, err = time.ParseDuration(*timeSchedule); err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	next := now.Add(dur)

	log.Println("update scheduled for:", next.String())

	api, err := http.NewSpotifyAPI(*clientID, *spdcCookieValue)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		t := api.BearerToken()
		if err := storage.WriteToken(*t, *tokenFileName); err != nil {
			log.Fatal(err)
		}
	}()

	token, err := storage.ReadToken(*tokenFileName)
	if err != nil {
		log.Println("can not read token:", err.Error())
		log.Println("trying to fetch new token ...")
		t, err := api.FetchBearerToken()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("token acquired:", t.TokenType)

		api.SetBearerToken(&t)
	} else {
		log.Println("using token:", *tokenFileName)
		api.SetBearerToken(&token)
	}

	prom := metrics.NewPrometheus()

	stdhttp.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Fatal(stdhttp.ListenAndServe(*listenAddr, nil))
	}()

	for {
		now = time.Now()

		m, err := api.FetchMetaData(*showID)
		if err != nil {
			log.Println("spotify-api: ", err)
		} else {
			prom.SetEpisodes(float64(m.TotalEpisodes))
			prom.SetFollowers(float64(m.Followers))
			prom.SetListeners(float64(m.Listeners))
			prom.SetStarts(float64(m.Starts))
			prom.SetStreams(float64(m.Streams))
		}

		next = time.Now().Add(dur)
		log.Println("update scheduled for:", next.String())
		time.Sleep(time.Until(next))
	}

}
