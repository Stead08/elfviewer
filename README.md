# ELF Viewer - Web Application

A modern React-based web application for analyzing ELF (Executable and Linkable Format) files directly in your browser. Upload ELF files and explore their structure through an intuitive web interface powered by WebAssembly.

## Features

- **File Upload Interface**: Drag and drop ELF files directly into the browser
- **Interactive Header Display**: View ELF file headers with detailed information
- **Section Headers Table**: Sortable and filterable section headers with properties
- **Program Headers Viewer**: Examine program headers (segments) with visual representation
- **Symbol Table Explorer**: Browse symbol tables (both static and dynamic) with search functionality
- **Dynamic Section Information**: Display dynamic section data in structured format
- **Hex Dump Viewer**: Interactive hex dump with section selection
- **Cross-Platform**: Runs in any modern web browser, no installation required
- **Full ELF Support**: Both 32-bit and 64-bit ELF files, little-endian and big-endian formats

## Technology Stack

- **Frontend**: React 19.1.0 with TypeScript
- **Build Tool**: Vite for fast development and optimized builds
- **Package Manager**: pnpm for efficient dependency management
- **Testing**: Vitest with React Testing Library
- **Linting**: Biome for code quality and formatting
- **Backend**: Go compiled to WebAssembly for ELF parsing

## Getting Started

### Prerequisites

- Node.js (latest LTS version recommended)
- pnpm package manager

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Stead08/elfviewer
cd elfviewer
```

2. Install frontend dependencies:
```bash
cd front
pnpm install
```

3. Start the development server:
```bash
pnpm dev
```

4. Open your browser and navigate to `http://localhost:5173`

### Building for Production

```bash
cd front
pnpm build
```

The built application will be available in the `dist` directory.

## Usage

1. **Upload an ELF File**: Click the upload area or drag and drop an ELF file
2. **Explore Headers**: View the ELF header information in the first tab
3. **Browse Sections**: Use the Section Headers tab to examine all sections
4. **Check Segments**: View program headers in the Program Headers tab
5. **Analyze Symbols**: Explore symbol tables in the Symbols tab
6. **Hex Dump**: Select any section for a detailed hex dump view

## Architecture

The application uses a dual-architecture approach:

- **Web Frontend**: React application providing the user interface
- **WebAssembly Backend**: Go code compiled to WASM for efficient ELF parsing
- **Shared Types**: TypeScript interfaces that match Go structures for type safety

### Key Components

- `FileUpload.tsx`: Handles ELF file uploads and processing
- `ELFHeader.tsx`: Displays parsed ELF header information
- `SectionHeaders.tsx`: Interactive section headers table
- `ProgramHeaders.tsx`: Program segments visualization
- `Symbols.tsx`: Symbol table explorer with search
- `HexDump.tsx`: Hex dump viewer with section selection

## Development

### Running Tests

```bash
cd front
pnpm test          # Watch mode
pnpm test -- --run # Run once
pnpm test:ui       # Run with UI
```

### Linting

```bash
cd front
pnpm run lint
```

### WebAssembly Build

The Go backend is compiled to WebAssembly for browser execution:

```bash
GOOS=js GOARCH=wasm go build -o front/public/elfviewer.wasm
```

## CLI Alternative

While this repository focuses on the web application, a command-line interface is also available:

```bash
go build
./elfviewer [options] <elf-file>
```

See `CLAUDE.md` for detailed CLI usage instructions.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

MIT License