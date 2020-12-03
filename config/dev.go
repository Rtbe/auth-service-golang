package config

import "os"

//Dev sets development environment variables.
func dev() {
	//API environment variables
	os.Setenv("PORT", "8080")
	//JWT secret
	os.Setenv("TOKEN_SECRET", "tokensecrettokensecret")
	//Database environment variables
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "testTask")
	os.Setenv("DB_PORT", "27017")
}
