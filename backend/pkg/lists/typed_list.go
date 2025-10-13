package lists

import (
	"container/list"
	"errors"
)

type TypedList[T comparable] struct {
	list *list.List
}

func New[T comparable]() *TypedList[T] {
	return &TypedList[T]{
		list: list.New(),
	}
}

func (tl *TypedList[T]) Len() int {
	return tl.list.Len()
}

func (tl *TypedList[T]) PushBack(value T) {
	tl.list.PushBack(value)
}

func (tl *TypedList[T]) PushFront(value T) {
	tl.list.PushFront(value)
}

func (tl *TypedList[T]) Remove(value T) (T, error) {
	if value == tl.list.Back().Value.(T) {
		tl.list.Remove(tl.list.Back())
		return value, nil
	}

	var assertedValue T
	for e := tl.list.Front(); e != nil; e = e.Next() {
		var ok bool
		assertedValue, ok = e.Value.(T)
		if !ok {
			panic("lists: Remove(): type assertion failed")
		}
		if assertedValue == value {
			tl.list.Remove(e)
			return assertedValue, nil
		}
	}
	return assertedValue, errors.New("lists: Remove(): value not found")
}

func (tl *TypedList[T]) Get(i int) (T, error) {
	var assertedValue T
	if i < 0 || i >= tl.list.Len() {
		return assertedValue, errors.New("lists: Get(): index out of range")
	}

	if i == tl.list.Len()-1 {
		return tl.list.Back().Value.(T), nil
	}

	e := tl.list.Front()
	for j := 0; j < i; j++ {
		e = e.Next()
	}
	assertedValue, ok := e.Value.(T)
	if !ok {
		panic("lists: Get(): type assertion failed")
	}
	return assertedValue, nil
}

func (tl *TypedList[T]) GetList() *list.List {
	return tl.list
}

func (tl *TypedList[T]) Contains(value T) bool {
	for e := tl.list.Front(); e != nil; e = e.Next() {
		if e.Value == value {
			return true
		}
	}
	return false
}
