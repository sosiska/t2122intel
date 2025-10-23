package records

// Trading212Record represents a record from Trading 212 CSV export
type Trading212Record struct {
	Action             string
	Time               string
	ISIN               string
	Ticker             string
	Name               string
	Notes              string
	ID                 string
	NumberOfShares     string
	PricePerShare      string
	CurrencyPriceShare string
	ExchangeRate       string
	CurrencyResult     string
	Total              string
	CurrencyTotal      string
}

// IntelinvestRecord represents a record in Intelinvest format
type IntelinvestRecord struct {
	Type          string
	Date          string
	TickerISIN    string
	Quantity      string
	Price         string
	Fee           string
	NKD           string
	Nominal       string
	Currency      string
	FeeCurrency   string
	Note          string
	TradeSystemID string
}

// AssetDefinition represents an asset definition for Intelinvest
type AssetDefinition struct {
	Type     string // ETF, STOCK, etc.
	Ticker   string
	Name     string
	Price    string
	Currency string
}

