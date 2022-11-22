package storage

import (
	"crypto/sha1"
	"example.com/m/v2/lib/e"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	isExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	const textError = "can't calculate hash"

	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap(textError, err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap(textError, err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
