package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/taofit/GPA-calculator/database"
	"github.com/taofit/GPA-calculator/types"
)

type APIServer struct {
	listenAdr string
	dbHandler database.DBHandle
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func NewAPIServer(listenAdr string, dbHandle database.DBHandle) *APIServer {
	return &APIServer{
		listenAdr: listenAdr,
		dbHandler: dbHandle,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

func (s *APIServer) handleRetrieveStudentsGPA(w http.ResponseWriter, r *http.Request) error {
	stdntsGPA, err := s.dbHandler.GetStudentsGPA()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, stdntsGPA)
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/students/gpa", makeHTTPHandleFunc(s.handleRetrieveStudentsGPA))

	log.Printf("API server is funning on port: %s", s.listenAdr)
	log.Fatal(http.ListenAndServe(s.listenAdr, router))
}

func (s *APIServer) SeedDB() {
	gradeScales := generateGradeScales()
	for _, gs := range gradeScales {
		err := s.dbHandler.CreateGradeScale(gs)
		if err != nil {
			log.Fatal(err)
		}
	}
	grades := generateGrades()
	for _, g := range grades {
		err := s.dbHandler.CreateGrade(g)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateGradeScales() []types.GradeScale {
	content, err := ioutil.ReadFile("./api/gradeScales.json")
	if err != nil {
		log.Fatal("Error when opening the file: ", err)
	}
	var gradeScales = []types.GradeScale{}
	err = json.Unmarshal(content, &gradeScales)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return gradeScales
}

func generateGrades() []types.Grade {
	content, err := ioutil.ReadFile("./api/grades.json")
	if err != nil {
		log.Fatal("Error when opening the file: ", err)
	}
	var grades = []types.Grade{}
	err = json.Unmarshal(content, &grades)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return grades
}
