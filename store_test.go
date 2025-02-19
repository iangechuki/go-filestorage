package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpic"
	pathname := CASPathTransformFunc(key)
	expectedPathName := "1b150/aae86/eedae/268f6/589f4/0fb48/b2a0d/47ff4"
	expectedOriginal := "1b150aae86eedae268f6589f40fb48b2a0d47ff4"
	if pathname.Pathname != expectedPathName {
		t.Errorf("expected %s, got %s", expectedPathName, pathname)
	}
	if pathname.Filename != expectedOriginal {
		t.Errorf("expected %s, got %s", expectedOriginal, pathname)
	}
}
func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "myspecialkey"
	data := []byte("new data")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if !s.Has(key) {
		t.Error("File does not exist")
	}
	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}
	b, _ := io.ReadAll(r)
	if string(b) != string(data) {
		t.Errorf("expected %s, got %s", string(data), string(b))
	}

}
func TestStoreDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "myspecialkey"
	data := []byte("new data")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	err := s.Delete(key)
	if err != nil {
		t.Error(err)
	}
	if s.Has(key) {
		t.Log("File exists!")
	} else {
		t.Log("File does not exist.")
	}
	// if err := s.Delete(key); err != nil {
	// 	t.Error(err)
	// }
}
