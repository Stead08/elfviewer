import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import type { SectionHeader } from "../utils/wasm";
import { SectionHeaders } from "./SectionHeaders";

describe("SectionHeaders", () => {
	const mockSections: SectionHeader[] = [
		{
			Name: ".text",
			Type: 1,
			Flags: 6,
			Addr: 0x1000,
			Offset: 0x1000,
			Size: 0x500,
			Link: 0,
			Info: 0,
			AddrAlign: 16,
			EntSize: 0,
		},
		{
			Name: ".data",
			Type: 1,
			Flags: 3,
			Addr: 0x2000,
			Offset: 0x2000,
			Size: 0x200,
			Link: 0,
			Info: 0,
			AddrAlign: 4,
			EntSize: 0,
		},
		{
			Name: ".bss",
			Type: 8,
			Flags: 3,
			Addr: 0x3000,
			Offset: 0x2200,
			Size: 0x100,
			Link: 0,
			Info: 0,
			AddrAlign: 4,
			EntSize: 0,
		},
	];

	it("renders section headers table", () => {
		render(<SectionHeaders sections={mockSections} />);

		expect(screen.getByText("Section Headers")).toBeInTheDocument();
		expect(screen.getByText(".text")).toBeInTheDocument();
		expect(screen.getByText(".data")).toBeInTheDocument();
		expect(screen.getByText(".bss")).toBeInTheDocument();
	});

	it("displays section count information", () => {
		render(<SectionHeaders sections={mockSections} />);

		expect(screen.getByText("Showing 3 of 3 sections")).toBeInTheDocument();
	});

	it("sorts sections by address when checkbox is checked", () => {
		render(<SectionHeaders sections={mockSections} />);

		const sortCheckbox = screen.getByLabelText("Sort by Address");
		fireEvent.click(sortCheckbox);

		const rows = screen.getAllByRole("row");
		// Skip header row
		const dataRows = rows.slice(1);

		// Check that sections are sorted by address (0x1000, 0x2000, 0x3000)
		expect(dataRows[0]).toHaveTextContent(".text");
		expect(dataRows[1]).toHaveTextContent(".data");
		expect(dataRows[2]).toHaveTextContent(".bss");
	});

	it("filters sections with flags when checkbox is checked", () => {
		const sectionsWithEmptyFlags = [
			...mockSections,
			{
				Name: ".comment",
				Type: 1,
				Flags: 0,
				Addr: 0,
				Offset: 0x4000,
				Size: 0x50,
				Link: 0,
				Info: 0,
				AddrAlign: 1,
				EntSize: 0,
			},
		];

		render(<SectionHeaders sections={sectionsWithEmptyFlags} />);

		const filterCheckbox = screen.getByLabelText(
			"Show only sections with flags",
		);
		fireEvent.click(filterCheckbox);

		expect(screen.getByText("Showing 3 of 4 sections")).toBeInTheDocument();
		expect(screen.queryByText(".comment")).not.toBeInTheDocument();
	});

	it("displays flag key information", () => {
		render(<SectionHeaders sections={mockSections} />);

		expect(screen.getByText("Key to Flags:")).toBeInTheDocument();
		expect(
			screen.getByText("W (write), A (alloc), X (execute)"),
		).toBeInTheDocument();
	});
});
