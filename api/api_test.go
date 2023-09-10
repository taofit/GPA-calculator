package api

import (
	"bytes"
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

	testCases := []struct {
		name      string
		gradeJson string
		expected  []database.StudentGPA
	}{
		{
			name:      "just load the initial GPA",
			gradeJson: "",
			expected: []database.StudentGPA{
				{SchoolID: 1, StudentID: 1, GPA: 2},
				{SchoolID: 1, StudentID: 2, GPA: 2.4},
				{SchoolID: 2, StudentID: 1, GPA: 3.6},
				{SchoolID: 2, StudentID: 2, GPA: 2.3},
				{SchoolID: 2, StudentID: 3, GPA: 1.9},
			},
		},
		{
			name: "load the Grade",
			gradeJson: `[
				{
					"school_id":  1,
					"student_id": 4,
					"course_id":  1,
					"grade":     "A"
				},
				{
					"school_id":  1,
					"student_id": 4,
					"course_id":  2,
					"grade":     "D"
				},
				{
					"school_id":  1,
					"student_id": 4,
					"course_id":  3,
					"grade":     "B"
				},
				{
					"school_id":  1,
					"student_id": 4,
					"course_id":  4,
					"grade":     "B"
				},
				{
					"school_id":  1,
					"student_id": 4,
					"course_id":  5,
					"grade":     "D"
				}
			]`,
			expected: []database.StudentGPA{
				{SchoolID: 1, StudentID: 1, GPA: 2},
				{SchoolID: 1, StudentID: 2, GPA: 2.4},
				{SchoolID: 1, StudentID: 4, GPA: 2.4},
				{SchoolID: 2, StudentID: 1, GPA: 3.6},
				{SchoolID: 2, StudentID: 2, GPA: 2.3},
				{SchoolID: 2, StudentID: 3, GPA: 1.9},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.gradeJson != "" {
				res, err := createStudentsGradeRequest(server, tc.gradeJson)
				checkError(res, err, t)
				res.Result().Body.Close()
			}
			res, err := retrieveStudentGPARequest(server)
			checkError(res, err, t)
			defer res.Result().Body.Close()

			var gpaArr []database.StudentGPA
			json.NewDecoder(res.Result().Body).Decode(&gpaArr)

			if !reflect.DeepEqual(gpaArr, tc.expected) {
				t.Errorf("expected %v, but got %v", tc.expected, gpaArr)
			}
		})
	}
	defer db.CloseConn()
}

func retrieveStudentGPARequest(server *APIServer) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(http.MethodGet, "/students/gpa", nil)
	if err != nil {
		return nil, err
	}
	res := httptest.NewRecorder()
	err = server.handleRetrieveStudentsGPA(res, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func createStudentsGradeRequest(server *APIServer, payload string) (*httptest.ResponseRecorder, error) {
	jsonBody := []byte(payload)
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, "/students/gpa", bodyReader)
	if err != nil {
		return nil, err
	}
	res := httptest.NewRecorder()
	err = server.handleCreateStudentsGrade(res, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func checkError(res *httptest.ResponseRecorder, err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
	if res.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 ,but got %d", res.Result().StatusCode)
	}
}
