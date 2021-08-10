package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Backend-GoAPI-server/db"
	"github.com/Backend-GoAPI-server/model"
	"github.com/Backend-GoAPI-server/model/method"
	"github.com/Backend-GoAPI-server/utils"
	"github.com/gorilla/mux"
	"github.com/savsgio/go-logger/v2"
)

type url string

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

// Document API
func Documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "root",
		},
		{
			URL:         url("/login"),
			Method:      "POST",
			Description: "login API",
		},
		{
			URL:         url("/signup"),
			Method:      "POST",
			Description: "signup API",
		},
		{
			URL:         url("/updata"),
			Method:      "PUT",
			Description: "update user API",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

// Login API
func LoginHandle(rw http.ResponseWriter, r *http.Request) {
	// Get data from request body
	var data utils.LoginReq
	err := json.NewDecoder(r.Body).Decode(&data)

	// Body data validation
	if err != nil {
		utils.BadRequestException(rw)
		return
	}

	// Get gorm.DB
	DB, err := db.Start()
	defer DB.Close()
	utils.HandlePanic(err)

	// Find user by id from request body data
	user, err := method.GetUserWithId(DB, data.Id)

	if err != nil {
		utils.NotFoundException(rw)
		return
	}

	// Hashing password
	hashedPw := utils.Hash(data.Password)

	// Password validataion
	if user.Password != hashedPw {
		logger.Error(errors.New("Wrong PW id : " + data.Id))
		utils.ForbiddenException(rw)
		return
	}

	// Generate Access, Refresh Token
	access := utils.AccessToken(data.Id)   // 10 Mins
	refresh := utils.RefreshToken(data.Id) // 14 Days

	// Response Token
	res := utils.LoginRes{
		Status:       200,
		Accesstoken:  access,
		Refreshtoken: refresh,
	}

	utils.MarshalAndRW(200, res, rw)
}

// Signup API
func SignupHandle(rw http.ResponseWriter, r *http.Request) {
	// Get data from request body
	var data utils.SignupReq

	// Body data validation
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		utils.BadRequestException(rw)
		return
	}

	// Get gorm.DB
	DB, err := db.Start()
	defer DB.Close()
	utils.HandlePanic(err)

	// Find user by user ID
	_, err = method.GetUserWithId(DB, data.Id)
	if err == nil {
		utils.BadRequestException(rw)
		return
	}

	// Hash & save password
	data.Password = utils.Hash(data.Password)
	if data.Provider == "" {
		data.Provider = "default"
	}

	err = method.CreateUser(DB, data)
	if err != nil {
		utils.ForbiddenException(rw)
		return
	}
	rw.WriteHeader(201)
}

// Drop out API
func DropoutHandle(rw http.ResponseWriter, r *http.Request) {

	// Get gorm.DB
	DB, err := db.Start()
	defer DB.Close()
	utils.HandlePanic(err)

	val := mux.Vars(r)

	// Find user by id from request body data
	user, err := method.GetUserWithId(DB, val["id"])
	if err != nil {
		utils.NotFoundException(rw)
		return
	}

	err = method.DeleteUserWithId(DB, user.Id)
	if err != nil {
		utils.ForbiddenException(rw)
		return
	}

	rw.WriteHeader(200)
}

// Update User API
func UpdateUserHandle(rw http.ResponseWriter, r *http.Request) {
	// Get data from request body
	var data model.User

	// Body data validation
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		utils.BadRequestException(rw)
		return
	}

	// Get gorm.DB
	DB, err := db.Start()
	defer DB.Close()
	utils.HandlePanic(err)

	err = method.UpdateUser(DB, data)
	if err != nil {
		utils.ForbiddenException(rw)
		return
	}
	rw.WriteHeader(201)
}
