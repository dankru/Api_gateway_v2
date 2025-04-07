package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strings"
)

type user struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name,omitempty"`
	Age       string    `json:"age,omitempty"`
	Anonymous bool      `json:"anonymous,omitempty"`
}

var users = make(map[uuid.UUID]user, 5)

func main() {

	http.HandleFunc("/user/", handleUsers)
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal("Error on running server: %s", err.Error())
	}
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUser(w, r)
	case http.MethodPost:
		createUser(w, r)
	case http.MethodPut:
		replaceUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid uuid provided"))
	}

	u, ok := users[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var userData user
	_ = json.Unmarshal(body, &userData)

	userData.Id = uuid.New()
	users[userData.Id] = userData
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(userData.Id.String()))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid uuid provided"))
	}
	_, ok := users[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}

func replaceUser(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid uuid provided"))
	}

	_, ok := users[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var userData user
	_ = json.Unmarshal(body, &userData)

	userData.Id = uuid.New()
	users[userData.Id] = userData

	resp, err := json.Marshal(userData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func getIdFromBody(r *http.Request) (uuid.UUID, error) {
	id := strings.TrimPrefix(r.URL.Path, "/user/")
	idParsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return idParsed, err
}
