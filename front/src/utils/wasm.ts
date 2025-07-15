export interface ELFIdent {
	Magic: number[];
	Class: number;
	Data: number;
	Version: number;
	OSABI: number;
	Pad: number[];
}

export interface ELFInfo {
	ident: ELFIdent;
	class: number;
	type: number;
	machine: number;
	entry: number;
	sectionHeaders: SectionHeader[];
	programHeaders: ProgramHeader[];
	symbols: ELFSymbol[];
}

export interface SectionHeader {
	Name: string;
	Type: number;
	Flags: number;
	Addr: number;
	Offset: number;
	Size: number;
	Link: number;
	Info: number;
	AddrAlign: number;
	EntSize: number;
}

export interface ProgramHeader {
	Type: number;
	Flags: number;
	Offset: number;
	VAddr: number;
	PAddr: number;
	FileSz: number;
	MemSz: number;
	Align: number;
}

export interface ELFSymbol {
	Name: string;
	Value: number;
	Size: number;
	Info: number;
	Other: number;
	Shndx: number;
}

declare global {
	interface Window {
		Go: new () => {
			run: (instance: WebAssembly.Instance) => void;
			importObject: WebAssembly.Imports;
		};
		parseELF: (buffer: ArrayBuffer) => { data?: string; error?: string };
		getHexDump: (
			buffer: ArrayBuffer,
			sectionName: string,
		) => { data?: string; error?: string };
	}
}

let wasmReady = false;
let wasmReadyPromise: Promise<void> | null = null;

export async function initWasm(): Promise<void> {
	if (wasmReady) return;
	if (wasmReadyPromise) return wasmReadyPromise;

	wasmReadyPromise = (async () => {
		const go = new window.Go();
		const response = await fetch("/elfviewer.wasm");
		const buffer = await response.arrayBuffer();
		const result = await WebAssembly.instantiate(buffer, go.importObject);
		go.run(result.instance);
		wasmReady = true;
	})();

	return wasmReadyPromise;
}

export async function parseELF(buffer: ArrayBuffer): Promise<ELFInfo> {
	await initWasm();

	const result = window.parseELF(buffer);
	if (result.error) {
		throw new Error(result.error);
	}

	if (!result.data) {
		throw new Error("No data returned from parseELF");
	}
	return JSON.parse(result.data);
}

export async function getHexDump(
	buffer: ArrayBuffer,
	sectionName: string,
): Promise<string> {
	await initWasm();

	const result = window.getHexDump(buffer, sectionName);
	if (result.error) {
		throw new Error(result.error);
	}

	if (!result.data) {
		throw new Error("No data returned from getHexDump");
	}
	return result.data;
}

// Type helpers
export function getClassName(classValue: number): string {
	switch (classValue) {
		case 1:
			return "ELF32";
		case 2:
			return "ELF64";
		default:
			return `Unknown (${classValue})`;
	}
}

export function getDataEncoding(data: number): string {
	switch (data) {
		case 1:
			return "2's complement, little endian";
		case 2:
			return "2's complement, big endian";
		default:
			return `Invalid (${data})`;
	}
}

export function getTypeString(type: number): string {
	switch (type) {
		case 0:
			return "NONE (No file type)";
		case 1:
			return "REL (Relocatable file)";
		case 2:
			return "EXEC (Executable file)";
		case 3:
			return "DYN (Shared object file)";
		case 4:
			return "CORE (Core file)";
		default:
			return `Unknown (0x${type.toString(16)})`;
	}
}

export function getMachineString(machine: number): string {
	switch (machine) {
		case 0:
			return "None";
		case 3:
			return "Intel 80386";
		case 40:
			return "ARM";
		case 62:
			return "AMD x86-64";
		case 183:
			return "ARM AARCH64";
		default:
			return `Unknown (${machine})`;
	}
}

export function getSectionTypeString(type: number): string {
	switch (type) {
		case 0:
			return "NULL";
		case 1:
			return "PROGBITS";
		case 2:
			return "SYMTAB";
		case 3:
			return "STRTAB";
		case 4:
			return "RELA";
		case 5:
			return "HASH";
		case 6:
			return "DYNAMIC";
		case 7:
			return "NOTE";
		case 8:
			return "NOBITS";
		case 9:
			return "REL";
		case 10:
			return "SHLIB";
		case 11:
			return "DYNSYM";
		default:
			return `Unknown (0x${type.toString(16)})`;
	}
}

export function getSegmentTypeString(type: number): string {
	switch (type) {
		case 0:
			return "NULL";
		case 1:
			return "LOAD";
		case 2:
			return "DYNAMIC";
		case 3:
			return "INTERP";
		case 4:
			return "NOTE";
		case 5:
			return "SHLIB";
		case 6:
			return "PHDR";
		default:
			return `Unknown (0x${type.toString(16)})`;
	}
}

export function formatFlags(flags: number): string {
	let s = "";
	if (flags & 0x1) s += "W";
	if (flags & 0x2) s += "A";
	if (flags & 0x4) s += "X";
	return s || "  ";
}

export function formatSegmentFlags(flags: number): string {
	let s = "";
	if (flags & 0x4) s += "R";
	if (flags & 0x2) s += "W";
	if (flags & 0x1) s += "E";
	return s.padEnd(3);
}
