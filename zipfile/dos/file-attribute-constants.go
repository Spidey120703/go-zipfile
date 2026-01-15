package dos

// https://learn.microsoft.com/en-us/windows/win32/fileio/file-attribute-constants
const (
	FileAttributeReadonly           uint32 = 0x00000001
	FileAttributeHidden                    = 0x00000002
	FileAttributeSystem                    = 0x00000004
	FileAttributeDirectory                 = 0x00000010
	FileAttributeArchive                   = 0x00000020
	FileAttributeDevice                    = 0x00000040
	FileAttributeNormal                    = 0x00000080
	FileAttributeTemporary                 = 0x00000100
	FileAttributeSparseFile                = 0x00000200
	FileAttributeReparsePoint              = 0x00000400
	FileAttributeCompressed                = 0x00000800
	FileAttributeOffline                   = 0x00001000
	FileAttributeNotContentIndexed         = 0x00002000
	FileAttributeEncrypted                 = 0x00004000
	FileAttributeIntegrityStream           = 0x00008000
	FileAttributeVirtual                   = 0x00010000
	FileAttributeNoScrubData               = 0x00020000
	FileAttributeEA                        = 0x00040000
	FileAttributePinned                    = 0x00080000
	FileAttributeUnpinned                  = 0x00100000
	FileAttributeRecallOnOpen              = 0x00040000
	FileAttributeRecallOnDataAccess        = 0x00400000
)
