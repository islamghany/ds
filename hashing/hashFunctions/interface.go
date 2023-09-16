package hashfunctions

type HashFunction interface {
	Hash(key string) int
}
