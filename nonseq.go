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
	key    []byte
	cipher SimonSpeckCipher
	seq    func() (seqid uint64, err error)
}

// Secret key ensures a unique permutation of the input sequence, so that only someone who knows the key can guess nonseqid value
// Key length determines block size and by this the number of significant bits in returned nonseqid (also its max value).
func NewGenerator(key []byte, seq func() (seqid uint64, err error)) (*Generator, error) {
	var cipher SimonSpeckCipher
	if blocksize, pres := keylen2blocksize[len(key)]; pres {
		if blocksize == 4 {
			cipher = simonspeck.NewSpeck32(key)
		} else if blocksize == 6 {
			cipher = simonspeck.NewSpeck48(key)
		} else if blocksize == 8 {
			cipher = simonspeck.NewSpeck64(key)
		} else {
			// should not happen
			return nil, fmt.Errorf("Internal error. Blocksize is %d", blocksize)
		}
		return &Generator{key, cipher, seq}, nil
	} else {
		return nil, fmt.Errorf("Allowed key length is 8, 9 or 12 bytes")
	}
}

// nonseqid is []byte of blocksize length (4, 6 or 8 bytes)
func (g *Generator) Next() (seqid uint64, nonseqid []byte, err error) {
	seqid, err = g.seq()
	// propagate the source sequence error
	if err != nil {
		return 0, nil, err
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, seqid)
	// trim b to blocksize
	b = b[(8 - g.cipher.BlockSize()):]
	nonseqid = make([]byte, g.cipher.BlockSize())
	g.cipher.Encrypt(nonseqid, b)
	return seqid, nonseqid, nil
}

func (g *Generator) Decode(nonseqid []byte) (seqid uint64, err error) {
	if len(nonseqid) != g.cipher.BlockSize() {
		return 0, fmt.Errorf("Wrong length of nonseqid. Actual=%d, Expected=%d", len(nonseqid), g.cipher.BlockSize())
	}
	block := make([]byte, g.cipher.BlockSize())
	g.cipher.Decrypt(block, nonseqid)
	b8 := make([]byte, 8)
	copy(b8[(8-len(block)):], block)
	seqid = binary.BigEndian.Uint64(b8)
	return seqid, nil

}
