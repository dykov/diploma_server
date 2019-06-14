package main

import (
	"./auth"
	"./domain"
	"./util"
	"github.com/gorilla/mux"
	"net/http"
)

var router *mux.Router

func getRouter() {

	router = mux.NewRouter()

	router.Path("/auth/login").Handler(
		auth.Middleware(http.HandlerFunc(auth.Login), 0),
	).Methods("POST")

	router.Path("/auth/refresh").Handler(
		auth.Middleware(http.HandlerFunc(auth.RefreshToken), util.GrantUser, auth.AuthMiddleware),
	).Methods("POST")

	router.Path("/user/{id}/verification").Queries("code", "{code_value}").Handler(
		auth.Middleware(http.HandlerFunc(domain.UserVerification), 0),
	).Methods("GET")

	router.Path("/user/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.GetUser), util.GrantUser, auth.AuthMiddleware),
	).Methods("GET")

	router.Path("/user").Handler(
		auth.Middleware(http.HandlerFunc(domain.CreateUser), 0),
	).Methods("POST")

	router.Path("/user/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.UpdateUser), util.GrantUser, auth.AuthMiddleware),
	).Methods("PUT")

	router.Path("/user/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.DeleteUser), util.GrantUser, auth.AuthMiddleware),
	).Methods("DELETE")

	router.Path("/user/{id}/answer").Handler(
		auth.Middleware(http.HandlerFunc(domain.GetUsersTests), util.GrantUser, auth.AuthMiddleware),
	).Methods("GET")

	router.Path("/user/{id}/answer").Handler(
		auth.Middleware(http.HandlerFunc(domain.AddUsersTest), util.GrantUser, auth.AuthMiddleware),
	).Methods("POST")

	router.Path("/course").Handler(
		auth.Middleware(http.HandlerFunc(domain.GetCourses), util.GrantUser, auth.AuthMiddleware),
	).Methods("GET")

	router.Path("/course").Handler(
		auth.Middleware(http.HandlerFunc(domain.CreateCourse), util.GrantModer, auth.AuthMiddleware),
	).Methods("POST")

	router.Path("/course/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.UpdateCourse), util.GrantModer, auth.AuthMiddleware),
	).Methods("PUT")

	router.Path("/course/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.DeleteCourse), util.GrantModer, auth.AuthMiddleware),
	).Methods("DELETE")

	router.Path("/course/{id}/section").Handler(
		auth.Middleware(http.HandlerFunc(domain.GetSections), util.GrantUser, auth.AuthMiddleware),
	).Methods("GET")

	router.Path("/section").Handler(
		auth.Middleware(http.HandlerFunc(domain.CreateSection), util.GrantModer, auth.AuthMiddleware),
	).Methods("POST")

	router.Path("/section/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.UpdateSection), util.GrantModer, auth.AuthMiddleware),
	).Methods("PUT")

	router.Path("/section/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.DeleteSection), util.GrantModer, auth.AuthMiddleware),
	).Methods("DELETE")

	router.Path("/section/{id}/lesson").Handler(
		auth.Middleware(http.HandlerFunc(domain.GetLessons), util.GrantUser, auth.AuthMiddleware),
	).Methods("GET")

	router.Path("/lesson").Handler(
		auth.Middleware(http.HandlerFunc(domain.CreateLesson), util.GrantModer, auth.AuthMiddleware),
	).Methods("POST")

	router.Path("/lesson/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.UpdateLesson), util.GrantModer, auth.AuthMiddleware),
	).Methods("PUT")

	router.Path("/lesson/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.DeleteLesson), util.GrantModer, auth.AuthMiddleware),
	).Methods("DELETE")

	router.Path("/lesson/{id}/paragraph").Handler(
		auth.Middleware(http.HandlerFunc(domain.GetParagraphs), util.GrantUser, auth.AuthMiddleware),
	).Methods("GET")

	router.Path("/paragraph").Handler(
		auth.Middleware(http.HandlerFunc(domain.CreateParagraph), util.GrantModer, auth.AuthMiddleware),
	).Methods("POST")

	router.Path("/paragraph/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.UpdateParagraph), util.GrantModer, auth.AuthMiddleware),
	).Methods("PUT")

	router.Path("/paragraph/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.DeleteParagraph), util.GrantModer, auth.AuthMiddleware),
	).Methods("DELETE")

	router.Path("/paragraph/{id}/answer").Handler(
		auth.Middleware(http.HandlerFunc(domain.GetAnswers), util.GrantUser, auth.AuthMiddleware),
	).Methods("GET")

	router.Path("/answer").Handler(
		auth.Middleware(http.HandlerFunc(domain.CreateAnswer), util.GrantModer, auth.AuthMiddleware),
	).Methods("POST")

	router.Path("/answer/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.UpdateAnswer), util.GrantModer, auth.AuthMiddleware),
	).Methods("PUT")

	router.Path("/answer/{id}").Handler(
		auth.Middleware(http.HandlerFunc(domain.DeleteAnswer), util.GrantModer, auth.AuthMiddleware),
	).Methods("DELETE")

	router.Path("/admin/backup").Handler(
		auth.Middleware(http.HandlerFunc(util.Backup), util.GrantAdmin, auth.AuthMiddleware),
	).Methods("POST")

	router.Path("/admin/restore").Handler(
		auth.Middleware(http.HandlerFunc(util.Restore), util.GrantAdmin, auth.AuthMiddleware),
	).Methods("POST")

	router.Path("/admin/onaft_review").Handler(
		auth.Middleware(http.HandlerFunc(domain.OnaftReview), 0),
	).Methods("GET")

	router.Path("/admin/users").Handler(
		auth.Middleware(http.HandlerFunc(domain.GetAllUsers), 0),
	).Methods("GET")

	router.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "main/index.html")
	})

}