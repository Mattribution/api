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

type TrackService struct {
	DB *sqlx.DB
}

// Store stores the track in the db
func (s TrackService) Store(t api.Track) (int, error) {
	sqlStatement :=
		`INSERT INTO public.tracks (owner_id, user_id, anonymous_id, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at, extra)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id`

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
	WHERE campaign_name <> ''
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
