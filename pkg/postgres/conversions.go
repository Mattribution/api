package postgres

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattribution/api/pkg/api"
)

type ConversionService struct {
	DB *sqlx.DB
}

// Store stores the conversion in the db
func (s ConversionService) Store(conversion api.Conversion) (int, error) {
	sqlStatement :=
		`INSERT INTO public.conversions (owner_id, track_id, kpi_id, created_at)
	VALUES($1, $2, $3, $4)
	RETURNING id`

	id := 0
	err := s.DB.QueryRow(sqlStatement, conversion.OwnerID, conversion.TrackID, conversion.KPIID, time.Now().Format(time.RFC3339)).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// Find all conversions for a user
func (s ConversionService) Find(ownerID int) ([]api.Conversion, error) {
	sqlStatement :=
		`SELECT * FROM public.conversions
		WHERE owner_id = $1`

	conversions := []api.Conversion{}
	err := s.DB.Select(&conversions, sqlStatement, ownerID)
	if err != nil {
		log.Printf("ERROR: %s", err)
	}

	return conversions, err
}

// Delete removes a single conversion by id
func (s ConversionService) Delete(id int, ownerID int) (int64, error) {
	sqlStatement :=
		`DELETE 
		FROM public.conversions 
		WHERE id = $1
		AND owner_id = $2`

	// TODO: Get remove count from this?
	res, err := s.DB.Exec(sqlStatement, id, ownerID)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
