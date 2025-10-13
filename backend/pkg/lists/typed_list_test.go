package lists

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushBack(t *testing.T) {
	list := New[int]()
	list.PushBack(1)
	assert.Equal(t, 1, list.Len())
}

func TestPushFront(t *testing.T) {
	list := New[int]()
	list.PushFront(1)
	assert.Equal(t, 1, list.Len())
}

func TestGet(t *testing.T) {
	var value int
	var err error

	list := New[int]()
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	value, err = list.Get(0)
	assert.Nil(t, err)
	assert.Equal(t, 1, value)

	value, err = list.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, 2, value)

	value, err = list.Get(2)
	assert.Nil(t, err)
	assert.Equal(t, 3, value)

	_, err = list.Get(3)
	assert.NotNil(t, err)

	_, err = list.Get(-1)
	assert.NotNil(t, err)
}

func TestRemove(t *testing.T) {
	var value int
	var err error

	list := New[int]()
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	_, err = list.Remove(2)
	assert.Nil(t, err)

	value, err = list.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, 3, value)
}

func TestRemoveBack(t *testing.T) {
	var value int
	var err error

	list := New[int]()
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	_, err = list.Remove(3)
	assert.Nil(t, err)

	value, err = list.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, 2, value)

	_, err = list.Get(2)
	assert.NotNil(t, err)
}

func TestContains(t *testing.T) {
	list := New[int]()
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	assert.True(t, list.Contains(1))
	assert.False(t, list.Contains(4))
}
