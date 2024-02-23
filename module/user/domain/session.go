package domain

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	id           uuid.UUID
	userId       uuid.UUID
	refreshToken string
	accessExpAt  time.Time
	refreshExpAt time.Time
}

func NewSession(id uuid.UUID, userId uuid.UUID, refreshToken string, accessExpAt time.Time, refreshExpAt time.Time) *Session {
	return &Session{id: id, userId: userId, refreshToken: refreshToken, accessExpAt: accessExpAt, refreshExpAt: refreshExpAt}
}

func (s Session) Id() uuid.UUID {
	return s.id
}

func (s Session) UserId() uuid.UUID {
	return s.userId
}

func (s Session) RefreshToken() string {
	return s.refreshToken
}

func (s Session) AccessExpAt() time.Time {
	return s.accessExpAt
}

func (s Session) RefreshExpAt() time.Time {
	return s.refreshExpAt
}
