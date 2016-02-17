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

func NewQuestion(b *Box) *Question {
	return &Question{
		Box:               b,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
	}
}

func (q *Question) calculateNext() {
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

func (q *Question) CorrectAnswerGiven() {
	q.CorrectlyAnswered++
	fmt.Println("Correct: ", q.CorrectlyAnswered)
	q.calculateNext()
}

func (q *Question) WrongAnswerGiven() {
	q.CorrectlyAnswered = 0
	fmt.Println("Wrong: ", q.CorrectlyAnswered)
	q.calculateNext()
}
