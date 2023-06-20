package local

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gogrlx/grlx/ingredients/file"
	"github.com/gogrlx/grlx/ingredients/file/hashers"
	"github.com/gogrlx/grlx/types"
)

type LocalFile struct {
	ID          string
	Source      string
	Destination string
	Hash        string
	Props       map[string]interface{}
}

func (lf LocalFile) Download(ctx context.Context) error {
	ok, err := lf.Verify(ctx)
	// if verification failed because the file doesn't exist,
	// that's ok. Otherwise, return the error.
	if !errors.Is(err, types.ErrFileNotFound) {
		return err
	}
	// if the file exists and the hash matches, we're done.
	if ok {
		return nil
	}
	// otherwise, "download" the file.
	f, err := os.Open(lf.Source)
	if err != nil {
		return err
	}
	defer f.Close()
	dest, err := os.Create(lf.Destination)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, f)
	return err
}

func (lf LocalFile) Properties() (map[string]interface{}, error) {
	return lf.Props, nil
}

func (lf LocalFile) Parse(id, source, destination, hash string, properties map[string]interface{}) (types.FileProvider, error) {
	// TODO make this properties nil check in other places
	if properties == nil {
		properties = make(map[string]interface{})
	}
	return LocalFile{ID: id, Source: source, Destination: destination, Hash: hash, Props: properties}, nil
}

func (lf LocalFile) Protocols() []string {
	return []string{"file"}
}

func (lf LocalFile) Verify(ctx context.Context) (bool, error) {
	_, err := os.Stat(lf.Destination)
	if err != nil {
		if os.IsNotExist(err) {
			return false, errors.Join(err, types.ErrFileNotFound)
		}
	}
	f, err := os.Open(lf.Destination)
	if err != nil {
		return false, err
	}
	defer f.Close()
	hashType := ""
	if lf.Props["hashType"] == nil {
		hashType = hashers.GuessHashType(lf.Hash)
	} else if ht, ok := lf.Props["hashType"].(string); !ok {
		hashType = hashers.GuessHashType(lf.Hash)
	} else {
		hashType = ht
	}
	hf, err := hashers.GetHashFunc(hashType)
	if err != nil {
		return false, err
	}
	hash, matches, err := hf(f, lf.Hash)
	if err != nil {
		return false, errors.Join(err, fmt.Errorf("recipe step %s: hash for %s failed: expected %s but found %s", lf.ID, lf.Destination, lf.Hash, hash))
	}
	return matches, err
}

func init() {
	file.RegisterProvider(LocalFile{})
}
