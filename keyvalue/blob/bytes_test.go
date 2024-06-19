package blob

import (
	"testing"

	"github.com/hack-pad/hackpadfs/internal/assert"
)

func TestBytes_Bytes(t *testing.T) {
	t.Parallel()

	b := NewBytes([]byte("hello"))
	assert.Equal(t, []byte("hello"), b.Bytes())
}

func TestBytes_Len(t *testing.T) {
	t.Parallel()

	b := NewBytes([]byte("hello"))
	assert.Equal(t, 5, b.Len())
}

func TestBytes_View(t *testing.T) {
	t.Parallel()

	b := NewBytes([]byte("hello"))
	v, err := b.View(1, 3)
	assert.NoError(t, err)
	assert.Equal(t, []byte("el"), v.(*Bytes).Bytes())

	_, err = b.View(1, 6)
	assert.Error(t, err)

	_, err = b.View(6, 1)
	assert.Error(t, err)
}

func TestBytes_Slice(t *testing.T) {
	t.Parallel()

	b := NewBytes([]byte("hello"))
	v, err := b.Slice(1, 3)
	assert.NoError(t, err)
	assert.Equal(t, "he", string(v.Bytes()))

	_, err = b.Slice(1, 6)
	assert.Error(t, err)

	_, err = b.Slice(6, 1)
	assert.Error(t, err)
}

func TestBytes_Set(t *testing.T) {
	t.Parallel()

	b := NewBytes([]byte("hello"))
	n, err := b.Set(NewBytes([]byte("world")), 0)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "world", string(b.Bytes()))

	n, err = b.Set(NewBytes([]byte("world")), 1)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, "wworl", string(b.Bytes()))

	n, err = b.Set(NewBytes([]byte("world")), 5)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, "wworl", string(b.Bytes()))
}

func TestBytes_Grow(t *testing.T) {
	t.Parallel()

	b := NewBytes([]byte("hello"))
	err := b.Grow(5)
	assert.NoError(t, err)
	assert.Equal(t, "hello\x00\x00\x00\x00\x00", string(b.Bytes()))

	err = b.Grow(0)
	assert.NoError(t, err)
	assert.Equal(t, "hello\x00\x00\x00\x00\x00", string(b.Bytes()))
}

func TestBytes_Truncate(t *testing.T) {
	t.Parallel()

	b := NewBytes([]byte("hello"))
	err := b.Truncate(3)
	assert.NoError(t, err)
	assert.Equal(t, "hel", string(b.Bytes()))

	err = b.Truncate(5)
	assert.NoError(t, err)
	assert.Equal(t, "hel", string(b.Bytes()))

	err = b.Truncate(0)
	assert.NoError(t, err)
	assert.Equal(t, "", string(b.Bytes()))
}
