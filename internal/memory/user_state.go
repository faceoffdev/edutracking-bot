package memory

import (
	"github.com/faceoffdev/edutracking-bot/internal/db"
	"strings"
)

type UserState struct {
	Questions []db.Question
	Index     int
}

var state = make(map[int64]*UserState)

func StartQuiz(userID int64, questions []db.Question) {
	state[userID] = &UserState{Questions: questions, Index: 0}
}

func GetCurrentQuestion(userID int64) *db.Question {
	s := state[userID]
	if s == nil || s.Index >= len(s.Questions) {
		return nil
	}
	return &s.Questions[s.Index]
}

func SubmitAnswer(userID int64, answer string) (bool, *db.Question) {
	s := state[userID]
	if s == nil || s.Index >= len(s.Questions) {
		return false, nil
	}

	curr := s.Questions[s.Index]

	s.Index++
	var next *db.Question
	if s.Index < len(s.Questions) {
		next = &s.Questions[s.Index]
	}

	for _, a := range strings.Split(curr.Answer, "/") {
		if strings.ToUpper(a) == strings.ToUpper(answer) {
			return true, next
		}
	}

	return false, next
}
