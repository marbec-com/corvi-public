package controllers

import (
	"errors"
	"fmt"
	"log"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/events"
)

var QuestionControllerSingleton *QuestionController

func QuestionCtrl() *QuestionController {
	return QuestionControllerSingleton
}

type QuestionController struct {
	db *DBController
}

func NewQuestionController(db *DBController) (*QuestionController, error) {
	c := &QuestionController{
		db: db,
	}
	err := c.createTables()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *QuestionController) createTables() error {

	// Create table for questions, only if it not already exists
	sql := "CREATE TABLE IF NOT EXISTS Question (ID INTEGER PRIMARY KEY ASC NOT NULL, Question VARCHAR (255) NOT NULL, Answer TEXT NOT NULL, BoxID INTEGER REFERENCES Box (ID) ON DELETE CASCADE NOT NULL, Next DATETIME NOT NULL, CorrectlyAnswered INTEGER NOT NULL, CreatedAt DATETIME NOT NULL);"
	_, err := c.db.Connection().Exec(sql)
	if err != nil {
		return err
	}

	// Create table for learnunits, only if it not already exists
	sql = "CREATE TABLE IF NOT EXISTS LearnUnit (QuestionID INTEGER REFERENCES Question (ID) ON DELETE CASCADE NOT NULL, BoxID INTEGER REFERENCES Box (ID) ON DELETE CASCADE NOT NULL, Correct BOOLEAN NOT NULL, PrevCorrect BOOLEAN NOT NULL, CreatedAt DATETIME NOT NULL);"
	_, err = c.db.Connection().Exec(sql)
	if err != nil {
		return err
	}

	// Create view for questions learned today
	sql = "CREATE VIEW IF NOT EXISTS QuestionsLearnedToday AS SELECT * FROM LearnUnit WHERE date(CreatedAt) = date('now') AND Correct = 1"
	_, err = c.db.Connection().Exec(sql)
	if err != nil {
		return err
	}

	// Create view for questions due
	sql = "CREATE VIEW IF NOT EXISTS QuestionsDue AS SELECT * FROM Question WHERE datetime(Next) < datetime('now', 'start of day', '+1 day') AND ID NOT IN (SELECT ID FROM QuestionsLearnedToday)"
	_, err = c.db.Connection().Exec(sql)

	return err
}

func (c *QuestionController) LoadQuestions() ([]*models.Question, error) {

	// Select all questions
	sql := "SELECT ID, Question, Answer, BoxID, Next, CorrectlyAnswered, CreatedAt FROM Question;"
	rows, err := c.db.Connection().Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create empty result set
	result := make([]*models.Question, 0)

	for rows.Next() {
		// Create new Box object
		newQuestion := models.NewQuestion()
		// Populate
		err = rows.Scan(&newQuestion.ID, &newQuestion.Question, &newQuestion.Answer, &newQuestion.BoxID, &newQuestion.Next, &newQuestion.CorrectlyAnswered, &newQuestion.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Append to result set
		result = append(result, newQuestion)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (c *QuestionController) LoadQuestionsOfBox(id uint) ([]*models.Question, error) {

	// Select all questions
	sql := "SELECT ID, Question, Answer, BoxID, Next, CorrectlyAnswered, CreatedAt FROM Question WHERE BoxID = ?;"
	rows, err := c.db.Connection().Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create empty result set
	result := make([]*models.Question, 0)

	for rows.Next() {
		// Create new Box object
		newQuestion := models.NewQuestion()
		// Populate
		err = rows.Scan(&newQuestion.ID, &newQuestion.Question, &newQuestion.Answer, &newQuestion.BoxID, &newQuestion.Next, &newQuestion.CorrectlyAnswered, &newQuestion.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Append to result set
		result = append(result, newQuestion)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (c *QuestionController) LoadQuestion(id uint) (*models.Question, error) {

	// Select question with matching ID
	sql := "SELECT ID, Question, Answer, BoxID, Next, CorrectlyAnswered, CreatedAt FROM Question WHERE ID = ?;"
	row := c.db.Connection().QueryRow(sql, id)

	// Create new Category object
	newQuestion := models.NewQuestion()

	// Populate
	err := row.Scan(&newQuestion.ID, &newQuestion.Question, &newQuestion.Answer, &newQuestion.BoxID, &newQuestion.Next, &newQuestion.CorrectlyAnswered, &newQuestion.CreatedAt)
	if err != nil {
		return nil, err
	}

	return newQuestion, nil

}

func (c *QuestionController) GiveAnswer(id uint, correct bool) error {

	// Begin Transaction
	tx, err := c.db.Connection().Begin()
	if err != nil {
		return err
	}

	// Rollback in case of an error
	defer tx.Rollback()

	// Load Question
	sql := "SELECT ID, BoxID, Next, CorrectlyAnswered, CreatedAt FROM Question WHERE ID = ?;"
	row := tx.QueryRow(sql, id)

	question := models.NewQuestion()

	err = row.Scan(&question.ID, &question.BoxID, &question.Next, &question.CorrectlyAnswered, &question.CreatedAt)
	if err != nil {
		return err
	}

	// Create Learn Unit
	unit := models.NewLearnUnit()
	unit.QuestionID = id
	unit.BoxID = question.BoxID
	unit.Correct = correct
	if question.CorrectlyAnswered == 0 {
		unit.PrevCorrect = false
	} else {
		unit.PrevCorrect = true
	}

	// Save Learn Unit
	sql = "INSERT INTO LearnUnit(QuestionID, BoxID, Correct, PrevCorrect, CreatedAt) VALUES (?, ?, ?, ?, ?);"
	_, err = tx.Exec(sql, unit.QuestionID, unit.BoxID, unit.Correct, unit.PrevCorrect, unit.CreatedAt)
	if err != nil {
		return err
	}

	// Increase CorrectlyAnswered if correct, else set to 0
	if correct {
		question.CorrectlyAnswered++
	} else {
		question.CorrectlyAnswered = 0
	}

	// Calculate Next
	question.CalculateNext()

	// Update Question
	sql = "UPDATE Question SET Next = ?, CorrectlyAnswered = ? WHERE ID = ?;"
	res, err := tx.Exec(sql, question.Next, question.CorrectlyAnswered, question.ID)
	if err != nil {
		return nil
	}

	// Check if update was performed
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Return error if no object was updated
	if rows == 0 {
		return errors.New("Qeustion to update was not found.")
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	// If answer was incorrect and RelearnUntilAccomplished is set, readd to heap
	if SettingsCtrl().Get().RelearnUntilAccomplished && !correct {
		BoxCtrl().reAddQuestionFromHeap(question.BoxID, question.ID)
	} else { // else remove from heap
		BoxCtrl().removeQuestionFromHeap(question.BoxID, question.ID)
	}

	// Publish Notification
	events.Events().Publish(events.Topic("questions"), c)
	events.Events().Publish(events.Topic(fmt.Sprintf("box-%d", question.BoxID)), c)
	events.Events().Publish(events.Topic("stats"), c)

	return nil
}

func (c *QuestionController) UpdateQuestion(qID uint, question *models.Question) error {

	// Begin Transaction
	tx, err := c.db.Connection().Begin()
	if err != nil {
		return err
	}

	// Rollback in case of an error
	defer tx.Rollback()

	// Get BoxID of original question
	sql := "SELECT BoxID FROM Question WHERE ID = ?;"
	row := tx.QueryRow(sql, qID)
	var originalBoxID uint
	err = row.Scan(&originalBoxID)
	if err != nil {
		return err
	}

	// Update category
	sql = "UPDATE Question SET Question = ?, Answer = ?, BoxID = ?, Next = ?, CorrectlyAnswered = ?, CreatedAt = ? WHERE ID = ?;"
	res, err := tx.Exec(sql, question.Question, question.Answer, question.BoxID, question.Next, question.CorrectlyAnswered, question.CreatedAt, qID)
	if err != nil {
		return err
	}

	// Check if update was performed
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Return error if no object was updated
	if rows == 0 {
		return errors.New("Question to update was not found.")
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Check if question was moved into another Box
	if question.BoxID != originalBoxID {
		BoxCtrl().buildHeap(question.BoxID)
		BoxCtrl().buildHeap(originalBoxID)
		events.Events().Publish(events.Topic("boxes"), c)
	}

	// Publish event to force client refresh
	events.Events().Publish(events.Topic("questions"), c)

	return nil

}

func (c *QuestionController) AddQuestion(q *models.Question) (*models.Question, error) {

	// Begin Transaction
	tx, err := c.db.Connection().Begin()
	if err != nil {
		return nil, err
	}

	// Rollback in case of an error
	defer tx.Rollback()

	// Execute insert statement
	sql := "INSERT INTO Question (Question, Answer, BoxID, Next, CorrectlyAnswered, CreatedAt) VALUES (?, ?, ?, ?, ?, ?);"
	res, err := tx.Exec(sql, q.Question, q.Answer, q.BoxID, q.Next, q.CorrectlyAnswered, q.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Update objects ID
	newID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	q.ID = uint(newID)

	// Commit
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Rebuild Heap for Box
	log.Println(BoxCtrl().buildHeap(q.BoxID))

	// Publish events to force client refresh
	events.Events().Publish(events.Topic(fmt.Sprintf("box-%d", q.BoxID)), c)
	events.Events().Publish(events.Topic("questions"), c)
	events.Events().Publish(events.Topic("stats"), c)

	// Return inserted object
	return q, nil

}

func (c *QuestionController) DeleteQuestion(qID uint) error {

	// Begin Transaction
	tx, err := c.db.Connection().Begin()
	if err != nil {
		return err
	}

	// Rollback in case of an error
	defer tx.Rollback()

	// Get BoxID of original question
	sql := "SELECT BoxID FROM Question WHERE ID = ?;"
	row := tx.QueryRow(sql, qID)
	var boxID uint
	err = row.Scan(&boxID)
	if err != nil {
		return err
	}

	// Execute delete statement
	// Because of foreign key contraints: deletes all answers of that Question
	sql = "DELETE FROM Question WHERE ID = ?;"
	res, err := tx.Exec(sql, qID)
	if err != nil {
		return err
	}

	// Check if delete was performed
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Return error if no object was deleted
	if rows <= 0 {
		return errors.New("Question could not be deleted.")
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Rebuild Heap of box
	BoxCtrl().buildHeap(boxID)

	// Publish events to force client refresh
	events.Events().Publish(events.Topic(fmt.Sprintf("box-%d", boxID)), c)
	events.Events().Publish(events.Topic("questions"), c)
	events.Events().Publish(events.Topic("stats"), c)

	return nil

}
