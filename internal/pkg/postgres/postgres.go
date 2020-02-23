package postgres

import (
	"fmt"
	"time"

	"github.com/mattribution/api/internal/app"

	// Import Postgres SQL driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewCloudSQLClient(dbUser, dbPass, dbName, dbHost string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("user=%s password='%s' host=%s dbname=%s sslmode=disable", dbUser, dbPass, dbHost, dbName)
	return sqlx.Open("postgres", connStr)
}

// ~=~=~=~=~=~=~=~=
// Tracks
// ~=~=~=~=~=~=~=~=

// TracksDAO handles Track data
type TracksDAO struct {
	DB *sqlx.DB
}

func (dao *TracksDAO) Store(t app.Track) (int64, error) {
	sqlStatement :=
		`INSERT INTO public.tracks (owner_id, user_id, anonymous_id, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id`

	var id int64
	err := dao.DB.QueryRow(sqlStatement, t.OwnerID, t.UserID, t.AnonymousID, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, t.SentAt, time.Now().Format(time.RFC3339)).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (dao *TracksDAO) GetNormalizedJourneyAggregate(ownerID string, columnName, conversionColumnName, conversionRowValue string) ([]app.PosAggregate, error) {
	sqlStatement :=
		fmt.Sprintf(
			`SELECT *, count(*)
		FROM (
			SELECT %s as value,
			ROW_NUMBER() OVER (PARTITION BY anonymous_id ORDER BY sent_at) AS position
			FROM tracks AS t
			WHERE %s <> ''
			AND owner_id = $1
			AND t.sent_at < (SELECT sent_at FROM tracks t2 WHERE t2.%s = '%s' AND t.anonymous_id = t2.anonymous_id)
		) as tracks
		GROUP BY position, value
		ORDER BY position;`, columnName, columnName, conversionColumnName, conversionRowValue)

	var posAggregates []app.PosAggregate
	err := dao.DB.Select(&posAggregates, sqlStatement, ownerID)
	if err != nil {
		return nil, err
	}

	return posAggregates, err
}

// ~=~=~=~=~=~=~=~=
// Kpis
// ~=~=~=~=~=~=~=~=

// Kpis handles Kpi data
type KpisDAO struct {
	DB *sqlx.DB
}

func (dao *KpisDAO) Store(kpi app.Kpi) (int64, error) {
	sqlStatement :=
		`INSERT INTO public.kpis (owner_id, name, target, pattern_match_column_name, pattern_match_row_value,  created_at)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id`

	var id int64
	err := dao.DB.QueryRow(sqlStatement, kpi.OwnerID, kpi.Name, kpi.Target, kpi.PatternMatchColumnName, kpi.PatternMatchRowValue, time.Now().Format(time.RFC3339)).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (dao *KpisDAO) FindByOwnerID(ownerID string) ([]app.Kpi, error) {
	sqlStatement :=
		`SELECT * FROM public.kpis 
		WHERE owner_id = $1`

	var kpis []app.Kpi

	err := dao.DB.Select(&kpis, sqlStatement, ownerID)
	if err != nil {
		return nil, err
	}

	return kpis, nil
}

func (dao *KpisDAO) Delete(id int64, ownerID string) (int64, error) {
	sqlStatement :=
		`DELETE FROM public.kpis 
		WHERE id = $1
		AND owner_id = $2`

	res, err := dao.DB.Exec(sqlStatement, id, ownerID)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
