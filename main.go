package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"

	"github.com/yuanying/crapi/pkgs/books"
)

func main() {
	var (
		log    logr.Logger
		root   string
		dbPath string
	)
	flag.StringVar(&root, "root", root, "Root directory")
	flag.StringVar(&dbPath, "db-path", dbPath, "DB path")
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
	if err := bookList.Register(ctx, root); err != nil {
		panic(err)
	}

	books, err := bookList.All()
	if err != nil {
		panic(err)
	}
	for i := range books {
		b := books[i]
		// fmt.Println(fmt.Sprintf("Time: %v, Artists: %v, Category: %v, Title: %v", b.ModTime, b.Artists, b.Categories, b.Filename))
		fmt.Println(fmt.Sprintf("Categories: %v", b.Categories))
	}

	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
