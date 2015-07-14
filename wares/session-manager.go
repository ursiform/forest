// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

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
