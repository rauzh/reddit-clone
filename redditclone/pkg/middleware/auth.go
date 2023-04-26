package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"redditclone/pkg/session"
	"redditclone/pkg/user"
)

type AuthURL struct {
	URL    *regexp.Regexp
	Method string
}

func Auth(sm *session.SessionsManager, next http.Handler, userRepo *user.UserRepo, authUrls []AuthURL) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, item := range authUrls {
			matched := item.URL.MatchString(r.URL.Path)
			if !matched {
				continue
			}
			if item.Method != r.Method {
				continue
			}
			sess, err := sm.Check(r, userRepo)
			if err != nil {
				fmt.Println("no auth")
				w.WriteHeader(http.StatusUnauthorized)
				resp, errMarshal := json.Marshal(map[string]interface{}{"message": "no auth"})
				if errMarshal != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Print("auth: marsh error")
					return
				}

				_, errWrite := w.Write(resp)
				if errWrite != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Print("auth: body write error")
					return
				}
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			ctx := context.WithValue(r.Context(), session.SessionKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}
