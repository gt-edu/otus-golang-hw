package hw04lrucache

import (
	"fmt"
	"strings"
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newFront := &ListItem{
		Value: v,
	}

	l.pushItemToFront(newFront)

	return newFront
}

func (l *list) pushItemToFront(newFront *ListItem) {
	if l.front != nil {
		l.front.Prev = newFront
		newFront.Next = l.front
	} else if l.back != nil {
		l.back.Prev = newFront
		newFront.Next = l.back
	}

	if l.len == 1 && l.back == nil {
		l.back = l.front
	}

	newFront.Prev = nil

	l.front = newFront

	l.len++
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBack := &ListItem{
		Value: v,
	}

	if l.back != nil {
		l.back.Next = newBack
		newBack.Prev = l.back
	} else if l.front != nil {
		l.front.Next = newBack
		newBack.Prev = l.front
	}

	if l.len == 1 && l.front == nil {
		l.front = l.back
	}

	l.back = newBack

	l.len++

	return newBack
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		// it is front
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		// it is back
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.pushItemToFront(i)
}

func (l *list) String() string {
	sb := strings.Builder{}
	for curr := l.front; curr != nil; curr = curr.Next {
		sb.WriteString(fmt.Sprintf("%v ", curr.Value))
	}

	return sb.String()
}

func ListToIntArray(l List) []int {
	arr := make([]int, 0, l.Len())
	for curr := l.Front(); curr != nil; curr = curr.Next {
		arr = append(arr, curr.Value.(int))
	}

	return arr
}

func NewList() List {
	return new(list)
}
