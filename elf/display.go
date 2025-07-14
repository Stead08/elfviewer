package elf

import (
	"fmt"
	"io"
	"text/tabwriter"
)

func (f *File) DisplayHeader(w io.Writer) {
	fmt.Fprintf(w, "ELF Header:\n")
	fmt.Fprintf(w, "  Magic:   ")
	for _, b := range f.Ident.Magic {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "\n")
	
	fmt.Fprintf(w, "  Class:                             ")
	switch f.Ident.Class {
	case ELFCLASS32:
		fmt.Fprintf(w, "ELF32\n")
	case ELFCLASS64:
		fmt.Fprintf(w, "ELF64\n")
	default:
		fmt.Fprintf(w, "Invalid (%d)\n", f.Ident.Class)
	}
	
	fmt.Fprintf(w, "  Data:                              ")
	switch f.Ident.Data {
	case ELFDATA2LSB:
		fmt.Fprintf(w, "2's complement, little endian\n")
	case ELFDATA2MSB:
		fmt.Fprintf(w, "2's complement, big endian\n")
	default:
		fmt.Fprintf(w, "Invalid (%d)\n", f.Ident.Data)
	}
	
	fmt.Fprintf(w, "  Version:                           %d (current)\n", f.Ident.Version)
	fmt.Fprintf(w, "  OS/ABI:                            %d\n", f.Ident.OSABI)
	fmt.Fprintf(w, "  Type:                              %s\n", TypeString(f.Type))
	fmt.Fprintf(w, "  Machine:                           %s\n", MachineString(f.Machine))
	fmt.Fprintf(w, "  Entry point address:               0x%x\n", f.Entry)
	fmt.Fprintf(w, "  Start of program headers:          %d (bytes into file)\n", f.phoff)
	fmt.Fprintf(w, "  Start of section headers:          %d (bytes into file)\n", f.shoff)
	fmt.Fprintf(w, "  Size of program headers:           %d (bytes)\n", f.phentsize)
	fmt.Fprintf(w, "  Number of program headers:         %d\n", f.phnum)
	fmt.Fprintf(w, "  Size of section headers:           %d (bytes)\n", f.shentsize)
	fmt.Fprintf(w, "  Number of section headers:         %d\n", f.shnum)
	fmt.Fprintf(w, "  Section header string table index: %d\n", f.shstrndx)
}

func (f *File) DisplaySectionHeaders(w io.Writer) {
	fmt.Fprintf(w, "Section Headers:\n")
	fmt.Fprintf(w, "  [Nr] Name              Type            Address          Offset\n")
	fmt.Fprintf(w, "       Size              EntSize         Flags  Link  Info  Align\n")
	
	for i, sh := range f.SectionHeaders {
		fmt.Fprintf(w, "  [%2d] %-16s %-15s %016x %08x\n",
			i, sh.Name, SectionTypeString(sh.Type), sh.Addr, sh.Offset)
		fmt.Fprintf(w, "       %016x %016x %3s %5d %5d %5d\n",
			sh.Size, sh.EntSize, formatFlags(sh.Flags), sh.Link, sh.Info, sh.AddrAlign)
	}
	
	fmt.Fprintf(w, "\nKey to Flags:\n")
	fmt.Fprintf(w, "  W (write), A (alloc), X (execute), M (merge), S (strings), I (info),\n")
	fmt.Fprintf(w, "  L (link order), O (extra OS processing required), G (group), T (TLS),\n")
	fmt.Fprintf(w, "  C (compressed), x (unknown), o (OS specific), E (exclude),\n")
	fmt.Fprintf(w, "  D (mbind), l (large), p (processor specific)\n")
}

func (f *File) DisplayProgramHeaders(w io.Writer) {
	fmt.Fprintf(w, "\nProgram Headers:\n")
	fmt.Fprintf(w, "  Type           Offset             VirtAddr           PhysAddr\n")
	fmt.Fprintf(w, "                 FileSiz            MemSiz              Flags  Align\n")
	
	for _, ph := range f.ProgramHeaders {
		fmt.Fprintf(w, "  %-14s 0x%016x 0x%016x 0x%016x\n",
			SegmentTypeString(ph.Type), ph.Offset, ph.VAddr, ph.PAddr)
		fmt.Fprintf(w, "                 0x%016x 0x%016x  %s 0x%x\n",
			ph.FileSz, ph.MemSz, formatSegmentFlags(ph.Flags), ph.Align)
	}
	
	fmt.Fprintf(w, "\n Section to Segment mapping:\n")
	fmt.Fprintf(w, "  Segment Sections...\n")
	
	for i, ph := range f.ProgramHeaders {
		fmt.Fprintf(w, "   %02d     ", i)
		for _, sh := range f.SectionHeaders {
			if sh.Offset >= ph.Offset && sh.Offset < ph.Offset+ph.FileSz {
				fmt.Fprintf(w, "%s ", sh.Name)
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

func (f *File) DisplaySymbols(w io.Writer) {
	if len(f.Symbols) == 0 {
		fmt.Fprintf(w, "\nNo symbols found.\n")
		return
	}
	
	fmt.Fprintf(w, "\nSymbol table contains %d entries:\n", len(f.Symbols))
	
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "   Num:\tValue\tSize\tType\tBind\tVis\tNdx\tName\n")
	
	for i, sym := range f.Symbols {
		symType := sym.Info & 0xf
		symBind := sym.Info >> 4
		symVis := sym.Other & 0x3
		
		fmt.Fprintf(tw, "%6d:\t%016x\t%d\t%s\t%s\t%s\t%s\t%s\n",
			i,
			sym.Value,
			sym.Size,
			symbolTypeString(symType),
			symbolBindString(symBind),
			symbolVisString(symVis),
			sectionIndexString(sym.Shndx),
			sym.Name)
	}
	tw.Flush()
}

func (f *File) DisplayDynamic(w io.Writer) {
	dynSec := f.GetSection(".dynamic")
	if dynSec == nil {
		fmt.Fprintf(w, "\nNo dynamic section found.\n")
		return
	}
	
	fmt.Fprintf(w, "\nDynamic section at offset 0x%x contains entries:\n", dynSec.Offset)
}

func formatFlags(flags uint64) string {
	s := ""
	if flags&SHF_WRITE != 0 {
		s += "W"
	}
	if flags&SHF_ALLOC != 0 {
		s += "A"
	}
	if flags&SHF_EXECINSTR != 0 {
		s += "X"
	}
	if s == "" {
		s = "  "
	}
	return fmt.Sprintf("%3s", s)
}

func formatSegmentFlags(flags uint32) string {
	s := ""
	if flags&0x4 != 0 {
		s += "R"
	}
	if flags&0x2 != 0 {
		s += "W"
	}
	if flags&0x1 != 0 {
		s += "E"
	}
	return fmt.Sprintf("%-3s", s)
}

func symbolTypeString(t uint8) string {
	switch t {
	case 0:
		return "NOTYPE"
	case 1:
		return "OBJECT"
	case 2:
		return "FUNC"
	case 3:
		return "SECTION"
	case 4:
		return "FILE"
	case 5:
		return "COMMON"
	case 6:
		return "TLS"
	default:
		return fmt.Sprintf("<%d>", t)
	}
}

func symbolBindString(b uint8) string {
	switch b {
	case 0:
		return "LOCAL"
	case 1:
		return "GLOBAL"
	case 2:
		return "WEAK"
	default:
		return fmt.Sprintf("<%d>", b)
	}
}

func symbolVisString(v uint8) string {
	switch v {
	case 0:
		return "DEFAULT"
	case 1:
		return "INTERNAL"
	case 2:
		return "HIDDEN"
	case 3:
		return "PROTECTED"
	default:
		return fmt.Sprintf("<%d>", v)
	}
}

func sectionIndexString(ndx uint16) string {
	switch ndx {
	case 0:
		return "UND"
	case 0xfff1:
		return "ABS"
	case 0xfff2:
		return "COMMON"
	default:
		return fmt.Sprintf("%d", ndx)
	}
}

func (f *File) DisplayHexDump(w io.Writer, sectionName string) error {
	sh := f.GetSection(sectionName)
	if sh == nil {
		return fmt.Errorf("section %s not found", sectionName)
	}
	
	data, err := f.GetSectionData(sh)
	if err != nil {
		return err
	}
	
	fmt.Fprintf(w, "\nHex dump of section '%s':\n", sectionName)
	
	for i := 0; i < len(data); i += 16 {
		fmt.Fprintf(w, "  0x%08x ", sh.Addr+uint64(i))
		
		for j := 0; j < 16; j++ {
			if i+j < len(data) {
				fmt.Fprintf(w, "%02x", data[i+j])
			} else {
				fmt.Fprintf(w, "  ")
			}
			if j%4 == 3 {
				fmt.Fprintf(w, " ")
			}
		}
		
		fmt.Fprintf(w, " ")
		for j := 0; j < 16 && i+j < len(data); j++ {
			c := data[i+j]
			if c >= 32 && c < 127 {
				fmt.Fprintf(w, "%c", c)
			} else {
				fmt.Fprintf(w, ".")
			}
		}
		fmt.Fprintf(w, "\n")
	}
	
	return nil
}