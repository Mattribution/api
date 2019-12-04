package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mattribution/api/pkg/api"
)

type UserService struct {
	DB *sqlx.DB
}

// NewUserService Creates a new User Service object for interracting with user data in Postgres
func NewUserService(host, username, password, dbName string, port int) (*TrackService, error) {
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

// // Store stores the object in the db
// func (s *UserService) Store(u api.User) (int, error) {

// 	sqlStatement :=
// 		`INSERT INTO public.users ()
// 		OUTPUT Inserted.ID
// 		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

// 	// Set default json value (so postgres doesn't get mad)
// 	if t.Extra == "" {
// 		t.Extra = "{}"
// 	}

// 	id := 0
// 	// err := s.DB.QueryRow(sqlStatement, t.OwnerID, t.UserID, t.FpHash, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, sentAt, sentAt, t.Extra).Scan(&id)
// 	err := s.DB.QueryRow(sqlStatement, t.OwnerID, t.UserID, t.FpHash, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, t.SentAt.Format(time.RFC3339), time.Now().Format(time.RFC3339), t.Extra).Scan(&id)
// 	if err != nil {
// 		return id, err
// 	}

// 	return id, nil
// }

// FindByID finds all objects by owner id
func (s UserService) FindByID(id int) (api.User, error) {
	sqlStatement :=
		`SELECT * from public.users
		WHERE id = $1`

	var t api.User

	switch err := s.DB.Get(&t, sqlStatement, id); err {
	case sql.ErrNoRows:
		return api.User{}, errors.New("Not found")
	case nil:
	default:
		panic(err)
	}

	return t, nil
}
