package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"marb.ec/corvi-backend/models"
	"time"
)

var StatsControllerSingleton *StatsController

func StatsCtrl() *StatsController {
	return StatsControllerSingleton
}

type StatsController struct {
	db       DatabaseService
	settings SettingsService
}

func NewStatsController(db DatabaseService, settings SettingsService) *StatsController {
	c := &StatsController{
		db:       db,
		settings: settings,
	}
	return c
}

func (c *StatsController) LoadStats(from, to time.Time) (*models.Stats, error) {

	if to.Before(from) {
		return nil, errors.New("Invalid range.")
	}

	stats := models.NewStats()
	stats.RangeFrom = from
	stats.RangeTo = to

	err := c.fillQuestionStats(stats)
	if err != nil {
		fmt.Println("Error while fillQuestionStats", err)
		return nil, err
	}
	err = c.fillLearnUnitStats(stats)
	if err != nil {
		fmt.Println("Error while fillLearnUnitStats", err)
		return nil, err
	}
	err = c.fillCountStats(stats)
	if err != nil {
		fmt.Println("Error while fillCountStats", err)
		return nil, err
	}
	err = c.fillBoxStats(stats)
	if err != nil {
		fmt.Println("Error while fillBoxStats", err)
		return nil, err
	}

	return stats, nil
}

func (c *StatsController) fillQuestionStats(stats *models.Stats) error {

	stats.WorstQuestionAnswers = 0

	sqlStmt := `SELECT Question.ID, Question.Question, Question.Answer, Question.BoxID, Question.Next, Question.CorrectlyAnswered, Question.CreatedAt,
				(
					SELECT COUNT(*) 
					FROM LearnUnit 
					WHERE datetime(LearnUnit.CreatedAt) > datetime(?, 'start of day') 
					AND datetime(LearnUnit.CreatedAt) < datetime(?, 'start of day', '+1 day')
					AND LearnUnit.Correct = 0
					AND LearnUnit.QuestionID = Question.ID
				) AS WrongAnswers
				FROM Question 
				WHERE datetime(Question.CreatedAt) > datetime(?, 'start of day') 
				  AND datetime(Question.CreatedAt) < datetime(?, 'start of day', '+1 day')`
	rows, err := c.db.Connection().Query(sqlStmt, stats.RangeFrom, stats.RangeTo, stats.RangeFrom, stats.RangeTo)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		// Create new Question object
		q := models.NewQuestion()
		var wrongAnswers uint
		// Populate
		err = rows.Scan(&q.ID, &q.Question, &q.Answer, &q.BoxID, &q.Next, &q.CorrectlyAnswered, &q.CreatedAt, &wrongAnswers)
		if err != nil {
			return err
		}
		if q.CorrectlyAnswered > 0 { // Question learned
			stats.TotalLearned++
		} else if q.CorrectlyAnswered == 0 && q.CreatedAt.Equal(q.Next) { // Question not learned, and no answer given yet
			stats.TotalUntouched++
		} else {
			stats.TotalUnlearned++ // Question not learned, was at least one time answered
		}
		// Best question: max CorrectlyAnswered value
		if stats.BestQuestion == nil || stats.BestQuestion.CorrectlyAnswered < q.CorrectlyAnswered {
			stats.BestQuestion = q
		}
		// Worst question: Most wrong answers
		if stats.WorstQuestionAnswers == 0 || wrongAnswers > stats.WorstQuestionAnswers {
			stats.WorstQuestion = q
			stats.WorstQuestionAnswers = wrongAnswers
		}
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil

}

func (c *StatsController) fillLearnUnitStats(stats *models.Stats) error {

	sqlStmt := `SELECT LearnUnit.Correct, LearnUnit.PrevCorrect, LearnUnit.CreatedAt
		   		FROM LearnUnit 
		  		WHERE datetime(LearnUnit.CreatedAt) > datetime(?, 'start of day') 
		 		  AND datetime(LearnUnit.CreatedAt) < datetime(?, 'start of day', '+1 day')`
	rows, err := c.db.Connection().Query(sqlStmt, stats.RangeFrom, stats.RangeTo)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {

		var correct, prevCorrect bool
		var createdAt time.Time
		err = rows.Scan(&correct, &prevCorrect, &createdAt)
		if err != nil {
			return err
		}

		// Answers by Weekday (0-6), Day in month (0-30) and Month (0-11)
		stats.LearnUnitsGroupByWeekday[(int(createdAt.Weekday())+6)%7]++ // 0 = Monday, 6 = Sunday
		stats.LearnUnitsGroupByMonthDay[int(createdAt.Day())-1]++
		stats.LearnUnitsGroupByMonth[int(createdAt.Month())-1]++

		// Count total
		stats.TotalLearnUnits++

		if correct {
			// Count total right answers
			stats.TotalNumberOfRightAnswers++
		} else {
			// Count total wrong answers
			stats.TotalNumberOfWrongAnswers++
		}

		if correct && prevCorrect {
			// Right answer on a previously learned question
			stats.KnowledgeRate++
		} else if correct && !prevCorrect {
			// Right answer on a previously unlearned question
			stats.LearnRate++
		} else if !correct && prevCorrect {
			// Wrong answer on a previously learned question
			stats.UnlearnRate++
		} else {
			// Wrong answer on a previously unlearned question
			stats.UnknowingRate++
		}

	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil

}

func (c *StatsController) fillCountStats(stats *models.Stats) error {
	sqlStmt := `SELECT (SELECT COUNT(*) FROM Box) AS TotalBoxes,
					   (SELECT COUNT(*) FROM Question) AS TotalQuestions;`
	row := c.db.Connection().QueryRow(sqlStmt)

	err := row.Scan(&stats.TotalBoxes, &stats.TotalQuestions)
	return err

}

func (c *StatsController) fillBoxStats(stats *models.Stats) error {

	sqlStmt := `SELECT max.ID, max.Name, max.Description, max.CategoryID, max.QuestionsTotal, max.QuestionsLearned, max.CreatedAt, 
	              	   min.ID, min.Name, min.Description, min.CategoryID, min.QuestionsTotal, min.QuestionsLearned, min.CreatedAt
				FROM   BoxWithMeta AS max, 
				       BoxWithMeta AS min
          	     WHERE max.QuestionsTotal > 0
				   AND (max.QuestionsLearned/max.QuestionsTotal) IN (SELECT MAX(QuestionsLearned/QuestionsTotal) FROM BoxWithMeta)
          	       AND min.QuestionsTotal > 0
				   AND (min.QuestionsLearned/min.QuestionsTotal) IN (SELECT MIN(QuestionsLearned/QuestionsTotal) FROM BoxWithMeta)
           	     LIMIT 1;`
	row := c.db.Connection().QueryRow(sqlStmt)

	bestBox := models.NewBox()
	worstBox := models.NewBox()

	err := row.Scan(&bestBox.ID, &bestBox.Name, &bestBox.Description, &bestBox.CategoryID, &bestBox.QuestionsTotal, &bestBox.QuestionsLearned, &bestBox.CreatedAt, &worstBox.ID, &worstBox.Name, &worstBox.Description, &worstBox.CategoryID, &worstBox.QuestionsTotal, &worstBox.QuestionsLearned, &worstBox.CreatedAt)

	if err == sql.ErrNoRows {
		stats.BestBox = nil
		stats.WorstBox = nil
		return nil
	} else if err != nil {
		return err
	}

	if bestBox.ID != 0 {
		stats.BestBox = bestBox
	}
	if worstBox.ID != 0 {
		stats.WorstBox = worstBox
	}

	return nil

}
