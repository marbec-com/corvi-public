package models

import (
	"sync"
)

type QuestionHeap struct {
	questions []*Question
	sync.Mutex
}

func NewQuestionHeap() *QuestionHeap {
	return &QuestionHeap{
		questions: make([]*Question, 0),
	}
}

func (h *QuestionHeap) Add(q *Question) {
	h.questions = append(h.questions, q)
}

func (h *QuestionHeap) Length() int {
	return len(h.questions)
}

func (h *QuestionHeap) Min() *Question {
	if len(h.questions) < 1 {
		return nil
	}

	first := h.questions[0]
	h.questions = h.questions[1:]

	return first
}

func (h *QuestionHeap) Clear() {
	h.questions = make([]*Question, 0)
}

func (h *QuestionHeap) Empty() bool {
	return h.Length() == 0
}

func (h *QuestionHeap) Peek() *Question {
	if len(h.questions) < 1 {
		return nil
	}

	return h.questions[0]
}

func (h *QuestionHeap) Equal(a *QuestionHeap) bool {
	if len(h.questions) != len(a.questions) {
		return false
	}
	for i := 0; i < len(h.questions); i++ {
		if !h.questions[i].Equal(a.questions[i]) {
			return false
		}
	}
	return true
}
