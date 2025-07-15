import { useCallback, useEffect, useState } from "react";
import "./App.css";
import { ELFHeader } from "./components/ELFHeader";
import { FileUpload } from "./components/FileUpload";
import { HexDump } from "./components/HexDump";
import { ProgramHeaders } from "./components/ProgramHeaders";
import { SectionHeaders } from "./components/SectionHeaders";
import { Symbols } from "./components/Symbols";
import { type ELFInfo, initWasm, parseELF } from "./utils/wasm";

function App() {
	const [elfData, setElfData] = useState<ELFInfo | null>(null);
	const [fileBuffer, setFileBuffer] = useState<ArrayBuffer | null>(null);
	const [error, setError] = useState<string | null>(null);
	const [activeTab, setActiveTab] = useState<
		"header" | "sections" | "segments" | "symbols" | "hex"
	>("header");
	const [wasmLoading, setWasmLoading] = useState(true);

	useEffect(() => {
		initWasm()
			.then(() => setWasmLoading(false))
			.catch((_err) => {
				setError("Failed to initialize WebAssembly");
				setWasmLoading(false);
			});
	}, []);

	const handleFileUpload = useCallback(async (file: File) => {
		try {
			setError(null);
			const buffer = await file.arrayBuffer();
			setFileBuffer(buffer);
			const data = await parseELF(buffer);
			setElfData(data);
		} catch (err) {
			setError(err instanceof Error ? err.message : "Failed to parse ELF file");
			setElfData(null);
			setFileBuffer(null);
		}
	}, []);

	return (
		<div className="app">
			<header className="app-header">
				<h1>ELF File Viewer</h1>
			</header>

			<main className="app-main">
				{wasmLoading ? (
					<div className="loading-message">
						<div style={{ fontSize: "2rem", marginBottom: "1rem" }}>⚡</div>
						<p>Initializing WebAssembly...</p>
					</div>
				) : (
					<FileUpload onFileSelect={handleFileUpload} />
				)}

				{error && <div className="error-message">{error}</div>}

				{elfData && (
					<div className="elf-content">
						<nav className="tabs">
							<button
								type="button"
								className={activeTab === "header" ? "active" : ""}
								onClick={() => setActiveTab("header")}
							>
								Header
							</button>
							<button
								type="button"
								className={activeTab === "sections" ? "active" : ""}
								onClick={() => setActiveTab("sections")}
							>
								Sections
							</button>
							<button
								type="button"
								className={activeTab === "segments" ? "active" : ""}
								onClick={() => setActiveTab("segments")}
							>
								Segments
							</button>
							<button
								type="button"
								className={activeTab === "symbols" ? "active" : ""}
								onClick={() => setActiveTab("symbols")}
							>
								Symbols
							</button>
							<button
								type="button"
								className={activeTab === "hex" ? "active" : ""}
								onClick={() => setActiveTab("hex")}
							>
								Hex Dump
							</button>
						</nav>

						<div className="tab-content">
							{activeTab === "header" && <ELFHeader data={elfData} />}
							{activeTab === "sections" && (
								<SectionHeaders sections={elfData.sectionHeaders} />
							)}
							{activeTab === "segments" && (
								<ProgramHeaders segments={elfData.programHeaders} />
							)}
							{activeTab === "symbols" && <Symbols symbols={elfData.symbols} sectionHeaders={elfData.sectionHeaders} />}
							{activeTab === "hex" && fileBuffer && (
								<HexDump
									buffer={fileBuffer}
									sections={elfData.sectionHeaders}
								/>
							)}
						</div>
					</div>
				)}
			</main>
		</div>
	);
}

export default App;
