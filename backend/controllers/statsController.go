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

	// TODO(mjb): Fill object

	return stats, nil
}
