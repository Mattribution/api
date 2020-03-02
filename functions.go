package functions

import (
	"net/http"
	"os"
	_ "github.com/joho/godotenv/autoload"

	"github.com/mattribution/api/internal/app"
	internal_http "github.com/mattribution/api/internal/pkg/http"
	"github.com/mattribution/api/internal/pkg/postgres"
)

var (
	dbUser      = getenv("DB_USER", "postgres")
	dbPass      = getenv("DB_PASS", "password")
	dbName      = getenv("DB_NAME", "mattribution")
	dbHost      = getenv("DB_HOST", "127.0.0.1")
	auth0ApiID  = getenv("AUTH0_API_ID", "")
	auth0Domain = getenv("AUTH0_DOMAIN", "")
	handler     *internal_http.Handler
)

func init() {
	// Setup db connection
	db, err := postgres.NewCloudSQLClient(dbUser, dbPass, dbName, dbHost)
	if err != nil {
		panic(err)
	}
	// Only allow 1 connection to the database to avoid overloading
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	tracksDAO := &postgres.TracksDAO{
		DB: db,
	}
	kpisDAO := &postgres.KpisDAO{
		DB: db,
	}

	// Setup services
	handler = internal_http.NewHandler(
		app.NewService(tracksDAO, kpisDAO),
		auth0Domain,
		auth0ApiID,
	)
}

// FunctionsEntrypoint represents cloud function entry point
func FunctionsEntrypoint(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
