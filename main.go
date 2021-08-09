package main

import (
	"admin-system-be/handlers"
	"admin-system-be/models"
	"log"
	"net/http"

	_ "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"admin-system-be/infrastructure"
)

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Connect database
	dbConn, err := infrastructure.ConnectDB()
	if err != nil || dbConn == nil {
		log.Fatal(err)
	}

	// Models registration
	userCrudOperation := &models.UserCRUDOperationsImpl{
		DbConn: dbConn,
	}

	//
	secretOperation := &models.SecretOperationsImpl{
		DbConn: dbConn,
	}

	// Handlers registration
	userHandlers := &handlers.UsersHandler{
		ModelsFunc: userCrudOperation,
		SecretFunc: secretOperation,
	}

	//sct, err := secretOperation.GetTokenKey()
	//if err != nil {
	//	return
	//} else {
	//	fmt.Println(*sct)
	//}

	// Router registration
	router := mux.NewRouter()
	handlers.AttachUsersHandler(router.PathPrefix("/users").Subrouter(), userHandlers)

	log.Fatal(http.ListenAndServe(":8080", router))
}
