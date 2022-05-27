package currency

func USDCentsToFloat64(usdCent int64) float64 {
	usdFloat := float64(usdCent)
	usdFloat = usdFloat / 100
	return usdFloat
}

func Float64ToUSDCent(f float64) int64 {
	usdCent := int64((f * 100) + 0.5)
	return usdCent
}
