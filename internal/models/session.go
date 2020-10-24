package models

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Session struct {
	ID        uint64
	Value     string
	UserID    uint64
	ExpiresAt time.Time
}

func NewSession(userID uint64) *Session {
	randValue := uuid.NewV4().String()
	expiresDur := consts.ExpiresDuration
	return &Session{
		Value:     randValue,
		UserID:    userID,
		ExpiresAt: time.Now().Add(expiresDur),
	}
}
