package mpt

// ToNibbles converts a byte slice to nibbles
func ToNibbles(bytes []byte) []byte {
	nibbles := make([]byte, len(bytes)*2)
	for i, b := range bytes {
		nibbles[i*2] = b >> 4
		nibbles[i*2+1] = b & 0x0f
	}
	return nibbles
}

/*// FromNibbles converts nibbles back to bytes
func FromNibbles(nibbles []byte) []byte {
	if len(nibbles)%2 != 0 {
		nibbles = append([]byte{0}, nibbles...)
	}
	
	bytes := make([]byte, len(nibbles)/2)
	for i := 0; i < len(nibbles); i += 2 {
		bytes[i/2] = byte(nibbles[i]<<4) | byte(nibbles[i+1])
	}
	return bytes
}

// ToPrefixed adds hex prefix to nibbles
func ToPrefixed(nibbles []byte, isLeaf bool) []byte {
	prefix := []byte{0}
	if len(nibbles)%2 != 0 {
		prefix = []byte{1}
	}
	
	if isLeaf {
		prefix[0] += 2
	}
	
	result := make([]byte, len(prefix)+len(nibbles))
	copy(result, prefix)
	copy(result[len(prefix):], nibbles)
	return result
}

// FromPrefixed removes hex prefix from nibbles
func FromPrefixed(prefixed []byte) ([]byte, bool) {
	if len(prefixed) == 0 {
		return nil, false
	}
	
	isLeaf := prefixed[0] >= 2
	odd := prefixed[0]%2 == 1
	
	if odd {
		return prefixed[1:], isLeaf
	}
	return prefixed[2:], isLeaf
}*/