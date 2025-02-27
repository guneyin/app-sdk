package store

import (
	"errors"
	"sync"

	"github.com/guneyin/app-sdk/utils"
)

var ErrRecordNotFound = errors.New("record not found")

type Store struct {
	store sync.Map
}

type Value[T any] struct {
	val T
}

func (v *Value[T]) Parse(dest T) error {
	_, err := utils.Convert(v.val, dest)
	if err != nil {
		return err
	}
	return nil
}

func New() *Store {
	return &Store{store: sync.Map{}}
}

func (s *Store) Set(key string, value any) {
	s.store.Store(key, value)
}

func (s *Store) Get(key string) Value[any] {
	val, _ := s.store.Load(key)
	return Value[any]{val: val}
}

func (s *Store) GetOnce(key string) (*Value[any], error) {
	val, ok := s.store.LoadAndDelete(key)
	if !ok {
		return nil, ErrRecordNotFound
	}
	return &Value[any]{val: val}, nil
}

func (s *Store) Delete(key string) {
	s.store.Delete(key)
}
