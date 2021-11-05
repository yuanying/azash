package caches

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-logr/logr"
	"github.com/yuanying/azash/pkgs/books"
)

type Cache struct {
	Filenames []string `json:"filenames"`
}

type Manager struct {
	Root string

	log logr.Logger
	mux sync.Mutex
}

func NewManager(log logr.Logger, root string) *Manager {
	return &Manager{
		log:  log,
		Root: root,
	}
}

func (c *Manager) Get(book *books.Book) (*Cache, error) {
	dir := c.Dir(book)
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return nil, nil
	}

	return c.Load(book)
}

func (c *Manager) Dir(book *books.Book) string {
	return filepath.Join(c.Root, string([]rune(book.ID)[:2]), book.ID)
}

func (c *Manager) IndexPath(book *books.Book) string {
	return filepath.Join(c.Dir(book), "index.json")
}

func (c *Manager) Load(book *books.Book) (cache *Cache, err error) {
	cache = &Cache{}
	index := c.IndexPath(book)
	c.mux.Lock()
	defer c.mux.Unlock()

	b, err := ioutil.ReadFile(index)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, cache)

	return
}

func (c *Manager) Generate(bookPath string, book *books.Book) (err error) {
	path := c.Dir(book)
	c.mux.Lock()
	defer c.mux.Unlock()

	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		if err != nil {
			return
		}
		return nil
	}

	if err = os.MkdirAll(path, 0755); err != nil {
		c.log.Error(err, "Failed to create cache dir", "path", path)
		return
	}
	defer func() {
		if err != nil {
			c.log.Error(err, "TODO: Try to remove all cache files")
		}
	}()

	cache := &Cache{}

	reader, err := zip.OpenReader(bookPath)
	if err != nil {
		return
	}
	defer reader.Close()

	for i, zf := range reader.File {
		zfinfo := zf.FileInfo()
		if zfinfo.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(zfinfo.Name()))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			continue
		}
		targetName := fmt.Sprintf("%v%v", i, ext)
		target := filepath.Join(path, targetName)
		err := c.createFileFromZipFile(target, zf)
		if err != nil {
			c.log.Error(err, "Failed to decompress zipped file", "target", target, "zipped", zf.Name)
			continue
		}
		cache.Filenames = append(cache.Filenames, targetName)
	}

	if err = c.createIndexFile(c.IndexPath(book), cache); err != nil {
		return
	}

	return
}

func (c *Manager) createIndexFile(target string, cache *Cache) error {
	var err error

	destFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destFile.Close()

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	_, err = destFile.Write(data)

	return err
}

func (c *Manager) createFileFromZipFile(target string, zf *zip.File) error {
	var err error

	reader, err := zf.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	destFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, reader)

	return err
}
