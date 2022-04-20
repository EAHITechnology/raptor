package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultQueueImpl(t *testing.T) {
	queue := NewDefaultqueue()
	err := queue.Put("test")
	assert.Nil(t, err)
	assert.Equal(t, int(queue.Len()), 1)

	item, err := queue.Pop()
	assert.Nil(t, err)
	assert.Equal(t, item.(string), "test")
	assert.Equal(t, int(queue.Len()), 0)

	item, err = queue.Pop()
	assert.Equal(t, err, ErrPopNil)
	assert.Equal(t, item, nil)
	assert.Equal(t, int(queue.Len()), 0)

	err = queue.Put("test again")
	assert.Nil(t, err)
	assert.Equal(t, int(queue.Len()), 1)
}
