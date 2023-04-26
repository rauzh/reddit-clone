package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/session"

	"redditclone/pkg/items"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ItemsHandler struct {
	Tmpl      *template.Template
	ItemsRepo *items.ItemsRepo
	Logger    *zap.SugaredLogger
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	elems, err := h.ItemsRepo.GetAll()
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(elems)
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

func (h *ItemsHandler) Add(w http.ResponseWriter, r *http.Request) {
	data, errReadBody := ioutil.ReadAll(r.Body)
	if errReadBody != nil {
		http.Error(w, `can't read req body`, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	item := items.NewItem()
	err := json.Unmarshal(data, item)
	if err != nil {
		http.Error(w, `Cant unmarshal JSON`, http.StatusBadRequest)
	}
	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	lastID, err := h.ItemsRepo.Add(sess, item)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, errMarshal := json.Marshal(item)
	if errMarshal != nil {
		http.Error(w, `db err`, http.StatusInternalServerError)
		return
	}
	_, errWrite := w.Write(resp)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Insert with id LastInsertId: %v", lastID)
}

func (h *ItemsHandler) Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, found := vars["id"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByItemID(id)
	if err != nil {
		http.Error(w, `db err`, http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, `no item`, http.StatusNotFound)
		return
	}
	item.Views++

	data, errMarshal := json.Marshal(item)
	if errMarshal != nil {
		http.Error(w, `db err`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, errWrite := w.Write(data)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, found := vars["id"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	sess, errSess := session.SessionFromContext(r.Context())
	if errSess != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err := h.ItemsRepo.Delete(sess, id)
	if errors.Is(err, items.ErrNotPostAuthor) {
		w.WriteHeader(http.StatusForbidden)
		resp, errMarshal := json.Marshal(map[string]interface{}{"message": "you're not the author"})
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, errWrite := w.Write(resp)
		if errWrite != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, errMarshal := json.Marshal(map[string]string{
		"message": "success",
	})
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, errWrite := w.Write(respJSON)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) Category(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	catName, found := vars["catName"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	items, err := h.ItemsRepo.GetByCategory(catName)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	data, errMarshal := json.Marshal(items)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, errWrite := w.Write(data)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) UserItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, found := vars["username"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	items, err := h.ItemsRepo.GetByUserID(username)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	data, errMarshal := json.Marshal(items)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, errWrite := w.Write(data)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
