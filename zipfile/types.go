package zipfile

import (
	"bytes"
	"compress/flate"
	"errors"
	"go-zipfile/crc"
	"go-zipfile/serial"
	"go-zipfile/zipfile/dos"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

var crc32 *crc.CyclicRedundancyCheck32

func init() {
	crc32 = crc.NewCRC32()
}

type FileEntry struct {
	FilePath          string
	CreationTime      windows.Filetime
	LastAccessTime    windows.Filetime
	LastWriteTime     windows.Filetime
	FileAttributes    uint32
	FileSize          uint32
	CRC32             uint32
	DataSize          uint32
	Data              FileData
	CompressionMethod uint16
}

func (e *FileEntry) Deflate(level int) (err error) {
	if e.FileSize < 16 {
		return
	}

	var buf bytes.Buffer

	writer, err := flate.NewWriter(&buf, level)
	if err != nil {
		return
	}

	if _, err = writer.Write(e.Data); err != nil {
		return
	}

	if err = writer.Close(); err != nil {
		return
	}

	e.CompressionMethod = CompressionMethodDeflated
	e.Data = buf.Bytes()
	e.DataSize = uint32(len(e.Data))

	return
}

func unixNanoseconds(nanoseconds int64) time.Time {
	return time.Unix(nanoseconds/int64(time.Second), nanoseconds%int64(time.Second))
}

func convertFiletime(filetime windows.Filetime) (*dos.Date, *dos.Time) {
	unix := unixNanoseconds(filetime.Nanoseconds())
	return &dos.Date{
			Year:  uint16(unix.Year()),
			Month: uint16(unix.Month()),
			Day:   uint16(unix.Day()),
		}, &dos.Time{
			Hour:   uint16(unix.Hour()),
			Minute: uint16(unix.Minute()),
			Second: uint16(unix.Second()),
		}
}

func NewFileEntry(path string) *FileEntry {
	entry := &FileEntry{
		FilePath: strings.ReplaceAll(path, "\\", "/"),
	}

	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			return nil
		}
	}

	attrs, err := windows.GetFileAttributes(windows.StringToUTF16Ptr(path))
	if err != nil {
		return nil
	}
	entry.FileAttributes = attrs

	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(path),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL|windows.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)

	var createTime, accessTime, writeTime windows.Filetime
	if err = windows.GetFileTime(handle, &createTime, &accessTime, &writeTime); err != nil {
		return nil
	}

	entry.CreationTime = createTime
	entry.LastAccessTime = accessTime
	entry.LastWriteTime = writeTime

	if stat.IsDir() {
		return entry
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	entry.CRC32 = crc32.Checksum(data)
	entry.Data = data
	entry.FileSize = uint32(stat.Size())
	entry.DataSize = entry.FileSize
	entry.CompressionMethod = CompressionMethodStored

	return entry
}

type Zip struct {
	CompressionMethod uint16
	CompressionLevel  int
	FileEntries       []*FileEntry
}

func NewZip() *Zip {
	return &Zip{
		CompressionMethod: CompressionMethodStored,
		CompressionLevel:  flate.DefaultCompression,
	}
}

func (z *Zip) SetCompressionMethod(method uint16) {
	z.CompressionMethod = method
}

func (z *Zip) SetDeflateLevel(level int) {
	z.CompressionLevel = level
}

func (z *Zip) Add(path string) (err error) {
	entry := NewFileEntry(path)

	switch z.CompressionMethod {
	case CompressionMethodStored:
	case CompressionMethodDeflated:
		if err = entry.Deflate(z.CompressionLevel); err != nil {
			return
		}
	default:
	}

	z.FileEntries = append(z.FileEntries, entry)
	return
}

func (z *Zip) Build() (FileFormat, error) {
	ff := FileFormat{}

	var offset, cdhSize uint32

	for _, entry := range z.FileEntries {
		LastModFileTime, LastModFileDate := convertFiletime(entry.LastWriteTime)
		FileNameLength := uint16(len(entry.FilePath))
		FileName := []byte(entry.FilePath)

		lfh := LocalFileHeader{
			Signature:         LocalFileHeaderSignature,
			Version:           DefaultVersion,
			Flags:             0,
			CompressionMethod: entry.CompressionMethod,
			LastModFileTime:   LastModFileDate,
			LastModFileDate:   LastModFileTime,
			CRC32:             entry.CRC32,
			CompressedSize:    entry.DataSize,
			UncompressedSize:  entry.FileSize,
			FileNameLength:    FileNameLength,
			ExtraFieldLength:  0,
			FileName:          FileName,
			ExtraField:        nil,
		}

		ff.LocalFileRecords = append(ff.LocalFileRecords, LocalFileRecord{
			LocalFileHeader: lfh,
			FileData:        entry.Data,
		})

		cdh := CentralDirectoryFileHeader{
			Signature:              CentralFileHeaderSignature,
			Version:                LatestVersion,
			VersionNeeded:          DefaultVersion,
			Flags:                  0,
			CompressionMethod:      entry.CompressionMethod,
			LastModFileTime:        LastModFileDate,
			LastModFileDate:        LastModFileTime,
			CRC32:                  entry.CRC32,
			CompressedSize:         entry.DataSize,
			UncompressedSize:       entry.FileSize,
			FileNameLength:         FileNameLength,
			ExtraFieldLength:       0,
			FileCommentLength:      0,
			DiskNumberStart:        0,
			InternalFileAttributes: 0,
			ExternalFileAttributes: entry.FileAttributes,
			OffsetOfLocalHeader:    offset,
			FileName:               FileName,
			ExtraField:             nil,
			FileComment:            nil,
		}
		ff.CentralDirectoryRecord.CentralDirectoryHeaders = append(ff.CentralDirectoryRecord.CentralDirectoryHeaders, cdh)
		offset += lfh.SizeOf() + entry.DataSize
		cdhSize += cdh.SizeOf()
	}

	TotalEntries := uint16(len(z.FileEntries))
	ff.EndOfCentralDirectoryRecord = EndOfCentralDirectoryRecord{
		Signature:                  EndOfCentralDirectorySignature,
		DiskNumber:                 0,
		StartingDiskNumber:         0,
		DiskTotalEntries:           TotalEntries,
		TotalEntries:               TotalEntries,
		CentralDirectorySize:       cdhSize,
		OffsetOfStartingDiskNumber: offset,
		ZIPFileCommentLength:       0,
		ZIPFileComment:             nil,
	}
	return ff, nil
}

func (z *Zip) Marshal(file *os.File) (err error) {
	var ff FileFormat
	if ff, err = z.Build(); err != nil {
		return
	}
	return serial.Marshal(file, ff)
}
