package extrafield

const (
	UNIXTagType uint16 = 0x000d
)

type UNIXExtraField struct {
	Tag   uint16
	TSize uint16
	Atime uint32
	Mtime uint32
	Uid   uint16
	Gid   uint16
}
