package settings

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	HasuraAdminSecret     string
	HasuraGraphQLEndpoint string
	Port                  string
	LogLevel              string
)

func init() {
	var ok bool

	godotenv.Load()

	variables := []string{"HASURA_GRAPHQL_ADMIN_SECRET", "HASURA_GRAPHQL_ENDPOINT"}

	for _, v := range variables {
		if _, ok := os.LookupEnv(v); !ok {
			logrus.Fatalf("%s variable is required\n", v)
		}
	}

	HasuraAdminSecret = os.Getenv(variables[0])
	HasuraGraphQLEndpoint = os.Getenv(variables[1])

	if Port, ok = os.LookupEnv("PORT"); !ok {
		Port = "9921"
	}

	if LogLevel, ok = os.LookupEnv("LOG_LEVEL"); !ok {
		LogLevel = "error"
	}
}
