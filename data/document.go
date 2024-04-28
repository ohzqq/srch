package data

import "github.com/bits-and-blooms/bloom"

type Doc struct {
	SrchFields map[string]*bloom.BloomFilter
	Facets     map[string]*bloom.BloomFilter
}
