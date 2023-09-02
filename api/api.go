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
	// grades := generateGrades()
}

func generateGradeScales() []types.GradeScale {
	schoolGradeScales := []types.GradeScale{
		{
			SchoolID:   1,
			Scale:      4.0,
			Grade:      "A",
			Percentage: 90,
		},
		{
			SchoolID:   1,
			Scale:      3.0,
			Grade:      "B",
			Percentage: 80,
		},
		{
			SchoolID:   1,
			Scale:      2.0,
			Grade:      "C",
			Percentage: 70,
		},
		{
			SchoolID:   1,
			Scale:      1.0,
			Grade:      "D",
			Percentage: 60,
		},
		{
			SchoolID:   1,
			Scale:      0.0,
			Grade:      "F",
			Percentage: 59,
		},
		{
			SchoolID:   2,
			Scale:      4.5,
			Grade:      "A+",
			Percentage: 95,
		},
		{
			SchoolID:   2,
			Scale:      4.0,
			Grade:      "A",
			Percentage: 90,
		},
		{
			SchoolID:   2,
			Scale:      3.5,
			Grade:      "B+",
			Percentage: 85,
		},
		{
			SchoolID:   2,
			Scale:      3.0,
			Grade:      "B",
			Percentage: 80,
		},
		{
			SchoolID:   2,
			Scale:      2.5,
			Grade:      "C+",
			Percentage: 75,
		},
		{
			SchoolID:   2,
			Scale:      2.0,
			Grade:      "C",
			Percentage: 70,
		},
		{
			SchoolID:   2,
			Scale:      1.5,
			Grade:      "D+",
			Percentage: 65,
		},
		{
			SchoolID:   2,
			Scale:      1.0,
			Grade:      "D",
			Percentage: 60,
		},
		{
			SchoolID:   2,
			Scale:      0.0,
			Grade:      "F",
			Percentage: 59,
		},
	}

	return schoolGradeScales
}

func generateGrades() []types.Grade {
	content, err := ioutil.ReadFile("./grades.json")
	if err != nil {
		log.Fatal("Error when opening the file: ", err)
	}
	var grades = []types.Grade{}
	err = json.Unmarshal(content, &grades)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	grades = []types.Grade{
		{
			SchoolID:  1,
			StudentID: 1,
			CourseID:  1,
			Grade:     "B",
		},
		{
			SchoolID:  1,
			StudentID: 1,
			CourseID:  2,
			Grade:     "A",
		},
		{
			SchoolID:  1,
			StudentID: 1,
			CourseID:  3,
			Grade:     "C",
		},
		{
			SchoolID:  1,
			StudentID: 1,
			CourseID:  4,
			Grade:     "D",
		},
	}

	return grades
}
