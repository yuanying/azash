package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"

	"github.com/yuanying/azash/pkgs/books"
	"github.com/yuanying/azash/pkgs/caches"
	books_handler "github.com/yuanying/azash/pkgs/handlers/books"
	caches_handler "github.com/yuanying/azash/pkgs/handlers/caches"
)

func main() {
	var (
		log      logr.Logger
		root     string
		cacheDir string
		dbPath   string
		wait     time.Duration
	)
	flag.StringVar(&root, "root", root, "Root directory")
	flag.StringVar(&cacheDir, "cache", "/tmp/azash", "Cache directory")
	flag.StringVar(&dbPath, "db-path", dbPath, "DB path")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Error(err, "Unable to open database", "path", dbPath)
		os.Exit(1)
	}
	defer db.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	bookList := books.NewBookList(log.WithName("books"), db)
	cache := caches.NewManager(log.WithName("cache"), cacheDir)

	go func() {
		if err := RegisterDir(ctx, &log, root, bookList, cache); err != nil {
			panic(err)
		}
	}()

	r := mux.NewRouter()
	booksHandler := books_handler.NewHandler(log, bookList)
	cachesHandler := caches_handler.NewHandler(log, bookList, cache)

	r.HandleFunc("/apis/books", booksHandler.All).Methods("GET")
	r.HandleFunc("/books/{books}", cachesHandler.Index).Methods("GET")
	r.HandleFunc("/books/{books}/thumbnail", cachesHandler.Thumbnail).Methods("GET")
	r.HandleFunc("/books/{books}/{filename}", cachesHandler.File).Methods("GET")

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()

		srv.Shutdown(ctx)
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Error(err, "Error")
	}
	os.Exit(0)
}

func RegisterDir(ctx context.Context, log *logr.Logger, root string, bookList *books.BookList, cache *caches.Manager) error {
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
				book := books.Book{
					ID:       hash,
					Filename: info.Name(),
					ModTime:  info.ModTime(),
				}
				book.ParseFilename()
				book.ParsePath(path)
				err = cache.Generate(path, &book)
				if err != nil {
					return err
				}
				return bookList.Register(book)
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err, "Something wrong")
		return err
	}

	return nil
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
	s := sha1.Sum(buf)

	return hex.EncodeToString(s[:]), nil
}
