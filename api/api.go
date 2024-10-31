package api

import (
	"GoCRUD/database"
	"GoCRUD/util"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", handlePostUsers)
			r.Get("/", handleGetUsers)
			r.Get("/{id}", handleGetUserByID)
			r.Delete("/{id}", handleDeleteUser)
			r.Put("/{id}", handlePutUser)
		})
	})

	return r
}

func handlePostUsers(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error when reading user json", "error", err)
		util.SendJson(w, util.Response{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}

	newUser, err := database.Insert(data)
	if err != nil {
		var invalidUserBody database.ErrorUserWrongData
		if errors.As(err, &invalidUserBody) {
			slog.Error("wrong user data", "error", err)
			util.SendJson(w, util.Response{Error: "Please provide FirstName LastName and bio for the user"}, http.StatusBadRequest)
			return
		}

		slog.Error("error when creating a new user", "error", err)
		util.SendJson(w, util.Response{Error: "There was an error while saving the user to the database"}, http.StatusInternalServerError)
		return
	}
	util.SendJson(w, util.Response{Data: newUser}, http.StatusCreated)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.FindAll()
	if err != nil {
		slog.Error("error when reading users database", "error", err)
		util.SendJson(w, util.Response{Error: "The users information could not be retrieved"}, http.StatusInternalServerError)
		return
	}

	util.SendJson(w, util.Response{Data: users}, http.StatusOK)
}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	user, err := database.FindByID(idStr)
	if err != nil {
		var usrErr database.ErrorUserNotFound
		if errors.As(err, &usrErr) {
			slog.Error("user not found", "error", err)
			util.SendJson(w, util.Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}
		slog.Error("error when getting user by id", "error", err)
		util.SendJson(w, util.Response{Error: "The user information could not be retrieved"}, http.StatusInternalServerError)
		return
	}
	util.SendJson(w, util.Response{Data: user}, http.StatusOK)
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if err := database.Delete(idStr); err != nil {
		var usrErr database.ErrorUserNotFound
		if errors.As(err, &usrErr) {
			slog.Error("user not found", "error", err)
			util.SendJson(w, util.Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		slog.Error("error when deleting user from database", "error", err)
		util.SendJson(w, util.Response{Error: "The user could not be removed"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handlePutUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error when reading user json", "error", err)
		util.SendJson(w, util.Response{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}

	updatedUser, err := database.Update(idStr, data)
	if err != nil {
		var invalidUserData database.ErrorUserWrongData
		if errors.As(err, &invalidUserData) {
			slog.Error("wrong user data", "error", err)
			util.SendJson(w, util.Response{Error: "Please provide name and bio for the user"}, http.StatusBadRequest)
			return
		}

		var usrErr database.ErrorUserNotFound
		if errors.As(err, &usrErr) {
			slog.Error("user not found", "error", err)
			util.SendJson(w, util.Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}
		slog.Error("error when updating the user information", "error", err)
		util.SendJson(w, util.Response{Error: "The user information could not be modified"}, http.StatusInternalServerError)
		return
	}
	util.SendJson(w, util.Response{Data: updatedUser}, http.StatusOK)
}
