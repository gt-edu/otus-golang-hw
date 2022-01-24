package hw04lrucache

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic capacity", func(t *testing.T) {
		c := NewCache(1)

		c.Set("test1", 1)
		c.Set("test2", 2)

		v, ok := c.Get("test1")
		require.False(t, ok)
		require.Nil(t, v)
	})

	t.Run("purge logic age", func(t *testing.T) {
		c := NewCache(2)

		c.Set("test1", 1)
		c.Set("test2", 2)

		v, ok := c.Get("test2")
		require.True(t, ok)
		require.Equal(t, 2, v)

		c.Set("test3", 3)

		v, ok = c.Get("test3")
		require.True(t, ok)
		require.Equal(t, 3, v)

		v, ok = c.Get("test2")
		require.True(t, ok)
		require.Equal(t, 2, v)

		v, ok = c.Get("test1")
		require.False(t, ok)
		require.Nil(t, v)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func TestCustomCache(t *testing.T) {
	t.Run("SetAndGet", func(t *testing.T) {
		c := NewCache(10)

		v, ok := c.Get("test1")
		require.False(t, ok)
		require.Equal(t, nil, v)

		exist := c.Set("test1", 1)
		require.False(t, exist)

		v, ok = c.Get("test1")
		require.True(t, ok)
		require.Equal(t, 1, v)

		exist = c.Set("test1", 2)
		require.True(t, exist)

		v, ok = c.Get("test1")
		require.True(t, ok)
		require.Equal(t, 2, v)

		v, ok = c.Get("test2")
		require.False(t, ok)
		require.Equal(t, nil, v)
	})

	t.Run("Clear", func(t *testing.T) {
		c := NewCache(10)

		c.Set("test1", 1)
		c.Set("test2", 2)

		v, ok := c.Get("test1")
		require.True(t, ok)
		require.Equal(t, 1, v)

		v, ok = c.Get("test2")
		require.True(t, ok)
		require.Equal(t, 2, v)

		c.Clear()
		v, ok = c.Get("test1")
		require.False(t, ok)
		require.Nil(t, v)

		v, ok = c.Get("test2")
		require.False(t, ok)
		require.Nil(t, v)
	})
}
