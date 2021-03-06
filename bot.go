package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/tidwall/gjson"

	"log"
	"time"

	tele "gopkg.in/tucnak/telebot.v2"
)

const dataFile = "data.json"
const timeout = 5

var isTest = false

func main() {
	t := os.Getenv("lgbtcntest")
	if t == "" {
		isTest = false
	}
	b, err := strconv.ParseBool(t)
	if err == nil {
		isTest = b
	}

	fmt.Printf("[I] TEST MODE: %v\n", isTest)

	bot, err := tele.NewBot(
		tele.Settings{
			Token: os.Getenv("Token"),
			Poller: &tele.LongPoller{Timeout: timeout * time.Second,
			},
		},
	)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println("[I] BOT CREATED WITHOUT ANY ERR.")

	month := time.Now().Format("01")
	day := time.Now().Format("02")

	c := os.Getenv("Chat_ID")
	ChatID, errC := strconv.ParseInt(c, 10, 64)

	if errC != nil {
		fmt.Println("[E] CAN NOT PARSE CHAR_ID TO INT64! EXIT!")
		os.Exit(1)
	}

	fmt.Println("[I] GET CHAT_ID SUCCESSFULLY")

	// bot.Start()
   	_, errS := bot.Send(tele.ChatID(ChatID), historyToday(month, day), tele.NoPreview, "Markdown")

	fmt.Println("[I] MSG SENT. EXITING.")

	// After sending msg, exit, cuz it boosts by CI/FaaS & Cron
   	if errS == nil {
		os.Exit(0)
	}
	fmt.Println(err)
	os.Exit(1)
}

func historyToday(month, day string) string {
	data, _ := ioutil.ReadFile(dataFile)
	tip := gjson.Get(string(data), month+".tip")
	today := month + "-" + day

	if tip.String() != "" {
		return tip.String() + "\n=====\n\n" + eventList(month, day) + "\n" + today
	}

	if isTest {
		return "[TEST]\n" + eventList(month, day) + "\n" + today
	}

	return eventList(month, day) + "\n" + today
}

func eventList(month, day string) string {
	data, _ := ioutil.ReadFile(dataFile)
	count, _ := strconv.Atoi((gjson.Get(string(data), month+"."+day+".#")).String())
	event := ""
	if (gjson.Get(string(data), month+"."+day+".0")).String() != "" {
		for i := 0; i < count; i++ {
			event = event + "\n" + (gjson.Get(string(data), month+"."+day+"."+strconv.Itoa(i))).String()
		}

		return event
	}
	event = "暂无历史今天的性少数群体历程\n你可以[前往 GitHub 提交数据](https://github.com/LGBT-CN/HistoryToday/edit/master/data.json)"
	return event
}
