package elf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func Open(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return Parse(data)
}

func Parse(data []byte) (*File, error) {
	if len(data) < EI_NIDENT {
		return nil, fmt.Errorf("file too small to be ELF")
	}

	f := &File{Raw: data}

	if err := f.parseIdent(data[:EI_NIDENT]); err != nil {
		return nil, err
	}

	switch f.Class {
	case ELFCLASS32:
		if err := f.parse32(data); err != nil {
			return nil, err
		}
	case ELFCLASS64:
		if err := f.parse64(data); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown ELF class: %d", f.Class)
	}

	if err := f.parseSectionHeaders(); err != nil {
		return nil, fmt.Errorf("failed to parse section headers: %w", err)
	}

	if err := f.parseProgramHeaders(); err != nil {
		return nil, fmt.Errorf("failed to parse program headers: %w", err)
	}

	if err := f.parseSymbols(); err != nil {
		return nil, fmt.Errorf("failed to parse symbols: %w", err)
	}

	return f, nil
}

func (f *File) parseIdent(ident []byte) error {
	copy(f.Ident.Magic[:], ident[:4])

	if f.Ident.Magic[0] != ELFMAG0 ||
		f.Ident.Magic[1] != ELFMAG1 ||
		f.Ident.Magic[2] != ELFMAG2 ||
		f.Ident.Magic[3] != ELFMAG3 {
		return fmt.Errorf("invalid ELF magic: %x", f.Ident.Magic)
	}

	f.Ident.Class = ident[EI_CLASS]
	f.Ident.Data = ident[EI_DATA]
	f.Ident.Version = ident[EI_VERSION]
	f.Ident.OSABI = ident[EI_OSABI]
	f.Class = f.Ident.Class

	switch f.Ident.Data {
	case ELFDATA2LSB:
		f.ByteOrder = binary.LittleEndian
	case ELFDATA2MSB:
		f.ByteOrder = binary.BigEndian
	default:
		return fmt.Errorf("invalid data encoding: %d", f.Ident.Data)
	}

	return nil
}

func (f *File) parse32(data []byte) error {
	if len(data) < 52 {
		return fmt.Errorf("file too small for ELF32 header")
	}

	var h Header32
	r := bytes.NewReader(data)
	r.Seek(0, io.SeekStart)

	if err := binary.Read(r, f.ByteOrder, &h); err != nil {
		return fmt.Errorf("failed to read ELF32 header: %w", err)
	}

	f.Type = h.Type
	f.Machine = h.Machine
	f.Entry = uint64(h.Entry)

	f.phoff = uint64(h.PhOff)
	f.phentsize = uint16(h.PhEntSize)
	f.phnum = uint16(h.PhNum)

	f.shoff = uint64(h.ShOff)
	f.shentsize = uint16(h.ShEntSize)
	f.shnum = uint16(h.ShNum)
	f.shstrndx = uint16(h.ShStrNdx)

	return nil
}

func (f *File) parse64(data []byte) error {
	if len(data) < 64 {
		return fmt.Errorf("file too small for ELF64 header")
	}

	var h Header64
	r := bytes.NewReader(data)
	r.Seek(0, io.SeekStart)

	if err := binary.Read(r, f.ByteOrder, &h); err != nil {
		return fmt.Errorf("failed to read ELF64 header: %w", err)
	}

	f.Type = h.Type
	f.Machine = h.Machine
	f.Entry = h.Entry

	f.phoff = h.PhOff
	f.phentsize = h.PhEntSize
	f.phnum = h.PhNum

	f.shoff = h.ShOff
	f.shentsize = h.ShEntSize
	f.shnum = h.ShNum
	f.shstrndx = h.ShStrNdx

	return nil
}

func (f *File) parseSectionHeaders() error {
	if f.shoff == 0 || f.shnum == 0 {
		return nil
	}

	f.SectionHeaders = make([]SectionHeader, f.shnum)

	for i := 0; i < int(f.shnum); i++ {
		offset := f.shoff + uint64(i)*uint64(f.shentsize)
		if offset+uint64(f.shentsize) > uint64(len(f.Raw)) {
			return fmt.Errorf("section header %d out of bounds", i)
		}

		r := bytes.NewReader(f.Raw[offset:])

		if f.Class == ELFCLASS32 {
			var sh SectionHeader32
			if err := binary.Read(r, f.ByteOrder, &sh); err != nil {
				return err
			}
			f.SectionHeaders[i] = SectionHeader{
				Type:      sh.Type,
				Flags:     uint64(sh.Flags),
				Addr:      uint64(sh.Addr),
				Offset:    uint64(sh.Offset),
				Size:      uint64(sh.Size),
				Link:      sh.Link,
				Info:      sh.Info,
				AddrAlign: uint64(sh.AddrAlign),
				EntSize:   uint64(sh.EntSize),
			}
			f.SectionHeaders[i].nameIdx = sh.Name
		} else {
			var sh SectionHeader64
			if err := binary.Read(r, f.ByteOrder, &sh); err != nil {
				return err
			}
			f.SectionHeaders[i] = SectionHeader{
				Type:      sh.Type,
				Flags:     sh.Flags,
				Addr:      sh.Addr,
				Offset:    sh.Offset,
				Size:      sh.Size,
				Link:      sh.Link,
				Info:      sh.Info,
				AddrAlign: sh.AddrAlign,
				EntSize:   sh.EntSize,
			}
			f.SectionHeaders[i].nameIdx = sh.Name
		}
	}

	if f.shstrndx < f.shnum {
		strTab := f.SectionHeaders[f.shstrndx]
		if strTab.Offset+strTab.Size <= uint64(len(f.Raw)) {
			strData := f.Raw[strTab.Offset : strTab.Offset+strTab.Size]
			for i := range f.SectionHeaders {
				if name := getString(strData, f.SectionHeaders[i].nameIdx); name != "" {
					f.SectionHeaders[i].Name = name
				}
			}
		}
	}

	return nil
}

func (f *File) parseProgramHeaders() error {
	if f.phoff == 0 || f.phnum == 0 {
		return nil
	}

	f.ProgramHeaders = make([]ProgramHeader, f.phnum)

	for i := 0; i < int(f.phnum); i++ {
		offset := f.phoff + uint64(i)*uint64(f.phentsize)
		if offset+uint64(f.phentsize) > uint64(len(f.Raw)) {
			return fmt.Errorf("program header %d out of bounds", i)
		}

		r := bytes.NewReader(f.Raw[offset:])

		if f.Class == ELFCLASS32 {
			var ph ProgramHeader32
			if err := binary.Read(r, f.ByteOrder, &ph); err != nil {
				return err
			}
			f.ProgramHeaders[i] = ProgramHeader{
				Type:   ph.Type,
				Flags:  ph.Flags,
				Offset: uint64(ph.Offset),
				VAddr:  uint64(ph.VAddr),
				PAddr:  uint64(ph.PAddr),
				FileSz: uint64(ph.FileSz),
				MemSz:  uint64(ph.MemSz),
				Align:  uint64(ph.Align),
			}
		} else {
			var ph ProgramHeader64
			if err := binary.Read(r, f.ByteOrder, &ph); err != nil {
				return err
			}
			f.ProgramHeaders[i] = ProgramHeader{
				Type:   ph.Type,
				Flags:  ph.Flags,
				Offset: ph.Offset,
				VAddr:  ph.VAddr,
				PAddr:  ph.PAddr,
				FileSz: ph.FileSz,
				MemSz:  ph.MemSz,
				Align:  ph.Align,
			}
		}
	}

	return nil
}

func getString(data []byte, offset uint32) string {
	if offset >= uint32(len(data)) {
		return ""
	}
	for i := offset; i < uint32(len(data)); i++ {
		if data[i] == 0 {
			return string(data[offset:i])
		}
	}
	return string(data[offset:])
}

func (f *File) GetSection(name string) *SectionHeader {
	for i := range f.SectionHeaders {
		if f.SectionHeaders[i].Name == name {
			return &f.SectionHeaders[i]
		}
	}
	return nil
}

func (f *File) GetSectionData(sh *SectionHeader) ([]byte, error) {
	if sh.Type == SHT_NOBITS {
		return nil, nil
	}
	if sh.Offset+sh.Size > uint64(len(f.Raw)) {
		return nil, fmt.Errorf("section data out of bounds")
	}
	return f.Raw[sh.Offset : sh.Offset+sh.Size], nil
}

func (f *File) parseSymbols() error {
	for _, sh := range f.SectionHeaders {
		if sh.Type != SHT_SYMTAB && sh.Type != SHT_DYNSYM {
			continue
		}

		data, err := f.GetSectionData(&sh)
		if err != nil {
			return err
		}

		if sh.Link >= uint32(len(f.SectionHeaders)) {
			continue
		}

		strTab, err := f.GetSectionData(&f.SectionHeaders[sh.Link])
		if err != nil {
			return err
		}

		var symSize int
		if f.Class == ELFCLASS32 {
			symSize = 16
		} else {
			symSize = 24
		}

		numSyms := int(sh.Size) / symSize
		for i := 0; i < numSyms; i++ {
			offset := i * symSize
			if offset+symSize > len(data) {
				break
			}

			r := bytes.NewReader(data[offset:])
			var sym Symbol

			if f.Class == ELFCLASS32 {
				var s Symbol32
				if err := binary.Read(r, f.ByteOrder, &s); err != nil {
					continue
				}
				sym = Symbol{
					Value: uint64(s.Value),
					Size:  uint64(s.Size),
					Info:  s.Info,
					Other: s.Other,
					Shndx: s.Shndx,
				}
				if s.Name < uint32(len(strTab)) {
					sym.Name = getString(strTab, s.Name)
				}
			} else {
				var s Symbol64
				if err := binary.Read(r, f.ByteOrder, &s); err != nil {
					continue
				}
				sym = Symbol{
					Value: s.Value,
					Size:  s.Size,
					Info:  s.Info,
					Other: s.Other,
					Shndx: s.Shndx,
				}
				if s.Name < uint32(len(strTab)) {
					sym.Name = getString(strTab, s.Name)
				}
			}

			f.Symbols = append(f.Symbols, sym)
		}
	}

	return nil
}

type SectionHeader struct {
	Name      string
	Type      uint32
	Flags     uint64
	Addr      uint64
	Offset    uint64
	Size      uint64
	Link      uint32
	Info      uint32
	AddrAlign uint64
	EntSize   uint64
	nameIdx   uint32
}

type File struct {
	Ident          Ident
	ByteOrder      binary.ByteOrder
	Class          uint8
	Type           uint16
	Machine        uint16
	Entry          uint64
	ProgramHeaders []ProgramHeader
	SectionHeaders []SectionHeader
	Symbols        []Symbol
	StringTable    []byte
	Raw            []byte
	
	phoff     uint64
	phentsize uint16
	phnum     uint16
	shoff     uint64
	shentsize uint16
	shnum     uint16
	shstrndx  uint16
}