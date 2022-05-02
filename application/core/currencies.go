package core

// amount, fromConversionRate, toConversionRate
func (core *CoreApplication) ConvertCurrency(amount, fromConversionRate, toConversionRate float64) float64 {
	if fromConversionRate == toConversionRate {
		return amount
	}

	// Convert amount to USD, then to target currency
	amountUSD := amount/fromConversionRate
	return amountUSD * toConversionRate
}