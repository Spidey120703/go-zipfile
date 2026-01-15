package posix

// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/uapi/linux/stat.h
const (
	StatIsExecutableOthers uint16 = 00001
	StatIsWriteableOthers         = 00002
	StatIsReadableOthers          = 00004
	StatIsExecutableGroup         = 00010
	StatIsWriteableGroup          = 00020
	StatIsReadableGroup           = 00040
	StatIsExecutableUser          = 00100
	StatIsWriteableUser           = 00200
	StatIsReadableUser            = 00400
	StatIsSticky                  = 0001000
	StatIsSetGroupID              = 0002000
	StatIsSetUserID               = 0004000
	StatIsNamedPipe               = 0010000
	StatIsCharacterDevice         = 0020000
	StatIsDirectory               = 0040000
	StatIsBlockDevice             = 0060000
	StatIsRegularFile             = 0100000
	StatIsSymbolicLink            = 0120000
	StatIsSocket                  = 0140000
)
