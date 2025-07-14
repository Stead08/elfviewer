package elf

import (
	"fmt"
)

const (
	EI_MAG0    = 0
	EI_MAG1    = 1
	EI_MAG2    = 2
	EI_MAG3    = 3
	EI_CLASS   = 4
	EI_DATA    = 5
	EI_VERSION = 6
	EI_OSABI   = 7
	EI_NIDENT  = 16

	ELFMAG0 = 0x7f
	ELFMAG1 = 'E'
	ELFMAG2 = 'L'
	ELFMAG3 = 'F'

	ELFCLASSNONE = 0
	ELFCLASS32   = 1
	ELFCLASS64   = 2

	ELFDATANONE = 0
	ELFDATA2LSB = 1
	ELFDATA2MSB = 2

	EV_NONE    = 0
	EV_CURRENT = 1
)

const (
	ET_NONE   = 0
	ET_REL    = 1
	ET_EXEC   = 2
	ET_DYN    = 3
	ET_CORE   = 4
	ET_LOOS   = 0xfe00
	ET_HIOS   = 0xfeff
	ET_LOPROC = 0xff00
	ET_HIPROC = 0xffff
)

const (
	EM_NONE    = 0
	EM_386     = 3
	EM_ARM     = 40
	EM_X86_64  = 62
	EM_AARCH64 = 183
)

const (
	SHT_NULL     = 0
	SHT_PROGBITS = 1
	SHT_SYMTAB   = 2
	SHT_STRTAB   = 3
	SHT_RELA     = 4
	SHT_HASH     = 5
	SHT_DYNAMIC  = 6
	SHT_NOTE     = 7
	SHT_NOBITS   = 8
	SHT_REL      = 9
	SHT_SHLIB    = 10
	SHT_DYNSYM   = 11
)

const (
	SHF_WRITE     = 0x1
	SHF_ALLOC     = 0x2
	SHF_EXECINSTR = 0x4
)

const (
	PT_NULL    = 0
	PT_LOAD    = 1
	PT_DYNAMIC = 2
	PT_INTERP  = 3
	PT_NOTE    = 4
	PT_SHLIB   = 5
	PT_PHDR    = 6
)

type Ident struct {
	Magic   [4]byte
	Class   uint8
	Data    uint8
	Version uint8
	OSABI   uint8
	Pad     [8]byte
}

type Header32 struct {
	Ident     Ident
	Type      uint16
	Machine   uint16
	Version   uint32
	Entry     uint32
	PhOff     uint32
	ShOff     uint32
	Flags     uint32
	EhSize    uint16
	PhEntSize uint16
	PhNum     uint16
	ShEntSize uint16
	ShNum     uint16
	ShStrNdx  uint16
}

type Header64 struct {
	Ident     Ident
	Type      uint16
	Machine   uint16
	Version   uint32
	Entry     uint64
	PhOff     uint64
	ShOff     uint64
	Flags     uint32
	EhSize    uint16
	PhEntSize uint16
	PhNum     uint16
	ShEntSize uint16
	ShNum     uint16
	ShStrNdx  uint16
}

type SectionHeader32 struct {
	Name      uint32
	Type      uint32
	Flags     uint32
	Addr      uint32
	Offset    uint32
	Size      uint32
	Link      uint32
	Info      uint32
	AddrAlign uint32
	EntSize   uint32
}

type SectionHeader64 struct {
	Name      uint32
	Type      uint32
	Flags     uint64
	Addr      uint64
	Offset    uint64
	Size      uint64
	Link      uint32
	Info      uint32
	AddrAlign uint64
	EntSize   uint64
}

type ProgramHeader32 struct {
	Type   uint32
	Offset uint32
	VAddr  uint32
	PAddr  uint32
	FileSz uint32
	MemSz  uint32
	Flags  uint32
	Align  uint32
}

type ProgramHeader64 struct {
	Type   uint32
	Flags  uint32
	Offset uint64
	VAddr  uint64
	PAddr  uint64
	FileSz uint64
	MemSz  uint64
	Align  uint64
}

type Symbol32 struct {
	Name  uint32
	Value uint32
	Size  uint32
	Info  uint8
	Other uint8
	Shndx uint16
}

type Symbol64 struct {
	Name  uint32
	Info  uint8
	Other uint8
	Shndx uint16
	Value uint64
	Size  uint64
}

type ProgramHeader struct {
	Type   uint32
	Flags  uint32
	Offset uint64
	VAddr  uint64
	PAddr  uint64
	FileSz uint64
	MemSz  uint64
	Align  uint64
}

type Symbol struct {
	Name  string
	Value uint64
	Size  uint64
	Info  uint8
	Other uint8
	Shndx uint16
}

func TypeString(t uint16) string {
	switch t {
	case ET_NONE:
		return "NONE (No file type)"
	case ET_REL:
		return "REL (Relocatable file)"
	case ET_EXEC:
		return "EXEC (Executable file)"
	case ET_DYN:
		return "DYN (Shared object file)"
	case ET_CORE:
		return "CORE (Core file)"
	default:
		return fmt.Sprintf("Unknown (%#x)", t)
	}
}

func MachineString(m uint16) string {
	switch m {
	case EM_NONE:
		return "None"
	case EM_386:
		return "Intel 80386"
	case EM_ARM:
		return "ARM"
	case EM_X86_64:
		return "AMD x86-64"
	case EM_AARCH64:
		return "ARM AARCH64"
	default:
		return fmt.Sprintf("Unknown (%d)", m)
	}
}

func SectionTypeString(t uint32) string {
	switch t {
	case SHT_NULL:
		return "NULL"
	case SHT_PROGBITS:
		return "PROGBITS"
	case SHT_SYMTAB:
		return "SYMTAB"
	case SHT_STRTAB:
		return "STRTAB"
	case SHT_RELA:
		return "RELA"
	case SHT_HASH:
		return "HASH"
	case SHT_DYNAMIC:
		return "DYNAMIC"
	case SHT_NOTE:
		return "NOTE"
	case SHT_NOBITS:
		return "NOBITS"
	case SHT_REL:
		return "REL"
	case SHT_SHLIB:
		return "SHLIB"
	case SHT_DYNSYM:
		return "DYNSYM"
	default:
		return fmt.Sprintf("Unknown (%#x)", t)
	}
}

func SegmentTypeString(t uint32) string {
	switch t {
	case PT_NULL:
		return "NULL"
	case PT_LOAD:
		return "LOAD"
	case PT_DYNAMIC:
		return "DYNAMIC"
	case PT_INTERP:
		return "INTERP"
	case PT_NOTE:
		return "NOTE"
	case PT_SHLIB:
		return "SHLIB"
	case PT_PHDR:
		return "PHDR"
	default:
		return fmt.Sprintf("Unknown (%#x)", t)
	}
}