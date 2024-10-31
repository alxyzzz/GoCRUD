package database

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type id uuid.UUID

func (i id) MarshalText() (text []byte, err error) {
	return []byte(uuid.UUID(i).String()), nil
}

type userErrorCode int

type ErrorUserNotFound struct {
	msg string
}

type ErrorUserWrongData struct {
	Code userErrorCode
	msg  string
}

func (e ErrorUserNotFound) Error() string { return e.msg }

func (e ErrorUserWrongData) Error() string { return e.msg }

type User struct {
	Id        id     `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Biography string `json:"biography"`
}

type application struct {
	data map[id]User
}

var Application application

func InitializeDatabase() {
	adminId, _ := uuid.NewRandom()
	initialData := map[id]User{
		id(adminId): {
			Id:        id(adminId),
			FirstName: "admin",
			LastName:  "admin",
			Biography: "soy admin",
		},
	}

	Application = application{data: initialData}
}

func Insert(data []byte) (User, error) {
	var newUser User
	if err := json.Unmarshal(data, &newUser); err != nil {
		return User{}, fmt.Errorf("error when unmarshal data: %w", err)
	}

	if err := checkUserData(newUser); err != nil {
		return User{}, fmt.Errorf("wrong user data: %w", err)
	}

	newId, err := uuid.NewRandom()

	if err != nil {
		return User{}, fmt.Errorf("failed to create UUID: %w", err)
	}

	newUser.Id = id(newId)
	Application.data[id(newId)] = newUser
	return newUser, nil
}

func FindAll() (map[id]User, error) {
	return Application.data, nil
}

func FindByID(idStr string) (User, error) {
	userId, err := uuid.Parse(idStr)
	if err != nil {
		return User{}, fmt.Errorf("user id format is invalid: %w", err)
	}
	user, ok := Application.data[id(userId)]
	if !ok {
		return User{}, ErrorUserNotFound{msg: "user not found"}
	}
	return user, nil
}

func Update(idStr string, data []byte) (User, error) {
	var newUserData User
	if err := json.Unmarshal(data, &newUserData); err != nil {
		return User{}, fmt.Errorf("error when unmarshal data: %w", err)
	}

	if err := checkUserData(newUserData); err != nil {
		return User{}, fmt.Errorf("wrong user data: %w", err)
	}

	userId, err := uuid.Parse(idStr)
	if err != nil {
		return User{}, fmt.Errorf("user id format is invalid: %w", err)
	}

	user, ok := Application.data[id(userId)]
	if !ok {
		return User{}, ErrorUserNotFound{msg: "user not found"}
	}
	Application.data[id(userId)] = newUserData
	return user, nil
}

func Delete(idStr string) error {
	userId, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("user id format is invalid: %w", err)
	}

	_, ok := Application.data[id(userId)]

	if !ok {
		return ErrorUserNotFound{msg: "user not found"}
	}

	delete(Application.data, id(userId))
	return nil
}

func checkUserData(user User) error {
	if user.FirstName == "" {
		return ErrorUserWrongData{msg: "empty first name", Code: 1}
	} else if user.LastName == "" {
		return ErrorUserWrongData{msg: "empty last name", Code: 2}
	} else if user.Biography == "" {
		return ErrorUserWrongData{msg: "empty biography", Code: 3}
	}

	if len(user.FirstName) < 2 || len(user.FirstName) > 20 {
		return ErrorUserWrongData{msg: "invalid first name", Code: 4}
	} else if len(user.LastName) < 2 || len(user.LastName) > 20 {
		return ErrorUserWrongData{msg: "invalid last name", Code: 5}
	} else if len(user.Biography) < 20 || len(user.Biography) > 450 {
		return ErrorUserWrongData{msg: "invalid biography", Code: 6}
	}

	return nil
}
