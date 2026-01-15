package extrafield

const (
	NTFSTagType       uint16 = 0x000a
	NTFSAttribute1Tag        = 0x0001
)

type NTFSExtraField struct {
	Tag      uint16
	TSize    uint16
	Reserved uint32
	Tag1     uint16
	Size1    uint16
	Mtime    uint64
	Atime    uint64
	Ctime    uint64
}
