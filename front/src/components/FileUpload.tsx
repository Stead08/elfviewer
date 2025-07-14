import type React from "react";
import { useCallback } from "react";

interface FileUploadProps {
	onFileSelect: (file: File) => void;
}

export const FileUpload: React.FC<FileUploadProps> = ({ onFileSelect }) => {
	const handleFileChange = useCallback(
		(event: React.ChangeEvent<HTMLInputElement>) => {
			const file = event.target.files?.[0];
			if (file) {
				onFileSelect(file);
			}
		},
		[onFileSelect],
	);

	const handleDrop = useCallback(
		(event: React.DragEvent<HTMLDivElement>) => {
			event.preventDefault();
			const file = event.dataTransfer.files[0];
			if (file) {
				onFileSelect(file);
			}
		},
		[onFileSelect],
	);

	const handleDragOver = useCallback(
		(event: React.DragEvent<HTMLDivElement>) => {
			event.preventDefault();
		},
		[],
	);

	return (
		// biome-ignore lint/a11y/useSemanticElements: <explanation>
		<div
			className="file-upload"
			onDrop={handleDrop}
			onDragOver={handleDragOver}
			role="button"
			tabIndex={0}
			onKeyDown={(e) => {
				if (e.key === "Enter" || e.key === " ") {
					e.preventDefault();
					document.getElementById("file-input")?.click();
				}
			}}
		>
			<input
				type="file"
				onChange={handleFileChange}
				accept=".elf,.so,.o,.out"
				style={{ display: "none" }}
				id="file-input"
			/>
			<label htmlFor="file-input" style={{ cursor: "pointer" }}>
				<div style={{ fontSize: "3rem", marginBottom: "1rem" }}>📦</div>
				<h3>Drop ELF file here or click to browse</h3>
				<p>
					Supported formats: ELF executables, shared libraries (.so), object
					files (.o)
				</p>
			</label>
		</div>
	);
};
