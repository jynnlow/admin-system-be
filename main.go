package main

import (
	"./database"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	db, err := database.NewDB()
	if err != nil {
		return 
	}

	mux := http.NewServeMux()

	//Open http server
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		response := &Response{}
		if r.Method == "POST" {
			//create a user instance of user struct
			user := &User{}
			//Parse and decode the request body into a new user instance
			if err := json.NewDecoder(r.Body).Decode(user); err != nil {
				response.Status = "FAIL"
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
			if user.Username == "" || user.Password == ""{
				response.Status = "FAIL"
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
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password),8)
			if err = db.InsertUsers(user.Username, string(hashedPassword)); err != nil {
				response.Status = "FAIL"
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

			response.Status = "SUCCESS"
			response.Message = "Insert successfully"
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				return
			}
			_, err = w.Write(jsonResponse)
			if err != nil {
				return
			}

		} else if r.Method == "GET" {
			//get the request token from the header authorization
			requestToken := r.Header.Get("Authorization")
			//check for empty request token
			if requestToken != "" {
				//split the token to remove the bearer string
				splitToken := strings.Split(requestToken, "Bearer ")
				requestToken = splitToken[1]
			}else{
				response.Status = "FAIL"
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
					return []byte(("my_secret_key")),nil
				},
			)

			if err != nil {
				response.Status = "FAIL"
				response.Message = err.Error() + ". Please log in again to access admin dashboard."
				jsonResponse, err := json.Marshal(response)
				if err != nil {
					return
				}
				_, err = w.Write(jsonResponse)
				if err != nil {
					return
				}
			}else{
				users, err := db.GetUsers()
				if err != nil {
					response.Status = "FAIL"
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
		}else if r.Method == "PATCH"{
			//create a user instance of user struct
			user := &User{}
			if err := json.NewDecoder(r.Body).Decode(user); err != nil {
				response.Status = "FAIL"
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
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password),8)
			if err = db.UpdateUser(user.ID,user.Username, string(hashedPassword)); err != nil {
				response.Status = "FAIL"
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

			response.Status = "SUCCESS"
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

		}else if r.Method == "DELETE" {
			//retrieve parameter from url
			param, ok := r.URL.Query()["id"]
			if !ok || len(param[0]) < 1 {
				response.Status = "FAIL"
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
			fmt.Println()

			//convert string id to int id
			intID, err := strconv.Atoi(stringID)
			if err != nil {
				response.Status = "FAIL"
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

			err = db.DeleteUser(intID)
			if err != nil {
				response.Status = "FAIL"
				response.Message = "Failed to delete"
				jsonResponse, err := json.Marshal(response)
				if err != nil {
					return
				}
				_, err = w.Write(jsonResponse)
				if err != nil {
					return
				}
				return
			}else {
				response.Status = "SUCCESS"
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
		}
	})

	mux.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request){
		response := &Response{}
		if r.Method == "POST"{
			//create a user instance of user struct
			loginUser := &User{}
			//Parse and decode the request body into a new user instance
			if err := json.NewDecoder(r.Body).Decode(loginUser); err != nil {
				response.Status = "FAIL"
				response.Message = "Something wrong"
				jsonResponse, err := json.Marshal(response)
				if err != nil {
					return
				}
				_, err = w.Write(jsonResponse)
				if err != nil {
					return
				}
			}

			//Using input username to get the hashed password in the database
			hashedPassword, err := db.GetUserPasswordByUsername(loginUser.Username)
			if err != nil {
				response.Status = "FAIL"
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
					response.Status = "FAIL"
					response.Message = "Incorrect Password"
					jsonResponse, err := json.Marshal(response)
					if err != nil {
						return
					}
					_, err = w.Write(jsonResponse)
					if err != nil {
						return
					}
				}else{
					//set the claims
					claim := &Claims{
						Username: loginUser.Username,
						StandardClaims: jwt.StandardClaims{
							ExpiresAt: time.Now().Add(100 * time.Minute).Unix(),
						},
					}
					//create a token
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
					//create the JWT token encoded string
					tokenEncodedString, err := token.SignedString([]byte("my_secret_key"))
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					//Encode with json
					//jsonTokenResponse, err := json.Marshal(tokenEncodedString)
					//if err != nil {
					//	return
					//}
					//return login status
					response.Status = "SUCCESS"
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
	})

	handler := cors.AllowAll().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

