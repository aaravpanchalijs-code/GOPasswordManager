package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aaravpanchalijs-code/secure-password-manager/database"
	"github.com/aaravpanchalijs-code/secure-password-manager/models"
	"github.com/aaravpanchalijs-code/secure-password-manager/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AddPasswordHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var vault models.Vault

	err := json.NewDecoder(r.Body).Decode(&vault)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if vault.Website == "" ||
		vault.LoginEmail == "" ||
		vault.Password == "" {

		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	userID, err := bson.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusInternalServerError)
		return
	}

	vault.UserID = userID

	encryptedPassword, err := utils.Encrypt(vault.Password)
	if err != nil {
		http.Error(w, "Encryption Error", http.StatusInternalServerError)
		return
	}

	vault.Password = encryptedPassword
	vault.CreatedAt = time.Now()
	vault.UpdatedAt = time.Now()

	_, err = database.VaultCollection.InsertOne(context.TODO(), vault)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password saved successfully",
	})

}

func GetPasswordsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := bson.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusInternalServerError)
		return
	}

	cursor, err := database.VaultCollection.Find(
		context.TODO(),
		bson.M{"user_id": userID})

	if err != nil {

		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	vaults := []models.Vault{}

	err = cursor.All(context.TODO(), &vaults)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	for i := range vaults {

		decryptedPassword, err := utils.Decrypt(vaults[i].Password)
		if err != nil {
			fmt.Println("Decrypt Error:", err)
			http.Error(w, "Decryption Error", http.StatusInternalServerError)
			return
		}

		vaults[i].Password = decryptedPassword
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(vaults)

}

func DeletePasswordHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert JWT UserID (string) to MongoDB ObjectID
	userID, err := bson.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusInternalServerError)
		return
	}

	passwordID := r.URL.Query().Get("id")
	if passwordID == "" {
		http.Error(w, "Password ID is required", http.StatusBadRequest)
		return
	}

	// Convert password ID to ObjectID
	passwordObjectID, err := bson.ObjectIDFromHex(passwordID)
	if err != nil {
		http.Error(w, "Invalid Password ID", http.StatusBadRequest)
		return
	}

	result, err := database.VaultCollection.DeleteOne(
		context.TODO(),
		bson.M{
			"_id":     passwordObjectID,
			"user_id": userID,
		},
	)

	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Password not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password deleted successfully",
	})
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Convert JWT UserID (string) to MongoDB ObjectID
	userID, err := bson.ObjectIDFromHex(claims.UserID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Invalid User ID")
		return
	}

	passwordID := r.URL.Query().Get("id")
	if passwordID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Password ID is required")
		return
	}

	// Convert password ID to ObjectID
	passwordObjectID, err := bson.ObjectIDFromHex(passwordID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid Password ID")
		return
	}

	var vault models.Vault

	err = json.NewDecoder(r.Body).Decode(&vault)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if vault.Website == "" ||
		vault.LoginEmail == "" ||
		vault.Password == "" {

		utils.ErrorResponse(w, http.StatusBadRequest, "All fields are required")
		return
	}

	encryptedPassword, err := utils.Encrypt(vault.Password)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Encryption Error")
		return
	}

	update := bson.M{
		"$set": bson.M{
			"website":     vault.Website,
			"login_email": vault.LoginEmail,
			"password":    encryptedPassword,
			"notes":       vault.Notes,
			"updated_at":  time.Now(),
		},
	}

	result, err := database.VaultCollection.UpdateOne(
		context.TODO(),
		bson.M{
			"_id":     passwordObjectID,
			"user_id": userID,
		},
		update,
	)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Database Error")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(w, http.StatusNotFound, "Password not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password updated successfully",
	})

}

func SharePasswordHandler(w http.ResponseWriter, r *http.Request) {

	type ShareRequest struct {
		PasswordID string `json:"password_id"`
		Username   string `json:"username"`
	}

	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var shareReq ShareRequest

	err := json.NewDecoder(r.Body).Decode(&shareReq)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if shareReq.PasswordID == "" || shareReq.Username == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "All fields are required")
		return
	}

	var receiver models.User
	err = database.UserCollection.FindOne(context.TODO(), bson.M{"username": shareReq.Username}).Decode(&receiver)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	senderID, err := bson.ObjectIDFromHex(claims.UserID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Invalid User ID")
		return
	}

	passwordObjectID, err := bson.ObjectIDFromHex(shareReq.PasswordID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid Password ID")
		return
	}

	var vault models.Vault

	err = database.VaultCollection.FindOne(
		context.TODO(),
		bson.M{
			"_id":     passwordObjectID,
			"user_id": senderID,
		},
	).Decode(&vault)

	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Password not found")
		return
	}

	vault.ID = bson.ObjectID{}
	vault.UserID = receiver.ID
	vault.CreatedAt = time.Now()
	vault.UpdatedAt = time.Now()

	_, err = database.VaultCollection.InsertOne(
		context.TODO(),
		vault,
	)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Database Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password shared successfully",
	})

}
