package serial

import "io"

type ISerializable interface {
	Marshal(io.WriteSeeker) error
}

type IDeserializable interface {
	Unmarshal(io.ReadSeeker) error
}

type ISizeOf interface {
	SizeOf() uint32
}
