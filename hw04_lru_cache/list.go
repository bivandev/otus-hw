package hw04lrucache

import "fmt"

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
	len   int
	front *ListItem
	back  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	var item = new(ListItem)
	item.Value = v
	item.Next = l.front

	if l.len == 0 {
		l.back = item
	} else {
		l.front.Prev = item
	}

	l.front = item
	l.len = l.len + 1

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	var item = new(ListItem)
	item.Value = v

	if l.len == 0 {
		l.front = item
	} else {
		l.back.Next = item
	}

	item.Prev = l.back
	l.back = item
	l.len = l.len + 1

	return item
}

func (l *list) Remove(i *ListItem) {
	if i == nil || (i.Prev == nil && i.Next == nil && i != l.front) {
		fmt.Println("Error while deleting cache item")
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	i.Next = nil
	i.Prev = nil
	l.len = l.len - 1
}

func (l *list) MoveToFront(i *ListItem) {
	exNext := i.Next
	exPrev := i.Prev

	if exPrev == nil {
		return
	}

	if exNext == nil {
		exPrev.Next = nil
		l.back = exPrev
	} else {
		exPrev.Next = i.Next
		exNext.Prev = i.Prev
	}

	exFront := l.Front()
	exFront.Prev = i

	i.Next = exFront

	i.Prev = nil

	l.front = i
}

func NewList() List {
	return new(list)
}
