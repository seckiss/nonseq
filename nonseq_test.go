package nonseq

import (
	"fmt"
	"testing"
)

func TestGenNext4(t *testing.T) {
	testGenNext(t, 4)
}
func TestGenNext6(t *testing.T) {
	testGenNext(t, 6)
}
func TestGenNext8(t *testing.T) {
	testGenNext(t, 8)
}
func TestGenNext12(t *testing.T) {
	testGenNext(t, 12)
}
func TestGenNext16(t *testing.T) {
	testGenNext(t, 16)
}

func TestGenNextNotExceed4(t *testing.T) {
	// 4 billion should not exceed 2^32
	seq := func() (uint64, error) {
		return 4000000000, nil
	}
	gen := NewGenerator(getKey(), seq)
	blocksize := 4
	nonseqid := make([]byte, blocksize)
	seqid, err := gen.Next(nonseqid)
	if err != nil {
		t.Fatalf("got error: %v, expected nil for seqid=%v", seqid)
	}
}

func TestGenNextExceed4(t *testing.T) {
	// 5 billion should exceed 2^32
	seq := func() (uint64, error) {
		return 5000000000, nil
	}
	gen := NewGenerator(getKey(), seq)
	blocksize := 4
	nonseqid := make([]byte, blocksize)
	seqid, err := gen.Next(nonseqid)
	//fmt.Printf("expected error: %v\n", err)
	if err == nil {
		t.Fatalf("got error nil while expected error exceeding for seqid=%v", seqid)
	}
}

func testGenNext(t *testing.T, blocksize int) {
	gen := getGenerator()
	for i := 0; i < 10; i++ {
		nonseqid := make([]byte, blocksize)
		seqid, err := gen.Next(nonseqid)
		//fmt.Printf("seqid=%d, nonseqid=%v, err=%v\n", seqid, nonseqid, err)
		_ = seqid
		_ = err
	}
}

func TestB58GenNext6(t *testing.T) {
	gen := getB58Generator()
	testB58GenNext(t, gen.Next6)
}
func TestB58GenNext9(t *testing.T) {
	gen := getB58Generator()
	testB58GenNext(t, gen.Next9)
}
func TestB58GenNext11(t *testing.T) {
	gen := getB58Generator()
	testB58GenNext(t, gen.Next11)
}
func TestB58GenNext17(t *testing.T) {
	gen := getB58Generator()
	testB58GenNext(t, gen.Next17)
}
func TestB58GenNext22(t *testing.T) {
	gen := getB58Generator()
	testB58GenNext(t, gen.Next22)
}

func testB58GenNext(t *testing.T, next func() (uint64, string, error)) {
	for i := 0; i < 10; i++ {
		seqid, cram, err := next()
		fmt.Printf("seqid=%d, cram=%v, err=%v\n", seqid, cram, err)
		_ = seqid
		_ = cram
		_ = err
	}
}

func getGenerator() *Generator {
	var counter uint64
	seq := func() (uint64, error) {
		counter++
		return counter, nil
	}
	return NewGenerator(getKey(), seq)
}

func getB58Generator() *B58Generator {
	var counter uint64
	seq := func() (uint64, error) {
		counter++
		return counter, nil
	}
	return NewB58Generator(getKey(), seq)
}

func testB58Decode(t *testing.T, cram string) (seqid uint64, err error) {
	gen := getB58Generator()
	seqid, err = gen.Decode(cram)
	//fmt.Printf("seqid=%d, err=%v\n", seqid, err)
	return seqid, err
}

func getKey() []byte {
	return []byte("0123456789abcdef")
}

func testDecode(t *testing.T, nonseqid []byte) (seqid uint64, err error) {
	gen := getGenerator()
	seqid, err = gen.Decode(nonseqid)
	//fmt.Printf("seqid=%d, err=%v\n", seqid, err)
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

func TestB58Decode6(t *testing.T) {
	inp, want := "4xnru9", uint64(3)
	got, err := testB58Decode(t, inp)
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
