package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/events"
	"sync"
)

var BoxControllerSingleton *BoxController

func BoxCtrl() *BoxController {
	return BoxControllerSingleton
}

type BoxController struct {
	db        DatabaseService
	settings  SettingsService
	heapCache map[uint]*models.QuestionHeap
	sync.Mutex
}

func NewBoxController(db DatabaseService, settings SettingsService) (*BoxController, error) {
	b := &BoxController{
		db:       db,
		settings: settings,
	}
	err := b.createTables()
	if err != nil {
		return nil, err
	}
	b.heapCache = make(map[uint]*models.QuestionHeap, 0)
	return b, nil
}

func (c *BoxController) createTables() error {

	// Create table, only if it not already exists
	// Includes foreign key constraint to Category table. We are not allowed to delete a Category that still has Boxes assigned.
	sql := "CREATE TABLE IF NOT EXISTS Box (ID INTEGER PRIMARY KEY ASC NOT NULL, Name VARCHAR (255) NOT NULL, Description TEXT NOT NULL, CategoryID INTEGER REFERENCES Category (ID) NOT NULL, CreatedAt DATETIME NOT NULL);"
	_, err := c.db.Connection().Exec(sql)
	if err != nil {
		return err
	}

	// Create additonal view with question meta data
	sql = "CREATE VIEW IF NOT EXISTS BoxWithMeta AS SELECT ID, Name, Description, CategoryID, (SELECT COUNT(*) FROM Question WHERE BoxID = Box.ID) AS QuestionsTotal, (SELECT COUNT(*) FROM Question WHERE BoxID = Box.ID AND CorrectlyAnswered > 0) AS QuestionsLearned, CreatedAt FROM Box;"
	_, err = c.db.Connection().Exec(sql)

	return err

}

func (c *BoxController) LoadBoxes() ([]*models.Box, error) {
	return c.LoadBoxesOfCategory(0)
}

func (c *BoxController) LoadBoxesOfCategory(id uint) ([]*models.Box, error) {

	// Select all boxes of category
	var rows *sql.Rows
	var err error
	if id == 0 {
		sql := "SELECT ID, Name, Description, CategoryID, QuestionsTotal, QuestionsLearned, CreatedAt FROM BoxWithMeta;"
		rows, err = c.db.Connection().Query(sql)
	} else {
		sql := "SELECT ID, Name, Description, CategoryID, QuestionsTotal, QuestionsLearned, CreatedAt FROM BoxWithMeta WHERE CategoryID = ?;"
		rows, err = c.db.Connection().Query(sql, id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create empty result set
	result := make([]*models.Box, 0)

	for rows.Next() {
		// Create new Box object
		newBox := models.NewBox()
		// Populate
		err = rows.Scan(&newBox.ID, &newBox.Name, &newBox.Description, &newBox.CategoryID, &newBox.QuestionsTotal, &newBox.QuestionsLearned, &newBox.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Questions To Learn from Heap
		heap, ok := c.heapCache[newBox.ID]
		if ok {
			newBox.QuestionsToLearn = uint(heap.Length())
		}

		// Append to result set
		result = append(result, newBox)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (c *BoxController) LoadBox(id uint) (*models.Box, error) {

	// Select box with matching ID
	sql := "SELECT ID, Name, Description, CategoryID, QuestionsTotal, QuestionsLearned, CreatedAt FROM BoxWithMeta WHERE ID = ?;"
	row := c.db.Connection().QueryRow(sql, id)

	// Create new Category object
	newBox := models.NewBox()

	// Populate
	err := row.Scan(&newBox.ID, &newBox.Name, &newBox.Description, &newBox.CategoryID, &newBox.QuestionsTotal, &newBox.QuestionsLearned, &newBox.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Questions To Learn from Heap
	heap, ok := c.heapCache[newBox.ID]
	if ok {
		newBox.QuestionsToLearn = uint(heap.Length())
	}

	return newBox, nil

}

func (c *BoxController) removeQuestionFromHeap(id, qID uint) error {

	// Get heap from cache
	c.Lock()
	heap, ok := c.heapCache[id]
	c.Unlock()
	if !ok {
		return errors.New("Heap not found.")
	}

	heap.Lock()

	// If top element in heap is given question
	if heap.Peek().ID == qID {
		// Remove top element
		heap.Min()
	}

	heap.Unlock()

	return nil

}

func (c *BoxController) reAddQuestionFromHeap(id, qID uint) error {

	// Get heap from cache
	c.Lock()
	heap, ok := c.heapCache[id]
	c.Unlock()
	if !ok {
		return errors.New("Heap not found.")
	}

	heap.Lock()

	// If top element in heap is given question
	if heap.Peek().ID == qID {
		// Readd question to heap
		heap.Add(heap.Min())
	}

	heap.Unlock()

	return nil

}

func (c *BoxController) GetQuestionToLearn(id uint) (*models.Question, error) {

	// Get heap from cache
	c.Lock()
	heap, ok := c.heapCache[id]
	c.Unlock()
	if !ok || heap.Length() == 0 { // Build new heap if non existent or empty
		c.BuildHeap(id)
	}

	// Return first question in heap without removing it, may be nil if no questions due
	return heap.Peek(), nil

}

func (c *BoxController) BuildHeap(id uint) error {

	// Create new heap
	heap := models.NewQuestionHeap()

	// Get capacity

	// If RelearnUntilAccomplished: capacity = max - correctly answered questions today
	// else: capacity = max - answered questions today
	sql := "SELECT COUNT(*) FROM LearnUnit WHERE date(CreatedAt) = date('now') AND BoxID = ?;"
	if c.settings.Get().RelearnUntilAccomplished {
		sql = "SELECT COUNT(*) FROM LearnUnit WHERE date(CreatedAt) = date('now') AND Correct = 1 AND BoxID = ?;"
	}
	row := c.db.Connection().QueryRow(sql, id)
	var learned int
	err := row.Scan(&learned)
	if err != nil {
		return err
	}

	capacity := int(c.settings.Get().MaxDailyQuestionsPerBox) - learned
	if capacity < 0 {
		capacity = 0
	}

	// Get due questions

	// If RelearnUntilAccomplished: Questions are due which were correctly answered today
	// else: Questions are due which were answered today
	sql = "SELECT ID, Question, Answer, BoxID, Next, CorrectlyAnswered, CreatedAt FROM Question WHERE datetime(Next) < datetime('now', 'start of day', '+1 day') AND BoxID = ? AND ID NOT IN(SELECT QuestionID FROM LearnUnit WHERE date(CreatedAt) = date('now') AND BoxID = ?) ORDER BY Next DESC LIMIT ?;"
	if c.settings.Get().RelearnUntilAccomplished {
		sql = "SELECT ID, Question, Answer, BoxID, Next, CorrectlyAnswered, CreatedAt FROM Question WHERE datetime(Next) < datetime('now', 'start of day', '+1 day') AND BoxID = ? AND ID NOT IN(SELECT QuestionID FROM LearnUnit WHERE date(CreatedAt) = date('now') AND Correct = 1 AND BoxID = ?) ORDER BY Next DESC LIMIT ?;"
	}

	qRows, err := c.db.Connection().Query(sql, id, id, capacity)
	if err != nil {
		return err
	}
	defer qRows.Close()

	for qRows.Next() {
		// Create new Box object
		newQuestion := models.NewQuestion()
		// Populate
		err = qRows.Scan(&newQuestion.ID, &newQuestion.Question, &newQuestion.Answer, &newQuestion.BoxID, &newQuestion.Next, &newQuestion.CorrectlyAnswered, &newQuestion.CreatedAt)
		if err != nil {
			return err
		}

		if capacity > 0 {
			heap.Add(newQuestion)
			capacity--
		} else { // Early break
			return nil
		}

	}

	err = qRows.Err()
	if err != nil {
		return err
	}

	// Overwrite old heap
	c.Lock()
	c.heapCache[id] = heap
	c.Unlock()

	return nil
}

func (c *BoxController) BuildHeaps() error {

	// Create new heap cache
	newHeapCache := make(map[uint]*models.QuestionHeap)

	// Get all capacities
	capacities := make(map[uint]uint, 0)

	// If RelearnUntilAccomplished: capacity = max - correctly answered questions today
	// else: capacity = max - answered questions today
	sql := "SELECT BoxID, COUNT(*) FROM LearnUnit WHERE date(CreatedAt) = date('now') GROUP BY BoxID;"
	if c.settings.Get().RelearnUntilAccomplished {
		sql = "SELECT BoxID, COUNT(*) FROM LearnUnit WHERE date(CreatedAt) = date('now') AND Correct = 1 GROUP BY BoxID;"
	}
	rows, err := c.db.Connection().Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var boxID uint
		var count int
		err = rows.Scan(&boxID, &count)
		if err != nil {
			return err
		}
		cap := int(c.settings.Get().MaxDailyQuestionsPerBox) - count
		if cap < 0 {
			cap = 0
		}
		capacities[boxID] = uint(cap)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	// Get all due questions
	// If RelearnUntilAccomplished: Questions are due which were correctly answered today
	// else: Questions are due which were answered today
	sql = "SELECT ID, Question, Answer, BoxID, Next, CorrectlyAnswered, CreatedAt FROM Question WHERE datetime(Next) < datetime('now', 'start of day', '+1 day') AND ID NOT IN(SELECT QuestionID FROM LearnUnit WHERE date(CreatedAt) = date('now'));"
	if c.settings.Get().RelearnUntilAccomplished {
		sql = "SELECT ID, Question, Answer, BoxID, Next, CorrectlyAnswered, CreatedAt FROM Question WHERE datetime(Next) < datetime('now', 'start of day', '+1 day') AND ID NOT IN(SELECT QuestionID FROM LearnUnit WHERE date(CreatedAt) = date('now') AND Correct = 1);"
	}

	qRows, err := c.db.Connection().Query(sql)
	if err != nil {
		return err
	}
	defer qRows.Close()

	for qRows.Next() {
		// Create new Box object
		newQuestion := models.NewQuestion()
		// Populate
		err = qRows.Scan(&newQuestion.ID, &newQuestion.Question, &newQuestion.Answer, &newQuestion.BoxID, &newQuestion.Next, &newQuestion.CorrectlyAnswered, &newQuestion.CreatedAt)
		if err != nil {
			return err
		}

		heap, ok := newHeapCache[newQuestion.BoxID]
		if !ok {
			heap = models.NewQuestionHeap()
			newHeapCache[newQuestion.BoxID] = heap
		}

		cap, ok := capacities[newQuestion.BoxID]
		if !ok {
			cap = c.settings.Get().MaxDailyQuestionsPerBox
			capacities[newQuestion.BoxID] = cap
		}

		if cap > 0 {
			heap.Add(newQuestion)
			capacities[newQuestion.BoxID]--
		}

	}

	err = qRows.Err()
	if err != nil {
		return err
	}

	c.Lock()
	c.heapCache = newHeapCache
	c.Unlock()

	return nil

}

func (c *BoxController) UpdateBox(boxID uint, box *models.Box) error {

	// Begin Transaction
	tx, err := c.db.Connection().Begin()
	if err != nil {
		return err
	}

	// Rollback in case of an error
	defer tx.Rollback()

	// Update category
	sql := "UPDATE Box SET Name = ?, Description = ?, CategoryID = ?, CreatedAt = ? WHERE ID = ?;"
	res, err := tx.Exec(sql, box.Name, box.Description, box.CategoryID, box.CreatedAt, boxID)
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
		return errors.New("Box to update was not found.")
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Publish event to force client refresh
	events.Events().Publish(events.Topic(fmt.Sprintf("box-%d", boxID)), c)

	return nil

}

func (c *BoxController) AddBox(box *models.Box) (*models.Box, error) {

	// Begin Transaction
	tx, err := c.db.Connection().Begin()
	if err != nil {
		return nil, err
	}

	// Rollback in case of an error
	defer tx.Rollback()

	// Execute insert statement
	sql := "INSERT INTO Box (Name, Description, CategoryID, CreatedAt) VALUES (?, ?, ?, ?);"
	res, err := tx.Exec(sql, box.Name, box.Description, box.CategoryID, box.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Update objects ID
	newID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	box.ID = uint(newID)

	// Commit
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Publish event to force client refresh
	events.Events().Publish(events.Topic("boxes"), c)
	events.Events().Publish(events.Topic("stats"), c)

	// Return inserted object
	return box, nil

}

func (c *BoxController) DeleteBox(boxID uint) error {

	// Begin Transaction
	tx, err := c.db.Connection().Begin()
	if err != nil {
		return err
	}

	// Rollback in case of an error
	defer tx.Rollback()

	// Execute delete statement
	// Because of foreign key contraints: deletes all quequestions of that Box
	// and all LearnUnits
	sql := "DELETE FROM Box WHERE ID = ?;"
	res, err := tx.Exec(sql, boxID)
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
		return errors.New("Box could not be deleted.")
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	events.Events().Publish(events.Topic("questions"), c)
	events.Events().Publish(events.Topic("boxes"), c)
	events.Events().Publish(events.Topic("stats"), c)

	return nil

}
