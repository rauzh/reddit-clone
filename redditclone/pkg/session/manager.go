package session

import (
	"fmt"
	"math/rand"
	"net/http"
	"redditclone/pkg/user"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET []byte = []byte("chickenNuggetWithNoSauce")

type SessionsManager struct {
	secret []byte
}

func NewSessionsMem() *SessionsManager {
	return &SessionsManager{
		secret: SECRET,
	}
}

func (sm *SessionsManager) GetSessSecret(token string) (secret []byte) {
	return sm.secret
}

func (sm *SessionsManager) Check(r *http.Request, userRepo *user.UserRepo) (*Session, error) {
	inToken := r.Header.Get("Authorization")
	if len(inToken) <= 7 {
		fmt.Println("session.manager: no token")
		return nil, ErrNoAuth
	}
	inToken = inToken[7:]

	secret := sm.GetSessSecret(inToken)
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return secret, nil
	}

	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err != nil || !token.Valid {
		fmt.Println("session.manager: bad token")
		return nil, fmt.Errorf("bad token")
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("session.manager: bad tokenClaims")
		return nil, fmt.Errorf("bad token")
	}

	userPayload := payload["user"].(map[string]interface{})
	sess := &Session{
		UserID:   userPayload["id"].(string),
		Username: userPayload["username"].(string),
	}

	if !userRepo.UserExist(sess.Username, sess.UserID) {
		return nil, fmt.Errorf("no user")
	}

	return sess, nil
}

func (sm *SessionsManager) Create(userID string, username string, secret []byte) (*Session, error) {
	t := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username": username,
			"id":       userID,
		},
		"iat": t.Unix(),
		"exp": t.Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("session.manager: cant create token")
		return nil, fmt.Errorf("bad token")
	}

	rand.Seed(time.Now().UnixNano())
	randID := make([]byte, 16)
	rand.Read(randID)
	sess := &Session{
		UserID:   userID,
		Username: username,
		Token:    tokenString,
	}

	return sess, nil
}
