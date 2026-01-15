package zipfile

import (
	"go-zipfile/zipfile/dos"
)

type Signature [4]byte

type LocalFileHeader struct {
	Signature         Signature
	Version           uint16
	Flags             uint16
	CompressionMethod uint16
	LastModFileTime   *dos.Time
	LastModFileDate   *dos.Date
	CRC32             uint32
	CompressedSize    uint32
	UncompressedSize  uint32
	FileNameLength    uint16
	ExtraFieldLength  uint16
	FileName          []byte `serial:"len=FileNameLength"`
	ExtraField        []byte `serial:"len=ExtraFieldLength"`
}

func (lfh *LocalFileHeader) SizeOf() uint32 {
	return 0 +
		4 + /* local file header signature     4 bytes  (0x04034b50) */
		2 + /* version needed to extract       2 bytes               */
		2 + /* general purpose bit flag        2 bytes               */
		2 + /* compression method              2 bytes               */
		2 + /* last mod file time              2 bytes               */
		2 + /* last mod file date              2 bytes               */
		4 + /* crc-32                          4 bytes               */
		4 + /* compressed size                 4 bytes               */
		4 + /* uncompressed size               4 bytes               */
		2 + /* file name length                2 bytes               */
		2 + /* extra field length              2 bytes               */
		uint32(lfh.FileNameLength) + /* file name (variable size)    */
		uint32(lfh.ExtraFieldLength) /* extra field (variable size)  */
}

/*
TODO: lazy reader/writer

type FileData struct {
	file   *os.File
	offset int64
	length int64
}

func (fd *FileData) Unmarshal(reader io.ReadSeeker) error {
	panic("implement me")
}

func (fd *FileData) Marshal(writer io.WriteSeeker) error {
	panic("implement me")
}
*/

type FileData []byte

type DataDescriptor struct {
	CRC32            uint32
	CompressedSize   uint32
	UncompressedSize uint32
}

type CentralDirectoryFileHeader struct {
	Signature              Signature
	Version                uint16
	VersionNeeded          uint16
	Flags                  uint16
	CompressionMethod      uint16
	LastModFileTime        *dos.Time
	LastModFileDate        *dos.Date
	CRC32                  uint32
	CompressedSize         uint32
	UncompressedSize       uint32
	FileNameLength         uint16
	ExtraFieldLength       uint16
	FileCommentLength      uint16
	DiskNumberStart        uint16
	InternalFileAttributes uint16
	ExternalFileAttributes uint32
	OffsetOfLocalHeader    uint32
	FileName               []byte `serial:"len=FileNameLength"`
	ExtraField             []byte `serial:"len=ExtraFieldLength"`
	FileComment            []byte `serial:"len=FileCommentLength"`
}

func (cdf *CentralDirectoryFileHeader) SizeOf() uint32 {
	return 0 +
		4 + /* central file header signature   4 bytes  (0x02014b50)   */
		2 + /* version made by                 2 bytes                 */
		2 + /* version needed to extract       2 bytes                 */
		2 + /* general purpose bit flag        2 bytes                 */
		2 + /* compression method              2 bytes                 */
		2 + /* last mod file time              2 bytes                 */
		2 + /* last mod file date              2 bytes                 */
		4 + /* crc-32                          4 bytes                 */
		4 + /* compressed size                 4 bytes                 */
		4 + /* uncompressed size               4 bytes                 */
		2 + /* file name length                2 bytes                 */
		2 + /* extra field length              2 bytes                 */
		2 + /* file comment length             2 bytes                 */
		2 + /* disk number start               2 bytes                 */
		2 + /* internal file attributes        2 bytes                 */
		4 + /* external file attributes        4 bytes                 */
		4 + /* relative offset of local header 4 bytes                 */
		uint32(cdf.FileNameLength) + /*   file name (variable size)    */
		uint32(cdf.ExtraFieldLength) + /* extra field (variable size)  */
		uint32(cdf.FileCommentLength) /*  file comment (variable size) */
}

type DigitalSignature struct {
	Signature     Signature `serial:"prefix='PK\x05\x05'"`
	DataSize      uint16
	SignatureData []byte
}

type EndOfCentralDirectoryRecord struct {
	Signature                  Signature
	DiskNumber                 uint16
	StartingDiskNumber         uint16
	DiskTotalEntries           uint16
	TotalEntries               uint16
	CentralDirectorySize       uint32
	OffsetOfStartingDiskNumber uint32
	ZIPFileCommentLength       uint16
	ZIPFileComment             []byte `serial:"len=ZIPFileCommentLength"`
}

func (eocdr *EndOfCentralDirectoryRecord) SizeOf() uint32 {
	return 0 +
		4 + /* end of central dir signature    4 bytes  (0x06054b50)                  */
		2 + /* number of this disk             2 bytes                                */
		0 + /* number of the disk with the                                            */
		2 + /* start of the central directory  2 bytes                                */
		0 + /* total number of entries in the                                         */
		2 + /* central directory on this disk  2 bytes                                */
		0 + /* total number of entries in                                             */
		2 + /* the central directory           2 bytes                                */
		4 + /* size of the central directory   4 bytes                                */
		0 + /* offset of start of central                                             */
		0 + /* directory with respect to                                              */
		4 + /* the starting disk number        4 bytes                                */
		2 + /* .ZIP file comment length        2 bytes                                */
		uint32(eocdr.ZIPFileCommentLength) /* .ZIP file comment       (variable size) */

}

type LocalFileRecord struct {
	LocalFileHeader LocalFileHeader
	FileData        FileData        `serial:"len=LocalFileHeader.CompressedSize,stream"`
	DataDescriptor  *DataDescriptor `serial:"condition=bit(LocalFileHeader.Flags,3)"`
}

type CentralDirectoryRecord struct {
	CentralDirectoryHeaders []CentralDirectoryFileHeader `serial:"prefix='PK\x01\x02'"`
	DigitalSignature        *DigitalSignature            `serial:"condition=false"`
}

type FileFormat struct {
	LocalFileRecords            []LocalFileRecord `serial:"prefix='PK\x03\x04'"`
	CentralDirectoryRecord      CentralDirectoryRecord
	EndOfCentralDirectoryRecord EndOfCentralDirectoryRecord
}
