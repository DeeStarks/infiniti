package config

func GetTokenSecret() string {
	return GetEnv("RESTAPI_SECRET")
}

func GetExchangeRateAPIKey() string {
	return GetEnv("EXCHANGE_RATE_API_KEY")
}