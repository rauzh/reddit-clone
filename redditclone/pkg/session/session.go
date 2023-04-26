package session

import (
	"context"
	"errors"
)

type Session struct {
	UserID   string
	Username string
	Token    string
}

type sessKey string

var SessionKey sessKey = "sessionKey"

var (
	ErrNoAuth = errors.New("no session found")
)

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}
