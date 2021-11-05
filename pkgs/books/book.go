package books

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"go.etcd.io/bbolt"
)

var (
	booksBucket = []byte("BOOKS")
)

type BookList struct {
	log   logr.Logger
	db    *bbolt.DB
	cache []Book
	mux   sync.Mutex
}

type Book struct {
	ID         string
	Filename   string
	Artists    []string
	Categories []string
	ModTime    time.Time
}

func NewBookList(log logr.Logger, db *bbolt.DB) *BookList {
	return &BookList{
		log: log,
		db:  db,
	}
}

func (b *BookList) Register(book Book) error {
	b.db.Update(func(tx *bbolt.Tx) error {
		if bucket, err := tx.CreateBucketIfNotExists(booksBucket); err != nil {
			return err
		} else {
			bid := []byte(book.ID)
			if bucket.Get(bid) == nil {
				if raw, err := json.Marshal(book); err != nil {
					return err
				} else {
					if err := bucket.Put(bid, raw); err != nil {
						return err
					}
					b.mux.Lock()
					b.cache = nil
					b.mux.Unlock()
				}
			} else {
				b.log.Info("Book is already registerd", "id", book.ID, "filename", book.Filename)
			}
		}
		return nil
	})

	return nil
}

func (b *BookList) Get(id string) (*Book, error) {
	book := &Book{}
	b.db.View(func(tx *bbolt.Tx) error {
		raw := tx.Bucket(booksBucket).Get([]byte(id))
		if err := json.Unmarshal(raw, book); err != nil {
			return err
		}
		return nil
	})
	return book, nil
}

func (b *BookList) All() ([]Book, error) {
	b.mux.Lock()
	defer b.mux.Unlock()
	if b.cache != nil {
		return b.cache, nil
	}
	var books []Book

	if err := b.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(booksBucket).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			book := Book{}
			if err := json.Unmarshal(v, &book); err != nil {
				b.log.Info("Failed to marshal book", "id", k)
			}
			books = append(books, book)
		}

		return nil
	}); err != nil {
		b.log.Error(err, "Failed to read db", "type", "Book")
		return nil, err
	}

	sort.Slice(books, func(i, j int) bool { return books[i].ModTime.After(books[j].ModTime) })
	b.cache = books

	return b.cache, nil
}

var titleRegExp = regexp.MustCompile(`^(\((?P<category>\S+)\))?\s*(\[(?P<rawArtists>[^\]]+)\])?\s*(.+)`)

func (book *Book) ParseFilename() error {
	g := titleRegExp.FindStringSubmatch(book.Filename)
	if g != nil {
		catIndex := titleRegExp.SubexpIndex("category")
		category := g[catIndex]
		if category != "" {
			book.Categories = []string{category}
		}
		artIndex := titleRegExp.SubexpIndex("rawArtists")
		rawArtists := g[artIndex]
		if rawArtists != "" {
			book.Artists = abstractArtists(rawArtists)
		}

	}
	return nil
}

var allNumber = regexp.MustCompile(`^\d+$`)

func (book *Book) ParsePath(path string) error {
	parentDir := filepath.Base(filepath.Dir(path))

	if !strings.Contains(parentDir, ".") && !strings.Contains(parentDir, "/") && parentDir != "temp" && !allNumber.MatchString(parentDir) {
		for _, artist := range book.Artists {
			if artist == parentDir {
				return nil
			}
		}
		book.Artists = append(book.Artists, parentDir)
	}
	return nil
}

var artistsRegExpWithGroup = regexp.MustCompile(`(?P<group>[^\(]+)\s*(\((?P<artist>.+)\))?`)

func abstractArtists(str string) []string {
	var artists []string
	g := artistsRegExpWithGroup.FindStringSubmatch(str)
	if g != nil {
		grpIndex := artistsRegExpWithGroup.SubexpIndex("group")
		group := g[grpIndex]
		if group != "" {
			artists = append(artists, strings.TrimSpace(group))
		}
		artIndex := artistsRegExpWithGroup.SubexpIndex("artist")
		artist := g[artIndex]
		if artist != "" {
			artists = append(artists, strings.TrimSpace(artist))
		}
	}
	return artists
}
