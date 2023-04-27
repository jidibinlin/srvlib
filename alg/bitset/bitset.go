package bitset

const (
	shift = 5
	mask  = 0x1F
)

/**
 * Returns the index of the lowest set bit in a 32-bit word. In other words,
 * counts the trailing zeroes on the word.
**/

// A list of bit indexes for every product of a word modulo 37.
// See http://graphics.stanford.edu/~seander/bithacks.html#ZerosOnRightModLookup

var mod37BitPosition = [37]int{
	32, 0, 1, 26, 2, 23, 27, 0, 3, 16,
	24, 30, 28, 11, 0, 13, 4, 7, 17, 0,
	25, 22, 31, 15, 29, 10, 12, 6, 0, 21,
	14, 9, 5, 20, 8, 19, 18,
}

func bitPosition(value int) int {
	return mod37BitPosition[(-value&value)%37]
}

type BitSet struct {
	_size uint32
	_bits []uint32
}

func New(max uint32) *BitSet {
	bs := &BitSet{}
	size := (max >> shift) + 1
	bs._size = size << shift
	bs._bits = make([]uint32, size)
	return bs
}

func (bs *BitSet) Set(bit uint32) {
	if bit >= bs._size {
		return
	}

	bs._bits[bit>>shift] |= 1 << (bit & mask)
}

func (bs *BitSet) Get(bit uint32) bool {
	if bit >= bs._size {
		return false
	}

	return (bs._bits[bit>>shift] & (1 << (bit & mask))) > 0
}

func (bs *BitSet) Unset(bit uint32) {
	if bit >= bs._size {
		return
	}

	bs._bits[bit>>shift] &= ^(1 << (bit & mask))
}

func (bs *BitSet) Test(bit uint32) bool {
	if bit >= bs._size {
		return false
	}

	return bs._bits[bit>>shift]&(1<<(bit&mask)) != 0
}

func (bs *BitSet) Range(fn func(uint32)) {
	for idx, value := range bs._bits {
		for value != 0 {
			key := idx<<shift + bitPosition(int(value))
			fn(uint32(key))
			value &= value - 1
		}
	}
}

func (bs *BitSet) Clear() {
	bs._bits = make([]uint32, bs._size)
}

func (bs *BitSet) GetSlice() []uint32 {
	return bs._bits
}

func (bs *BitSet) SetSlice(bits []uint32) {
	bs._bits = bits
}
