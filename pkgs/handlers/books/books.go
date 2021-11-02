package books

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	bo "github.com/yuanying/crapi/pkgs/books"
)

type Book struct {
	ID         string
	Filename   string
	Categories []string
	Artists    []string
	ModTime    time.Time
}

type Handlers struct {
	log      logr.Logger
	bookList *bo.BookList
}

func NewHandler(log logr.Logger, bookList *bo.BookList) Handlers {
	return Handlers{
		log:      log,
		bookList: bookList,
	}
}

func (h *Handlers) All(w http.ResponseWriter, r *http.Request) {
	rawBooks, _ := h.bookList.All()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	books := make([]Book, len(rawBooks))

	for i := range rawBooks {
		raw := rawBooks[i]
		book := Book{
			ID:         raw.ID,
			Filename:   raw.Filename,
			Categories: raw.Categories,
			Artists:    raw.Artists,
			ModTime:    raw.ModTime,
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}
