package skip_list

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultSkipListImpl(t *testing.T) {
	list, err := NewDefaultSkipListImpl()
	assert.Nil(t, err)
	assert.NotNil(t, list)
}

func TestDefaultSkipListImpl_PutGet(t *testing.T) {
	list, err := NewDefaultSkipListImpl()
	assert.Nil(t, err)

	err = list.Put([]byte("test_key"), []byte("test_value"))
	assert.Nil(t, err)

	value, ok, err := list.Get([]byte("test_key"))
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, []byte("test_value"))
}

func TestDefaultSkipListImpl_Del(t *testing.T) {
	list, err := NewDefaultSkipListImpl()
	assert.Nil(t, err)

	err = list.Put([]byte("test_key"), []byte("test_value"))
	assert.Nil(t, err)

	value, ok, err := list.Get([]byte("test_key"))
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, []byte("test_value"))

	err = list.Del([]byte("test_key"))
	assert.Nil(t, err)

	value, ok, err = list.Get([]byte("test_key"))
	assert.Nil(t, err)
	assert.Equal(t, false, ok)
	assert.Equal(t, []uint8([]byte(nil)), value)

	assert.Equal(t, int64(0), list.Len())
}

func TestDefaultSkipListImpl_Iter(t *testing.T) {
	list, err := NewDefaultSkipListImpl()
	assert.Nil(t, err)

	err = list.Put([]byte("test_key0"), []byte("test_value0"))
	assert.Nil(t, err)

	err = list.Put([]byte("test_key1"), []byte("test_value1"))
	assert.Nil(t, err)

	err = list.Put([]byte("test_key2"), []byte("test_value2"))
	assert.Nil(t, err)

	err = list.Put([]byte("test_key3"), []byte("test_value3"))
	assert.Nil(t, err)

	iter := list.GetIter()
	err = iter.Seek([]byte("test_key0"))
	assert.Nil(t, err)

	value, err := iter.Get()
	assert.Nil(t, err)
	assert.Equal(t, []byte("test_value0"), value)

	idx := 1
	for iter.HasNext() {
		value, err = iter.Next()
		assert.Nil(t, err)
		assert.Equal(t, []byte("test_value"+strconv.Itoa(idx)), value)
		idx++
	}

	iter.Close()
}
