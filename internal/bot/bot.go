package bot

import (
	"fmt"
	"github.com/faceoffdev/edutracking-bot/internal/config"
	"github.com/faceoffdev/edutracking-bot/internal/memory"
	"time"

	"gopkg.in/telebot.v3"
)

type Bot struct{}

func New() *Bot {
	return &Bot{}
}

var topicsBtn = telebot.InlineButton{Unique: "topics", Text: "📚 По темам"}

//var formatBtn = telebot.InlineButton{Unique: "oge_format", Text: "📝 По формату ОГЭ"}

func (h *Bot) Run(cfg *config.Config) error {
	pref := telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: 60 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return err
	}

	b.Handle("/start", handleStartCommand)
	//b.Handle(&formatBtn, func(c telebot.Context) error {
	//	return c.Send("Формат ОГЭ в разработке :)")
	//})
	b.Handle(&topicsBtn, func(c telebot.Context) error {
		msg, markup, err := getTopicMessage(b)
		if err != nil {
			return err
		}

		return c.Send(msg, markup)
	})
	b.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID
		curr := memory.GetCurrentQuestion(userID)
		if curr == nil {
			return c.Send("Введите команду /start, чтобы начать.")
		}

		answer := c.Text()
		correct, next := memory.SubmitAnswer(userID, answer)

		if correct {
			if next != nil {
				return c.Send(fmt.Sprintf("✅ Верно!\n\n%s\n\nСлово: %s", next.Question, next.Word))
			}
			return c.Send("✅ Верно!\n\n🎉 Вы прошли все вопросы этой темы!")
		} else {
			return c.Send("❌ Неверно. Попробуйте ещё раз.")
		}
	})

	b.Start()

	return nil
}
