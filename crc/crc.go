package crc

type UintLike interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type ICyclicRedundancyCheck[T UintLike] interface {
	Checksum(data []byte) T
}
