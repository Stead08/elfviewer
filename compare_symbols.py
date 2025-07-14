import re

# Read CLI labels
with open('cli_labels_raw.txt', 'r') as f:
    cli_lines = f.readlines()

# Extract and format
cli_labels = {}
for line in cli_lines:
    parts = line.strip().split()
    if len(parts) >= 9:
        addr = parts[1].lstrip('0') or '0'
        name = parts[8]
        cli_labels[name] = addr.upper()

# Read expected labels  
with open('expected_labels.txt', 'r') as f:
    expected_lines = f.readlines()

# Compare
print('Comparison Results:')
print('==================')
print()

missing_in_cli = []
for line in expected_lines:
    parts = line.strip().split()
    if len(parts) >= 2:
        addr = parts[0].upper()
        name = parts[1]
        if name in cli_labels:
            if cli_labels[name] != addr:
                print(f'MISMATCH: {name}')
                print(f'  Expected: {addr}')
                print(f'  CLI:      {cli_labels[name]}')
                print()
        else:
            missing_in_cli.append(f'{addr} {name}')

if missing_in_cli:
    print('Missing in CLI output:')
    for item in missing_in_cli:
        print(f'  {item}')
    print()

# Check for _start 
if '_start' not in [name for name in cli_labels.keys()]:
    print('Note: _start symbol missing from CLI output')
    # Look for it in the full symbol table
    with open('cli_output_all.txt', 'r') as f:
        for line in f:
            if '_start' in line and 'GLOBAL' in line:
                print(f'Found in full output: {line.strip()}')

print('\nSummary:')
print(f'Expected symbols: {len([l for l in expected_lines if l.strip()])}')
print(f'Found in CLI: {len([name for name in cli_labels.keys() if name in [l.split()[1] for l in expected_lines if len(l.split()) >= 2]])}')