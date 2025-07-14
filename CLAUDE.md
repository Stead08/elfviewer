# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Go Backend (CLI Tool)

**Build the CLI application:**
```bash
go build
```

**Build for WebAssembly:**
```bash
GOOS=js GOARCH=wasm go build -o front/public/elfviewer.wasm
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

### Frontend (React Application)

**Install dependencies:**
```bash
cd front
pnpm install
```

**Development server:**
```bash
pnpm dev
```

**Build for production:**
```bash
pnpm build
```

**Run linter (Biome):**
```bash
pnpm run lint
```

**Run tests:**
```bash
pnpm test          # Watch mode
pnpm test -- --run # Run once
```

**Run tests with UI:**
```bash
pnpm test:ui
```

## Architecture

The ELFViewer is a dual-mode application that can run as both a command-line tool and a web application through WebAssembly.

### Go Backend Structure

**`main.go`**: Entry point with build tags:
- Default build (`!wasm`): CLI version that delegates to cmd package
- WASM build: Excluded from compilation

**`wasm_exports.go`**: WebAssembly exports (build tag: `wasm`):
- `parseELF(buffer)`: Parses ELF file from ArrayBuffer and returns JSON
- `getHexDump(buffer, sectionName)`: Returns hex dump of specific section

**`cmd/` package**: Command-line interface logic:
- `root.go`: Defines CLI flags (-h, -S, -l, -s, -d, -a, -x) and orchestrates display flow
- Uses Go's standard `flag` package

**`elf/` package**: Core ELF parsing and display functionality:
- `types.go`: ELF constants, structures, and type definitions
- `parser.go`: Parsing logic for 32-bit and 64-bit ELF formats
  - `Open()`: Opens and reads ELF file from disk
  - `Parse()`: Main parsing logic with endianness handling
- `display.go`: Display methods for CLI output
  - `DisplayHeader()`, `DisplaySectionHeaders()`, `DisplayProgramHeaders()`
  - `DisplaySymbols()`, `DisplayDynamic()`, `DisplayHexDump()`

### Frontend Structure (React + TypeScript)

**Technology Stack:**
- React 19.1.0 with TypeScript
- Vite as build tool
- Biome for linting (replacing ESLint)
- Vitest for testing
- pnpm as package manager

**Key Components:**
- `App.tsx`: Main component with tab navigation
- `FileUpload.tsx`: Handles ELF file uploads
- `ELFHeader.tsx`: Displays parsed ELF header
- `SectionHeaders.tsx`: Section headers with sorting/filtering
- `ProgramHeaders.tsx`: Program segments display
- `Symbols.tsx`: Symbol table viewer
- `HexDump.tsx`: Hex dump with section selector

**WebAssembly Integration:**
- `utils/wasm.ts`: Initializes WASM module and wraps exported functions
- `public/wasm_exec.js`: Go's WebAssembly support runtime
- `public/elfviewer.wasm`: Compiled WebAssembly module

### Key Design Patterns

1. **Dual Build System**: Same Go code serves both CLI and web interfaces
2. **Separation of Concerns**: Parsing logic separated from display logic
3. **Type Safety**: TypeScript interfaces match Go structures
4. **Error Handling**: Go errors wrapped with context, displayed in UI

### Data Flow (Web Version)

1. User uploads ELF file → `FileUpload` component
2. File converted to ArrayBuffer → Passed to WASM `parseELF`
3. Go parses binary data → Returns structured JSON
4. React components render parsed data with interactive features

## Testing

**Frontend Tests:**
- Test framework: Vitest with React Testing Library
- Test files: `*.test.tsx` (e.g., `SectionHeaders.test.tsx`)
- Configuration: `vitest.config.ts`
- Setup file: `src/test/setup.ts`

## CI/CD

**GitHub Actions Workflows:**
- `test-frontend.yml`: Runs on PRs affecting `front/`
  - Uses pnpm for dependency management
  - Runs lint, tests, and build

## Important Notes

- イシューの要望は原則/frontでフロントアプリで実装して
- frontアプリのlintはpackage.jsonに定義されているlintコマンドを利用して
- When building WASM, output must go to `front/public/elfviewer.wasm`
- The frontend expects the WASM file to be available at runtime