# ELF Viewer

A command-line tool for displaying and analyzing ELF (Executable and Linkable Format) files on Unix-like systems.

## Features

- Display ELF file headers with detailed information
- View section headers and their properties
- Examine program headers (segments)
- List symbol tables (both static and dynamic)
- Display dynamic section information
- Hex dump of specific sections
- Support for both 32-bit and 64-bit ELF files
- Support for little-endian and big-endian formats

## Installation

```bash
go get github.com/elfviewer/elfviewer
```

Or build from source:

```bash
git clone https://github.com/elfviewer/elfviewer
cd elfviewer
go build
```

## Usage

```
elfviewer [options] <elf-file>

Options:
  -h, --header      Show ELF header (default)
  -S, --sections    Show section headers
  -l, --segments    Show program headers
  -s, --symbols     Show symbol table
  -d, --dynamic     Show dynamic section
  -a, --all         Show all information
  -x, --hex <section>  Dump section in hex
  --help           Show help message
```

## Examples

Display basic ELF header information:
```bash
elfviewer /bin/ls
```

Show all information about an ELF file:
```bash
elfviewer -a /usr/bin/gcc
```

Display section headers:
```bash
elfviewer -S /lib/libc.so.6
```

Show program headers (segments):
```bash
elfviewer -l /bin/bash
```

Display symbol table:
```bash
elfviewer -s /usr/lib/libm.so
```

Hex dump of .text section:
```bash
elfviewer -x .text /bin/echo
```

## Output Format

### ELF Header
Shows basic file information including:
- Magic bytes (7f 45 4c 46)
- Class (32-bit or 64-bit)
- Data encoding (little-endian or big-endian)
- File type (executable, shared object, etc.)
- Target architecture
- Entry point address

### Section Headers
Lists all sections with:
- Section name
- Type (PROGBITS, SYMTAB, STRTAB, etc.)
- Virtual address
- File offset
- Size
- Flags (Write, Alloc, Execute)

### Program Headers
Shows loadable segments with:
- Segment type (LOAD, DYNAMIC, INTERP, etc.)
- File offset and size
- Virtual and physical addresses
- Memory size
- Permissions (Read, Write, Execute)

### Symbol Table
Displays symbols with:
- Symbol value (address)
- Size
- Type (FUNC, OBJECT, etc.)
- Binding (LOCAL, GLOBAL, WEAK)
- Section index
- Symbol name

## Supported Architectures

- x86 (i386)
- x86-64 (amd64)
- ARM
- ARM64 (aarch64)
- And many others (displays numeric ID for unsupported architectures)

## Requirements

- Go 1.16 or later

## License

MIT License