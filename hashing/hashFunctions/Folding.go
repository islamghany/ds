package hashfunctions

import "strconv"

type Folding struct{}

func (f Folding) Hash(key string) int {
	hash := 0
	for i := 0; i < len(key); i += 4 {
		end := i + 4
		if end > len(key) {
			end = len(key)
		}
		num, _ := strconv.Atoi(key[i:end])
		hash += num
	}
	return hash
}

func NewFolding() *Folding {
	return &Folding{}
}
