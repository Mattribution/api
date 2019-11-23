package presto

import (
	"database/sql"
	_ "github.com/prestodb/presto-go-client/presto"
	"context"
	"log"
)

type TrackService struct {
	DB *sql.DB
	username string,
	password string,
	uri string
}

func NewTrackService(username, password, uri string) (TrackService, error) {
	dsn := fmt.Sprintf("http://%s:%s@%s:8080?catalog=postgres-dev&schema=public", username, password, uri)
	db, err := sql.Open("presto", dsn)	
	if err != nil {
		println(err)
		return
	}

	return TrackService{
		db
	}
}

func (s TrackService) NewTrack(t Track) error {
	sqlStatement :=
		`INSERT INTO tracks (owner_id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at, extra)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING id`

	id = 0
	// Set default json value (so postgres doesn't get mad)
	if t.Extra == "" {
		t.Extra = "{}"
	}

	err = s.DB.QueryRow(sqlStatement, ownerID, t.UserID, t.FpHash, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, t.SentAt.Format(time.RFC3339), time.Now().Format(time.RFC3339), t.Extra).Scan(&id)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return id, err
	}

	return id, nil
}

// FindByID finds all track objects by owner id
func (p TrackService) FindByID(id int) (Track, error) {
	sqlStatement :=
		`SELECT id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, extra from public.tracks
		WHERE id = $1`

	rows, err := p.db.Query(sqlStatement, id)
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	defer rows.Close()

	tracks := []Track{}
	for rows.Next() {
		var t Track
		err = rows.Scan(&t.ID, &t.UserID, &t.FpHash, &t.PageURL, &t.PagePath, &t.PageReferrer, &t.PageTitle, &t.Event, &t.CampaignSource, &t.CampaignMedium, &t.CampaignName, &t.CampaignContent, &t.SentAt, &t.Extra)
		if err != nil {
			// handle this error
			return nil, err
		}
		tracks = append(tracks, t)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}