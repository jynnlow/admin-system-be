package handlers

import (
	"admin-system-be/entity"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"admin-system-be/constants"
	"admin-system-be/dto"
	"admin-system-be/helpers"
	"admin-system-be/models"
)

type UsersHandler struct {
	ModelsFunc *models.UserCRUDOperationsImpl
	SecretFunc *models.SecretOperationsImpl
}

type Admin struct {
	LoginKey string `json:"login_key"`
}

func AttachUsersHandler(subRouter *mux.Router, handlers *UsersHandler) {
	subRouter.HandleFunc("/login", handlers.Login).Methods(constants.POST)
	subRouter.HandleFunc("", handlers.Signup).Methods(constants.POST)
	subRouter.HandleFunc("", handlers.GetAll).Methods(constants.GET)
	subRouter.HandleFunc("", handlers.Update).Methods(constants.PATCH)
	subRouter.HandleFunc("", handlers.Delete).Methods(constants.DELETE)
}

// Login ...
func (u *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginUser := &dto.UserRequest{}
	// Parse and decode the request body into a new user instance
	err := json.NewDecoder(r.Body).Decode(loginUser)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	// check if the username or password is empty
	err = loginUser.Validate()
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"username or password cannot be empty",
			w,
		)
		return
	}

	//get user from db
	user, err := u.ModelsFunc.GetByUsername(loginUser.Username)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"user not found in the system",
			w,
		)
		return
	}

	//check if the user is approved to login
	if user.Approved != true {
		err = helpers.JsonResponse(
			constants.FAIL,
			"please wait for admin to approve your registration",
			w,
		)
		return
	}

	// check if the input is matched with the password in db
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"incorrect password",
			w,
		)
		return
	}

	//retrieve token key from db
	tokenKey, err := u.SecretFunc.GetTokenKey()
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"token key does not exist",
			w,
		)
		return
	}

	//create token
	token := helpers.NewClaim(
		user.ID,
		user.Username,
		user.Role,
	)

	tokenEncodedString, err := token.CreateToken(tokenKey)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"failed to create token",
			w,
		)
		return
	}

	//respond when successful
	err = helpers.JsonResponse(
		constants.SUCCESS,
		tokenEncodedString,
		w,
	)
}

// Signup ...
func (u *UsersHandler) Signup(w http.ResponseWriter, r *http.Request) {
	registrationReq := &dto.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(registrationReq)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	err = registrationReq.Validate()
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"username or password cannot be empty",
			w,
		)
		return
	}

	//retrieve admin login key from db
	adminLoginKey, err := u.SecretFunc.GetAdminLoginKey()
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"admin login key not exists",
			w,
		)
		return
	}

	//determine the role of user
	user, err := entity.NewUserEntity(
		registrationReq.Username,
		registrationReq.Password,
		registrationReq.SecretKey,
		adminLoginKey,
	)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"wrong admin secret key",
			w,
		)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	userModel := &models.User{
		Username: user.Username,
		Password: string(hashedPassword),
		Role:     user.Role,
		Approved: user.Approved,
	}

	err = u.ModelsFunc.Insert(userModel)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	err = helpers.JsonResponse(
		constants.SUCCESS,
		"Insert successfully",
		w,
	)
}

// GetAll ...
func (u *UsersHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	claims := &helpers.Claims{}
	requestToken, err := claims.GetToken(r)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	// get token key from db
	tokenKey, err := u.SecretFunc.GetTokenKey()
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"token key does not exists",
			w,
		)
		return
	}

	//verify token
	verifiedToken, err := claims.VerifyToken(requestToken, tokenKey)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	if verifiedToken.Role != "admin" {
		err = helpers.JsonResponse(
			constants.FAIL,
			"please look for an admin to help out",
			w,
		)
		return
	}

	users, err := u.ModelsFunc.GetAll()
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	err = helpers.JsonUserList(users, w)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}
}

// Update ...
func (u *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := &helpers.Claims{}
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	var hashedPassword []byte
	if user.Password != "" {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			err = helpers.JsonResponse(
				constants.FAIL,
				err.Error(),
				w,
			)
			return
		}
		user.Password = string(hashedPassword)
	}

	requestToken, err := claims.GetToken(r)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	// get token key from db
	tokenKey, err := u.SecretFunc.GetTokenKey()
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"token key does not exists",
			w,
		)
		return
	}

	//verify token
	verifiedToken, err := claims.VerifyToken(requestToken, tokenKey)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	//if it is user, only retrieve its own data
	if verifiedToken.Role == "user" {
		user.ID = verifiedToken.Id
	}

	err = u.ModelsFunc.Update(user.ID, user.Username, user.Password)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	err = helpers.JsonResponse(
		constants.SUCCESS,
		"Update successfully",
		w,
	)
}

// Delete ...
func (u *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := &helpers.Claims{}
	requestToken, err := claims.GetToken(r)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	// get token key from db
	tokenKey, err := u.SecretFunc.GetTokenKey()
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			"token key does not exists",
			w,
		)
		return
	}

	//verify token
	verifiedToken, err := claims.VerifyToken(requestToken, tokenKey)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	if verifiedToken.Role != "admin" {
		err = helpers.JsonResponse(
			constants.FAIL,
			"only admin is allow to delete user data",
			w,
		)
		return
	}

	//retrieve parameter from url
	param, ok := r.URL.Query()["id"]
	if !ok || len(param[0]) < 1 {
		_ = helpers.JsonResponse(
			constants.FAIL,
			"Url Param 'key' is missing",
			w,
		)
		return
	}

	// convert id to uint64 type
	uintID, err := strconv.ParseUint(param[0], 10, 64)
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	err = u.ModelsFunc.Delete(uint(uintID))
	if err != nil {
		err = helpers.JsonResponse(
			constants.FAIL,
			err.Error(),
			w,
		)
		return
	}

	err = helpers.JsonResponse(
		constants.SUCCESS,
		"Delete successfully",
		w,
	)
}
