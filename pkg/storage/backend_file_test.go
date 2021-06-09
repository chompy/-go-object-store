package storage

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/contextualcode/storage-backend/pkg/types"
)

func TestStore(t *testing.T) {
	b := FileBackend{
		Path: os.TempDir(),
		GZIP: false,
	}
	// create
	o := types.NewObject()
	o.Type = "test"
	o.Data["name"] = "Test Person"
	o.Data["age"] = float64(32)
	if err := b.Put(o); err != nil {
		t.Error(err)
	}
	// fetch
	so, err := b.Get(o.UID)
	if err != nil {
		t.Error(err)
	}
	if so.UID != o.UID {
		t.Error("stored object uid does not match")
	}
	log.Println("DSFsd")
	if so.Data["name"] != o.Data["name"] || so.Data["age"] != o.Data["age"] {
		t.Error("stored data does not match")
	}
	// delete
	if err := b.Delete(o); err != nil {
		t.Error(err)
	}
	f, err := os.Open(filepath.Join(os.TempDir(), o.UID))
	if err == nil {
		f.Close()
		t.Error("expected object file to be deleted")
	}
}
