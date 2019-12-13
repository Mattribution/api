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
	VALUES(:user_id, :amount, :created_at)
	RETURNING id`

	id := 0
	billingEvent.CreatedAt = time.Now()
	err := s.DB.QueryRow(sqlStatement, billingEvent).Scan(&id)
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
