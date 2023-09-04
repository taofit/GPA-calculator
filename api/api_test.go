package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/taofit/GPA-calculator/database"
)

func TestHandleRetrieveStudentsGPA(t *testing.T) {
	db, err := database.NewDBInstance()
	if err != nil {
		log.Fatal(err)
	}
	server := NewAPIServer(":8080", db)
	db.ClearTables()
	server.SeedDB("gradeScales.json", "grades.json")
	t.Run("can fetch student gpa", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/students/gpa", nil)
		if err != nil {
			t.Error(err)
		}
		res := httptest.NewRecorder()
		err = server.handleRetrieveStudentsGPA(res, req)
		if err != nil {
			t.Error(err)
		}
		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected 200 ,but got %d", res.Result().StatusCode)
		}
		defer res.Result().Body.Close()
		var expectedRst = []database.StudentGPA{
			{SchoolID: 1, StudentID: 1, GPA: 2},
			{SchoolID: 1, StudentID: 2, GPA: 2.4},
			{SchoolID: 2, StudentID: 1, GPA: 3.6},
			{SchoolID: 2, StudentID: 2, GPA: 2.3},
			{SchoolID: 2, StudentID: 3, GPA: 1.9},
		}
		var gpaArr []database.StudentGPA
		json.NewDecoder(res.Result().Body).Decode(&gpaArr)

		if !reflect.DeepEqual(gpaArr, expectedRst) {
			t.Errorf("expected %v, but got %v", expectedRst, gpaArr)
		}
	})
	defer db.CloseConn()
}
