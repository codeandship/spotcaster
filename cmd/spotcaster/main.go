package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codeandship/spotcaster/http"
	"github.com/codeandship/spotcaster/storage"
	"github.com/peterbourgon/ff"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)

	fs := flag.NewFlagSet("my-program", flag.ExitOnError)
	var (
		clientID        = fs.String("client-id", "", "spotify client id for authentication (does not need to be set for now)")
		debug           = fs.Bool("debug", false, "enable debug output")
		showID          = fs.String("show-id", "", "podcast show id")
		spdcCookieValue = fs.String("spdc-cookie-value", "", "sp_dc cookie value from spotify's webapp when you are logged in")
		tgChatID        = fs.Int64("tg-chat-id", 0, "set telegram chat id")
		tgToken         = fs.String("tg-token", "", "telegram bot token")
		timeSchedule    = fs.String("time-sched", "8h", "set scheduled start (offset from midnight 0:00, defaults to 8 hours)")
		tokenFileName   = fs.String("token-file", "token.json", "access token file in json format")
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
	if *tgToken == "" {
		log.Fatal("tg-token not set")
	}

	var (
		dur time.Duration
		err error
	)
	if dur, err = time.ParseDuration(*timeSchedule); err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	next := time.Time{}
	if time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(dur).After(now) {
		next = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		next = next.Add(dur)
	} else {
		next = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		next = next.Add(dur)
	}

	log.Println("update scheduled for:", next.String())

	bot, err := tgbotapi.NewBotAPI(*tgToken)
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		bot.Debug = true
	}

	log.Printf("telegram bot account=%s", bot.Self.UserName)

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

	for {
		time.Sleep(time.Until(next))
		now = time.Now()
		msg := tgbotapi.NewMessage(*tgChatID, "")

		m, err := api.FetchMetaData(*showID)
		if err != nil {
			log.Println("spotify-api: ", err)
			msg.Text = fmt.Sprintf("spotify-api: %s", err.Error())
		} else {
			msg.Text = m.Markdown()
			msg.ParseMode = tgbotapi.ModeMarkdown
		}
		if _, err := bot.Send(msg); err != nil {
			log.Println("telegram bot:", err)
		}

		next = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()).Add(dur)
		log.Println("update scheduled for:", next.String())
	}

}
