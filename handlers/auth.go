package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aaravpanchalijs-code/secure-password-manager/database"
	"github.com/aaravpanchalijs-code/secure-password-manager/models"
	"github.com/aaravpanchalijs-code/secure-password-manager/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate required fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "All fields are required")
		return
	}

	// -----------------------------
	// CHECK IF EMAIL ALREADY EXISTS
	// -----------------------------
	var existingUser models.User

	err = database.UserCollection.FindOne(
		context.TODO(),
		bson.M{"email": user.Email},
	).Decode(&existingUser)

	if err == nil {
		// User found
		utils.ErrorResponse(w, http.StatusConflict, "Email already exists")
		return
	}

	// Any error other than "not found" is a database error
	if !errors.Is(err, mongo.ErrNoDocuments) {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Database Error")
		return
	}

	// -----------------------------
	// CHECK IF USERNAME EXISTS
	// -----------------------------

	err = database.UserCollection.FindOne(
		context.TODO(),
		bson.M{"username": user.Username},
	).Decode(&existingUser)

	if err == nil {
		utils.ErrorResponse(w, http.StatusConflict, "Username already exists")
		return
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Database Error")
		return
	}

	// -------------------------------------------------
	// Next step:
	// Hash password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	user.Password = string(hashedPassword)

	// Generate UserID

	// Set CreatedAt

	user.CreatedAt = time.Now()
	// Insert user into MongoDB

	_, err = database.UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error inserting user")
		return
	}

	// -------------------------------------------------
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginReq LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	//validate fields

	if loginReq.Email == "" || loginReq.Password == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "All fields are required")
		return
	}

	var user models.User

	err = database.UserCollection.FindOne(context.TODO(), bson.M{"email": loginReq.Email}).Decode(&user)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Database Error")
		return
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(loginReq.Password),
	)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	type LoginResponse struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}

	response := LoginResponse{
		Message: "Login successful",
		Token:   token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
