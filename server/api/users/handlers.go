package users

import (
	"errors"
	"github.com/raunofreiberg/kyrene/server/database"
	"github.com/raunofreiberg/kyrene/server/model"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-pg/pg"
)

func CreateUser(name string, email string, hashedPassword []uint8, isFacebook bool) (interface{}, error) {
	user := database.User{
		Name:					name,
		Email:				email,
		Password:			hashedPassword,
		IsFacebook:		isFacebook,
	}

	exists, err := database.DB.Model(&user).Where("email = ?", email).Exists()
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("Este email ya esta registrado, intenta con otro")
	}

	if _, err := database.DB.Model(&user).Returning("id").Insert(); err != nil {
		return nil, err
	}

	role := model.Role{
		UserId:		user.Id,
		Name:			user.Name,
	}

	if _, err := database.DB.Model(&role).Insert(); err != nil {
		return nil, err
	}

	return model.User{
		Id:						user.Id,
		Name:					user.Name,
		Email:				user.Email,
		IsFacebook:		user.IsFacebook,
	}, nil
}

func QueryUser(email string) (interface{}, error) {
	user := database.User{}

	_, err := database.DB.QueryOne(
		&user,
		"SELECT * FROM users WHERE email = ?", email,
	)

	if err != nil {
		return nil, err
	}

	return model.User{
		Id:						user.Id,
		Name:					user.Name,
		Email:				user.Email,
		IsFacebook:		user.IsFacebook,
		IsGoogle:			user.IsGoogle,
	}, nil
}

func QueryUserById(userID int) (interface{}, error) {
	if userID == 0 {
		return nil, errors.New("This query requires a userID param")
	}

	user := database.User{}

	_, err := database.DB.QueryOne(
		&user,
		"SELECT * FROM users WHERE id = ?", userID,
	)

	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New("This user doesn't exist")
		}
		return nil, err
	}

	return model.User{
		Id:				user.Id,
		Name:			user.Name,
		Email:		user.Email,
	}, nil
}

func QueryUsers() (interface{}, error) {
	var users []model.User
	var dbUsers []database.User

	err := database.DB.Model(&dbUsers).Select()

	if err != nil {
		return nil, err
	}

	for _, user := range dbUsers {
		users = append(users, model.User{
			Id:				user.Id,
			Name:			user.Name,
			Email:		user.Email,
		})
	}

	return users, nil
}

func QueryRole(userId int) (interface{}, error) {
	if userId == 0 {
		return nil, errors.New("This query requires a userID param")
	}

	role := database.Role{}

	_, err := database.DB.QueryOne(
		&role,
		"SELECT * FROM roles WHERE user_id = ?",
		userId,
	)

	if err != nil {
		if err == pg.ErrNoRows {
			return model.Role{}, nil
		}
		return nil, err
	}

	return model.Role{
		Id:					role.Id,
		UserId:			role.UserId,
		Name:				role.Name,
		IsAdmin:		role.IsAdmin,
		IsGod:			role.IsGod,
	}, nil
}

func IsAuthenticated(email string, password []byte) (bool, error) {
	user := database.User{}

	_, err := database.DB.QueryOne(
		&user,
		"SELECT password FROM users WHERE email = ?",
		email,
	)

	if err != nil {
		return false, nil
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, password); err != nil {
		return false, errors.New("Incorrect password")
	}

	return true, nil
}
