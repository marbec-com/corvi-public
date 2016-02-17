package models

type QuestionHeap struct {
}

func NewQuestionHeap() *QuestionHeap {
	return &QuestionHeap{}
}

func (h *QuestionHeap) Add(q *Question) {

}

func (h *QuestionHeap) Length() int {
	return 0
}

func (h *QuestionHeap) Min() *Question {
	return nil
}
