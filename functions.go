package functions

import (
	"net/http"
	"os"

	"github.com/mattribution/api/internal/app"
	internal_http "github.com/mattribution/api/internal/pkg/http"
	"github.com/mattribution/api/internal/pkg/postgres"
)

var (
	dbUser      = os.Getenv("DB_USER")
	dbPass      = os.Getenv("DB_PASS")
	dbName      = os.Getenv("DB_NAME")
	dbHost      = os.Getenv("DB_HOST")
	auth0ApiID  = os.Getenv("AUTH0_API_ID")
	auth0Domain = os.Getenv("AUTH0_DOMAIN")
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
