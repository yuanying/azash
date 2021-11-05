package caches

import (
	"encoding/json"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"

	bo "github.com/yuanying/azash/pkgs/books"
	ca "github.com/yuanying/azash/pkgs/caches"
)

type Handlers struct {
	log      logr.Logger
	bookList *bo.BookList
	cache    *ca.Manager
}

func NewHandler(log logr.Logger, bookList *bo.BookList, cache *ca.Manager) Handlers {
	return Handlers{
		log:      log,
		bookList: bookList,
		cache:    cache,
	}
}

func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	bookID := mux.Vars(r)["books"]
	book, err := h.bookList.Get(bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if book == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	c, err := h.cache.Get(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}
