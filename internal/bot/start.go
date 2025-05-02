package bot

import (
	"fmt"
	"github.com/faceoffdev/edutracking-bot/internal/memory"

	"gopkg.in/telebot.v3"
)

func handleStartCommand(c telebot.Context) error {
	defer c.Respond()

	memory.StartQuiz(c.Sender().ID, nil)

	markup := &telebot.ReplyMarkup{}
	markup.InlineKeyboard = [][]telebot.InlineButton{{topicsBtn}}

	tgErr := c.Send("Добро пожаловать!"+
		//"\n\n"+
		//"ТУТ ДОЛЖНО БЫТЬ КАКОЕ-ТО ОПИСАНИЕ. КТО МЫ И ЗАЧЕМ ЗДЕСЬ"+
		"\n\n"+
		"Чем больше практики, тем быстрее вы освоите иностранный язык и подготовитесь к экзамену.", markup)

	if tgErr != nil {
		fmt.Println("Ошибка отправки сообщения в Telegram: " + tgErr.Error())
		return tgErr
	}

	return nil
}
