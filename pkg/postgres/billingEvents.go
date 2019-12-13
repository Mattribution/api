package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/mattribution/api/pkg/api"
)

type BillingEventService struct {
	DB *sqlx.DB
}

// Store stores the track in the db
func (s BillingEventService) Store(billingEvent api.BillingEvent) (int, error) {
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
