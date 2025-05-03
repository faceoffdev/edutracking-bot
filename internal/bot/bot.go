package bot

import (
	"fmt"
	"github.com/faceoffdev/edutracking-bot/internal/config"
	"github.com/faceoffdev/edutracking-bot/internal/memory"
	"gopkg.in/telebot.v3/middleware"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

type Bot struct{}

func New() *Bot {
	return &Bot{}
}

var topicsBtn = telebot.InlineButton{Unique: "topics", Text: "📚 По темам"}

var logChatId int64

//var formatBtn = telebot.InlineButton{Unique: "oge_format", Text: "📝 По формату ОГЭ"}

func (h *Bot) Run(cfg *config.Config) error {
	logChatId = cfg.Telegram.LogChatID

	pref := telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: 60 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return err
	}

	b.Use(middleware.AutoRespond())

	b.Handle("/start", handleStartCommand)
	//b.Handle(&formatBtn, func(c telebot.Context) error {
	//	return c.Send("Формат ОГЭ в разработке :)")
	//})
	b.Handle(&topicsBtn, func(c telebot.Context) error {
		msg, markup, err := getTopicMessage(b)
		if err != nil {
			return err
		}

		return c.EditOrSend(msg, &telebot.SendOptions{ParseMode: telebot.ModeHTML, ReplyMarkup: markup})
	})
	b.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID
		curr := memory.GetCurrentQuestion(userID)
		if curr == nil {
			return c.Send("Введите команду /start, чтобы начать.")
		}

		answer := c.Text()
		correct, next := memory.SubmitAnswer(userID, answer)

		var msg strings.Builder

		if correct {
			msg.WriteString("✅ Верно!")
		} else {
			msgLog := fmt.Sprintf("Пользователь %s ответил неверно на вопрос под номером %s\n\n%s\n\nОтвет: %s\nПравильный ответ: %s", getUserString(c.Sender()), curr.ID, curr.Question, answer, curr.Answer)
			c.Bot().Send(telebot.ChatID(logChatId), msgLog)

			msg.WriteString("❌ Неверно.")
		}

		if next != nil {
			msg.WriteString(" Следующий вопрос:\n\n")
			msg.WriteString("<i>" + next.Question + "</i>\n\n")
			msg.WriteString("Слово: <b>" + next.Word + "</b>")
		} else {
			msg.WriteString(" Вы прошли все вопросы этой темы! 🎉")
		}

		return c.Send(msg.String(), &telebot.SendOptions{ParseMode: telebot.ModeHTML})
	})

	b.Start()

	return nil
}

func getUserString(user *telebot.User) string {
	var str strings.Builder

	if user.FirstName != "" {
		str.WriteString(user.FirstName + " ")
	}

	if user.LastName != "" {
		str.WriteString(user.LastName + " ")
	}

	if user.Username != "" {
		str.WriteString(fmt.Sprintf("(@%s)", user.Username))
	}

	return str.String()
}
