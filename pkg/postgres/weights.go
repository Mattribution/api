package postgres

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattribution/api/pkg/api"
)

type WeightService struct {
	DB *sqlx.DB
}

// Store stores into the db
func (s WeightService) Store(weight api.Weight) (int, error) {
	sqlStatement :=
		`INSERT INTO public.weights (owner_id, model_name, key, value, created_at)
	VALUES(:owner_id, :model_name, :key, :value, :created_at)
	RETURNING id`

	id := 0
	weight.CreatedAt = time.Now()
	err := s.DB.QueryRow(sqlStatement, weight).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// FindByKpiAndModelName finds all by model name for a user
func (s WeightService) FindByKpiAndModelName(ownerID int32, kpiID int32, modelName string) ([]api.Weight, error) {
	sqlStatement :=
		`SELECT * FROM public.weights
		WHERE owner_id = $1
		AND kpi_id = $3
		AND model_name = $2`

	weights := []api.Weight{}
	err := s.DB.Select(&weights, sqlStatement, ownerID, kpiID, modelName)
	if err != nil {
		log.Printf("ERROR: %s", err)
	}

	return weights, err
}
