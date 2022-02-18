package settings

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	HasuraAdminSecret     string
	HasuraGraphQLEndpoint string
	WebAddr               string
	LogLevel              string
)

func init() {

	var ok bool

	godotenv.Load()

	variables := []string{"HASURA_GRAPHQL_ADMIN_SECRET", "HASURA_GRAPHQL_ENDPOINT", "WEB_ADDR", "LOG_LEVEL"}

	for _, v := range variables {
		if _, ok := os.LookupEnv(v); !ok {
			logrus.Fatalf("%s variable is required\n", v)
		}
	}

	HasuraAdminSecret = os.Getenv(variables[0])

	HasuraGraphQLEndpoint = os.Getenv(variables[1])

	if WebAddr, ok = os.LookupEnv(variables[2]); !ok {
		WebAddr = ":9921"
	}

	if LogLevel, ok = os.LookupEnv(variables[2]); !ok {
		LogLevel = "info"
	}

}
