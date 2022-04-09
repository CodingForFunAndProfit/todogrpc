package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/CodingForFunAndProfit/todogrpc/internal/jobs/v1"
)

type TodoSvc struct {
	mu    sync.Mutex
	todos []*jobs.Todo
}

func (s *TodoSvc) Add(_ context.Context, todos *jobs.Todos) (
	*jobs.Response, error) {
	s.mu.Lock()
	s.todos = append(s.todos, jobs.Todos...)
	s.mu.Unlock()
	return &jobs.Response{Message: "ok"}, nil
}

func (s *TodoSvc) Completed(_ context.Context,
	req *jobs.CompletedRequest) (*jobs.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.todos == nil || req.TodoNumber < 1 ||
		int(req.TodoNumber) > len(s.todos) {
		return nil, fmt.Errorf("todo %d not found", req.TodoNumber)
	}
	s.todos[req.TodoNumber-1].Completed = true
	return &jobs.Response{Message: "ok"}, nil
}

func (s *TodoSvc) List(_ context.Context, _ *jobs.Empty) (
	*jobs.Todos, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.todos == nil {
		s.todos = make([]*jobs.Todo, 0)
	}
	return &jobs.Todos{Todos: s.todos}, nil
}

func (s *TodoSvc) Service() *jobs.TodoService {
	return &jobs.TodoService{
		Add:      s.Add,
		Complete: s.Completed,
		List:     s.List,
	}
}
