// Package session provides interfaces and implementations for session storage in DynamoDB.
package session

import (
	"context"
)

// Session represents an OIDC authentication session.
type Session struct {
	SessionID   string
	ClientID    string
	AccountID   string
	RoleName    string
	State       string
	Status      string // e.g., "pending", "complete"
	AccessToken string // Encrypted at rest
	ExpiresAt   int64
}

// Store defines the interface for session storage.
type Store interface {
	Put(ctx context.Context, s *Session) error
	Get(ctx context.Context, sessionID string) (*Session, error)
	Update(ctx context.Context, s *Session) error
}
