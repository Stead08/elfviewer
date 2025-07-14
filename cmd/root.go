package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/elfviewer/elfviewer/elf"
)

var (
	showHeader   bool
	showSections bool
	showSegments bool
	showSymbols  bool
	showDynamic  bool
	showAll      bool
	hexDump      string
	help         bool
)

func init() {
	flag.BoolVar(&showHeader, "h", true, "Show ELF header")
	flag.BoolVar(&showHeader, "header", true, "Show ELF header")
	flag.BoolVar(&showSections, "S", false, "Show section headers")
	flag.BoolVar(&showSections, "sections", false, "Show section headers")
	flag.BoolVar(&showSegments, "l", false, "Show program headers")
	flag.BoolVar(&showSegments, "segments", false, "Show program headers")
	flag.BoolVar(&showSymbols, "s", false, "Show symbol table")
	flag.BoolVar(&showSymbols, "symbols", false, "Show symbol table")
	flag.BoolVar(&showDynamic, "d", false, "Show dynamic section")
	flag.BoolVar(&showDynamic, "dynamic", false, "Show dynamic section")
	flag.BoolVar(&showAll, "a", false, "Show all information")
	flag.BoolVar(&showAll, "all", false, "Show all information")
	flag.StringVar(&hexDump, "x", "", "Dump section in hex")
	flag.StringVar(&hexDump, "hex", "", "Dump section in hex")
	flag.BoolVar(&help, "help", false, "Show help message")
}

func Execute() error {
	flag.Parse()

	if help || flag.NArg() == 0 {
		printUsage()
		if help {
			return nil
		}
		return fmt.Errorf("no ELF file specified")
	}

	filename := flag.Arg(0)
	
	file, err := elf.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", filename, err)
	}

	if showAll {
		showHeader = true
		showSections = true
		showSegments = true
		showSymbols = true
		showDynamic = true
	}

	if showHeader && hexDump == "" {
		file.DisplayHeader(os.Stdout)
		fmt.Println()
	}

	if showSections {
		file.DisplaySectionHeaders(os.Stdout)
		fmt.Println()
	}

	if showSegments {
		file.DisplayProgramHeaders(os.Stdout)
		fmt.Println()
	}

	if showSymbols {
		file.DisplaySymbols(os.Stdout)
		fmt.Println()
	}

	if showDynamic {
		file.DisplayDynamic(os.Stdout)
		fmt.Println()
	}

	if hexDump != "" {
		if err := file.DisplayHexDump(os.Stdout, hexDump); err != nil {
			return err
		}
		fmt.Println()
	}

	return nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: elfviewer [options] <elf-file>\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -h, --header      Show ELF header (default)\n")
	fmt.Fprintf(os.Stderr, "  -S, --sections    Show section headers\n")
	fmt.Fprintf(os.Stderr, "  -l, --segments    Show program headers\n")
	fmt.Fprintf(os.Stderr, "  -s, --symbols     Show symbol table\n")
	fmt.Fprintf(os.Stderr, "  -d, --dynamic     Show dynamic section\n")
	fmt.Fprintf(os.Stderr, "  -a, --all         Show all information\n")
	fmt.Fprintf(os.Stderr, "  -x, --hex <section>  Dump section in hex\n")
	fmt.Fprintf(os.Stderr, "  --help           Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  elfviewer /bin/ls               # Show ELF header\n")
	fmt.Fprintf(os.Stderr, "  elfviewer -a /bin/ls            # Show all information\n")
	fmt.Fprintf(os.Stderr, "  elfviewer -S /bin/ls            # Show section headers\n")
	fmt.Fprintf(os.Stderr, "  elfviewer -x .text /bin/ls      # Hex dump of .text section\n")
}