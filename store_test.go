package main

import (
	"log"
	"math/rand"
	"testing"
)

func TestGetSet(t *testing.T) {
	s := NewStore(nil)
	o := NewObject(nil)
	o.Data["test"] = "hello world"
	o.Data["test2"] = 123
	if err := s.Set(o, nil); err != nil {
		t.Error(err)
		return
	}
	storedObj, err := s.Get(o.UID, nil)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(o.Data, storedObj.Data)
	if storedObj.UID != o.UID || storedObj.Data["test"] != o.Data["test"] || storedObj.Data["test2"] != o.Data["test2"] {
		t.Error("stored object does not match")
		return
	}
}

func TestIndexSet(t *testing.T) {
	s := NewStore(nil)
	o := NewObject(nil)
	o.Data["test"] = "hello world"
	o.Data["test_long"] = ""
	for i := 0; i < 256; i++ {
		o.Data["test_long"] = o.Data["test_long"].(string) + "a"
	}
	if err := s.Set(o, nil); err != nil {
		t.Error(err)
		return
	}
	index, err := s.Index()
	if err != nil {
		t.Error(err)
		return
	}
	if index[0].UID != o.UID || index[0].Data["test"] != o.Data["test"] {
		t.Error("indexed object does not match")
	}
	if len(index[0].Data["test_long"].(string)) > indexValueMaxSize {
		t.Error("unexpected long string indexed")
	}
}

func TestQuery(t *testing.T) {

	s := NewStore(nil)
	o := NewObject(nil)
	o.Data["test_int"] = 123
	o.Data["test_float"] = 123.4
	o.Data["test_bool"] = false
	o.Data["test_string"] = "hello world"
	if err := s.Set(o, nil); err != nil {
		t.Error(err)
		return
	}
	if err := s.Set(o, nil); err != nil {
		t.Error(err)
		return
	}
	res, err := s.Query("test_int = 123", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res) == 0 {
		t.Error("unexpected empty index")
	}
	if res[0].UID != o.UID {
		t.Error("unexpected item in index")
	}

	res, err = s.Query("test_int > 64 and test_int < 128", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res) == 0 {
		t.Error("unexpected empty index")
		return
	}
	if res[0].UID != o.UID {
		t.Error("unexpected item in index")
		return
	}

	res, err = s.Query("test_int > 123", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res) != 0 {
		t.Error("unexpected index")
		return
	}

	res, err = s.Query("test_string = 'hello world'", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res) == 0 {
		t.Error("unexpected empty index")
		return
	}
	if res[0].UID != o.UID {
		t.Error("unexpected item in index")
		return
	}

}

func TestQueryMulti(t *testing.T) {

	s := NewStore(nil)
	o1 := NewObject(nil)
	o1.Data["test_str"] = "hello"
	o1.Data["test_int"] = 1
	s.Set(o1, nil)

	o2 := NewObject(nil)
	o2.Data["test_str"] = "world"
	o2.Data["test_int"] = 99
	s.Set(o2, nil)

	o3 := NewObject(nil)
	o2.Data["test_str"] = "world"
	o2.Data["test_float"] = 153.4
	s.Set(o3, nil)

	res, err := s.Query("test_int >= 1", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res) != 2 {
		t.Error("unexpected query results")
		return
	}
}

func TestLargeIndex(t *testing.T) {
	s := NewStore(nil)
	// build very large index
	for i := 0; i < 4096; i++ {
		o := NewObject(nil)
		o.Data["test_int"] = rand.Int()
		o.Data["test_float"] = rand.Float64()
		o.Data["test_letter"] = string(byte(65 + (i % 24)))
		s.Set(o, nil)
	}
	index, _ := s.Index()
	if len(index) != 4096 {
		t.Error("unexpected index size")
	}
	res, _ := s.Query("test_int > 0", nil)
	if len(res) == 0 {
		t.Error("expected at least one result from query")
	}
	res, _ = s.Query("test_letter = 'A'", nil)
	if len(res) == 0 || len(res) == 4096 {
		t.Error("expected more than one result from query but less than 4096")
	}
}

func TestSyncIndex(t *testing.T) {

	s := NewStore(nil)
	o := NewObject(nil)
	o.Data["test_int"] = 123
	o.Data["test_float"] = 123.4
	o.Data["test_bool"] = false
	o.Data["test_string"] = "hello world"
	if err := s.Set(o, nil); err != nil {
		t.Error(err)
		return
	}

	// store object and sync index
	s.Set(o, nil)
	if err := s.Sync(); err != nil {
		t.Error(err)
		return
	}

	// update object without sync
	o.Data["test_string"] = "hello world two"
	s.Set(o, nil)

	// fetch remote index prior to sync to ensure
	// old value still remains
	remoteIndex := make([]*IndexObject, 0)
	s.getRaw(indexName, &remoteIndex)
	if len(remoteIndex) == 0 {
		t.Error("unexpected remote index length")
	}
	if remoteIndex[0].Data["test_string"] != "hello world" {
		t.Error("unexpected value in remote index")
	}

	// sync and ensure remote index is now updated
	s.Sync()
	remoteIndex = make([]*IndexObject, 0)
	s.getRaw(indexName, &remoteIndex)
	if remoteIndex[0].Data["test_string"] != "hello world two" {
		t.Error("unexpected value in remote index")
	}

}
