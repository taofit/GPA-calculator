package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/taofit/GPA-calculator/types"
)

type StudentGPA struct {
	SchoolID int     `json:"schoolID"`
	GPA      float32 `json:"gpa"`
}

type DBInstance struct {
	db *sql.DB
}

type DBHandle interface {
	GetStudentsGPA() ([]StudentGPA, error)
	CreateGradeScale(types.GradeScale) error
	CreateGrade(types.Grade) error
}

func NewDBInstance() (*DBInstance, error) {
	dbConnStr := fmt.Sprintf("host=database port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	// dbConnStr := "postgres://admin:password@localhost/gpa_calculator?sslmode=verify-full"
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, err
	}

	return &DBInstance{
		db: db,
	}, nil
}

func (d *DBInstance) GetStudentsGPA() ([]StudentGPA, error) {
	rows, err := d.db.Query("select * from grade")
	if err != nil {
		return nil, err
	}
	stdntsGPA := []StudentGPA{}
	for rows.Next() {
		stdntGPA, err := getStudentGPA(rows)
		if err != nil {
			return nil, err
		}
		stdntsGPA = append(stdntsGPA, stdntGPA)
	}

	return stdntsGPA, nil
}

func (d *DBInstance) CloseConn() {
	defer d.db.Close()
}

func getStudentGPA(row *sql.Rows) (StudentGPA, error) {
	stdntGPA := StudentGPA{}
	err := row.Scan(
		&stdntGPA.SchoolID,
		&stdntGPA.GPA,
	)

	return stdntGPA, err
}

func (d *DBInstance) CreateGradeScale(gs types.GradeScale) error {
	query := `
		INSERT INTO grade_scale (school_id, scale, grade, percent)
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT ON CONSTRAINT grade_scale_school_id_grade_key DO NOTHING;`
	_, err := d.db.Query(query, gs.SchoolID, gs.Scale, gs.Grade, gs.Percentage)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBInstance) CreateGrade(g types.Grade) error {
	query := `
		INSERT INTO grade (school_id, student_id, course_id, grade)
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT ON CONSTRAINT grade_school_id_student_id_course_id_key DO NOTHING;`
	_, err := d.db.Query(query, g.SchoolID, g.StudentID, g.CourseID, g.Grade)
	if err != nil {
		return err
	}
	return nil
}
