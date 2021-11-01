package books

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"go.etcd.io/bbolt"
)

var (
	booksBucket = []byte("BOOKS")
)

type BookList struct {
	log logr.Logger
	db  *bbolt.DB
}

type Book struct {
	ID         string
	Filename   string
	Path       string
	Artists    []string
	Categories []string
	ModTime    time.Time
}

func NewBookList(log logr.Logger, db *bbolt.DB) BookList {
	return BookList{
		log: log,
		db:  db,
	}
}

func (b *BookList) Register(ctx context.Context, root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return errors.New("Canceled")
		default:
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") && (strings.HasSuffix(info.Name(), ".zip") || strings.HasPrefix(info.Name(), ".cbr")) {
				hash, err := fileHash(path)
				if err != nil {
					return err
				}
				book := Book{
					ID:       hash,
					Filename: info.Name(),
					Path:     path,
					ModTime:  info.ModTime(),
				}
				book.parseFilename()
				book.parsePath()
				return b.register(book)
			}
		}

		return nil
	})

	if err != nil {
		b.log.Error(err, "Something wrong")
		return err
	}

	return nil
}

func (b *BookList) register(book Book) error {
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
				}
			} else {
				b.log.Info("Book is already registerd", "id", book.ID, "filename", book.Filename)
			}
		}
		return nil
	})

	return nil
}

func (b *BookList) All() ([]Book, error) {
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

	return books, nil
}

var titleRegExp = regexp.MustCompile(`^(\((?P<category>\S+)\))?\s*(\[(?P<rawArtists>[^\]]+)\])?\s*(.+)`)

func (book *Book) parseFilename() error {
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

func (book *Book) parsePath() error {
	parentDir := filepath.Base(filepath.Dir(book.Path))

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

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	_, err = f.Read(buf)
	if err != nil {
		return "", err
	}
	s := sha256.Sum256(buf)

	return hex.EncodeToString(s[:]), nil
}
