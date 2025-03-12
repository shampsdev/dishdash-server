package session

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"dishdash.ru/pkg/domain"
)

type Session struct {
	lock         sync.Mutex
	recordEvents map[string]struct{}
	shorteners   map[string]Shortener

	Steps []*Step
}

type Step struct {
	Name       string
	Events     map[string][]EventData
	respAmount atomic.Uint32
}

type EventData struct {
	Event string
	Data  interface{}
}

func New() *Session {
	return &Session{
		recordEvents: make(map[string]struct{}),
		shorteners:   make(map[string]Shortener),
	}
}

func (s *Session) NewStep(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	step := &Step{
		Name:   name,
		Events: make(map[string][]EventData),
	}
	s.Steps = append(s.Steps, step)
}

// RecordEvent records an event for the current step
// Even if the event should not be recorded, the counter will increase
func (s *Session) RecordEvent(user *domain.User, event string, data interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	step := s.Steps[len(s.Steps)-1]

	step.respAmount.Add(1)
	if _, ok := s.recordEvents[event]; ok {
		step.Events[user.Name] = append(step.Events[user.Name], EventData{
			Event: event,
			Data:  data,
		})
	}
}

func (s *Session) SetRecordEvents(events ...string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.recordEvents = make(map[string]struct{})
	for _, event := range events {
		s.recordEvents[event] = struct{}{}
	}
}

func (s *Session) UseShortener(event string, sh Shortener) {
	s.shorteners[event] = sh
}

// WaitNResponses waits for n responses to be recorded with [Session.RecordEvent]
func (s *Session) WaitNResponses(n uint32) {
	timeout := time.NewTimer(time.Second * 10).C
	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return
		case <-ticker.C:
			actual := s.Steps[len(s.Steps)-1].respAmount.Load()
			if actual == n {
				return
			}
		}
	}
}

func (s *Session) SaveToFile(file string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.saveStepsToFile(s.Steps, file)
}

func (s *Session) SaveToFileShortened(file string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	stepsShortened := make([]*Step, len(s.Steps))
	for i, step := range s.Steps {
		stepsShortened[i] = &Step{
			Name:   step.Name,
			Events: make(map[string][]EventData),
		}
		for user, events := range step.Events {
			for _, event := range events {
				data := event.Data
				if sh, ok := s.shorteners[event.Event]; ok {
					data = sh(data)
				}
				stepsShortened[i].Events[user] = append(stepsShortened[i].Events[user], EventData{
					Event: event.Event,
					Data:  data,
				})
			}
		}
	}
	return s.saveStepsToFile(stepsShortened, file)
}

func (s *Session) saveStepsToFile(steps []*Step, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(steps)
	if err != nil {
		return err
	}
	return nil
}

func LoadFromFile(file string) (*Session, error) {
	s := &Session{}

	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&s.Steps)
	if err != nil {
		return nil, err
	}
	return s, nil
}
