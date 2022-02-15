package settings

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	HasuraAdminSecret string
	HasuraGraphQLURL  string
	WebPath           string
	WebAddr           string
	LogLevel          string
	Debug             bool
)

func init() {

	var ok bool

	godotenv.Load()

	variables := []string{"HASURA_GRAPHQL_ADMIN_SECRET", "HASURA_GRAPHQL_URL", "WEB_PATH", "WEB_ADDR", "LOG_LEVEL"}

	for _, v := range variables {
		if _, ok := os.LookupEnv(v); !ok {
			logrus.Fatalf("%s variable is required\n", v)
		}
	}

	HasuraAdminSecret = os.Getenv(variables[0])

	HasuraGraphQLURL = os.Getenv(variables[1])

	WebPath = os.Getenv(variables[2])

	WebAddr = os.Getenv(variables[3])

	if LogLevel, ok = os.LookupEnv(variables[4]); !ok {
		LogLevel = "info"
	}

}
