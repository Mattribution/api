package auth0

import (
	"errors"
	"fmt"
	"log"

	"github.com/mattribution/api/internal/app"
	"gopkg.in/auth0.v3/management"
)

type UsersDAO struct {
	Manager *management.Management
}

func (dao *UsersDAO) FindBySecret(secret string) (*app.User, error) {
	queryStr := fmt.Sprintf(`app_metadata.secret:"%s"`, secret)
	query := management.Query(queryStr)
	users, err := dao.Manager.User.List(query)
	if err != nil {
		return nil, err
	}

	// TODO: Make this print a 4xx error instead of flowing up to a 500
	// this will probably involve creating a custom error object
	if len(users.Users) == 0 {
		return nil, errors.New("No user was found for that secret")
	}

	if len(users.Users) > 1 {
		errStr := "Found multiple users for one secret key"
		// Note: This error is serious af... idk how this could happen
		log.Println(errStr)
		return nil, errors.New(errStr)
	}

	user := users.Users[0]

	return &app.User{
		Name:  *user.Name,
		Email: *user.Email,
		UUID:  user.AppMetadata["uuid"].(string),
	}, nil
}
