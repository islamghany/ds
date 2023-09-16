package hashtables

import hashfunctions "hashing/hashFunctions"

const (
	LOAD_FACTOR = 0.75
)

type Node[T any] struct {
	key   string
	value T
	next  *Node[T]
}

type LinkedList[T any] struct {
	head *Node[T]
	size int
}

func (l *LinkedList[T]) Insert(key string, value T) {
	node := &Node[T]{key, value, nil}
	if l.head == nil {
		l.head = node
	} else {
		node.next = l.head
		l.head = node
	}
	l.size++
}

func (l *LinkedList[T]) Search(key string) (*T, bool) {
	if l.head == nil {
		return nil, false
	}
	current := l.head
	for current != nil {
		if current.key == key {
			return &current.value, true
		}
		current = current.next
	}
	return nil, false
}

func (l *LinkedList[T]) Delete(key string) {
	if l.head == nil {
		return
	}
	if l.head.key == key {
		l.head = l.head.next
		l.size--
		return
	}
	current := l.head
	for current.next != nil {
		if current.next.key == key {
			current.next = current.next.next
			l.size--
			return
		}
		current = current.next
	}
}

type SeperateChaining[T any] struct {
	table        []LinkedList[T]
	tsize        int // size of the table
	count        int // number of elements in the table
	hashFunction hashfunctions.HashFunction
}

func NewSeperateChaining[T any](size int, hashFunction hashfunctions.HashFunction) *SeperateChaining[T] {
	return &SeperateChaining[T]{make([]LinkedList[T], size), size, 0, hashFunction}
}

func (s *SeperateChaining[T]) rehash() {
	oldTable := s.table
	s.tsize = s.tsize * 2
	s.table = make([]LinkedList[T], s.tsize)
	for _, list := range oldTable {
		current := list.head
		for current != nil {
			s.Insert(current.key, current.value)
			current = current.next
		}
	}

}
func (s *SeperateChaining[T]) Insert(key string, value T) {
	index := s.hashFunction.Hash(key) % s.tsize
	s.table[index].Insert(key, value)
	s.count++
	if float64(s.count)/float64(s.tsize) > LOAD_FACTOR {
		s.rehash()
	}
}

func (s *SeperateChaining[T]) Search(key string) (*T, bool) {
	index := s.hashFunction.Hash(key) % s.tsize
	return s.table[index].Search(key)
}

func (s *SeperateChaining[T]) Delete(key string) {
	index := s.hashFunction.Hash(key) % s.tsize
	s.table[index].Delete(key)
}

func (s *SeperateChaining[T]) IsEmpty() bool {
	return s.count == 0
}

func (s *SeperateChaining[T]) Size() int {
	return s.count
}
