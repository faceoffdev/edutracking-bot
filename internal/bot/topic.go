package bot

import (
	"fmt"
	"github.com/faceoffdev/edutracking-bot/internal/db"
	"github.com/faceoffdev/edutracking-bot/internal/memory"
	"gopkg.in/telebot.v3"
	"strings"
)

func getTopicMessage(b *telebot.Bot) (string, *telebot.ReplyMarkup, error) {
	topics, err := db.GetTopics()
	if err != nil {
		return "", nil, err
	}

	markup := &telebot.ReplyMarkup{}
	var rows [][]telebot.InlineButton

	var messageText strings.Builder
	messageText.WriteString("Список доступных тем\n\n")

	for i, topic := range topics {
		messageText.WriteString(fmt.Sprintf("%s — %s\n", topic.ID, topic.Name))

		btn := telebot.InlineButton{
			Unique: "topic_" + strings.ReplaceAll(topic.ID, ".", "_"),
			Text:   topic.ID,
			Data:   topic.ID,
		}
		if i%2 == 0 {
			rows = append(rows, []telebot.InlineButton{btn})
		} else {
			rows[len(rows)-1] = append(rows[len(rows)-1], btn)
		}

		b.Handle(&btn, func(c telebot.Context) error {
			defer c.Respond()

			topicID := c.Callback().Data
			questions, err := db.GetQuestionsByTopic(topicID)
			if err != nil || len(questions) == 0 {
				return c.Send("Нет вопросов по выбранной теме.")
			}

			memory.StartQuiz(c.Sender().ID, questions)

			q := questions[0]
			return c.Send(fmt.Sprintf("Прочитайте приведенный ниже текст. Преобразуйте слово, напечатанное заглавными буквами так, чтобы оно грамматически соответствовало содержанию текста. Заполните пропуск полученным словом.\n\n%s\n\nСлово: %s", q.Question, q.Word))
		})
	}

	markup.InlineKeyboard = rows

	return messageText.String(), markup, nil
}
