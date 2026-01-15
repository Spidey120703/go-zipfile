package crc

type CyclicRedundancyCheck32 struct {
	poly  uint32
	table [256]uint32
}

func (crc32 *CyclicRedundancyCheck32) Checksum(data []byte) uint32 {
	var crc uint32 = 0xffffffff

	for _, b := range data {
		crc = (crc >> 8) ^ crc32.table[(crc^uint32(b))&0xff]
	}

	return ^crc
}

func NewCRC32() *CyclicRedundancyCheck32 {
	var crc32 = &CyclicRedundancyCheck32{0xedb88320, [256]uint32{}}

	for i := 0; i < 256; i++ {
		crc := uint32(i)
		for j := 0; j < 8; j++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ crc32.poly
			} else {
				crc >>= 1
			}
		}
		crc32.table[i] = crc
	}

	return crc32
}

func ChecksumCRC32(data []byte) uint32 {
	const poly uint32 = 0xedb88320

	var crc uint32 = 0xffffffff

	for _, b := range data {
		crc ^= uint32(b)
		for i := 0; i < 8; i++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ poly
			} else {
				crc >>= 1
			}
		}
	}

	return ^crc
}
