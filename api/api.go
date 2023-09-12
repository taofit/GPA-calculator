package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/taofit/GPA-calculator/database"
	"github.com/taofit/GPA-calculator/types"
)

type APIServer struct {
	listenAdr string
	dbHandler database.DBHandle
}

type apiFunc func(http.ResponseWriter, *http.Request) error
type page_key string
type perPage_key string

const pageKey page_key = "page"
const perPageKey perPage_key = "perPage"
const per_page = 10

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

func paginationMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		page, err := strconv.ParseInt(r.URL.Query().Get(string(pageKey)), 10, 0)
		if err != nil {
			page = 1
		}
		perPage, err := strconv.ParseInt(r.URL.Query().Get(string(perPageKey)), 10, 0)
		if err != nil {
			perPage = per_page
		}
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = per_page
		}
		ctx = context.WithValue(ctx, pageKey, page)
		ctx = context.WithValue(ctx, perPageKey, perPage)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

// func (s *APIServer) handleStudentsGPA(w http.ResponseWriter, r *http.Request) error {
// 	if r.Method == http.MethodGet {
// 		return s.handleRetrieveStudentsGPA(w, r)
// 	}

// 	return fmt.Errorf("method not allowed %s", r.Method)
// }

func (s *APIServer) handleRetrieveStudentsGrade(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page := ctx.Value(pageKey).(int64)
	perPage := ctx.Value(perPageKey).(int64)
	stdGrade, err := s.dbHandler.GetStudentsGrade(page, perPage)
	if err != nil {
		return err
	}
	if len(stdGrade) == 0 {
		msg := "There is no record"
		return writeJSON(w, http.StatusOK, msg)
	}
	return writeJSON(w, http.StatusOK, stdGrade)
}

func (s *APIServer) handleRetrieveStudentsGPA(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page := ctx.Value(pageKey).(int64)
	perPage := ctx.Value(perPageKey).(int64)
	stdGPA, err := s.dbHandler.GetStudentsGPA(page, perPage)
	if err != nil {
		return err
	}

	if len(stdGPA) == 0 {
		msg := "There is no record"
		return writeJSON(w, http.StatusOK, msg)
	}
	return writeJSON(w, http.StatusOK, stdGPA)
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
	router.HandleFunc("/students/gpa", paginationMiddleware(makeHTTPHandleFunc(s.handleRetrieveStudentsGPA))).Methods("GET")
	router.HandleFunc("/students/grade", makeHTTPHandleFunc(s.handleCreateStudentsGrade)).Methods("POST")
	router.HandleFunc("/students/grade", paginationMiddleware(makeHTTPHandleFunc(s.handleRetrieveStudentsGrade))).Methods("GET")

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
