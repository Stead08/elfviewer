import type React from "react";
import { useState, useMemo } from "react";
import {
	formatFlags,
	getSectionTypeString,
	type SectionHeader,
} from "../utils/wasm";

interface SectionHeadersProps {
	sections: SectionHeader[];
}

export const SectionHeaders: React.FC<SectionHeadersProps> = ({ sections }) => {
	const [sortByAddress, setSortByAddress] = useState(false);
	const [showFlaggedOnly, setShowFlaggedOnly] = useState(false);

	const processedSections = useMemo(() => {
		let filtered = sections;

		// Apply flag filtering
		if (showFlaggedOnly) {
			filtered = filtered.filter(section => formatFlags(section.Flags).trim() !== "");
		}

		// Apply address sorting
		if (sortByAddress) {
			filtered = [...filtered].sort((a, b) => a.Addr - b.Addr);
		}

		return filtered;
	}, [sections, sortByAddress, showFlaggedOnly]);

	return (
		<div className="section-headers">
			<h2>Section Headers</h2>

			{/* Controls */}
			<div style={{ marginBottom: "1rem", display: "flex", gap: "1rem", flexWrap: "wrap" }}>
				<label style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
					<input
						type="checkbox"
						checked={sortByAddress}
						onChange={(e) => setSortByAddress(e.target.checked)}
					/>
					Sort by Address
				</label>
				<label style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
					<input
						type="checkbox"
						checked={showFlaggedOnly}
						onChange={(e) => setShowFlaggedOnly(e.target.checked)}
					/>
					Show only sections with flags
				</label>
			</div>

			<div className="table-container">
				<table>
					<thead>
						<tr>
							<th>[Nr]</th>
							<th>Name</th>
							<th>Type</th>
							<th>Address</th>
							<th>Offset</th>
							<th>Size</th>
							<th>EntSize</th>
							<th>Flags</th>
							<th>Link</th>
							<th>Info</th>
							<th>Align</th>
						</tr>
					</thead>
					<tbody>
						{processedSections.map((section, index) => (
							<tr key={`section-${section.Name || section.Addr}-${index}`}>
								<td>[{section.originalIndex}]</td>
								<td>{section.Name || "<no-name>"}</td>
								<td>{getSectionTypeString(section.Type)}</td>
								<td className="mono">
									0x{section.Addr.toString(16).padStart(16, "0")}
								</td>
								<td className="mono">
									0x{section.Offset.toString(16).padStart(8, "0")}
								</td>
								<td className="mono">
									0x{section.Size.toString(16).padStart(16, "0")}
								</td>
								<td className="mono">
									0x{section.EntSize.toString(16).padStart(16, "0")}
								</td>
								<td>{formatFlags(section.Flags)}</td>
								<td>{section.Link}</td>
								<td>{section.Info}</td>
								<td>{section.AddrAlign}</td>
							</tr>
						))}
					</tbody>
				</table>
			</div>

			{/* Results info */}
			<div style={{ marginTop: "1rem", fontSize: "0.875rem", color: "#666" }}>
				<div>
					Showing {processedSections.length} of {sections.length} sections
				</div>
				<div style={{ marginTop: "0.5rem" }}>
					<strong>Key to Flags:</strong>
					<br />W (write), A (alloc), X (execute)
				</div>
			</div>
		</div>
	);
};
