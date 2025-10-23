package main

import (
	"fmt"
	"os"

	"github.com/sosiska/t2122intel/converter"
	"github.com/sosiska/t2122intel/exporter"
	"github.com/sosiska/t2122intel/parser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: t2122intel <input.csv> <output.csv>")
		fmt.Println("Example: t2122intel trading212.csv intelinvest.csv")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	p := parser.New()
	records, err := p.ParseFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	c := converter.New()
	assets := c.ExtractAssets(records)
	converted := c.Convert(records)

	e := exporter.New()
	err = e.WriteFile(outputFile, assets, converted)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Converted %d records to %s\n", len(records), outputFile)
}
