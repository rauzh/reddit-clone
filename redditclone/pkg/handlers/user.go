package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"go.uber.org/zap"
)

type UserHandler struct {
	Logger   *zap.SugaredLogger
	UserRepo *user.UserRepo
	Sessions *session.SessionsManager
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	username, password, err := getUserParams(r)
	if err != nil {
		http.Error(w, `Incorrect JSON`, http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.Authorize(username, password)
	if err == user.ErrBadPass {
		w.WriteHeader(http.StatusUnauthorized)
		resp, errMarshal := json.Marshal(map[string]interface{}{"message": "invalid password"})
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Logger.Infof("create session FAILED, JSON marshal error")
			return
		}

		_, errWrite := w.Write(resp)
		if errWrite != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Logger.Infof("create session FAILED, body write error")
			return
		}
		return
	}
	if err == user.ErrNoUser {
		w.WriteHeader(http.StatusUnauthorized)
		resp, errMarshal := json.Marshal(map[string]interface{}{"message": "user not found"})
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Logger.Infof("create session FAILED, JSON marshal error")
			return
		}

		_, errWrite := w.Write(resp)
		if errWrite != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Logger.Infof("create session FAILED, body write error")
			return
		}
		return
	}

	secret := h.Sessions.GetSessSecret("")
	sess, err := h.Sessions.Create(fmt.Sprint(u.ID), u.Username, secret)
	if err != nil {
		http.Error(w, `Cant create session`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, errMarshal := json.Marshal(map[string]interface{}{"token": sess.Token})
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("create session FAILED, JSON marshal error")
		return
	}

	_, errWrite := w.Write(resp)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("create session FAILED, body write error")
		return
	}
	h.Logger.Infof("created session for %v", sess.UserID)
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	username, password, err := getUserParams(r)
	if err != nil {
		http.Error(w, `Incorrect JSON`, http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.Registration(username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp, errMarshal := json.Marshal(map[string]interface{}{
			"errors": []map[string]interface{}{{
				"location": "body",
				"param":    "username",
				"value":    username,
				"msg":      "already exist",
			}}})
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Logger.Infof("create session FAILED, JSON marshal error")
			return
		}

		_, errWrite := w.Write(resp)
		if errWrite != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Logger.Infof("create session FAILED, body write error")
			return
		}
		return
	}
	// fmt.Print("\n\n\n", u.ID, "\n\n\n")
	secret := h.Sessions.GetSessSecret("")
	sess, err := h.Sessions.Create(fmt.Sprintf("%d", u.ID), u.Username, secret)
	if err != nil {
		http.Error(w, `Cant create session`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	resp, errMarshal := json.Marshal(map[string]interface{}{"token": sess.Token})
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("created session for %x  FAILED, JSON marshal error", sess.UserID)
		return
	}

	_, errWrite := w.Write(resp)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("created session for %x  FAILED, body write error", sess.UserID)
		return
	}

	h.Logger.Infof("created session for %x", sess.UserID)
}

type UserJSON struct {
	Username string
	Password string
}

func getUserParams(r *http.Request) (string, string, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", "", fmt.Errorf("body reading error")
	}
	defer r.Body.Close()
	userJSON := &UserJSON{}
	err = json.Unmarshal(body, userJSON)
	if err != nil {
		return "", "", fmt.Errorf("incorrect JSON")
	}
	return userJSON.Username, userJSON.Password, nil
}
