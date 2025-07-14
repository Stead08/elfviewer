# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

**Build the application:**

```bash
go build
```

**Run with a sample ELF file:**

```bash
./elfviewer <elf-file>
# Example: ./elfviewer /bin/ls
```

**Run directly without building:**

```bash
go run main.go <elf-file>
```

**Format code:**

```bash
go fmt ./...
```

**Run static analysis:**

```bash
go vet ./...
```

## Architecture

The ELFViewer is a command-line tool for analyzing ELF (Executable and Linkable Format) files. The codebase follows a clean package structure:

### Core Packages

**`main.go`**: Entry point that delegates to the cmd package.

**`cmd/` package**: Contains command-line interface logic using Go's standard `flag` package. The `root.go` file:

- Defines all command-line flags (-h, -S, -l, -s, -d, -a, -x)
- Orchestrates the parsing and display flow
- Calls appropriate display methods based on flags

**`elf/` package**: Core ELF parsing and display functionality:

- `types.go`: Defines ELF constants, structures, and type definitions (headers, sections, segments, symbols)
- `parser.go`: Contains the parsing logic for reading and interpreting ELF binary data
  - `Open()`: Opens and reads an ELF file from disk
  - `Parse()`: Main parsing logic that handles both 32-bit and 64-bit ELF formats
  - Separate parsing methods for headers, sections, and program segments
- `display.go`: Implements display methods for outputting parsed ELF data
  - `DisplayHeader()`: Shows basic ELF header information
  - `DisplaySectionHeaders()`: Lists all sections
  - `DisplayProgramHeaders()`: Shows loadable segments
  - `DisplaySymbols()`: Displays symbol table entries
  - `DisplayDynamic()`: Shows dynamic linking information
  - `DisplayHexDump()`: Provides hex dump of specific sections

### Key Design Patterns

1. **Separation of Concerns**: Parsing logic is separated from display logic and CLI handling
2. **Dual Architecture Support**: Code handles both 32-bit and 64-bit ELF files transparently
3. **Endianness Handling**: Supports both little-endian and big-endian formats
4. **Error Propagation**: Uses Go's error wrapping for better error context

The tool reads the entire ELF file into memory, parses its structure based on the ELF specification, and provides various display options for different ELF components.

## イシュー

イシューの要望は原則/frontでフロントアプリで実装して。

## lint

frontアプリのlintはpackage.jsonに定義されているlintコマンドを利用して
