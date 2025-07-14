import type React from 'react';
import { type ELFInfo, getClassName, getDataEncoding, getTypeString, getMachineString } from '../utils/wasm';

interface ELFHeaderProps {
  data: ELFInfo;
}

export const ELFHeader: React.FC<ELFHeaderProps> = ({ data }) => {
  return (
    <div className="elf-header">
      <h2>ELF Header</h2>
      <table>
        <tbody>
          <tr>
            <td><strong>Magic:</strong></td>
            <td className="mono">
              {data.ident.Magic.map(b => b.toString(16).padStart(2, '0')).join(' ')}
            </td>
          </tr>
          <tr>
            <td><strong>Class:</strong></td>
            <td>{getClassName(data.class)}</td>
          </tr>
          <tr>
            <td><strong>Data:</strong></td>
            <td>{getDataEncoding(data.ident.Data)}</td>
          </tr>
          <tr>
            <td><strong>Version:</strong></td>
            <td>{data.ident.Version} (current)</td>
          </tr>
          <tr>
            <td><strong>OS/ABI:</strong></td>
            <td>{data.ident.OSABI}</td>
          </tr>
          <tr>
            <td><strong>Type:</strong></td>
            <td>{getTypeString(data.type)}</td>
          </tr>
          <tr>
            <td><strong>Machine:</strong></td>
            <td>{getMachineString(data.machine)}</td>
          </tr>
          <tr>
            <td><strong>Entry point address:</strong></td>
            <td className="mono">0x{data.entry.toString(16)}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};