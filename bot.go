package main

import (
	"io/ioutil"

	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"

	"log"
	"time"

	tele "gopkg.in/tucnak/telebot.v2"
)

const (
	spec        = "0 0 0 * * ?"
	telegramAPI = ""
	botToken    = "1111111111:***********************************"
	timeout     = 10
	gid         = 1010001100
	dataFile    = "data.json"
)

func main() {
	bot, err := tele.NewBot(tele.Settings{URL: telegramAPI, Token: botToken, Poller: &tele.LongPoller{Timeout: timeout * time.Second}})
	if err != nil {
		log.Fatal(err)
		return
	}

	month := time.Now().Format("01")
	day := time.Now().Format("02")

	bot.Handle("/today", func(msg *tele.Message) {
		bot.Send(msg.Chat, historyToday(month, day))
	})

	// 定时
	crontab := cron.New(cron.WithSeconds())
	task := func() {
		bot.Send(tele.ChatID(gid), historyToday(month, day), "Markdown")
	}
	crontab.AddFunc(spec, task)

	crontab.Start()
	bot.Start()
}

func historyToday(month, day string) string {
	data, _ := ioutil.ReadFile(dataFile)
	tip := gjson.Get(string(data), month+".tip")
	event := gjson.Get(string(data), month+"."+day)
	today := month + "-" + day

	if tip.String() != "" {
		return tip.String() + "\n=====\n\n" + event.String() + "\n\n" + today
	}
	return event.String() + "\n\n" + today
}
