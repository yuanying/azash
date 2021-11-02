package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"

	"github.com/yuanying/crapi/pkgs/books"
	books_handler "github.com/yuanying/crapi/pkgs/handlers/books"
)

func main() {
	var (
		log    logr.Logger
		root   string
		dbPath string
		wait   time.Duration
	)
	flag.StringVar(&root, "root", root, "Root directory")
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
	if err := bookList.RegisterDir(ctx, root); err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	booksHandler := books_handler.NewHandler(log, bookList)

	r.HandleFunc("/apis/books", booksHandler.All).Methods("GET")

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
