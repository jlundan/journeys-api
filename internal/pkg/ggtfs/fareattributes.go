package ggtfs

type FareAttributes struct {
	Id               string
	Price            float64
	CurrencyType     string
	PaymentMethod    int
	Transfers        int
	AgencyId         *string
	TransferDuration *uint
	LineNumber       int
}
