import React from 'react';
import { Symbol } from '../utils/wasm';

interface SymbolsProps {
  symbols: Symbol[];
}

function getSymbolType(info: number): string {
  const type = info & 0xf;
  switch (type) {
    case 0: return 'NOTYPE';
    case 1: return 'OBJECT';
    case 2: return 'FUNC';
    case 3: return 'SECTION';
    case 4: return 'FILE';
    case 5: return 'COMMON';
    case 6: return 'TLS';
    default: return `<${type}>`;
  }
}

function getSymbolBind(info: number): string {
  const bind = info >> 4;
  switch (bind) {
    case 0: return 'LOCAL';
    case 1: return 'GLOBAL';
    case 2: return 'WEAK';
    default: return `<${bind}>`;
  }
}

function getSymbolVisibility(other: number): string {
  const vis = other & 0x3;
  switch (vis) {
    case 0: return 'DEFAULT';
    case 1: return 'INTERNAL';
    case 2: return 'HIDDEN';
    case 3: return 'PROTECTED';
    default: return `<${vis}>`;
  }
}

function getSectionIndex(shndx: number): string {
  switch (shndx) {
    case 0: return 'UND';
    case 0xfff1: return 'ABS';
    case 0xfff2: return 'COMMON';
    default: return shndx.toString();
  }
}

export const Symbols: React.FC<SymbolsProps> = ({ symbols }) => {
  if (symbols.length === 0) {
    return (
      <div className="symbols">
        <h2>Symbols</h2>
        <p>No symbols found.</p>
      </div>
    );
  }

  return (
    <div className="symbols">
      <h2>Symbols</h2>
      <div className="table-container">
        <table>
          <thead>
            <tr>
              <th>Num</th>
              <th>Value</th>
              <th>Size</th>
              <th>Type</th>
              <th>Bind</th>
              <th>Vis</th>
              <th>Ndx</th>
              <th>Name</th>
            </tr>
          </thead>
          <tbody>
            {symbols.map((symbol, index) => (
              <tr key={index}>
                <td>{index}:</td>
                <td className="mono">0x{symbol.Value.toString(16).padStart(16, '0')}</td>
                <td>{symbol.Size}</td>
                <td>{getSymbolType(symbol.Info)}</td>
                <td>{getSymbolBind(symbol.Info)}</td>
                <td>{getSymbolVisibility(symbol.Other)}</td>
                <td>{getSectionIndex(symbol.Shndx)}</td>
                <td>{symbol.Name || '<no-name>'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};