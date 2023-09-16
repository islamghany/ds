package hashtables

type HashTable interface {
	Insert(key string, value interface{})
	Search(key string) (interface{}, bool)
	Delete(key string)
	IsEmpty() bool
	Size() int
}
