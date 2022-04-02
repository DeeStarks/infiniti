package config

func GetTokenSecret() string {
	return GetEnv("RESTAPI_SECRET")
}