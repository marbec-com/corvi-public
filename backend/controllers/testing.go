package controllers

import (
	"log"
	"marb.ec/maf/events"
	"marb.ec/maf/interfaces"
	"os"
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

func setupTestDBController(path string) *DBController {
	c, err := NewDBController(path)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil
	}
	return c
}

func tearDownTestDBController(db *DBController) {
	// Close database connection
	err := db.Close()
	if err != nil {
		log.Fatal("Error in Teardown", err)
		return
	}
	// Remove file
	err = os.Remove(db.databasePath)
	if err != nil {
		log.Fatal("Error in Teardown", err)
		return
	}
}
