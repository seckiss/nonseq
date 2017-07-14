package nonseq

import (
	"fmt"
	"testing"
)

func TestGeneratorKey4(t *testing.T) {
	testGenKey(t, 4)
}
func TestGeneratorKey6(t *testing.T) {
	testGenKey(t, 6)
}
func TestGeneratorKey8(t *testing.T) {
	testGenKey(t, 8)
}
func TestGeneratorKey12(t *testing.T) {
	testGenKey(t, 12)
}
func TestGeneratorKey16(t *testing.T) {
	testGenKey(t, 16)
}

func testGenKey(t *testing.T, blocksize int) {
	gen := getGenerator()
	for i := 0; i < 10; i++ {
		nonseqid := make([]byte, blocksize)
		seqid, _ := gen.Next(nonseqid)
		fmt.Printf("seqid=%d, nonseqid=%v\n", seqid, nonseqid)
	}
}

func getGenerator() *Generator {
	var counter uint64
	seq := func() (uint64, error) {
		counter++
		return counter, nil
	}
	key := []byte("0123456789abcdef")
	return NewGenerator(key, seq)
}

func testDecode(t *testing.T, nonseqid []byte) (seqid uint64, err error) {
	gen := getGenerator()
	seqid, err = gen.Decode(nonseqid)
	fmt.Printf("seqid=%d, err=%v\n", seqid, err)
	return seqid, err
}

func TestDecode4(t *testing.T) {
	inp, want := []byte{154, 255, 88, 12}, uint64(3)
	got, err := testDecode(t, inp)
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
