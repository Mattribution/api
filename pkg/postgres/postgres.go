package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/mattribution/api/pkg/api"
)

const (
	missingDataErr = "The payload is missing required data"
)

type TrackService struct {
	DB *sql.DB
}

type KPIService struct {
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

// GetTopValuesFromColumn gets a theGetTopValuesFromColumn top values from a column along with their counts
// TODO: implement daily limit
func (s TrackService) GetTopValuesFromColumn(days int, column, table string) ([]api.ValueCount, error) {
	sqlStatement :=
		fmt.Sprintf(`SELECT %s, count(*) count FROM %s
		GROUP BY 1
		ORDER BY 2 DESC
		LIMIT 10;`, column, table)

	rows, err := s.DB.Query(sqlStatement)
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	defer rows.Close()

	vCounts := []api.ValueCount{}
	for rows.Next() {
		var vCount api.ValueCount
		err = rows.Scan(&vCount.Value, &vCount.Count)
		if err != nil {
			// handle this error
			return nil, err
		}
		vCounts = append(vCounts, vCount)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return vCounts, nil
}

func (s TrackService) GetCountsFromColumn(days int, column, table string) ([]api.ValueCount, error) {
	sqlStatement :=
		fmt.Sprintf(`SELECT %s, count(*) count FROM %s
		GROUP BY 1
		ORDER by 1 ASC`, column, table)

	rows, err := s.DB.Query(sqlStatement)
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	defer rows.Close()

	vCounts := []api.ValueCount{}
	for rows.Next() {
		var vCount api.ValueCount
		err = rows.Scan(&vCount.Value, &vCount.Count)
		if err != nil {
			// handle this error
			return nil, err
		}
		vCounts = append(vCounts, vCount)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return vCounts, nil
}

// // StoreKPI stores a new KPI object in the DB
// func (s TrackService) StoreKPI(id int) (api.Track, error) {
// 	sqlStatement :=
// 		`SELECT id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, extra from public.tracks
// 		WHERE id = $1`

// 	var t api.Track

// 	row := s.DB.QueryRow(sqlStatement, id)
// 	switch err := row.Scan(&t.ID, &t.UserID, &t.FpHash, &t.PageURL, &t.PagePath, &t.PageReferrer, &t.PageTitle, &t.Event, &t.CampaignSource, &t.CampaignMedium, &t.CampaignName, &t.CampaignContent, &t.SentAt, &t.Extra); err {
// 	case sql.ErrNoRows:
// 		log.Println("No rows were returned!")
// 	case nil:
// 		log.Printf("FOUND NIL: %v", t)
// 	default:
// 		panic(err)
// 	}

// 	return t, nil
// }

// NewKPIService Creates a new KPIService object
func NewKPIService(host, username, password, dbName string, port int) (*KPIService, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		println(err)
		return nil, err
	}

	return &KPIService{
		db,
	}, nil
}

// StoreKPI stores the track in the db
func (s *KPIService) StoreKPI(kpi api.KPI) (int, error) {
	sqlStatement :=
		`INSERT INTO public.kpis (column_name, value, name, created_at)
	VALUES($1, $2, $3, $4)
	RETURNING id`

	if !kpi.IsValid() {
		return 0, errors.New(missingDataErr)
	}

	id := 0
	// err := s.DB.QueryRow(sqlStatement, t.OwnerID, t.UserID, t.FpHash, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, sentAt, sentAt, t.Extra).Scan(&id)
	err := s.DB.QueryRow(sqlStatement, kpi.Column, kpi.Value, kpi.Name, time.Now().Format(time.RFC3339)).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}
