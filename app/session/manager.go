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

func (sm *SessionManager) Get(sessID string) (*Session, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sess, has := sm.data[sessID]
	return sess, has
}

func (sm *SessionManager) IsValid(session *Session) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	hasExpired := session.ExpiresAt.Before(time.Now())
	if hasExpired {
		return false
	}
	return true
}
