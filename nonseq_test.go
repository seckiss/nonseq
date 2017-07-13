package nonseq

import (
	"fmt"
	"testing"
)

func TestGeneratorKey8(t *testing.T) {
	testGenKey(t, 8)
}
func TestGeneratorKey9(t *testing.T) {
	testGenKey(t, 9)
}
func TestGeneratorKey12(t *testing.T) {
	testGenKey(t, 12)
}

func testGenKey(t *testing.T, keylen int) {
	gen := getGenerator(t, keylen)
	for i := 0; i < 10; i++ {
		seqid, nonseqid, _ := gen.Next()
		fmt.Printf("seqid=%d, nonseqid=%v\n", seqid, nonseqid)
	}
}

func getGenerator(t *testing.T, keylen int) *Generator {
	var counter uint64
	seq := func() (uint64, error) {
		counter++
		return counter, nil
	}
	longkey := []byte("123456789abcdefghij")
	gen, err := NewGenerator(longkey[:keylen], seq)
	if err != nil {
		t.Fatalf("NewGenerator error=%v", err)
	}
	return gen
}

func TestDecode8(t *testing.T) {
	nonseqid := []byte{42, 5, 18, 59}
	gen := getGenerator(t, 8)
	seqid, err := gen.Decode(nonseqid)
	fmt.Printf("seqid=%d, err=%v\n", seqid, err)

}
