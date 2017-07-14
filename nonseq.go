package nonseq

import (
	"encoding/binary"
	"fmt"

	"github.com/ankitkalbande/simonspeck"
)

// Simon/Speck cipher key size mapping to block size. There are many variants
// Block size/key length in bits:
// - 32/64
// - 48/72, 48/96
// - 64/96, 64/128
// - 96/96, 96/144
// - 128/128, 128/192, 128/256
// We chose 32/64, 48/72 and 64/96 which corresponds to the reversed mapping in bytes:
var keylen2blocksize = map[int]int{8: 4, 9: 6, 12: 8}

// missing interface in simonspec library
type SimonSpeckCipher interface {
	Encrypt(dst, src []byte)
	Decrypt(dst, src []byte)
	BlockSize() int
}

type Generator struct {
	key []byte
	// block size to corresponding cipher mapping
	cipher map[int]SimonSpeckCipher
	seq    func() (seqid uint64, err error)
}

// Secret key ensures a unique permutation of the input sequence, so that only someone who knows the key can guess nonseqid value
// The key should be 16 bytes however only part of it is used for weaker ciphers
func NewGenerator(key []byte, seq func() (seqid uint64, err error)) *Generator {
	if len(key) != 16 {
		// wrong key length means developer error
		panic("Key length should be 16 bytes")
	}
	g := &Generator{key, make(map[int]SimonSpeckCipher), seq}
	g.cipher[4] = simonspeck.NewSpeck32(key[:8])
	g.cipher[6] = simonspeck.NewSpeck48(key[:9])
	g.cipher[8] = simonspeck.NewSpeck64(key[:12])
	g.cipher[12] = simonspeck.NewSpeck96(key[:12])
	g.cipher[16] = simonspeck.NewSpeck128(key)
	return g
}

// nonseqid is []byte of blocksize length (4, 6, 8, 12 or 16)
// it will be filled with nonseqid generated from seqid which is also returned
func (g *Generator) Next(nonseqid []byte) (seqid uint64, err error) {
	blocksize := len(nonseqid)
	c := g.cipher[blocksize]
	if c == nil {
		// wrong block size means developer error
		panic("Block size should be 4, 6, 8, 12 or 16 bytes")
	}
	seqid, err = g.seq()
	// propagate the source sequence error
	if err != nil {
		return 0, err
	}
	// convert seqid to []byte of same length as nonseqid
	bytes, err := toBytes(seqid, blocksize)
	if err != nil {
		return seqid, err
	}
	c.Encrypt(nonseqid, bytes)
	return seqid, nil
}

func isZeroSlice(b []byte) bool {
	for _, bb := range b {
		if bb != 0 {
			return false
		}
	}
	return true
}

// Convert uint64 into []byte
// Also for size < 8 return error if given uint64 id exceeds maximum number encodable in size bytes
func toBytes(id uint64, size int) ([]byte, error) {
	b8 := make([]byte, 8)
	binary.BigEndian.PutUint64(b8, id)
	if size < 8 {
		if !isZeroSlice(b8[:(8 - size)]) {
			return nil, fmt.Errorf("id %d exceeds maximum number encodeable in %d bytes", id, size)
		}
		// trim to blocksize
		return b8[(8 - size):], nil
	} else if size > 8 {
		// expand to blocksize
		bytes := make([]byte, size)
		copy(bytes[(size-8):], b8)
		return bytes, nil
	} else {
		return b8, nil
	}
}

// Convert []byte to uint64
// Also for len(b) > 8 return error if the trimmed MSB bytes are non-zero
// It may work as a correctness checksum for 12 and 16 byte blocksize
func fromBytes(b []byte) (id uint64, err error) {
	size := len(b)
	b8 := make([]byte, 8)
	if size < 8 {
		copy(b8[(8-size):], b)
	} else if size >= 8 {
		copy(b8, b[(size-8):])
		if !isZeroSlice(b[0:(size - 8)]) {
			err = fmt.Errorf("Trimmed MSB bytes are non-zero in %v", b)
		}
	}
	return binary.BigEndian.Uint64(b8), err
}

func (g *Generator) Decode(nonseqid []byte) (seqid uint64, err error) {
	blocksize := len(nonseqid)
	c := g.cipher[blocksize]
	if c == nil {
		return 0, fmt.Errorf("Block size should be 4, 6, 8, 12 or 16 bytes")
	}
	block := make([]byte, blocksize)
	c.Decrypt(block, nonseqid)
	seqid, err = fromBytes(block)
	// rewrite error to be more informative
	if err != nil {
		err = fmt.Errorf("The nonseq %v is decodable but does not come from this generator", nonseqid)
	}
	return seqid, err
}
