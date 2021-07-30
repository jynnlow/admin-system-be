package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"admin-system-be/constants"
	"admin-system-be/models"
)

type UsersHandler struct {
	ModelsFunc *models.UserCRUDOperationsImpl
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func AttachUsersHandler(subRouter *mux.Router, handlers *UsersHandler) {
	subRouter.HandleFunc("/login", handlers.Login).Methods("POST")
	subRouter.HandleFunc("", handlers.Signup).Methods("POST")
	subRouter.HandleFunc("", handlers.GetAll).Methods("GET")
	subRouter.HandleFunc("", handlers.Update).Methods("PATCH")
	subRouter.HandleFunc("", handlers.Delete).Methods("DELETE")
}

//Login ...
func (u *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	response := &constants.Response{}
	//create a user instance of user struct
	loginUser := &models.User{}
	//Parse and decode the request body into a new user instance

	if err := json.NewDecoder(r.Body).Decode(loginUser); err != nil {
		response.Status = constants.FAIL
		response.Message = "Something wrong"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Println(err)
		}
	}

	//Using input username to get the hashed password in the database
	hashedPassword, err := u.ModelsFunc.GetPwdByUsername(loginUser.Username)
	if err != nil {
		response.Status = constants.FAIL
		response.Message = "Users not found in the system"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	}

	if hashedPassword != nil {
		hashedPasswordByte := []byte(*hashedPassword)
		err = bcrypt.CompareHashAndPassword(hashedPasswordByte, []byte(loginUser.Password))
		if err != nil {
			response.Status = constants.FAIL
			response.Message = "Incorrect Password"
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				return
			}
			_, err = w.Write(jsonResponse)
			if err != nil {
				return
			}
		} else {
			//set the claims
			claim := &Claims{
				Username: loginUser.Username,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(100 * time.Minute).Unix(),
				},
			}
			//create a token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
			tokenKey := fmt.Sprintf("%s", os.Getenv("TOKEN_KEY"))
			//create the JWT token encoded string
			tokenEncodedString, err := token.SignedString([]byte(tokenKey))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//return login status
			response.Status = constants.SUCCESS
			response.Message = tokenEncodedString
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				return
			}
			_, err = w.Write(jsonResponse)
			if err != nil {
				return
			}
		}
	}

}

// Signup ...
func (u *UsersHandler) Signup(w http.ResponseWriter, r *http.Request) {
	response := &constants.Response{}
	//create a user instance of user struct
	user := &models.User{}
	//Parse and decode the request body into a new user instance
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		response.Status = constants.FAIL
		response.Message = "Something went wrong"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	}

	//check if body object is empty, if yes then return
	if user.Username == "" || user.Password == "" {
		response.Status = constants.FAIL
		response.Message = "Username or password is empty"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
		return
	}

	//Salt and hash the password using the bcrypt algorithm
	//The second argument is the cost of hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err = u.ModelsFunc.InsertUsers(user.Username, string(hashedPassword)); err != nil {
		response.Status = constants.FAIL
		response.Message = "Username already exits"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	} else {
		response.Status = constants.SUCCESS
		response.Message = "Insert successfully"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	}

}

// GetAll ...
func (u *UsersHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	response := constants.Response{}
	//get the request token from the header authorization
	requestToken := r.Header.Get("Authorization")
	//check for empty request token
	if requestToken != "" {
		//split the token to remove the bearer string
		splitToken := strings.Split(requestToken, "Bearer ")
		requestToken = splitToken[1]
	} else {
		response.Status = constants.FAIL
		response.Message = "Cannot get request token"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	}

	//parse the token with the claims
	_, err := jwt.ParseWithClaims(
		requestToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(fmt.Sprintf("%s", os.Getenv("TOKEN_KEY"))), nil
		},
	)

	if err != nil {
		response.Status = constants.FAIL
		response.Message = err.Error() + ". Please log in again to access admin dashboard."
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	} else {
		users, err := u.ModelsFunc.GetAllUsers()
		if err != nil {
			response.Status = constants.FAIL
			response.Message = "No user"
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				return
			}
			_, err = w.Write(jsonResponse)
			if err != nil {
				return
			}
		}
		jsonUsers, err := json.Marshal(users)
		if err != nil {
			errString := fmt.Sprintf("%e", err)
			_, err = w.Write([]byte(errString))
			if err != nil {
				return
			}
		}
		_, err = w.Write(jsonUsers)
		if err != nil {
			return
		}
	}
}

// Update ...
func (u *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	response := constants.Response{}
	//create a user instance of user struct
	user := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		response.Status = constants.FAIL
		response.Message = "Something went wrong"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	}

	//hash password before update user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err = u.ModelsFunc.UpdateUsers(user.ID, user.Username, string(hashedPassword)); err != nil {
		response.Status = constants.FAIL
		response.Message = "Something went wrong"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	}

	response.Status = constants.SUCCESS
	response.Message = "Updated Successfully"
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
	return
}

// Delete ...
func (u *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	response := constants.Response{}
	//retrieve parameter from url
	param, ok := r.URL.Query()["id"]
	if !ok || len(param[0]) < 1 {
		response.Status = constants.FAIL
		response.Message = "Url Param 'key' is missing"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
		return
	}

	//retrieve first string element from param(array)
	stringID := param[0]

	//convert string id to uint64 type
	uintID, err := strconv.ParseUint(stringID, 10, 64)
	if err != nil {
		response.Status = constants.FAIL
		response.Message = "Fail to convert string ID to int ID"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
		return
	}

	err = u.ModelsFunc.DeleteUsers(uint(uintID))
	if err != nil {
		response.Status = constants.FAIL
		response.Message = "User not found"
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
		return
	}
	response.Status = constants.SUCCESS
	response.Message = "Deleted Successfully"
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
	return
}
