package trie

type Trie interface {
	Insert(key []byte, value []byte) error
	Has(key []byte) (bool, error)
}