package postgres

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"

// 	"github.com/jmoiron/sqlx"
// 	"github.com/mattribution/api/pkg/api"
// )

// type UserService struct {
// 	DB *sqlx.DB
// }

// // NewUserService Creates a new User Service object for interracting with user data in Postgres
// func NewUserService(host, username, password, dbName string, port int) (*TrackService, error) {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
// 		"password=%s dbname=%s sslmode=disable",
// 		host, port, username, password, dbName)
// 	db, err := sqlx.Open("postgres", psqlInfo)
// 	if err != nil {
// 		println(err)
// 		return nil, err
// 	}

// 	return &TrackService{
// 		db,
// 	}, nil
// }

// // FindByID finds all objects by owner id
// func (s UserService) FindByID(id int) (api.User, error) {
// 	sqlStatement :=
// 		`SELECT * from public.users
// 		WHERE id = $1`

// 	var t api.User

// 	switch err := s.DB.Get(&t, sqlStatement, id); err {
// 	case sql.ErrNoRows:
// 		return api.User{}, errors.New("Not found")
// 	case nil:
// 	default:
// 		panic(err)
// 	}

// 	return t, nil
// }
