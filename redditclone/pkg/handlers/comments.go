package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/items"
	"redditclone/pkg/session"

	"github.com/gorilla/mux"
)

type CommentJSON struct {
	Message string `json:"comment"`
}

func (h *ItemsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, found := vars["id"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	item, err := h.ItemsRepo.GetByItemID(itemID)
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	data, errReadBody := ioutil.ReadAll(r.Body)
	if errReadBody != nil {
		http.Error(w, `can't read req body`, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	message := new(CommentJSON)
	err = json.Unmarshal(data, message)
	if err != nil {
		http.Error(w, `Cant unmarshal JSON`, http.StatusBadRequest)
	}

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	lastID, err := h.ItemsRepo.AddComment(sess, item, message.Message)
	if err == items.ErrEmptyComm {
		errMsg(w, "you can't post empty comment")
		return
	}
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	resp, errMarshal := json.Marshal(item)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("Insert with id LastInsertId: %v FAILED: JSON marshal error", lastID)
		return
	}
	_, errWrite := w.Write(resp)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("Insert with id LastInsertId: %v FAILED: body write error", lastID)
		return
	}
	h.Logger.Infof("Insert with id LastInsertId: %v", lastID)
}

func (h *ItemsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, found := vars["postID"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	comID, found := vars["comID"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByItemID(postID)
	if (err != nil) || (item == nil) {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ItemsRepo.DeleteComment(sess, comID, item)
	if errors.Is(err, items.ErrNotCommAuthor) {
		errMsg(w, "you're not the author")
		return
	}
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, errWrite := w.Write(data)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func errMsg(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusForbidden)
	resp, errMarshal := json.Marshal(map[string]interface{}{"message": msg})
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, errWrite := w.Write(resp)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
