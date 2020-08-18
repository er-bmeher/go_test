package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"CRUD_POC/dbconn"
	"CRUD_POC/student"

	"github.com/gorilla/mux"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type stud struct {
	//id      string `json:"id"`
	name     string `json:"name"`
	subject  string `json:"subject"`
	location string `json:"location"`
	pin      string `json:"pin"`
}

// Get all Students
func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	students, err := student.GetStudents()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, students)
}

// Get single Student
func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	student, err := student.GetStudent(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, student)
}

// Add new Student
func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.FormValue("name")
	sub := r.FormValue("subject")
	loc := r.FormValue("location")
	pin := r.FormValue("pin")
	//fmt.Println(name, loc, sub, pin)

	student, err := student.CreateStudent(name, sub, loc, pin)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, student)
}

// Update Student
func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	name := r.FormValue("name")
	sub := r.FormValue("subject")
	loc := r.FormValue("location")
	pin := r.FormValue("pin")
	//fmt.Println(name, loc, sub, pin)

	status, err := student.UpdateStudent(params["id"], name, sub, loc, pin)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, status)
}

// Delete Student
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println("in Server deleteStudent Params - ", params)
	fmt.Println("in Server deleteStudent Params - ", params["id"])
	status, err := student.DeleteStudent(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	//json.NewEncoder(w).Encode(books)
	respondWithJSON(w, http.StatusOK, status)
}

func mainStart() {

	// var err error
	// conn, err = dbconn.GetDBConnection()

	// defer conn.Close()
	// if err != nil {
	// 	fmt.Println("Error in db connection.", err)
	// 	return
	// }

	student.InitStudent()
	defer student.ShutDownStudent()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students/{id}", getStudent).Methods("GET")
	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":7777", r))
}

// Main function
func main() {
	mainStart()
}
