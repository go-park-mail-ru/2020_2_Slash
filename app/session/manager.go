package session

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2020_2_Slash/app/user"
)

var (
	ErrNoAuth = errors.New("No session found")
)

type Session struct {
	ID        string
	UserID    uint64
	ExpiresAt time.Time
}

func NewSession(userID uint64) *Session {
	randID := uuid.NewV4().String()
	return &Session{
		ID:        randID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(10 * time.Hour),
	}
}

type SessionManager struct {
	data map[string]*Session
	mu   *sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		data: make(map[string]*Session, 10),
		mu:   &sync.Mutex{},
	}
}

func (sm *SessionManager) Create(user *user.User) *Session {
	session := NewSession(user.ID)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[session.ID] = session
	return session
}

func (sm *SessionManager) Check(cookieValue string) (*Session, error) {
	sm.mu.Lock()
	sess, has := sm.data[cookieValue]
	sm.mu.Unlock()
	if !has || sess.ExpiresAt.After(time.Now()) {
		// Doesn't exist or expired
		return nil, ErrNoAuth
	}
	return sess, nil
}
