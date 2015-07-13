package forest

import (
	"github.com/ursiform/bear"
	"time"
)

type SessionManager interface {
	Create(sessionID string, userID string, userJSON string, ctx *bear.Context) error
	CreateEmpty(sessionID string, ctx *bear.Context)
	Delete(sessionID string, userID string) error
	Marshal(ctx *bear.Context) ([]byte, error)
	Read(sessionID string) (userID string, userJSON string, err error)
	Revoke(userID string) error
	Update(sessionID string, userID string, userJSON string, duration time.Duration) error
}
