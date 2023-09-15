package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aethiopicuschan/voicevox"
	gt "github.com/bas24/googletranslatefree"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func loadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	loadConfig()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			client := voicevox.NewClient("http", "localhost:50021")
			result, _ := gt.Translate(update.Message.Text, "id", "ja")
			query, err := client.CreateQuery(1, result)
			if err != nil {
				fmt.Printf("CreateQuery error: %v\n", err)
				return
			}
			wavAudio, err := client.CreateVoice(1, true, query)
			if err != nil {
				fmt.Printf("CreateVoice error: %v\n", err)
				return
			}

			CreateAudioFile(wavAudio, result)

			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg := tgbotapi.NewAudio(update.Message.From.ID, tgbotapi.FilePath("audio/"+result+".mp3"))

			bot.Send(msg)
		}
	}
}

func CreateAudioFile(wavAudio []byte, filename string) {
	file, err := os.Create("audio/" + filename + ".mp3")
	if err != nil {
		fmt.Printf("CreateAudioFile error ; %v\n ", err)
	}
	bytesWritten, err := file.Write(wavAudio)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
}
