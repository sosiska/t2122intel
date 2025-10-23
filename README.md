# t2122intel

Convert Trading 212 CSV export files to Intelinvest format.

## Installation

```bash
go build -o t2122intel ./cmd/t2122intel
```

## Usage

```bash
./t2122intel input.csv output.csv
```

Example:
```bash
./t2122intel ~/Downloads/trading212_export.csv intelinvest_import.csv
```

## Features

- Converts Trading 212 CSV exports to Intelinvest CSV format
- Automatic GBX to GBP conversion (divides by 100)
- Converts Interest on cash to INCOME + MONEYDEPOSIT
- Converts Market buy/sell to SHARE_BUY/SHARE_SELL + MONEYWITHDRAW/MONEYDEPOSIT
- Zero external dependencies
- Preserves Trading 212 transaction IDs

## Supported Operations

- **Interest on cash** → INCOME + MONEYDEPOSIT
- **Market buy** → SHARE_BUY + MONEYWITHDRAW
- **Market sell** → SHARE_SELL + MONEYDEPOSIT

## Project Structure

```
cmd/t2122intel/    - CLI entry point
records/           - Data structures
parser/            - Trading 212 CSV parser
exporter/          - Intelinvest CSV exporter
converter/         - Format conversion logic
```

## Example

Input (Trading 212):
```csv
Action,Time,ISIN,Ticker,Name,Notes,ID,No. of shares,Price / share,Currency (Price / share),Exchange rate,Currency (Result),Total,Currency (Total)
Interest on cash,2025-10-17 00:05:53,,,,"Interest on cash",b4cf857f,,,,,,,3.62,"GBP"
Market buy,2025-10-22 14:45:46,IE00BKW9SV11,IGBE,"Invesco Bond",,EOF407,2.0000000000,3339.5000000000,GBX,100.00000000,"GBP",66.79,"GBP"
```

Output (Intelinvest):
```csv
#CsvFormatVersion:v1

INCOME;17.10.2025 00:05:53;;;3.62;;;;GBP;;;;b4cf857f
MONEYDEPOSIT;17.10.2025 00:05:53;;;3.62;;;;GBP;;Interest deposit from 17.10.2025 00:05:53;;b4cf857f
SHARE_BUY;22.10.2025 14:45:46;IGBE:IE00BKW9SV11;2.0000000000;33.395;0;;;GBP;GBP;;;EOF407
MONEYWITHDRAW;22.10.2025 14:45:46;;;66.79;;;;GBP;;Payment for IGBE:IE00BKW9SV11;;EOF407
```

## License

MIT
