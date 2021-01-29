package main

import (
	"io/ioutil"
	"strconv"

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
		bot.Send(msg.Chat, historyToday(month, day), tele.NoPreview, "Markdown")
	})

	// 定时
	crontab := cron.New(cron.WithSeconds())
	task := func() {
		bot.Send(tele.ChatID(gid), historyToday(month, day), tele.NoPreview, "Markdown")
	}
	crontab.AddFunc(spec, task)

	crontab.Start()
	bot.Start()
}

func historyToday(month, day string) string {
	data, _ := ioutil.ReadFile(dataFile)
	tip := gjson.Get(string(data), month+".tip")
	today := month + "-" + day

	if tip.String() != "" {
		return tip.String() + "\n=====\n\n" + eventList(month, day) + "\n" + today
	}

	return eventList(month, day) + "\n" + today
}

func eventList(month, day string) string {
	data, _ := ioutil.ReadFile("data.json")
	count, _ := strconv.Atoi((gjson.Get(string(data), month+"."+day+".#")).String())
	event := ""
	if (gjson.Get(string(data), month+"."+day+".0")).String() != "" {
		for i := 0; i < count; i++ {
			event = event + "\n" + (gjson.Get(string(data), month+"."+day+"."+strconv.Itoa(i))).String()
		}
		return event
	} else {
		event = "暂无历史今天的性少数群体历程\n你可以前往 GitHub (https://github.com/LGBT-CN/HistoryToday/edit/master/data.json) 提交数据"
		return event
	}
}
