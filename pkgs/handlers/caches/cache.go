package caches

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"

	bo "github.com/yuanying/azash/pkgs/books"
	ca "github.com/yuanying/azash/pkgs/caches"
	"github.com/yuanying/azash/pkgs/thumbnail"
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
	_, c, err := h.handleBook(w, r)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func (h *Handlers) Thumbnail(w http.ResponseWriter, r *http.Request) {
	book, c, err := h.handleBook(w, r)
	if err != nil {
		return
	}
	thumbnailPath := filepath.Join(h.cache.Dir(book), "thumbnail.jpg")
	titlePath := filepath.Join(h.cache.Dir(book), c.Filenames[0])

	thumbnail.Generate(thumbnailPath, titlePath)

	//FIXME
	http.ServeFile(w, r, thumbnailPath)
}

func (h *Handlers) File(w http.ResponseWriter, r *http.Request) {
	book, _, err := h.handleBook(w, r)
	if err != nil {
		return
	}

	//FIXME
	http.ServeFile(w, r, filepath.Join(h.cache.Dir(book), mux.Vars(r)["filename"]))
}

func (h *Handlers) handleBook(w http.ResponseWriter, r *http.Request) (*bo.Book, *ca.Cache, error) {
	bookID := mux.Vars(r)["books"]
	book, err := h.bookList.Get(bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, nil, err
	}
	if book == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return nil, nil, errors.New("Not Found")
	}

	c, err := h.cache.Get(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, nil, err
	}

	return book, c, nil
}
