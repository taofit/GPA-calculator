package api

import (
	"encoding/json"
	"fmt"
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

func (s *APIServer) handleStudentsGPA(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleRetrieveStudentsGPA(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateStudentsGrade(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleRetrieveStudentsGPA(w http.ResponseWriter, r *http.Request) error {
	stdntsGPA, err := s.dbHandler.GetStudentsGPA()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, stdntsGPA)
}

func (s *APIServer) handleCreateStudentsGrade(w http.ResponseWriter, r *http.Request) error {
	var gradeArr = []types.Grade{}
	if err := json.NewDecoder(r.Body).Decode(&gradeArr); err != nil {
		return err
	}
	for _, g := range gradeArr {
		err := s.dbHandler.CreateGrade(g)
		if err != nil {
			return err
		}
	}
	return writeJSON(w, http.StatusOK, "Students grades are generated successfully!!!")
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/students/gpa", makeHTTPHandleFunc(s.handleStudentsGPA))

	log.Printf("API server is running on port: %s", s.listenAdr)
	log.Fatal(http.ListenAndServe(s.listenAdr, router))
}

func (s *APIServer) SeedDB(args ...string) {
	gsJsonPath := "./api/gradeScales.json"
	gJsonPath := "./api/grades.json"
	if len(args) == 2 {
		gsJsonPath = args[0]
		gJsonPath = args[1]
	}
	gradeScales, err := generateGradeScales(gsJsonPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, gs := range gradeScales {
		err := s.dbHandler.CreateGradeScale(gs)
		if err != nil {
			log.Fatal(err)
		}
	}
	grades, err := generateGrades(gJsonPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, g := range grades {
		err := s.dbHandler.CreateGrade(g)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateGradeScales(gsJsonPath string) ([]types.GradeScale, error) {
	content, err := ioutil.ReadFile(gsJsonPath)
	if err != nil {
		return nil, fmt.Errorf("error when opening the file: %w", err)
	}
	var gradeScales = []types.GradeScale{}
	err = json.Unmarshal(content, &gradeScales)
	if err != nil {
		return nil, fmt.Errorf("error during Unmarshal(): %w", err)
	}

	return gradeScales, nil
}

func generateGrades(gJsonPath string) ([]types.Grade, error) {
	content, err := ioutil.ReadFile(gJsonPath)
	if err != nil {
		return nil, fmt.Errorf("error when opening the file: %w", err)
	}
	var grades = []types.Grade{}
	err = json.Unmarshal(content, &grades)
	if err != nil {
		return nil, fmt.Errorf("error during Unmarshal(): %w", err)
	}

	return grades, nil
}
