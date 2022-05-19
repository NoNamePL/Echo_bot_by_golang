package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	BOTTOKEN = "5300131496:AAFb9regAsI0lmI1cS3qh5CEE13KvC86PJE"
	BOTAPI   = "https://api.telegram.org/bot"
	URL      = BOTAPI + BOTTOKEN
)

func main() {
	var offset = 0
	for {
		updates, err := getUpdates(URL, offset)
		if err != nil {
			log.Fatal(err)
		}
		for _, update := range updates {
			err = respond(URL, update)
			offset = update.UpdateID + 1
		}
		fmt.Println(updates)
	}
}

// запрос обновления
func getUpdates(URL string, offset int) ([]Update, error) {
	resp, err := http.Get(URL + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var restResponce RestResponce
	err = json.Unmarshal(text, &restResponce)
	if err != nil {
		log.Fatal(err)
	}

	return restResponce.Result, nil
}

// отвечает на обновления
func respond(URL string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatID = update.Message.Chat.ChatID
	botMessage.Text = update.Message.Text
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(URL + "/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}
