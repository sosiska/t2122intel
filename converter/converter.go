package converter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sosiska/t2122intel/records"
)

// Converter converts Trading 212 records to Intelinvest format
type Converter struct{}

// New creates a new converter instance
func New() *Converter {
	return &Converter{}
}

// Convert converts a slice of Trading 212 records to Intelinvest records
func (c *Converter) Convert(recs []records.Trading212Record) []records.IntelinvestRecord {
	var result []records.IntelinvestRecord

	for _, record := range recs {
		converted := c.convertRecord(record)
		result = append(result, converted...)
	}

	return result
}

// ExtractAssets extracts unique asset definitions from Trading 212 records
// Returns only custom assets that need to be defined in AssetsDefinitions section
func (c *Converter) ExtractAssets(recs []records.Trading212Record) []records.AssetDefinition {
	// Custom assets that are not in Intelinvest database
	customAssets := map[string]records.AssetDefinition{
		"XMWX": {
			Type:     "ETF",
			Ticker:   "XMWX",
			Name:     "Xtrackers MSCI World ex USA UCITS ETF 1C",
			Price:    "29.54",
			Currency: "GBP",
		},
		"HEMC": {
			Type:     "ETF",
			Ticker:   "HEMC",
			Name:     "HSBC MSCI Emerging Markets UCITS ETF USD (Acc)",
			Price:    "11.12",
			Currency: "GBP",
		},
		"EMUG": {
			Type:     "ETF",
			Ticker:   "EMUG",
			Name:     "L&G ESG Emerging Markets Corporate Bond UCITS ETF",
			Price:    "6.61",
			Currency: "GBP",
		},
		"V3GS": {
			Type:     "ETF",
			Ticker:   "V3GS",
			Name:     "Vanguard ESG Global Corporate Bond UCITS ETF GBP Hedged Accumulating",
			Price:    "5.13",
			Currency: "GBP",
		},
	}

	// Collect only custom assets that appear in transactions
	usedAssets := make(map[string]bool)
	for _, record := range recs {
		if record.Ticker == "" {
			continue
		}
		if record.Action == "Market buy" || record.Action == "Market sell" {
			usedAssets[record.Ticker] = true
		}
	}

	var assets []records.AssetDefinition
	for ticker := range usedAssets {
		if asset, ok := customAssets[ticker]; ok {
			assets = append(assets, asset)
		}
	}

	return assets
}

func (c *Converter) convertRecord(record records.Trading212Record) []records.IntelinvestRecord {
	switch record.Action {
	case "Interest on cash":
		return c.convertInterestOnCash(record)
	case "Market buy":
		return c.convertMarketBuy(record)
	case "Market sell":
		return c.convertMarketSell(record)
	default:
		return nil
	}
}

func (c *Converter) convertInterestOnCash(record records.Trading212Record) []records.IntelinvestRecord {
	date := c.convertDate(record.Time)
	amount := record.Total
	currency := record.CurrencyTotal

	// Generate unique LinkID for related operations
	linkID := c.generateLinkID(record.ID)

	return []records.IntelinvestRecord{
		{
			Type:          "INCOME",
			Date:          date,
			TickerISIN:    "",
			Quantity:      "",
			Price:         amount,
			Fee:           "",
			NKD:           "",
			Nominal:       "",
			Currency:      currency,
			FeeCurrency:   "",
			Note:          "",
			LinkID:        linkID,
			TradeSystemID: record.ID,
		},
		{
			Type:          "MONEYDEPOSIT",
			Date:          date,
			TickerISIN:    "",
			Quantity:      "",
			Price:         amount,
			Fee:           "",
			NKD:           "",
			Nominal:       "",
			Currency:      currency,
			FeeCurrency:   "",
			Note:          fmt.Sprintf("Income deposit from %s", date),
			LinkID:        linkID,
			TradeSystemID: record.ID,
		},
	}
}

func (c *Converter) convertMarketBuy(record records.Trading212Record) []records.IntelinvestRecord {
	date := c.convertDate(record.Time)
	ticker := record.Ticker
	isin := record.ISIN
	quantity := record.NumberOfShares
	price := c.normalizePrice(record.PricePerShare, record.CurrencyPriceShare)
	total := record.Total
	currency := record.CurrencyTotal

	tickerISIN := c.formatTickerISIN(ticker, isin)
	name := record.Name
	if name == "" {
		name = ticker
	}

	// Generate unique LinkID for related operations
	linkID := c.generateLinkID(record.ID)

	return []records.IntelinvestRecord{
		{
			Type:          "SHARE_BUY",
			Date:          date,
			TickerISIN:    tickerISIN,
			Quantity:      quantity,
			Price:         price,
			Fee:           "0",
			NKD:           "",
			Nominal:       "",
			Currency:      currency,
			FeeCurrency:   currency,
			Note:          "",
			LinkID:        linkID,
			TradeSystemID: record.ID,
		},
		{
			Type:          "MONEYWITHDRAW",
			Date:          date,
			TickerISIN:    "",
			Quantity:      "",
			Price:         total,
			Fee:           "",
			NKD:           "",
			Nominal:       "",
			Currency:      currency,
			FeeCurrency:   "",
			Note:          fmt.Sprintf("Payment for %s", name),
			LinkID:        linkID,
			TradeSystemID: record.ID,
		},
	}
}

func (c *Converter) convertMarketSell(record records.Trading212Record) []records.IntelinvestRecord {
	date := c.convertDate(record.Time)
	ticker := record.Ticker
	isin := record.ISIN
	quantity := record.NumberOfShares
	price := c.normalizePrice(record.PricePerShare, record.CurrencyPriceShare)
	total := record.Total
	currency := record.CurrencyTotal

	tickerISIN := c.formatTickerISIN(ticker, isin)
	name := record.Name
	if name == "" {
		name = ticker
	}

	// Generate unique LinkID for related operations
	linkID := c.generateLinkID(record.ID)

	return []records.IntelinvestRecord{
		{
			Type:          "SHARE_SELL",
			Date:          date,
			TickerISIN:    tickerISIN,
			Quantity:      quantity,
			Price:         price,
			Fee:           "0",
			NKD:           "",
			Nominal:       "",
			Currency:      currency,
			FeeCurrency:   currency,
			Note:          "",
			LinkID:        linkID,
			TradeSystemID: record.ID,
		},
		{
			Type:          "MONEYDEPOSIT",
			Date:          date,
			TickerISIN:    "",
			Quantity:      "",
			Price:         total,
			Fee:           "",
			NKD:           "",
			Nominal:       "",
			Currency:      currency,
			FeeCurrency:   "",
			Note:          fmt.Sprintf("Proceeds from %s", name),
			LinkID:        linkID,
			TradeSystemID: record.ID,
		},
	}
}

func (c *Converter) convertDate(dateStr string) string {
	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("02.01.2006 15:04:05")
}

func (c *Converter) normalizePrice(price, currency string) string {
	if currency == "GBX" {
		if priceFloat, err := strconv.ParseFloat(price, 64); err == nil {
			return c.formatNumber(priceFloat / 100.0)
		}
	}
	return price
}

func (c *Converter) formatTickerISIN(ticker, isin string) string {
	// Map Trading 212 tickers to exchange-specific tickers for Intelinvest
	tickerMapping := map[string]string{
		"IGBE": "IGBE.L",  // London Stock Exchange
		"EPRA": "EPRA.PA", // Euronext Paris
		"XSGI": "DX2E.DE", // Deutsche Börse
	}

	// Custom ETFs use TICKER:TICKER format instead of TICKER:ISIN
	customTickers := map[string]bool{
		"EMUG": true,
		"HEMC": true,
		"V3GS": true,
		"XMWX": true,
	}

	if mappedTicker, ok := tickerMapping[ticker]; ok {
		ticker = mappedTicker
	}

	if ticker == "" {
		return isin
	}
	if isin == "" {
		return ticker
	}

	// For custom ETFs use TICKER:TICKER format
	if customTickers[ticker] {
		return fmt.Sprintf("%s:%s", ticker, ticker)
	}

	return fmt.Sprintf("%s:%s", ticker, isin)
}

func (c *Converter) formatNumber(num float64) string {
	str := strconv.FormatFloat(num, 'f', 10, 64)
	str = strings.TrimRight(str, "0")
	str = strings.TrimRight(str, ".")
	return str
}

// generateLinkID generates a unique numeric LinkID from TradeSystemID
// Extracts numeric part or generates hash to create numeric ID
func (c *Converter) generateLinkID(tradeSystemID string) string {
	// Extract only digits from TradeSystemID
	var digits strings.Builder
	for _, ch := range tradeSystemID {
		if ch >= '0' && ch <= '9' {
			digits.WriteByte(byte(ch))
		}
	}

	numericID := digits.String()
	// Take last 8 digits to create shorter ID
	if len(numericID) > 8 {
		return numericID[len(numericID)-8:]
	}
	if len(numericID) > 0 {
		return numericID
	}

	// Fallback: if no digits, return "1"
	return "1"
}
