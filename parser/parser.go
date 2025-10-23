package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/sosiska/t2122intel/records"
)

// Trading212Parser reads and parses Trading 212 CSV files
type Trading212Parser struct{}

// New creates a new parser instance
func New() *Trading212Parser {
	return &Trading212Parser{}
}

// ParseFile reads a Trading 212 CSV file and returns parsed records
func (p *Trading212Parser) ParseFile(filename string) ([]records.Trading212Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return p.Parse(file)
}

// Parse reads Trading 212 CSV data from a reader
func (p *Trading212Parser) Parse(r io.Reader) ([]records.Trading212Record, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("file is empty or contains only header")
	}

	var recs []records.Trading212Record
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 14 {
			continue
		}

		record := records.Trading212Record{
			Action:             row[0],
			Time:               row[1],
			ISIN:               row[2],
			Ticker:             row[3],
			Name:               row[4],
			Notes:              row[5],
			ID:                 row[6],
			NumberOfShares:     row[7],
			PricePerShare:      row[8],
			CurrencyPriceShare: row[9],
			ExchangeRate:       row[10],
			CurrencyResult:     row[11],
			Total:              row[12],
			CurrencyTotal:      row[13],
		}
		recs = append(recs, record)
	}

	return recs, nil
}

