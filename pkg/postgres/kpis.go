package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/mattribution/api/pkg/api"
)

type KPIService struct {
	DB *sqlx.DB
}

// Store stores the track in the db
func (s KPIService) Store(kpi api.KPI) (int64, error) {
	sqlStatement :=
		`INSERT INTO public.kpis (column_name, value, name, owner_id, target, created_at)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id`

	if !kpi.IsValid() {
		return 0, errors.New(missingDataErr)
	}

	var id int64
	err := s.DB.QueryRow(sqlStatement, kpi.Column, kpi.Value, kpi.Name, kpi.OwnerID, kpi.Target, time.Now().Format(time.RFC3339)).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// FindByID finds all track objects by owner id
func (s KPIService) FindByID(id int64) (api.KPI, error) {
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
func (s KPIService) Find(ownerID int64) ([]api.KPI, error) {
	sqlStatement :=
		`SELECT * 
		FROM public.kpis
		WHERE owner_id = $1
		`

	kpis := []api.KPI{}
	err := s.DB.Select(&kpis, sqlStatement, ownerID)
	if err != nil {
		return nil, err
	}

	return kpis, nil
}

// Delete removes a single KPI by id
func (s KPIService) Delete(id int64) (int64, error) {
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
