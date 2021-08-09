package helpers

import (
	"admin-system-be/models"
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserList struct {
	users []*models.User
}

func JsonResponse(status, message string, w http.ResponseWriter) error {
	response := &Response{
		Status:  status,
		Message: message,
	}

	jsonRes, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonRes)
	if err != nil {
		return err
	}

	return nil
}

func JsonUserList(users []*models.User, w http.ResponseWriter) error {
	jsonRes, err := json.Marshal(users)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonRes)
	if err != nil {
		return err
	}

	return nil

}
