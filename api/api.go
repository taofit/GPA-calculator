package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/taofit/GPA-calculator/database"
)

type APIServer struct {
	listenAdr string
	dbHandle  database.DBHandle
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func NewAPIServer(listenAdr string, dbHandle database.DBHandle) *APIServer {
	return &APIServer{
		listenAdr: listenAdr,
		dbHandle:  dbHandle,
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
	stdntsGPA, err := s.dbHandle.GetStudentsGPA()
	if err != nil {
		return nil
	}
	return writeJSON(w, http.StatusOK, stdntsGPA)
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/students/gpa", makeHTTPHandleFunc(s.handleRetrieveStudentsGPA))

	log.Printf("API server is funning on port: %s", s.listenAdr)
	log.Fatal(http.ListenAndServe(s.listenAdr, router))
}
