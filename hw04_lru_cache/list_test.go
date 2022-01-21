package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestListFuncs(t *testing.T) {
	t.Run("PushFront", func(t *testing.T) {
		l := NewList()
		l.PushFront(1)
		require.Equal(t, 1, l.Front().Value)
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Front().Next)
		require.Nil(t, l.Back())

		l.PushFront(2)
		require.Equal(t, 2, l.Front().Value)
		require.Nil(t, l.Front().Prev)
		require.NotNil(t, l.Front().Next)
		require.NotNil(t, l.Back())
		require.Equal(t, l.Front().Next, l.Back())

		l = NewList()
		l.PushBack(3)
		l.PushFront(2)

		require.Equal(t, 2, l.Front().Value)
		require.Nil(t, l.Front().Prev)
		require.NotNil(t, l.Front().Next)
		require.NotNil(t, l.Back())
		require.Equal(t, l.Front().Next, l.Back())
		require.Equal(t, 3, l.Back().Value)

		l.PushFront(1)
		require.Equal(t, []int{1, 2, 3}, ListToIntArray(l))
	})
	t.Run("PushBack", func(t *testing.T) {
		l := NewList()
		l.PushBack(1)
		require.Equal(t, 1, l.Back().Value)
		require.Equal(t, 1, l.Len())
		require.Nil(t, l.Back().Prev)
		require.Nil(t, l.Back().Next)
		require.Nil(t, l.Front())

		l.PushBack(2)
		require.Equal(t, 2, l.Back().Value)
		require.Equal(t, 2, l.Len())
		require.NotNil(t, l.Back().Prev)
		require.Nil(t, l.Back().Next)
		require.NotNil(t, l.Front())

		l.PushBack(3)
		require.Equal(t, 3, l.Back().Value)
		require.Equal(t, 3, l.Len())
		require.NotNil(t, l.Back().Prev)
		require.Nil(t, l.Back().Next)
		require.NotNil(t, l.Front())

		l = NewList()
		l.PushFront(1)
		l.PushBack(2)
		require.Equal(t, 2, l.Back().Value)
		require.Equal(t, 2, l.Len())
		require.NotNil(t, l.Back().Prev)
		require.Nil(t, l.Back().Next)
		require.NotNil(t, l.Front())
		require.Equal(t, 1, l.Front().Value)

		l.PushBack(3)
		require.Equal(t, []int{1, 2, 3}, ListToIntArray(l))
	})

	t.Run("Remove", func(t *testing.T) {
		l := NewList()
		l.PushFront(2)
		l.PushBack(3)
		l.PushFront(1)
		require.Equal(t, []int{1, 2, 3}, ListToIntArray(l))
		require.Equal(t, 3, l.Len())

		l.Remove(l.Back())
		require.Equal(t, []int{1, 2}, ListToIntArray(l))
		require.Equal(t, 2, l.Len())

		l.Remove(l.Front())
		require.Equal(t, []int{2}, ListToIntArray(l))
		require.Equal(t, 1, l.Len())

		l.Remove(l.Front())
		require.Equal(t, []int{}, ListToIntArray(l))
		require.Equal(t, 0, l.Len())

		l = NewList()
		l.PushFront(2)
		l.PushBack(3)
		l.PushFront(1)
		require.Equal(t, []int{1, 2, 3}, ListToIntArray(l))

		l.Remove(l.Front())
		require.Equal(t, []int{2, 3}, ListToIntArray(l))

		l = NewList()
		l.PushBack(2)
		l.PushFront(1)
		l.PushBack(3)
		require.Equal(t, []int{1, 2, 3}, ListToIntArray(l))

		l.Remove(l.Front())
		require.Equal(t, []int{2, 3}, ListToIntArray(l))
	})

	t.Run("MoveToFront", func(t *testing.T) {
		l := NewList()
		l.PushFront(2)
		l.PushBack(3)
		l.PushFront(1)
		require.Equal(t, []int{1, 2, 3}, ListToIntArray(l))
		require.Equal(t, 3, l.Back().Value)

		l.MoveToFront(l.Front())
		require.Equal(t, []int{1, 2, 3}, ListToIntArray(l))

		l.MoveToFront(l.Back())
		require.Equal(t, []int{3, 1, 2}, ListToIntArray(l))
		require.Equal(t, 2, l.Back().Value)

		l.MoveToFront(l.Front().Next)
		require.Equal(t, []int{1, 3, 2}, ListToIntArray(l))

		l = NewList()
		l.PushFront(2)
		l.PushBack(3)
		l.PushFront(1)
		require.Equal(t, []int{1, 2, 3}, ListToIntArray(l))

		l.MoveToFront(l.Front())
		require.Equal(t, []int{1, 2, 3}, ListToIntArray(l))

		l.MoveToFront(l.Back().Prev)
		require.Equal(t, []int{2, 1, 3}, ListToIntArray(l))
	})

	t.Run("String", func(t *testing.T) {
		l := NewList()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)
		require.Equal(t, "1 2 3 ", l.(*list).String())
	})
}
