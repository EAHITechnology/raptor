package utils

const (
	BIG_M = 0xc6a4a7935bd1e995
	BIG_R = 47
	SEED  = 0x1234ABCD
)

/*
murmur hash function.
When the number of unique elements is 10000000, the hash distribution is as follows:
|--------------------||--------------------|
| bucket |    num    || bucket |    num    |
|--------------------||--------------------|
|  0     |  1000095  ||  5     |  999697   |
|--------------------||--------------------|
|  1     |  1000449  ||  6     |  1000249  |
|--------------------||--------------------|
|  2     |  999580   ||  7     |  1000917  |
|--------------------||--------------------|
|  3     |  999182   ||  8     |  1000388  |
|--------------------||--------------------|
|  4     |  998817   ||  9     |  1000626  |
|--------------------||--------------------|

E(X) = 389275.8
*/
func MurmurHash64A(data []byte) (h uint64) {
	var k uint64
	h = SEED ^ uint64(len(data))*BIG_M

	var ubigm uint64 = BIG_M
	var ibigm = ubigm
	for l := len(data); l >= 8; l -= 8 {
		k = uint64(uint64(data[0]) | uint64(data[1])<<8 | uint64(data[2])<<16 | uint64(data[3])<<24 |
			uint64(data[4])<<32 | uint64(data[5])<<40 | uint64(data[6])<<48 | uint64(data[7])<<56)

		k = k * ibigm
		k ^= uint64(k) >> BIG_R
		k = k * ibigm

		h = h ^ k
		h = h * ibigm
		data = data[8:]
	}

	switch len(data) {
	case 7:
		h ^= uint64(data[6]) << 48
		fallthrough
	case 6:
		h ^= uint64(data[5]) << 40
		fallthrough
	case 5:
		h ^= uint64(data[4]) << 32
		fallthrough
	case 4:
		h ^= uint64(data[3]) << 24
		fallthrough
	case 3:
		h ^= uint64(data[2]) << 16
		fallthrough
	case 2:
		h ^= uint64(data[1]) << 8
		fallthrough
	case 1:
		h ^= uint64(data[0])
		h *= ibigm
	}

	h ^= uint64(h >> BIG_R)
	h *= ibigm
	h ^= uint64(h >> BIG_R)
	return
}
