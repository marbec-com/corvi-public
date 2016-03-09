package controllers

import (
	"log"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/events"
	"marb.ec/maf/interfaces"
)

type MockSubscriber struct {
	Notifications map[string]uint
}

func NewMockSubscriber(topics []string) *MockSubscriber {
	s := &MockSubscriber{
		Notifications: make(map[string]uint, len(topics)),
	}

	for _, topic := range topics {
		events.Events().Subscribe(events.Topic(topic), s)
	}
	return s
}

func (s *MockSubscriber) Notify(t interfaces.Topic, p interfaces.Publisher) {
	s.Notifications[t.Topic()]++
}

func (s *MockSubscriber) Assert(topic string, count uint) bool {
	c, ok := s.Notifications[topic]
	if !ok {
		return false
	}
	return c == count
}

type MockSettingsService struct {
	settings *models.Settings
}

func NewMockSettingsService(settings *models.Settings) *MockSettingsService {
	return &MockSettingsService{
		settings: settings,
	}
}

func (s *MockSettingsService) Update() error {
	return nil
}

func (s *MockSettingsService) Get() *models.Settings {
	return s.settings
}

func setupTestDBController(path string) DatabaseService {
	c, err := NewSQLiteDBService(path)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil
	}
	return c
}

func tearDownTestDBController(db DatabaseService) {
	// Close database connection
	err := db.Close()
	if err != nil {
		log.Fatal("Error in Teardown", err)
		return
	}
	// Remove file
	err = db.Destroy()
	if err != nil {
		log.Fatal("Error in Teardown", err)
		return
	}
}
