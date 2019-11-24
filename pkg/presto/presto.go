package presto

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mattribution/api/pkg/api"
	_ "github.com/prestodb/presto-go-client/presto"
)

type TrackService struct {
	DB       *sql.DB
	username string
	password string
	uri      string
}

// NewTrackService Creates a new Trackservice object
func NewTrackService(username, password, uri string) (*TrackService, error) {
	dsn := fmt.Sprintf("http://%s:%s@%s:8080?catalog=postgres-dev&schema=public", username, password, uri)
	db, err := sql.Open("presto", dsn)
	if err != nil {
		println(err)
		return nil, err
	}

	return &TrackService{
		db,
		username,
		password,
		uri,
	}, nil
}

// StoreTrack stores the track in the db
func (s *TrackService) StoreTrack(t api.Track) (int, error) {
	// sqlStatement := `INSERT INTO tracks (owner_id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at, extra)
	// VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlStatement := `INSERT INTO tracks (owner_id)
	VALUES(?)`

	log.Println(sqlStatement)
	// Set default json value (so postgres doesn't get mad)
	if t.Extra == "" {
		t.Extra = "{}"
	}

	id := 0
	// err := s.DB.QueryRow(sqlStatement, t.OwnerID, t.UserID, t.FpHash, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, sentAt, sentAt, t.Extra).Scan(&id)
	err := s.DB.QueryRow(sqlStatement, t.OwnerID).Scan(&id)
	if err != nil {
		log.Printf("ERROR: %v", err)
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

// func (s TrackService) DeleteTrack(id int) error {

// }
