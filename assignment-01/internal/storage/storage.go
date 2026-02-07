package storage

import (
	"sync"
	"github.com/ameiswq/Golang2026/assignment-01/internal/models"
)

type Store struct {
	mu     sync.Mutex
	nextID int
	tasks  map[int]models.Task
}

func NewStore() *Store {
	return &Store{
		nextID: 1,
		tasks:  make(map[int]models.Task),
	}
}

func (s *Store) Create(title string) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := models.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[t.ID] = t
	s.nextID++
	return t
}

func (s *Store) Get(id int) (models.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *Store) List() []models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	return out
}

func (s *Store) UpdateDone(id int, done bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return false
	}
	t.Done = done
	s.tasks[id] = t
	return true
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[id]
	if !ok {
		return false
	}
	delete(s.tasks, id)
	return true
}
