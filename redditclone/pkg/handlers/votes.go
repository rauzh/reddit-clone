package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"redditclone/pkg/session"

	"github.com/gorilla/mux"
)

func (h *ItemsHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, found := vars["id"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByItemID(id)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, `no item`, http.StatusNotFound)
		return
	}

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ItemsRepo.Upvote(sess, item)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	data, errMarshal := json.Marshal(item)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("Update for item with id FAILED, JSON marshall error")
		return
	}

	w.WriteHeader(http.StatusOK)

	_, errWrite := w.Write(data)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("Update for item with id FAILED, body write error")
		return
	}
	h.Logger.Infof("Update for item with id : %v", id)
}

func (h *ItemsHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, found := vars["id"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByItemID(id)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, `no item`, http.StatusNotFound)
		return
	}

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("FLAG:DV(1);\tsess =", sess)
	err = h.ItemsRepo.Downvote(sess, item)
	fmt.Println("FLAG:DV(2)")
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	data, errMarshal := json.Marshal(item)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("Update for item with id FAILED, JSON marshall error")
		return
	}

	w.WriteHeader(http.StatusOK)

	_, errWrite := w.Write(data)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("Update for item with id FAILED, body write error")
		return
	}
	h.Logger.Infof("Update for item with id : %v", id)
}

func (h *ItemsHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, found := vars["id"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByItemID(id)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, `no item`, http.StatusNotFound)
		return
	}

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ItemsRepo.Unvote(sess, item)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	data, errMarshal := json.Marshal(item)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("Update for item with id FAILED, JSON marshall error")
		return
	}

	w.WriteHeader(http.StatusOK)

	_, errWrite := w.Write(data)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Infof("Update for item with id FAILED, body write error")
		return
	}
	h.Logger.Infof("Update for item with id : %v", id)
}
