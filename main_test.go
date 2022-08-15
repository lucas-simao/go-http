package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostUser(t *testing.T) {
	cases := map[string]struct {
		body       string
		statusCode int
	}{
		"should register user": {
			body:       `{"id":"1","name":"lucas"}`,
			statusCode: http.StatusCreated,
		},
		"should return error 400": {
			body:       `{"id":"1","name":"lucas"}`,
			statusCode: http.StatusBadRequest,
		},
		"should return error 400, empty body": {
			body:       `{}`,
			statusCode: http.StatusBadRequest,
		},
	}

	store := make(store)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(test.body))
			responseRecorder := httptest.NewRecorder()

			service := Service{
				Store: store,
			}
			service.PostUser(responseRecorder, request)

			if responseRecorder.Code != test.statusCode {
				t.Errorf("Want status '%d', got '%d'", test.statusCode, responseRecorder.Code)
			}

			if test.statusCode == http.StatusCreated {
				if strings.TrimSpace(string(responseRecorder.Body.String())) != test.body {
					t.Errorf("Want user '%v', got '%v'", test.body, string(responseRecorder.Body.String()))
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	cases := map[string]struct {
		id         string
		want       string
		statusCode int
	}{
		"should return user id 1": {
			id:         "1",
			want:       `{"id":"1","name":"lucas"}`,
			statusCode: http.StatusOK,
		},
		"should return error 400": {
			id:         "2",
			want:       `{"id":"1","name":"lucas"}`,
			statusCode: http.StatusBadRequest,
		},
		"should return error 400, without query id": {
			id:         "",
			want:       `{"id":"1","name":"lucas"}`,
			statusCode: http.StatusBadRequest,
		},
	}

	user := User{
		Id:   "1",
		Name: "lucas",
	}

	store := make(store)

	store[user.Id] = user

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user?id=%v", test.id), nil)
			responseRecorder := httptest.NewRecorder()

			service := Service{
				Store: store,
			}
			service.GetUser(responseRecorder, request)

			if responseRecorder.Code != test.statusCode {
				t.Errorf("Want status '%d', got '%d'", test.statusCode, responseRecorder.Code)
			}

			if test.statusCode == http.StatusOK {
				if strings.TrimSpace(string(responseRecorder.Body.String())) != test.want {
					t.Errorf("Want user '%v', got '%v'", test.want, string(responseRecorder.Body.String()))
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	cases := map[string]struct {
		id         string
		statusCode int
	}{
		"should delete user id 1": {
			id:         "1",
			statusCode: http.StatusOK,
		},
		"should return error 400": {
			id:         "2",
			statusCode: http.StatusBadRequest,
		},
		"should return error 400, without query id": {
			id:         "",
			statusCode: http.StatusBadRequest,
		},
	}

	user := User{
		Id:   "1",
		Name: "lucas",
	}

	store := make(store)

	store[user.Id] = user

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user?id=%v", test.id), nil)
			responseRecorder := httptest.NewRecorder()

			service := Service{
				Store: store,
			}
			service.DeleteUser(responseRecorder, request)

			if responseRecorder.Code != test.statusCode {
				t.Errorf("Want status '%d', got '%d'", test.statusCode, responseRecorder.Code)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	cases := map[string]struct {
		body       string
		statusCode int
	}{
		"should register user": {
			body:       `{"id":"1","name":"lucas simão"}`,
			statusCode: http.StatusOK,
		},
		"should return error 400": {
			body:       `{"id":"2","name":"lucas"}`,
			statusCode: http.StatusBadRequest,
		},
		"should return error 400, empty body": {
			body:       `{}`,
			statusCode: http.StatusBadRequest,
		},
	}

	user := User{
		Id:   "1",
		Name: "lucas",
	}

	store := make(store)

	store[user.Id] = user

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(test.body))
			responseRecorder := httptest.NewRecorder()

			service := Service{
				Store: store,
			}
			service.UpdateUser(responseRecorder, request)

			if responseRecorder.Code != test.statusCode {
				t.Errorf("Want status '%d', got '%d'", test.statusCode, responseRecorder.Code)
			}

			if test.statusCode == http.StatusOK {
				if strings.TrimSpace(string(responseRecorder.Body.String())) != test.body {
					t.Errorf("Want user '%v', got '%v'", test.body, string(responseRecorder.Body.String()))
				}
			}
		})
	}
}
