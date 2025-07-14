import type React from 'react';
import { type ProgramHeader, getSegmentTypeString, formatSegmentFlags } from '../utils/wasm';

interface ProgramHeadersProps {
  segments: ProgramHeader[];
}

export const ProgramHeaders: React.FC<ProgramHeadersProps> = ({ segments }) => {
  return (
    <div className="program-headers">
      <h2>Program Headers</h2>
      <div className="table-container">
        <table>
          <thead>
            <tr>
              <th>Type</th>
              <th>Offset</th>
              <th>VirtAddr</th>
              <th>PhysAddr</th>
              <th>FileSiz</th>
              <th>MemSiz</th>
              <th>Flags</th>
              <th>Align</th>
            </tr>
          </thead>
          <tbody>
            {segments.map((segment, index) => (
              <tr key={`segment-${index}`}>
                <td>{getSegmentTypeString(segment.Type)}</td>
                <td className="mono">0x{segment.Offset.toString(16).padStart(16, '0')}</td>
                <td className="mono">0x{segment.VAddr.toString(16).padStart(16, '0')}</td>
                <td className="mono">0x{segment.PAddr.toString(16).padStart(16, '0')}</td>
                <td className="mono">0x{segment.FileSz.toString(16).padStart(16, '0')}</td>
                <td className="mono">0x{segment.MemSz.toString(16).padStart(16, '0')}</td>
                <td>{formatSegmentFlags(segment.Flags)}</td>
                <td className="mono">0x{segment.Align.toString(16)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};