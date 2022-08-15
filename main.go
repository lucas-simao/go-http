package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var errorResponse ErrorResponse

type store map[string]User

type Service struct {
	Store store
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	service := Service{}

	service.Store = make(store)

	mux := http.NewServeMux()

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			service.PostUser(w, r)
		case http.MethodGet:
			service.GetUser(w, r)
		case http.MethodDelete:
			service.DeleteUser(w, r)
		case http.MethodPut:
			service.UpdateUser(w, r)
		default:
			errorResponse.Message = "the request can't be process"
			responseJson(w, http.StatusNotFound, errorResponse)
		}

	})

	if err := http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), mux); err != nil {
		log.Fatal(err)
	}
}

func (s *Service) PostUser(w http.ResponseWriter, r *http.Request) {
	var u User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		errorResponse.Message = err.Error()
		responseJson(w, http.StatusBadRequest, errorResponse)
		return
	}

	if u.Id == "" || u.Name == "" {
		errorResponse.Message = "the body is empty"
		responseJson(w, http.StatusBadRequest, errorResponse)
		return
	}

	if _, ok := s.Store[u.Id]; ok {
		errorResponse.Message = fmt.Sprintf("the id:%v exist", u.Id)
		responseJson(w, http.StatusBadRequest, errorResponse)
		return
	}

	s.Store[u.Id] = u

	responseJson(w, http.StatusCreated, u)
}

func (s *Service) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		errorResponse.Message = "query id is missing"
		responseJson(w, http.StatusBadRequest, errorResponse)
		return
	}

	user, ok := s.Store[id]
	if !ok {
		errorResponse.Message = fmt.Sprintf("user with id: %v don't exist", id)
		responseJson(w, http.StatusBadRequest, errorResponse)
		return
	}

	responseJson(w, http.StatusOK, user)
}

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		errorResponse.Message = "query id is missing"
		responseJson(w, http.StatusBadRequest, errorResponse)
		return
	}

	_, ok := s.Store[id]
	if !ok {
		errorResponse.Message = fmt.Sprintf("user with id: %v don't exist", id)
		responseJson(w, http.StatusBadRequest, errorResponse)
		return
	}

	delete(s.Store, id)
	responseJson(w, http.StatusOK, nil)
}

func (s *Service) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var u User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		errorResponse.Message = "error decode body update"
		responseJson(w, http.StatusBadRequest, errorResponse)
		return
	}

	_, ok := s.Store[u.Id]
	if !ok {
		errorResponse.Message = fmt.Sprintf("user with id: %v don't exist", u.Id)
		responseJson(w, http.StatusBadRequest, errorResponse)
		return
	}

	s.Store[u.Id] = u

	responseJson(w, http.StatusOK, u)
}

func responseJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
