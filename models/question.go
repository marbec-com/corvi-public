package models

import (
	"fmt"
	"math"
	"time"
)

const (
	nextFormulaA float64 = 180.0
	nextFormulaB float64 = -6.0
	nextFormulaC float64 = -0.1
)

type Question struct {
	ID                uint
	Question          string
	Answer            string
	Box               *Box
	Next              time.Time
	CorrectlyAnswered int
}

func NewQuestion(question, answer string, b *Box) *Question {
	return &Question{
		Question:          question,
		Answer:            answer,
		Box:               b,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
	}
}

func (q *Question) CalculateNext() {
	result := 0.0

	if q.CorrectlyAnswered < 7 { // Linear increase
		result = float64(q.CorrectlyAnswered)

	} else if q.CorrectlyAnswered < 80 { // Saturation function
		result = nextFormulaA * math.Pow(math.E, nextFormulaB*math.Pow(math.E, nextFormulaC*float64(q.CorrectlyAnswered)))

	} else { // Saturation limit
		result = 180
	}

	days := int(result + 0.5)

	newTime := time.Now().AddDate(0, 0, days).Truncate(time.Hour * 24)
	fmt.Println(q.CorrectlyAnswered, days, newTime)

	q.Next = newTime
}
