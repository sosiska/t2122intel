package exporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/sosiska/t2122intel/records"
)

// Exporter writes records in Intelinvest CSV format
type Exporter struct{}

// New creates a new exporter instance
func New() *Exporter {
	return &Exporter{}
}

// WriteFile writes Intelinvest records to a CSV file
func (e *Exporter) WriteFile(filename string, assets []records.AssetDefinition, recs []records.IntelinvestRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return e.Write(file, assets, recs)
}

// Write writes Intelinvest records to a writer
func (e *Exporter) Write(wr io.Writer, assets []records.AssetDefinition, recs []records.IntelinvestRecord) error {
	writer := csv.NewWriter(wr)
	writer.Comma = ';'
	defer writer.Flush()

	// Header
	if err := writer.Write([]string{"#CsvFormatVersion:v1"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Assets definitions section
	if len(assets) > 0 {
		if err := writer.Write([]string{"#AssetsDefinitionsStart.v1"}); err != nil {
			return fmt.Errorf("failed to write assets header: %w", err)
		}

		for _, asset := range assets {
			row := []string{
				asset.Type,
				asset.Ticker,
				asset.Name,
				asset.Price,
				asset.Currency,
				"", // URL
				"", // Selector
				"", // Extra field 1
				"", // Extra field 2
			}
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write asset: %w", err)
			}
		}

		if err := writer.Write([]string{"#AssetsDefinitionsEnd"}); err != nil {
			return fmt.Errorf("failed to write assets end: %w", err)
		}

		if err := writer.Write([]string{""}); err != nil {
			return fmt.Errorf("failed to write empty line: %w", err)
		}
	}

	// Records
	for _, record := range recs {
		row := []string{
			record.Type,
			record.Date,
			record.TickerISIN,
			record.Quantity,
			record.Price,
			record.Fee,
			record.NKD,
			record.Nominal,
			record.Currency,
			record.FeeCurrency,
			record.Note,
			record.LinkID,
			record.TradeSystemID,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}
