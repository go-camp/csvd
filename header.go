package csvd

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// CanonicalHeaderKey returns the canonical format of the header key.
// For example, the canonical key for " User\t  Id  " is "user id".
func CanonicalHeaderKey(key string) string {
	var hasSpace bool
	var start, size int
	var r rune
	var b strings.Builder

	for start < len(key) {
		r, size = utf8.DecodeRuneInString(key[start:])
		if !unicode.IsSpace(r) {
			break
		}
		start += size
	}
	key = key[start:]
	start = 0

	if start < len(key) {
		for {
			if unicode.IsUpper(r) {
				goto canonicalHeaderKey
			}

			if unicode.IsSpace(r) {
				if hasSpace || r != ' ' {
					goto canonicalHeaderKey
				}
				hasSpace = true
			} else {
				hasSpace = false
			}

			start += size
			if start < len(key) {
				r, size = utf8.DecodeRuneInString(key[start:])
			} else {
				break
			}
		}
	}

	if hasSpace {
		return key[:start-size]
	}

	return key

canonicalHeaderKey:
	b.Grow(len(key))
	b.WriteString(key[:start])

	for {
		if unicode.IsSpace(r) {
			hasSpace = true
		} else {
			if hasSpace {
				b.WriteByte(' ')
			}
			b.WriteRune(unicode.ToLower(r))
			hasSpace = false
		}

		start += size
		if start < len(key) {
			r, size = utf8.DecodeRuneInString(key[start:])
		} else {
			break
		}
	}

	return b.String()
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
