package session

import (
	"context"
	"errors"
)

// MockStore is a mock implementation of Store for testing.
type MockStore struct {
	Sessions   map[string]*Session
	PutFunc    func(ctx context.Context, s *Session) error
	GetFunc    func(ctx context.Context, sessionID string) (*Session, error)
	UpdateFunc func(ctx context.Context, s *Session) error
}

func (m *MockStore) Put(ctx context.Context, s *Session) error {
	if m.PutFunc != nil {
		return m.PutFunc(ctx, s)
	}
	m.Sessions[s.SessionID] = s
	return nil
}

func (m *MockStore) Get(ctx context.Context, sessionID string) (*Session, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, sessionID)
	}
	s, ok := m.Sessions[sessionID]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *MockStore) Update(ctx context.Context, s *Session) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, s)
	}
	m.Sessions[s.SessionID] = s
	return nil
}
