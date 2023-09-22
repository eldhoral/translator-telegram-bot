package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

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
	StartMessage     = "/start"
	VoiceCallback    = "Voice"
	LanguageCallback = "Language"
)

type User struct {
	TelegramId      int64
	SpeakerSelected int
	SourceLanguage  string
	TargetLanguage  string
}

var onlyNumeric = regexp.MustCompile(`[^0-9]+`)

func TelegramBot(client *redis.Client) {
	DeleteAllCreatedAudio()
	voiceStyle := InlineVoiceStyle()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = cast.ToBool(os.Getenv("BOT_DEBUG"))
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.Text == StartMessage {
				SendMessageChooseVoice(update.Message.From.ID, bot, voiceStyle)
				continue
			} else {
				data, err := redisClient.NewRedisClient(client).Get(context.TODO(), CacheKeyTelegram+cast.ToString(update.Message.From.ID))
				if err != nil || data == "" {
					SendMessageChooseVoiceAgain(update.Message.From.ID, bot, voiceStyle)
					continue
				}
				translationResult, err := QueryToVoiceVox(data, bot, update)
				if err != nil {
					SendMessageError(update.Message.From.ID, bot)
					continue
				}
				SendMessageAudioWithTranslation(update.Message.From.ID, bot, translationResult, update)
				continue
			}

		}

		if strings.Contains(update.CallbackData(), VoiceCallback) {
			text := "Voice is choosed. You can start typing any word ٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ"
			err := SaveUserPreferencesTelegram(text, bot, client, update)
			if err != nil {
				SendMessageError(update.CallbackQuery.Message.Chat.ID, bot)
			}
			continue
		}

		if strings.Contains(update.CallbackData(), LanguageCallback) {
			// TODO
			continue
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

func SendMessageError(chatMessageId int64, bot *tgbotapi.BotAPI) {
	text := "Oops, voice styling is in development. So many error (つ﹏⊂)"
	msgText := tgbotapi.NewMessage(chatMessageId, text)
	bot.Send(msgText)
}

func SendMessageChooseVoice(chatMessageId int64, bot *tgbotapi.BotAPI, inlineKeyboardMessage tgbotapi.InlineKeyboardMarkup) {
	text := "Please choose voice style (づ ◕‿◕ )づ"
	msgCallBack := tgbotapi.NewMessage(chatMessageId, text)
	msgCallBack.ReplyMarkup = inlineKeyboardMessage
	bot.Send(msgCallBack)
}

func SendMessageChooseVoiceAgain(chatMessageId int64, bot *tgbotapi.BotAPI, inlineKeyboardMessage tgbotapi.InlineKeyboardMarkup) {
	text := "Please choose voice again ヽ༼ ಠ益ಠ ༽ﾉ"
	msgCallBack := tgbotapi.NewMessage(chatMessageId, text)
	msgCallBack.ReplyMarkup = inlineKeyboardMessage
	bot.Send(msgCallBack)
}

func SendMessageAudioWithTranslation(chatMessageId int64, bot *tgbotapi.BotAPI, translationResult string, update tgbotapi.Update) {
	msg := tgbotapi.NewAudio(update.Message.From.ID, tgbotapi.FilePath(os.Getenv("AUDIO_PATH")+cast.ToString(update.Message.From.ID)+"_"+translationResult+".mp3"))
	bot.Send(msg)
	msgText := tgbotapi.NewMessage(update.Message.From.ID, translationResult)
	bot.Send(msgText)
}

func QueryToVoiceVox(userDataTelegram string, bot *tgbotapi.BotAPI, update tgbotapi.Update) (translationResult string, err error) {
	userTelegram := &User{}
	err = jsoniter.Unmarshal([]byte(userDataTelegram), userTelegram)

	client := voicevox.NewClient("http", os.Getenv("VOICEVOX_HOST")+":"+os.Getenv("VOICEVOX_PORT"))
	translationResult, _ = gt.Translate(update.Message.Text, "id", "ja")
	query, err := client.CreateQuery(userTelegram.SpeakerSelected, translationResult)
	if err != nil {
		log.Printf("CreateVoice error: %v\n", err)
		return
	}

	wavAudio, err := client.CreateVoice(userTelegram.SpeakerSelected, true, query)
	if err != nil {
		log.Printf("CreateVoice error: %v\n", err)
		return
	}

	CreateAudioFile(wavAudio, translationResult, update.Message.From.ID)
	return
}

func SaveUserPreferencesTelegram(text string, bot *tgbotapi.BotAPI, client *redis.Client, update tgbotapi.Update) (err error) {
	msgText := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	bot.Send(msgText)

	userTelegram := User{
		TelegramId:      update.CallbackQuery.Message.Chat.ID,
		SpeakerSelected: cast.ToInt(removeExceptNumber(update.CallbackQuery.Data)),
	}

	serviceBytes, err := jsoniter.Marshal(userTelegram)
	if err != nil {
		log.Printf("Error marshal with error: %v\n", err)
	}
	_, err = redisClient.NewRedisClient(client).Set(context.TODO(), CacheKeyTelegram+cast.ToString(update.CallbackQuery.Message.Chat.ID), serviceBytes)
	if err != nil {
		log.Printf("[Redis] Error set cache to redis when adding data(): %v\n", err)
	}
	return
}

func InlineVoiceStyle() (msgButtonMarkup tgbotapi.InlineKeyboardMarkup) {
	msgButton1 := tgbotapi.NewInlineKeyboardButtonData("1.Normal", "Voice 2")
	msgButton2 := tgbotapi.NewInlineKeyboardButtonData("2.Sweet", "Voice 0")
	msgButton3 := tgbotapi.NewInlineKeyboardButtonData("3.Tsun Tsun", "Voice 6")
	msgButton4 := tgbotapi.NewInlineKeyboardButtonData("4.Whisper", "Voice 26")
	msgButtonRow1 := tgbotapi.NewInlineKeyboardRow(msgButton1, msgButton2)
	msgButtonRow2 := tgbotapi.NewInlineKeyboardRow(msgButton3, msgButton4)
	msgButtonMarkup = tgbotapi.NewInlineKeyboardMarkup(msgButtonRow1, msgButtonRow2)
	return
}

func InlineLanguage() (msgButtonMarkup tgbotapi.InlineKeyboardMarkup) {
	msgButton1 := tgbotapi.NewInlineKeyboardButtonData("Indonesia to Japanese", "Language 1")
	msgButton2 := tgbotapi.NewInlineKeyboardButtonData("Japanese to Indonesia", "Language 2")
	msgButton3 := tgbotapi.NewInlineKeyboardButtonData("Indonesia to English", "Language 3")
	msgButton4 := tgbotapi.NewInlineKeyboardButtonData("English to Indonesia", "Language 4")
	msgButtonRow1 := tgbotapi.NewInlineKeyboardRow(msgButton1)
	msgButtonRow2 := tgbotapi.NewInlineKeyboardRow(msgButton2)
	msgButtonRow3 := tgbotapi.NewInlineKeyboardRow(msgButton3)
	msgButtonRow4 := tgbotapi.NewInlineKeyboardRow(msgButton4)
	msgButtonMarkup = tgbotapi.NewInlineKeyboardMarkup(msgButtonRow1, msgButtonRow2, msgButtonRow3, msgButtonRow4)
	return
}

func removeExceptNumber(str string) (result string) {
	result = onlyNumeric.ReplaceAllString(str, "")
	return
}
