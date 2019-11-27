package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattribution/api/pkg/api"
)

const (
	missingDataErr = "The payload is missing required data"
	// For SQLx helper
	schema = `
	CREATE TABLE tracks (
		id SERIAL PRIMARY KEY,
		owner_id integer,
		user_id character varying,
		fp_hash character varying,
		page_url character varying,
		page_path character varying,
		page_referrer character varying,
		extra json,
		event character varying,
		ip character varying,
		page_title character varying,
		campaign_source character varying,
		campaign_medium character varying,
		campaign_name character varying,
		campaign_content character varying,
		received_at timestamp without time zone,
		sent_at timestamp without time zone
	);
	
	CREATE TABLE kpis (
		id SERIAL PRIMARY KEY,
		column_name varchar,
		value varchar,
		name varchar,
		created_at timestamp
	);`
)

type TrackService struct {
	DB *sqlx.DB
}

type KPIService struct {
	DB *sqlx.DB
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~
// =~ Tacks
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~

// NewTrackService Creates a new Trackservice object
func NewTrackService(host, username, password, dbName string, port int) (*TrackService, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
	db, err := sqlx.Open("postgres", psqlInfo)
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
		OUTPUT Inserted.ID
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

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
		`SELECT * from public.tracks
		WHERE id = $1`

	var t api.Track

	switch err := s.DB.Get(&t, sqlStatement, id); err {
	case sql.ErrNoRows:
		return api.Track{}, errors.New("Not found")
	case nil:
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

	vCounts := []api.ValueCount{}
	err := s.DB.Select(&vCounts, sqlStatement)
	if err != nil {
		// handle this error better than this
		return nil, err
	}

	return vCounts, nil
}

// GetCountsFromColumn will group a column and get the count of each unique value
func (s TrackService) GetCountsFromColumn(days int, column, table string) ([]api.ValueCount, error) {
	sqlStatement :=
		fmt.Sprintf(`SELECT %s, count(*) count FROM %s
		GROUP BY 1
		ORDER by 1 ASC`, column, table)

	vCounts := []api.ValueCount{}
	err := s.DB.Select(&vCounts, sqlStatement)
	if err != nil {
		// handle this error better than this
		return nil, err
	}

	return vCounts, nil
}

// GetDailyConversionCountForKPI looks through tracks to find daily conversion counts for a KPI
func (s TrackService) GetDailyConversionCountForKPI(kpi api.KPI) ([]api.ValueCount, error) {
	sqlStatement :=
		fmt.Sprintf(`SELECT received_at, count(*) FROM tracks 
	WHERE %s = $1
	GROUP BY 1
	ORDER BY 1 asc;`, kpi.Column)

	vCounts := []api.ValueCount{}
	err := s.DB.Select(&vCounts, sqlStatement)
	if err != nil {
		// handle this error better than this
		return nil, err
	}

	return vCounts, nil
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~
// =~ KPIs
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~

// NewKPIService Creates a new KPIService object
func NewKPIService(host, username, password, dbName string, port int) (*KPIService, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
	db, err := sqlx.Open("postgres", psqlInfo)
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

// FindByID finds all track objects by owner id
func (s KPIService) FindByID(id int) (api.KPI, error) {
	sqlStatement :=
		`SELECT * FROM public.kpis
		WHERE id = $1`

	var kpi api.KPI

	switch err := s.DB.Get(&kpi, sqlStatement, id); err {
	case sql.ErrNoRows:
		return api.KPI{}, errors.New("Not found")
	case nil:
	default:
		return api.KPI{}, err
	}

	return kpi, nil
}

// Find queries for all kpis (no filter)
func (s KPIService) Find() ([]api.KPI, error) {
	sqlStatement :=
		`SELECT * FROM public.kpis`

	var kpis []api.KPI
	err := s.DB.Select(&kpis, sqlStatement)
	if err != nil {
		// handle this error better than this
		return nil, err
	}

	return kpis, nil
}
