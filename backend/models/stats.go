package models

import (
	"time"
)

type Stats struct {
	RangeFrom                 time.Time
	RangeTo                   time.Time
	TotalQuestions            uint // Where CreatedAt > RangeFrom and < RangeTo
	TotalLearnUnits           uint
	LearnUnitsGroupByWeekday  []uint // Key is weekday 0-6
	LearnUnitsGroupByMonthDay []uint // Key is day in month 1-31
	LearnUnitsGroupByMonth    []uint // Key is month 1-12
	TotalBoxes                uint
	TotalNumberOfRightAnswers uint
	TotalNumberOfWrongAnswers uint
	TotalLearned              uint // number of total questions that are learned
	TotalUnlearned            uint // number of total questions that are not learned
	TotalUntouched            uint // number of total questions that were not yet learned (Next == CreatedAt)
	LearnRate                 uint // number of total LearnUnits were a unlearned question was answered correctly
	KnowledgeRate             uint // number of total LearnUnits were a learned question was answered correctly
	UnlearnRate               uint // number of total LearnUnits were a learned question was answered incorrectly
	UnknowingRate             uint // number of total LearnUnits were a unlearned question was answered incorrectly
	BestBox                   *Box
	WorstBox                  *Box
	BestQuestion              *Question
	WorstQuestion             *Question
}

func NewStats() *Stats {
	return &Stats{
		LearnUnitsGroupByWeekday:  make([]uint, 7),
		LearnUnitsGroupByMonthDay: make([]uint, 31),
		LearnUnitsGroupByMonth:    make([]uint, 12),
	}
}
