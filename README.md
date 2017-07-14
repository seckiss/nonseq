# nonseq

## Non sequential unique id generator from sequential source



Generate non-sequential unique IDs from a serial uint64 sequence (like from Postgres bigserial sequence or from a simple non-durable counter)

```
var counter uint64
// sequence generator
seq := func() (uint64, error) {
	counter++
	return counter, nil
}
gen := NewGenerator(getKey(), seq)
nonseqid := make([]byte, 6)
// 6 byte non-sequential ID
seqid, err := gen.Next(nonseqid)
```
