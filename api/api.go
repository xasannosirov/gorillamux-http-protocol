package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"gorillaMux/models"
	"gorillaMux/storage"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/user/create", CreateUser).Methods("POST")
	r.HandleFunc("/user/update", UpdateUser).Methods("PUT")
	r.HandleFunc("/user/delete", DeleteUser).Methods("DELETE")
	r.HandleFunc("/user/get", GetUser).Methods("GET")
	r.HandleFunc("/user/all", GetAllUsers).Methods("GET")
	log.Println("Server is running...")
	if err := http.ListenAndServe("localhost:8088", r); err != nil {
		log.Println("Error sever is running!")
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while getting body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user *models.User
	if err = json.Unmarshal(bodyByte, &user); err != nil {
		log.Println("error while unmarshalling body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := uuid.NewString()
	user.Id = id

	respUser, err := storage.CreateUser(user)
	if err != nil {
		log.Println("error while creating user", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respBody, err := json.Marshal(respUser)
	if err != nil {
		log.Println("error while marshalling to response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while getting body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user *models.User
	if err = json.Unmarshal(bodyByte, &user); err != nil {
		log.Println("error while unmarshalling body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user_id := r.URL.Query().Get("id")

	respUser, err := storage.UpdateUser(user_id, user)
	if err != nil {
		log.Println("error while updating user", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respBody, err := json.Marshal(respUser)
	if err != nil {
		log.Println("error while marshalling to response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("id")

	if err := storage.DeleteUser(user_id); err != nil {
		log.Println("error while deleting user", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted User"))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("id")

	respUser, err := storage.GetUser(user_id)
	if err != nil {
		log.Println("Error while getting user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(respUser)
	if err != nil {
		log.Println("error while marshalling body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		log.Println("Error while converting page")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("Error while converting limit")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := storage.GetAllUsers(intPage, intLimit)
	if err != nil {
		log.Println("Error while getting all users", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(users)
	if err != nil {
		log.Println("error while marshalling body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
