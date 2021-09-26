package csvd

import (
	"unicode"
)

// CanonicalHeaderKey returns the canonical format of the header key.
// For example, the canonical key for " User\t  Id  " is "user id".
func CanonicalHeaderKey(key string) string {
	rk := []rune(key)

	i := len(rk)
	for ; i > 0; i-- {
		if !unicode.IsSpace(rk[i-1]) {
			break
		}
	}

	rk = rk[:i]
	removeSpace := true
	i = 0
	for _, r := range rk {
		if unicode.IsSpace(r) {
			if !removeSpace {
				rk[i] = ' '
				i++
				removeSpace = true
			}
		} else {
			rk[i] = unicode.ToLower(r)
			i++
			removeSpace = false
		}
	}

	return string(rk[:i])
}

// ParseHeader parses a csv file header.
func ParseHeader(record []string) Header {
	h := make(Header, len(record))
	for i, key := range record {
		h.Set(key, i)
	}
	return h
}

// Header represents the key-index pairs in the csv file header.
// The keys should be in canonical form, as returned by CanonicalHeaderKey.
type Header map[string]int

// Set sets the header entries associated with key to the single index.
// It replaces any existing values associated with key.
//
// The key is canonicalized by CanonicalHeaderKey.
func (h Header) Set(key string, index int) {
	h[CanonicalHeaderKey(key)] = index
}

// Get gets the index associated with the given key.
// If there are no index associated with the key, Get returns -1.
//
// The key is canonicalized by CanonicalHeaderKey.
func (h Header) Get(key string) int {
	idx, ok := h[CanonicalHeaderKey(key)]
	if !ok {
		return -1
	}
	return idx
}

// Get checks if has a index associated with the given key.
// If there are no index associated with the key, Get returns -1.
//
// The key is canonicalized by CanonicalHeaderKey.
func (h Header) Has(key string) bool {
	_, ok := h[CanonicalHeaderKey(key)]
	return ok
}

// Del deletes the index associated with the given key.
//
// The key is canonicalized by CanonicalHeaderKey.
func (h Header) Del(key string) {
	delete(h, CanonicalHeaderKey(key))
}
