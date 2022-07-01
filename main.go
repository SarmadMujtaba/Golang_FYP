package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Users struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Pass string `json:"pass"`
}

func main() {
	fmt.Println("Program Running...")
	API()
}

func getData(w http.ResponseWriter, r *http.Request) {
	var user []Users
	fileData, err := ioutil.ReadFile("users.json")
	json.Unmarshal(fileData, &user)
	var count int = 0
	id := r.URL.Query().Get("id")
	if len(id) > 0 {
		for _, v := range user {
			if v.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(user[count])
				return
			}
			count++
		}
	}
	if err != nil {
		fmt.Fprintf(w, "File not found!!")
		w.WriteHeader(404)
		return
	}
	json.Unmarshal(fileData, &user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func postData(w http.ResponseWriter, r *http.Request) {

	var user []Users
	var add Users
	fileData, err := ioutil.ReadFile("users.json")
	if err != nil {
		fmt.Fprintf(w, "File not found!!")
		w.WriteHeader(404)
		return
	}
	json.Unmarshal(fileData, &user)
	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)
	for _, v := range user {
		if v.ID == dataToCompare["id"] {
			w.WriteHeader(409)
			fmt.Fprintf(w, "ID already in use!!")
			return
		}
	}
	add.ID = dataToCompare["id"]
	add.Name = dataToCompare["name"]
	add.Pass = dataToCompare["pass"]
	user = append(user, add)
	out, _ := json.MarshalIndent(user, "", "\t")
	_ = ioutil.WriteFile("users.json", out, 0755)
	w.WriteHeader(201)
	fmt.Fprintf(w, "User Created!!")
	return
}

func API() {
	r := mux.NewRouter()
	r.HandleFunc("/users", getData).Methods("GET")
	r.HandleFunc("/users", postData).Methods("POST")
	log.Fatal(http.ListenAndServe(":5005", r))
}
