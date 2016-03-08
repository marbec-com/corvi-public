package controllers

import (
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
