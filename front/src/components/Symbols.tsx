import type React from "react";
import { useState, useMemo } from "react";
import type { Symbol as ELFSymbol, SectionHeader } from "../utils/wasm";

interface SymbolsProps {
	symbols: ELFSymbol[];
	sectionHeaders: SectionHeader[];
}

function getSymbolType(info: number): string {
	const type = info & 0xf;
	switch (type) {
		case 0:
			return "NOTYPE";
		case 1:
			return "OBJECT";
		case 2:
			return "FUNC";
		case 3:
			return "SECTION";
		case 4:
			return "FILE";
		case 5:
			return "COMMON";
		case 6:
			return "TLS";
		default:
			return `<${type}>`;
	}
}

function getSymbolBind(info: number): string {
	const bind = info >> 4;
	switch (bind) {
		case 0:
			return "LOCAL";
		case 1:
			return "GLOBAL";
		case 2:
			return "WEAK";
		default:
			return `<${bind}>`;
	}
}

function getSymbolVisibility(other: number): string {
	const vis = other & 0x3;
	switch (vis) {
		case 0:
			return "DEFAULT";
		case 1:
			return "INTERNAL";
		case 2:
			return "HIDDEN";
		case 3:
			return "PROTECTED";
		default:
			return `<${vis}>`;
	}
}

function getSectionIndex(shndx: number, sectionHeaders: SectionHeader[]): string {
	switch (shndx) {
		case 0:
			return "UND";
		case 0xfff1:
			return "ABS";
		case 0xfff2:
			return "COMMON";
		default:
			// If shndx is a valid section index, return the section name
			if (shndx > 0 && shndx < sectionHeaders.length) {
				return sectionHeaders[shndx].Name || shndx.toString();
			}
			return shndx.toString();
	}
}

export const Symbols: React.FC<SymbolsProps> = ({ symbols, sectionHeaders }) => {
	const [hideNotype, setHideNotype] = useState(false);
	const [sortBy, setSortBy] = useState<"num" | "name" | "size">("num");
	const [sortOrder, setSortOrder] = useState<"asc" | "desc">("asc");

	const filteredAndSortedSymbols = useMemo(() => {
		let filtered = symbols;
		
		// Filter out NOTYPE symbols if checkbox is checked
		if (hideNotype) {
			filtered = symbols.filter(symbol => getSymbolType(symbol.Info) !== "NOTYPE");
		}
		
		// Sort symbols
		const sorted = [...filtered].sort((a, b) => {
			let compareValue = 0;
			
			switch (sortBy) {
				case "num":
					compareValue = symbols.indexOf(a) - symbols.indexOf(b);
					break;
				case "name":
					compareValue = (a.Name || "").localeCompare(b.Name || "");
					break;
				case "size":
					compareValue = a.Size - b.Size;
					break;
			}
			
			return sortOrder === "asc" ? compareValue : -compareValue;
		});
		
		return sorted;
	}, [symbols, hideNotype, sortBy, sortOrder]);

	if (symbols.length === 0) {
		return (
			<div className="symbols">
				<h2>Symbols</h2>
				<p>No symbols found.</p>
			</div>
		);
	}

	const handleSortChange = (newSortBy: "num" | "name" | "size") => {
		if (sortBy === newSortBy) {
			setSortOrder(sortOrder === "asc" ? "desc" : "asc");
		} else {
			setSortBy(newSortBy);
			setSortOrder("asc");
		}
	};

	return (
		<div className="symbols">
			<h2>Symbols</h2>
			
			<div style={{ marginBottom: "1rem", display: "flex", gap: "1rem", alignItems: "center" }}>
				<label style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
					<input
						type="checkbox"
						checked={hideNotype}
						onChange={(e) => setHideNotype(e.target.checked)}
					/>
					Hide NOTYPE symbols
				</label>
				
				<div style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
					<span>Sort by:</span>
					<select
						value={sortBy}
						onChange={(e) => setSortBy(e.target.value as "num" | "name" | "size")}
					>
						<option value="num">Num</option>
						<option value="name">Name</option>
						<option value="size">Size</option>
					</select>
					<button
						type="button"
						onClick={() => setSortOrder(sortOrder === "asc" ? "desc" : "asc")}
						style={{ 
							background: "none", 
							border: "1px solid var(--border-color)", 
							borderRadius: "4px",
							padding: "0.25rem 0.5rem",
							cursor: "pointer",
							color: "var(--text-secondary)"
						}}
					>
						{sortOrder === "asc" ? "↑" : "↓"}
					</button>
				</div>
			</div>

			<div className="table-container">
				<table>
					<thead>
						<tr>
							<th 
								style={{ cursor: "pointer" }}
								onClick={() => handleSortChange("num")}
							>
								Num {sortBy === "num" && (sortOrder === "asc" ? "↑" : "↓")}
							</th>
							<th>Value</th>
							<th 
								style={{ cursor: "pointer" }}
								onClick={() => handleSortChange("size")}
							>
								Size {sortBy === "size" && (sortOrder === "asc" ? "↑" : "↓")}
							</th>
							<th>Type</th>
							<th 
								style={{ cursor: "pointer" }}
								onClick={() => handleSortChange("name")}
							>
								Name {sortBy === "name" && (sortOrder === "asc" ? "↑" : "↓")}
							</th>
							<th>Ndx</th>
							<th>Vis</th>
							<th>Bind</th>
						</tr>
					</thead>
					<tbody>
						{filteredAndSortedSymbols.map((symbol, index) => (
							<tr key={`symbol-${symbols.indexOf(symbol)}`}>
								<td>{symbols.indexOf(symbol)}:</td>
								<td className="mono">
									0x{symbol.Value.toString(16).padStart(16, "0")}
								</td>
								<td>{symbol.Size}</td>
								<td>{getSymbolType(symbol.Info)}</td>
								<td>{symbol.Name || "<no-name>"}</td>
								<td>{getSectionIndex(symbol.Shndx, sectionHeaders)}</td>
								<td>{getSymbolVisibility(symbol.Other)}</td>
								<td>{getSymbolBind(symbol.Info)}</td>
							</tr>
						))}
					</tbody>
				</table>
			</div>
		</div>
	);
};
