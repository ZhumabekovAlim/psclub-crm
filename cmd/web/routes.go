package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, makeResponseJSON)

	dynamicMiddleware := alice.New()

	mux := pat.New()

	//// USERS
	//mux.Post("/users/signup", dynamicMiddleware.ThenFunc(app.userHandler.SignUp))           // Sign up user / work
	//mux.Post("/users/login", dynamicMiddleware.ThenFunc(app.userHandler.LogIn))             // Log in /work
	//mux.Get("/users", standardMiddleware.ThenFunc(app.userHandler.GetAllUsers))             // Get all users /work
	//mux.Get("/users/details/:id", standardMiddleware.ThenFunc(app.userHandler.GetUserByID)) // Get user by ID /work
	//mux.Del("/users/:id", standardMiddleware.ThenFunc(app.userHandler.DeleteUserByID))      // Delete user by ID /work
	//mux.Put("/users/:id", standardMiddleware.ThenFunc(app.userHandler.UpdateUser))          // Update user /work

	return standardMiddleware.Then(mux)
}
