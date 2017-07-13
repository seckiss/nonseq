package nonseq

import (
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

type Generator struct {
	key       []byte
	blocksize int
	cipher    *simonspeck.Speck48Cipher
	seq       func() (seqid uint64, err error)
}

// Secret key ensures a unique permutation of the input sequence, so that only someone who knows the key can guess nonseqid value
// Key length determines block size and by this the number of significant bits in returned nonseqid (also its max value).
// nonseqid is uint64 so:
// - for 32-bit block size the 32 most significant bits of nonseqid should be zero
// - for 48-bit block size the 16 most significant bits of nonseqid should be zero
func NewGenerator(key []byte, seq func() (seqid uint64, err error)) (*Generator, error) {
	if blocksize, pres := keylen2blocksize[len(key)]; pres {
		cipher := simonspeck.NewSpeck48(key)
		return &Generator{key, blocksize, cipher, seq}, nil
	} else {
		return nil, fmt.Errorf("Allowed key length is 8, 9 or 12 bytes")
	}
}

func (g *Generator) Next() (seqid uint64, nonseqid uint64, err error) {
	return 0, 0, nil
}

func (g *Generator) Decode(nonseqid uint64) (seqid uint64, err error) {
	return 0, nil
}
