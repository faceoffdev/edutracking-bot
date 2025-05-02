package db

type Question struct {
	ID       string
	TopicID  string
	Question string
	Word     string
	Answer   string
}

func GetQuestionsByTopic(topicID string) ([]Question, error) {
	rows, err := doQuery[Question](`SELECT id, topic_id, question, word, answer
        FROM grammar_questions
        WHERE topic_id = $1
        ORDER BY id ASC`, topicID)
	if err != nil {
		return nil, err
	}

	return getList[Question](rows), nil
}
