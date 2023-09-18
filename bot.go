package main

import (
	"context"
	"fmt"
	"log"
	"os"

	redisClient "go-telegram-bot/redis"

	"github.com/aethiopicuschan/voicevox"
	gt "github.com/bas24/googletranslatefree"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

const (
	CacheKeyTelegram = "go_telegram_bot:chat_id:"
)

type User struct {
	TelegramId      int64
	SpeakerSelected int
}

func TelegramBot(client *redis.Client) {
	// new chat
	msgButton1 := tgbotapi.NewInlineKeyboardButtonData("1.Normal", "2")
	msgButton2 := tgbotapi.NewInlineKeyboardButtonData("2.Sweet", "0")
	msgButton3 := tgbotapi.NewInlineKeyboardButtonData("3.Tsun Tsun", "6")
	msgButton4 := tgbotapi.NewInlineKeyboardButtonData("4.Whisper", "26")
	msgButtonRow1 := tgbotapi.NewInlineKeyboardRow(msgButton1, msgButton2)
	msgButtonRow2 := tgbotapi.NewInlineKeyboardRow(msgButton3, msgButton4)
	msgButtonMarkup := tgbotapi.NewInlineKeyboardMarkup(msgButtonRow1, msgButtonRow2)

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
			if update.Message.Text == "/start" {
				text := "Please choose voice style (づ ◕‿◕ )づ"
				msgText := tgbotapi.NewMessage(update.Message.From.ID, text)
				bot.Send(msgText)

				msgCallBack := tgbotapi.NewMessage(update.Message.From.ID, update.Message.Text)
				msgCallBack.ReplyMarkup = msgButtonMarkup
				bot.Send(msgCallBack)
				continue
			} else {
				data, err := redisClient.NewRedisClient(client).Get(context.TODO(), CacheKeyTelegram+cast.ToString(update.Message.From.ID))
				if err != nil || data == "" {
					text := "Please choose voice again ヽ༼ ಠ益ಠ ༽ﾉ"
					msgText := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
					bot.Send(msgText)

					msgCallBack := tgbotapi.NewMessage(update.Message.From.ID, update.Message.Text)
					msgCallBack.ReplyMarkup = msgButtonMarkup
					bot.Send(msgCallBack)
				}
				userTelegram := &User{}
				err = jsoniter.Unmarshal([]byte(data), userTelegram)

				client := voicevox.NewClient("http", "localhost:50021")
				result, _ := gt.Translate(update.Message.Text, "id", "ja")
				query, err := client.CreateQuery(userTelegram.SpeakerSelected, result)
				if err != nil {
					fmt.Printf("CreateQuery error: %v\n", err)
					return
				}
				wavAudio, err := client.CreateVoice(userTelegram.SpeakerSelected, true, query)
				if err != nil {
					fmt.Printf("CreateVoice error: %v\n", err)
					return
				}

				CreateAudioFile(wavAudio, result, update.Message.From.ID)

				// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg := tgbotapi.NewAudio(update.Message.From.ID, tgbotapi.FilePath(os.Getenv("AUDIO_PATH")+cast.ToString(update.Message.From.ID)+"_"+result+".mp3"))
				bot.Send(msg)
				msgText := tgbotapi.NewMessage(update.Message.From.ID, result)
				bot.Send(msgText)
				continue
			}

		}

		if update.CallbackData() != "" {
			text := "Voice is choosed. You can start typing any word ٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ"
			msgText := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
			bot.Send(msgText)

			userTelegram := User{
				TelegramId:      update.CallbackQuery.Message.Chat.ID,
				SpeakerSelected: cast.ToInt(update.CallbackQuery.Data),
			}
			serviceBytes, err := jsoniter.Marshal(userTelegram)
			if err != nil {
				log.Printf("Error marshal with error: %v\n", err)
			}
			_, err = redisClient.NewRedisClient(client).Set(context.TODO(), CacheKeyTelegram+cast.ToString(update.CallbackQuery.Message.Chat.ID), serviceBytes)
			if err != nil {
				text = "Oops, voice styling is in development. So many error (つ﹏⊂)"
				msgText = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
				bot.Send(msgText)
				log.Printf("[Redis] Error set cache to redis when adding data(): %v\n", err)
			}
		}
	}
}

func CreateAudioFile(wavAudio []byte, filename string, teleId int64) {
	file, err := os.Create(os.Getenv("AUDIO_PATH") + cast.ToString(teleId) + "_" + filename + ".mp3")
	if err != nil {
		fmt.Printf("CreateAudioFile error ; %v\n ", err)
	}
	bytesWritten, err := file.Write(wavAudio)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
}
