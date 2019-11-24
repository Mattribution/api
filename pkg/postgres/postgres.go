package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/mattribution/api/pkg/api"
)

type TrackService struct {
	DB *sql.DB
}

// NewTrackService Creates a new Trackservice object
func NewTrackService(host, username, password, dbName string, port int) (*TrackService, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		println(err)
		return nil, err
	}

	return &TrackService{
		db,
	}, nil
}

// StoreTrack stores the track in the db
func (s *TrackService) StoreTrack(t api.Track) (int, error) {
	// sqlStatement := `INSERT INTO tracks (owner_id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at, extra)
	// VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlStatement :=
		`INSERT INTO public.tracks (owner_id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at, extra)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING id`

	log.Println(sqlStatement)
	// Set default json value (so postgres doesn't get mad)
	if t.Extra == "" {
		t.Extra = "{}"
	}

	id := 0
	// err := s.DB.QueryRow(sqlStatement, t.OwnerID, t.UserID, t.FpHash, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, sentAt, sentAt, t.Extra).Scan(&id)
	err := s.DB.QueryRow(sqlStatement, t.OwnerID, t.UserID, t.FpHash, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, t.SentAt.Format(time.RFC3339), time.Now().Format(time.RFC3339), t.Extra).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// FindByID finds all track objects by owner id
func (s TrackService) FindByID(id int) (api.Track, error) {
	sqlStatement :=
		`SELECT id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, extra from public.tracks
		WHERE id = $1`

	var t api.Track

	row := s.DB.QueryRow(sqlStatement, id)
	switch err := row.Scan(&t.ID, &t.UserID, &t.FpHash, &t.PageURL, &t.PagePath, &t.PageReferrer, &t.PageTitle, &t.Event, &t.CampaignSource, &t.CampaignMedium, &t.CampaignName, &t.CampaignContent, &t.SentAt, &t.Extra); err {
	case sql.ErrNoRows:
		log.Println("No rows were returned!")
	case nil:
		log.Printf("FOUND NIL: %v", t)
	default:
		panic(err)
	}

	return t, nil
}
