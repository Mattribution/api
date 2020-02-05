package postgres

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/mattribution/api/internal/app"

	// Import Postgres SQL driver
	_ "github.com/lib/pq"
)

// Tracks handles Track data
type Tracks struct {
	DB *sql.DB
}

func NewCloudSQLClient(dbUser, dbPass, dbName, instanceName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password='%s' host=/cloudsql/%s dbname=%s sslmode=disable", dbUser, dbPass, instanceName, dbName)
	// connStr := fmt.Sprintf("postgres://%s:%s@/%s?unix_socket=/cloudsql/%s/.s.PGSQL.5432", dbUser, dbPass, dbName, instanceName)
	log.Println(connStr)
	return sql.Open("postgres", connStr)
}

func (s *Tracks) Store(t app.Track) (int64, error) {
	sqlStatement :=
		`INSERT INTO public.tracks (owner_id, user_id, anonymous_id, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at, extra)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING id`

	var id int64
	err := s.DB.QueryRow(sqlStatement, t.OwnerID, t.UserID, t.AnonymousID, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, t.SentAt, time.Now().Format(time.RFC3339), t.Extra).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil

}
