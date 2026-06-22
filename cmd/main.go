package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aaravpanchalijs-code/secure-password-manager/database"
	"github.com/aaravpanchalijs-code/secure-password-manager/handlers"
	"github.com/aaravpanchalijs-code/secure-password-manager/middleware"
)

// ============================
// CORS Middleware
// ============================

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	database.ConnectDB()

	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	http.Handle(
		"/vault/add",
		middleware.AuthMiddleware(http.HandlerFunc(handlers.AddPasswordHandler)),
	)

	http.Handle(
		"/vault/get",
		middleware.AuthMiddleware(http.HandlerFunc(handlers.GetPasswordsHandler)),
	)

	http.Handle(
		"/vault/delete",
		middleware.AuthMiddleware(http.HandlerFunc(handlers.DeletePasswordHandler)),
	)

	http.Handle(
		"/vault/update",
		middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdatePasswordHandler)),
	)

	http.Handle(
		"/vault/share",
		middleware.AuthMiddleware(http.HandlerFunc(handlers.SharePasswordHandler)),
	)

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	fmt.Println("Server started on port 8080")

	err := http.ListenAndServe(
		":8080",
		enableCORS(http.DefaultServeMux),
	)

	if err != nil {
		log.Fatal(err)
	}
}