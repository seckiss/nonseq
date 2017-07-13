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
	for i := 0; i < 10; i++ {
		seqid, nonseqid, _ := gen.Next()
		fmt.Printf("seqid=%d, nonseqid=%v\n", seqid, nonseqid)
	}
}
