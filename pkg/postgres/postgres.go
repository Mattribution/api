package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattribution/api/pkg/api"
)

const (
	missingDataErr = "The payload is missing required data"
	mockOwnerID    = 1
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

type BillingEventService struct {
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

// Store stores the track in the db
func (s *TrackService) Store(t api.Track) (int, error) {
	// sqlStatement := `INSERT INTO tracks (owner_id, user_id, anonymous_id, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at, extra)
	// VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlStatement :=
		`INSERT INTO public.tracks (owner_id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at, extra)
		OUTPUT Inserted.ID
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	// Set default json value (so postgres doesn't get mad)
	if t.Extra == "" {
		t.Extra = "{}"
	}

	id := 0
	err := s.DB.QueryRow(sqlStatement, t.OwnerID, t.UserID, t.AnonymousID, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, t.SentAt.Format(time.RFC3339), time.Now().Format(time.RFC3339), t.Extra).Scan(&id)
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

// GetCampaigns infers campaign names from tracking data
// 	TODO: When adding custom campaigns, this should return an api.Campaign object
//  that signifies it is an inferred campaign
func (s TrackService) GetCampaigns() ([]string, error) {
	sqlStatement := `SELECT campaign_name FROM tracks
	GROUP BY 1`

	campaigns := []string{}
	err := s.DB.Select(&campaigns, sqlStatement)
	if err != nil {
		// handle this error better than this
		return nil, err
	}

	return campaigns, nil
}

// GetTopValuesFromColumn gets a theGetTopValuesFromColumn top values from a column along with their counts
// TODO: implement daily limit
func (s TrackService) GetTopValuesFromColumn(days int, column, table string, extraWheres string) ([]api.ValueCount, error) {
	sqlStatement :=
		fmt.Sprintf(`SELECT %s as value, count(*) count FROM %s
		%s
		GROUP BY 1
		ORDER BY 2 DESC
		LIMIT 10;`, column, table, extraWheres)

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
		fmt.Sprintf(`SELECT %s as value, count(*) count FROM %s
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
		fmt.Sprintf(`SELECT date_trunc('day', received_at) as value, count(*) FROM tracks 
	WHERE %s = $1
	GROUP BY 1
	ORDER BY 1 asc;`, kpi.Column)

	vCounts := []api.ValueCount{}
	err := s.DB.Select(&vCounts, sqlStatement, kpi.Value)
	if err != nil {
		return nil, err
	}

	return vCounts, nil
}

func (s TrackService) GetFirstTouchForKPI(kpi api.KPI) ([]api.ValueCount, error) {
	sqlStatement := fmt.Sprintf(`
	SELECT campaign_name as value, count(*) as count
	FROM tracks t
	WHERE EXISTS (
		SELECT 1 
		FROM tracks t2
		WHERE t2.%s = $1
		AND t.anonymous_id = t2.anonymous_id
	)
	AND NOT EXISTS (
		SELECT * 
		FROM tracks t2
		WHERE t.anonymous_id = t2.anonymous_id
		AND t2.received_at < t.received_at
	)
	GROUP BY 1
	ORDER BY 2 DESC;`, kpi.Column)

	vCounts := []api.ValueCount{}
	err := s.DB.Select(&vCounts, sqlStatement, kpi.Value)
	if err != nil {
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

// Store stores the track in the db
func (s *KPIService) Store(kpi api.KPI) (int, error) {
	// TODO
	kpi.OwnerID = mockOwnerID

	sqlStatement :=
		`INSERT INTO public.kpis (column_name, value, name, owner_id, target, created_at)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id`

	if !kpi.IsValid() {
		return 0, errors.New(missingDataErr)
	}

	id := 0
	err := s.DB.QueryRow(sqlStatement, kpi.Column, kpi.Value, kpi.Name, kpi.OwnerID, kpi.Target, time.Now().Format(time.RFC3339)).Scan(&id)
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

	kpis := []api.KPI{}
	err := s.DB.Select(&kpis, sqlStatement)
	if err != nil {
		return nil, err
	}

	return kpis, nil
}

// Delete removes a single KPI by id
func (s KPIService) Delete(id int) (int64, error) {
	sqlStatement :=
		`DELETE FROM public.kpis WHERE id = $1`

	// TODO: Get remove count from this?
	res, err := s.DB.Exec(sqlStatement, id)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~
// =~ BillingEvents
// =~ Note: this section may be removed
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~

// NewBillingEventService Creates a new BillingEventService object
func NewBillingEventService(host, username, password, dbName string, port int) (*BillingEventService, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		println(err)
		return nil, err
	}

	return &BillingEventService{
		db,
	}, nil
}

// Store stores the track in the db
func (s *BillingEventService) Store(billingEvent api.BillingEvent) (int, error) {
	sqlStatement :=
		`INSERT INTO public.billing_events (user_id, amount, created_at)
	VALUES($1, $2, $3)
	RETURNING id`

	id := 0
	err := s.DB.QueryRow(sqlStatement, billingEvent.UserID, billingEvent.Amount, time.Now().Format(time.RFC3339)).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// FindByUserID finds all billing events by user_id
func (s BillingEventService) FindByUserID(userID int) (api.BillingEvent, error) {
	sqlStatement :=
		`SELECT * FROM public.billing_events
		WHERE user_id = $1`

	var billingEvent api.BillingEvent

	switch err := s.DB.Get(&billingEvent, sqlStatement, userID); err {
	case sql.ErrNoRows:
		return api.BillingEvent{}, errors.New("Not found")
	case nil:
	default:
		return api.BillingEvent{}, err
	}

	return billingEvent, nil
}
