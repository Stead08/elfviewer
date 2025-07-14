import type React from 'react';
import { useState, useCallback, useEffect } from 'react';
import { type SectionHeader, getHexDump } from '../utils/wasm';

interface HexDumpProps {
  buffer: ArrayBuffer;
  sections: SectionHeader[];
}

export const HexDump: React.FC<HexDumpProps> = ({ buffer, sections }) => {
  const [selectedSection, setSelectedSection] = useState<string>('');
  const [hexDump, setHexDump] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Select first non-null section by default
    const firstSection = sections.find(s => s.Name && s.Size > 0);
    if (firstSection) {
      setSelectedSection(firstSection.Name);
    }
  }, [sections]);

  const loadHexDump = useCallback(async () => {
    if (!selectedSection || !buffer) return;

    setLoading(true);
    setError(null);
    try {
      const dump = await getHexDump(buffer, selectedSection);
      setHexDump(dump);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load hex dump');
      setHexDump('');
    } finally {
      setLoading(false);
    }
  }, [buffer, selectedSection]);

  useEffect(() => {
    if (selectedSection) {
      loadHexDump();
    }
  }, [selectedSection, loadHexDump]);

  const handleSectionChange = useCallback((event: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedSection(event.target.value);
  }, []);

  return (
    <div className="hex-dump">
      <h2>Hex Dump</h2>
      <div style={{ marginBottom: '1rem' }}>
        <label htmlFor="section-select" style={{ marginRight: '0.5rem' }}>
          Select Section:
        </label>
        <select
          id="section-select"
          value={selectedSection}
          onChange={handleSectionChange}
          style={{
            padding: '0.5rem',
            fontSize: '0.875rem',
            borderRadius: '4px',
            border: '1px solid #ccc',
          }}
        >
          {sections
            .filter(s => s.Name && s.Size > 0)
            .map((section, _index) => (
              <option key={section.Name} value={section.Name}>
                {section.Name} (0x{section.Size.toString(16)} bytes)
              </option>
            ))}
        </select>
      </div>

      {loading && <p>Loading hex dump...</p>}
      {error && <p style={{ color: '#c00' }}>Error: {error}</p>}
      {!loading && !error && hexDump && (
        <div className="hex-dump-container">
          {hexDump}
        </div>
      )}
    </div>
  );
};