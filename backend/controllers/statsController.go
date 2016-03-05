package controllers

import (
	"errors"
	"marb.ec/corvi-backend/models"
	"time"
)

var StatsControllerSingleton *StatsController

func StatsControllerInstance() *StatsController {
	if StatsControllerSingleton == nil {
		StatsControllerSingleton = NewStatsController()
	}
	return StatsControllerSingleton
}

type StatsController struct{}

func NewStatsController() *StatsController {
	return &StatsController{}
}

func (c *StatsController) LoadStats(from, to time.Time) (*models.Stats, error) {

	if to.Before(from) {
		return nil, errors.New("Invalid range.")
	}

	stats := &models.Stats{
		RangeFrom: from,
		RangeTo:   to,
	}

	// TODO(mjb): Replace with correct SQL

	for _, question := range mockQuestions {
		if question.CreatedAt.After(from) && question.CreatedAt.Before(to) {
			stats.TotalQuestions++
			if question.CorrectlyAnswered > 0 {
				stats.TotalLearned++
			} else if question.CorrectlyAnswered == 0 && question.CreatedAt == question.Next {
				stats.TotalUntouched++
			} else {
				stats.TotalUnlearned++
			}
		}
	}

	for _, lu := range mockAnswers {
		if lu.CreatedAt.After(from) && lu.CreatedAt.Before(to) {
			stats.TotalLearnUnits++
			if lu.Correct {
				stats.TotalNumberOfRightAnswers++
			} else {
				stats.TotalNumberOfWrongAnswers++
			}
			if lu.Correct && lu.PrevCorrect {
				stats.KnowledgeRate++
			} else if lu.Correct && !lu.PrevCorrect {
				stats.LearnRate++
			} else if !lu.Correct && lu.PrevCorrect {
				stats.UnlearnRate++
			} else {
				stats.UnknowingRate++
			}
		}
	}

	for _, box := range mockBoxes {
		if box.CreatedAt.After(from) && box.CreatedAt.Before(to) {
			stats.TotalBoxes++
		}
	}

	stats.BestBox = mockBoxes[0]
	stats.WorstBox = mockBoxes[1]
	stats.BestQuestion = mockQuestions[0]
	stats.WorstQuestion = mockQuestions[4]

	return stats, nil
}
