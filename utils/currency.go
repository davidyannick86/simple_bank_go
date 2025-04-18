package utils

const (
	USD = "USD"
	CAD = "CAD"
	EUR = "EUR"
	GBP = "GBP"
	CHF = "CHF"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, CAD, EUR, GBP, CHF:
		return true
	default:
		return false
	}
}
