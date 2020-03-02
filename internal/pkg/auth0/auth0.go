package auth0

import (
	"log"
	"errors"
	"gopkg.in/auth0.v3/management"
	"fmt"
	"github.com/mattribution/api/internal/app"
)

type UsersDAO struct {
	UserManager *management.UserManager
}

func (dao *UsersDAO) FindBySecret(secret string) (*app.User, error) {
	users, err := dao.UserManager.Search(management.ListOption(management.Query(fmt.Sprintf(`app_metadata.secret:"%s"`, secret))))
	if err != nil {
		return nil, err
	}
	
	if users.Length > 1 {
		errStr := "Found multiple users for one secret key"
		// Note: This error is serious af... idk how this could happen
		log.Println(errStr)
		return nil, errors.New(errStr)
	}
	
	user := users.Users[0]
	log.Printf("%+v\n", user)

	return &app.User{
		Name: *user.Name,
		Email: *user.Email,
		UUID: user.AppMetadata["uuid"].(string),
	}, nil
}
