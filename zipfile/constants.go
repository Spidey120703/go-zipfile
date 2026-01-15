package zipfile

var (
	LocalFileHeaderSignature       Signature = [4]byte{0x50, 0x4b, 0x03, 0x04}
	CentralFileHeaderSignature               = [4]byte{0x50, 0x4b, 0x01, 0x02}
	EndOfCentralDirectorySignature           = [4]byte{0x50, 0x4b, 0x05, 0x06}
	DigitalHeaderSignature                   = [4]byte{0x50, 0x4b, 0x05, 0x05}
)

const (
	VersionMadeByMS_DOS_and_OS_2 uint8 = iota
	VersionMadeByAmiga
	VersionMadeByOpenVMS
	VersionMadeByUNIX
	VersionMadeByVMC_MS
	VersionMadeByAtariST
	VersionMadeByOS_2_HPFS
	VersionMadeByMacintosh
	VersionMadeByZSystem
	VersionMadeByCP_M
	VersionMadeByWindowsNTFS
	VersionMadeByMVS
	VersionMadeByVSE
	VersionMadeByAcornRisc
	VersionMadeByVFAT
	VersionMadeByAlternativeMVS
	VersionMadeByBeOS
	VersionMadeByTandem
	VersionMadeByOS400
	VersionMadeByOSX_Darwin
)

var MapOfVersionMadeBy = map[uint8]string{
	VersionMadeByMS_DOS_and_OS_2: "MS-DOS and OS/2 (FAT / VFAT / FAT32 file systems)",
	VersionMadeByAmiga:           "Amiga",
	VersionMadeByOpenVMS:         "OpenVMS",
	VersionMadeByUNIX:            "UNIX",
	VersionMadeByVMC_MS:          "VM/CMS",
	VersionMadeByAtariST:         "Atari ST",
	VersionMadeByOS_2_HPFS:       "OS/2 H.P.F.S.",
	VersionMadeByMacintosh:       "Macintosh",
	VersionMadeByZSystem:         "Z-System",
	VersionMadeByCP_M:            "CP/M",
	VersionMadeByWindowsNTFS:     "Windows NTFS",
	VersionMadeByMVS:             "MVS (OS/390 - Z/OS)",
	VersionMadeByVSE:             "VSE",
	VersionMadeByAcornRisc:       "Acorn Risc",
	VersionMadeByVFAT:            "VFAT",
	VersionMadeByAlternativeMVS:  "alternate MVS",
	VersionMadeByBeOS:            "BeOS",
	VersionMadeByTandem:          "Tandem",
	VersionMadeByOS400:           "OS/400",
	VersionMadeByOSX_Darwin:      "OS X (Darwin)",
}

const DefaultVersion uint16 = 10
const LatestVersion uint16 = 63

var MinimumFeatureVersions = map[uint8][]string{
	10: {"1.0 - Default value"},
	11: {"1.1 - File is a volume label"},
	20: {
		"2.0 - File is a folder (directory)",
		"2.0 - File is compressed using Deflate compression",
		"2.0 - File is encrypted using traditional PKWARE encryption"},
	21: {"2.1 - File is compressed using Deflate64(tm)"},
	25: {"2.5 - File is compressed using PKWARE DCL Implode"},
	27: {"2.7 - File is a patch data set"},
	45: {"4.5 - File uses ZIP64 format extensions"},
	46: {"4.6 - File is compressed using BZIP2 compression*"},
	50: {"5.0 - File is encrypted using DES",
		"5.0 - File is encrypted using 3DES",
		"5.0 - File is encrypted using original RC2 encryption",
		"5.0 - File is encrypted using RC4 encryption"},
	51: {
		"5.1 - File is encrypted using AES encryption",
		"5.1 - File is encrypted using corrected RC2 encryption**"},
	52: {"5.2 - File is encrypted using corrected RC2-64 encryption**"},
	61: {"6.1 - File is encrypted using non-OAEP key wrapping***"},
	62: {"6.2 - Central directory encryption"},
	63: {"6.3 - File is compressed using LZMA",
		"6.3 - File is compressed using PPMd+",
		"6.3 - File is encrypted using Blowfish",
		"6.3 - File is encrypted using Twofish"},
}

const (
	EncryptedFlag uint16 = 1 << iota
	CompressionOption1
	CompressionOption2
	DataDescriptorFlag
	EnhancedDeflationFlag
	CompressedPatchedDataFlag
	StrongEncryptionFlag
	UnusedFlag1
	UnusedFlag2
	UnusedFlag3
	UnusedFlag4
	LanguageEncodingFlag
	ReservedFlag1
	MaskHeaderValuesFlag
	ReservedFlag2
	ReservedFlag3
)

const (
	CompressionMethodStored = iota
	CompressionMethodShrunk
	CompressionMethodReducedWithCompressionFactor1
	CompressionMethodReducedWithCompressionFactor2
	CompressionMethodReducedWithCompressionFactor3
	CompressionMethodReducedWithCompressionFactor4
	CompressionMethodImploded
	CompressionMethodTokenized
	CompressionMethodDeflated
	CompressionMethodDeflate64
	CompressionMethodPKWARE_DCL_Imploded
	CompressionMethodReserved1
	CompressionMethodBZIP2
	CompressionMethodReserved2
	CompressionMethodLZMA
	CompressionMethodReserved3
	CompressionMethodReserved34
	CompressionMethodIBM_zOS_CMPSC
	CompressionMethodIBM_TERSE
	CompressionMethodIBM_LZ77_z_Architecture
	CompressionMethodDeprecatedZSTD
	CompressionMethodZSTD = iota + 72
	CompressionMethodMP3
	CompressionMethodXZ
	CompressionMethodJPEG
	CompressionMethodWavPack
	CompressionMethodPPMd
	CompressionMethodAEx
)
