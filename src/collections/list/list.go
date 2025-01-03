package list

import (
	"fmt"
	"strings"
)

type (
	List[T comparable] struct {
		length int
		head *listItem[T]
		tail *listItem[T]
	}

	listItem[T comparable] struct {
		value T
		prev *listItem[T]
		next *listItem[T]
		owner *List[T]
	}
)

func New[T comparable](values ...T) *List[T] {
	l := Empty[T]()
	for _, val := range values {
		l.Add(val)
	}
	return l
}

func NewFor[T comparable](_ T) *List[T] {
	return Empty[T]()
}

func Empty[T comparable]() *List[T] {
	return &List[T]{0, nil, nil}
}

func All[T comparable](all func (T)) func (T) bool {
	return func(v T) bool {
		all(v)
		return true
	}
}
func AllIndex[T comparable](all func(int, T)) func (int, T) bool {
	return func(idx int, v T) bool {
		all(idx, v)
		return true
	}
}

func First[T comparable](first func(T)) func (T) bool {
	return func(v T) bool {
		first(v)
		return false
	}
}

func (l *List[T]) newItem(value T) *listItem[T] {
	return &listItem[T]{value, nil, nil, l}
}

func (l *List[T]) Get(idx int) (T, bool) {
	value := *new(T)
	if idx < 0 || idx >= l.length {
		return value, false
	}

	finder := func (i int, v T) bool {
		if i == idx {
			value = v
			return false
		}
		return true
	}

	if idx <= l.length {
		l.ForEachIndex(finder)
	} else {
		l.ForEachIndexReverse(finder)
	}
	return value, true
}

func (l *List[T]) Add(value T) {
	l.Append(value)
}

func (l *List[T]) Append(value T) *List[T] {
	item := l.newItem(value)
	item.prev = l.tail
	if l.length == 0 {
		l.head = item
		l.tail = item
	} else {
		l.tail.next = item
	}
	l.tail = item
	l.length++
	return l
}

func (l *List[T]) Prepend(value T) *List[T] {
	item := l.newItem(value)
	item.next = l.head
	if l.length == 0{
		l.head = item
		l.tail = item
	} else {
		l.head.prev = item
	}
	l.head = item
	l.length++
	return l
}

func (l *List[T]) PopFront() (T, bool) {
	item := remove(l.head)
	if item == nil {
		return *new(T), false
	}
	return item.value, true
}

func (l *List[T]) PeekFront() (T, bool) {
	if l.head == nil {
		return *new(T), false
	}
	return l.head.value, true
}

func (l *List[T]) PopBack() (T, bool) {
	item := remove(l.tail)
	if item == nil {
		return *new(T), false
	}
	return item.value, true
}

func (l *List[T]) PeekBack() (T, bool) {
	if l.tail == nil {
		return *new(T), false
	}
	return l.tail.value, true
}

func (l *List[T]) ForEachIndex(foreach func(idx int, value T) bool) {
	item := l.head
	idx := 0
	for item != nil {
		if !foreach(idx, item.value) {
			break
		}
		item = item.next
		idx++
	}
}

func (l *List[T]) ForEachIndexReverse(foreach func(idx int, value T) bool) {
	item := l.tail
	idx := 0
	for item != nil {
		if !foreach(idx, item.value) {
			break
		}
		item = item.prev
		idx++
	}
}

func (l *List[T]) ForEach(foreach func(value T) bool) {
	l.ForEachIndex(func (_ int, v T) bool {
		return foreach(v)
	})
}

func (l *List[T]) ForEachReverse(foreach func(value T) bool) {
	l.ForEachIndexReverse(func (_ int, v T) bool {
		return foreach(v)
	})
}

func (l *List[T]) Has(value T) bool {
	found := false
	l.ForEach(func(v T) bool {
		if v == value {
			found = true
		}
		return !found
	})
	return found
}

func (l *List[T]) Len() int {
	return l.length
}

func (l *List[T]) IsEmpty() bool {
	return l.head == nil
}

func (l *List[T]) String() string {
	var b strings.Builder
	b.WriteString("[")
	l.ForEach(All(func (v T) {
		b.WriteString(fmt.Sprintf(" %v", v))
	}))
	b.WriteString(" ]")
	return b.String()
}

func (l *List[T]) Values() []T {
	values := make([]T, l.length)
	l.ForEachIndex(AllIndex(func (idx int, v T) {
		values[idx] = v
	}))
	return values
}

func (l *List[T]) Equals(other *List[T]) bool {
	if other == nil {
		return false
	}
	
	if l.Len() != other.Len() {
		return false
	}

	left := l.head
	right := other.head
	for left != nil {
		if left.value != right.value {
			return false
		}
		left = left.next
		right = right.next
	}
	return true
}

func remove[T comparable](item *listItem[T]) *listItem[T] {
	if item == nil {
		return nil
	}

	owner := item.owner
	if owner == nil {
		item.prev = nil
		item.next = nil
		return item
	}

	prev := item.prev
	next := item.next
	
	if prev != nil {
		prev.next = next
	} else {
		owner.head = next
	}
	
	if next != nil {
		next.prev = prev
	} else {
		owner.tail = prev
	}

	item.prev = nil
	item.next = nil
	owner.length--
	return item
}