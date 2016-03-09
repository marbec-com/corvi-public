package controllers

import (
	"log"
	"marb.ec/corvi-backend/models"
	"os"
	"path"
	"time"
)

func InitControllerSingletons(db DatabaseService, s SettingsService) {

	b, err := NewBoxController(db, s)
	if err != nil {
		log.Fatal(err)
	}
	BoxControllerSingleton = b

	c, err := NewCategoryController(db, s)
	if err != nil {
		log.Fatal(err)
	}
	CategoryControllerSingleton = c

	q, err := NewQuestionController(db, s)
	if err != nil {
		log.Fatal(err)
	}
	QuestionControllerSingleton = q

	StatsControllerSingleton = NewStatsController(db, s)

	populateDummyData()

}

func GenerateUserDataPath(fileName string) string {
	userPath := os.Getenv("USER_DATA")
	if userPath != "" {
		fileName = path.Join(userPath, fileName)
	}
	return fileName
}

func populateDummyData() {
	dummyData := os.Getenv("DUMMY_DATA")
	if dummyData != "1" {
		return
	}

	var mockCategories = []*models.Category{
		&models.Category{
			Name:      "Computer Science",
			CreatedAt: time.Now(),
		},
		&models.Category{
			Name:      "Vocabulary",
			CreatedAt: time.Now(),
		},
	}

	catCtrl := CategoryCtrl()
	for k, cat := range mockCategories {
		newCat, err := catCtrl.AddCategory(cat)
		if err == nil {
			mockCategories[k] = newCat
		}
	}

	var mockBoxes = []*models.Box{
		&models.Box{
			Name:             "SQL Statements",
			CategoryID:       mockCategories[0].ID,
			QuestionsToLearn: 2,
			QuestionsTotal:   2,
			QuestionsLearned: 0,
			CreatedAt:        time.Now(),
		},
		&models.Box{
			Name:             "English - Kitchen",
			CategoryID:       mockCategories[1].ID,
			QuestionsToLearn: 1,
			QuestionsTotal:   1,
			QuestionsLearned: 0,
			CreatedAt:        time.Now(),
		},
		&models.Box{
			Name:             "French - Cuisine",
			CategoryID:       mockCategories[1].ID,
			QuestionsToLearn: 0,
			QuestionsTotal:   1,
			QuestionsLearned: 1,
			CreatedAt:        time.Now(),
		},
	}

	boxCtrl := BoxCtrl()
	for k, box := range mockBoxes {
		newBox, err := boxCtrl.AddBox(box)
		if err == nil {
			mockBoxes[k] = newBox
		}
	}

	var mockQuestions = []*models.Question{
		&models.Question{
			Question:          "A Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 5,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 20,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Update Statement?",
			Answer:            "UPDATE table SET key = value",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Insert Statement?",
			Answer:            "INSERT INTO table (key, key) VALUES (value, value)",
			BoxID:             mockBoxes[0].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 1,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Küche",
			Answer:            "Kitchen",
			BoxID:             mockBoxes[1].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 0,
			CreatedAt:         time.Now(),
		},
		&models.Question{
			Question:          "Küche",
			Answer:            "Cuisine",
			BoxID:             mockBoxes[2].ID,
			Next:              time.Now(),
			CorrectlyAnswered: 3,
			CreatedAt:         time.Now(),
		},
	}

	questionCtrl := QuestionCtrl()
	for k, q := range mockQuestions {
		newQ, err := questionCtrl.AddQuestion(q)
		if err == nil {
			mockQuestions[k] = newQ
		}
	}

	mockCategories = nil
	mockBoxes = nil
	mockQuestions = nil

}
