import type React from 'react';
import { type SectionHeader, getSectionTypeString, formatFlags } from '../utils/wasm';

interface SectionHeadersProps {
  sections: SectionHeader[];
}

export const SectionHeaders: React.FC<SectionHeadersProps> = ({ sections }) => {
  return (
    <div className="section-headers">
      <h2>Section Headers</h2>
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
            {sections.map((section, index) => (
              <tr key={`section-${index}`}>
                <td>[{index}]</td>
                <td>{section.Name || '<no-name>'}</td>
                <td>{getSectionTypeString(section.Type)}</td>
                <td className="mono">0x{section.Addr.toString(16).padStart(16, '0')}</td>
                <td className="mono">0x{section.Offset.toString(16).padStart(8, '0')}</td>
                <td className="mono">0x{section.Size.toString(16).padStart(16, '0')}</td>
                <td className="mono">0x{section.EntSize.toString(16).padStart(16, '0')}</td>
                <td>{formatFlags(section.Flags)}</td>
                <td>{section.Link}</td>
                <td>{section.Info}</td>
                <td>{section.AddrAlign}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div style={{ marginTop: '1rem', fontSize: '0.875rem', color: '#666' }}>
        <strong>Key to Flags:</strong><br />
        W (write), A (alloc), X (execute)
      </div>
    </div>
  );
};