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

	"github.com/go-logr/logr"
	"github.com/yuanying/crapi/pkgs/books"
)

type Cache struct {
	Filenames []string `json:"filenames"`
}

type Manager struct {
	Root string

	log logr.Logger
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
		return c.Generate(dir, book)
	}

	return c.Load(c.IndexPath(book), book)
}

func (c *Manager) Dir(book *books.Book) string {
	return filepath.Join(c.Root, string([]rune(book.ID)[:2]), book.ID)
}

func (c *Manager) IndexPath(book *books.Book) string {
	return filepath.Join(c.Dir(book), "index.json")
}

func (c *Manager) Load(index string, book *books.Book) (cache *Cache, err error) {
	cache = &Cache{}

	b, err := ioutil.ReadFile(index)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, cache)

	return
}

func (c *Manager) Generate(path string, book *books.Book) (cache *Cache, err error) {
	if err = os.MkdirAll(path, 0755); err != nil {
		c.log.Error(err, "Failed to create cache dir", "path", path)
		return
	}

	cache = &Cache{}

	reader, err := zip.OpenReader(book.Path)
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

	defer func() {
		if err != nil {
			c.log.Error(err, "Try to remove all cache files")
		}
	}()

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
