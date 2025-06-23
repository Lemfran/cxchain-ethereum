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

// ToBytes converts a nibbles slice to bytes
func ToBytes(nibbles []byte) []byte {
	bytes := make([]byte, len(nibbles)/2)
	for i := 0; i < len(nibbles)/2; i++ {
		bytes[i] = nibbles[i*2]<<4 | nibbles[i*2+1]
	}
	return bytes
}
