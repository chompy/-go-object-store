package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"gitlab.com/contextualcode/storage-backend/pkg/types"

	"github.com/pkg/errors"
)

const indexFilename = "__index"

// FileBackend handles file system backend.
type FileBackend struct {
	Path  string
	GZIP  bool
	index []*types.Object
	sync  sync.Mutex
}

func (b *FileBackend) pathTo(uid string) string {
	return filepath.Join(b.Path, uid)
}

func (b *FileBackend) initIndex() error {
	b.index = make([]*types.Object, 0)
	data, err := ioutil.ReadFile(filepath.Join(b.Path, indexFilename))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.WithStack(err)
	}
	if b.GZIP {
		data, err = uncompressBytes(data)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	if err := yaml.Unmarshal(data, &b.index); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *FileBackend) updateIndex(o *types.Object) error {
	if b.index == nil {
		if err := b.initIndex(); err != nil {
			return errors.WithStack(err)
		}
	}
	for i := range b.index {
		if b.index[i].UID == o.UID {
			b.index[i] = o.Index()
			return errors.WithStack(b.saveIndex())
		}
	}
	b.index = append(b.index, o)
	return errors.WithStack(b.saveIndex())
}

func (b *FileBackend) deleteIndex(o *types.Object) error {
	if b.index == nil {
		return nil
	}
	for i := range b.index {
		if b.index[i].UID == o.UID {
			b.index = append(b.index[:i], b.index[i+1:]...)
			return errors.WithStack(b.saveIndex())
		}
	}
	return nil
}

func (b *FileBackend) saveIndex() error {
	data, err := json.Marshal(b.index)
	if err != nil {
		return errors.WithStack(err)
	}
	if b.GZIP {
		data, err = compressBytes(data)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	b.sync.Lock()
	defer b.sync.Unlock()
	if err := ioutil.WriteFile(
		filepath.Join(b.Path, indexFilename),
		data,
		0660,
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Put uploads given object.
func (b *FileBackend) Put(o *types.Object) error {
	o.Modified = time.Now()
	if err := b.updateIndex(o); err != nil {
		return errors.WithStack(err)
	}
	data, err := o.Serialize()
	if err != nil {
		return errors.WithStack(err)
	}
	if b.GZIP {
		data, err = compressBytes(data)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	b.sync.Lock()
	defer b.sync.Unlock()
	if err := ioutil.WriteFile(
		b.pathTo(o.UID), data, 0660,
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete deletes given object.
func (b *FileBackend) Delete(o *types.Object) error {
	if err := os.Remove(b.pathTo(o.UID)); err != nil {
		if os.IsNotExist(err) {
			return errors.WithStack(ErrNotExists)
		}
		return errors.WithStack(err)
	}
	return errors.WithStack(b.deleteIndex(o))
}

// Get downloads object.
func (b *FileBackend) Get(uid string) (*types.Object, error) {
	b.sync.Lock()
	defer b.sync.Unlock()
	data, err := ioutil.ReadFile(b.pathTo(uid))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.WithStack(ErrNotExists)
		}
		return nil, errors.WithStack(err)
	}
	if b.GZIP {
		data, err = uncompressBytes(data)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	o := &types.Object{}
	if err := o.Unserialize(data); err != nil {
		return nil, errors.WithStack(err)
	}
	return o, nil
}

// Query fetches objects based on provided query.
func (b *FileBackend) Query(q types.Query) ([]*types.Object, error) {
	return nil, nil
}
