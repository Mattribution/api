package auth0

import (
	"errors"
	"fmt"

	"github.com/mattribution/api/internal/app"
	"gopkg.in/auth0.v3/management"
)

type UsersDAO struct {
	Manager *management.Management
}

func (dao *UsersDAO) FindBySecret(secret string) ([]app.User, error) {
	queryStr := fmt.Sprintf(`app_metadata.secret:"%s"`, secret)
	query := management.Query(queryStr)
	usersList, err := dao.Manager.User.List(query)
	if err != nil {
		return nil, err
	}

	users := []app.User{}
	for _, user := range usersList.Users {
		newUser := app.User{}
		if user.Name != nil {
			newUser.Name = *user.Name
		}
		if user.Email != nil {
			newUser.Email = *user.Email
		}
		uuid, ok := user.AppMetadata["uuid"].(string)
		if !ok {
			return nil, errors.New("Error parsing user data from backend")
		}
		newUser.UUID = uuid
		users = append(users, newUser)
	}

	return users, nil
}
