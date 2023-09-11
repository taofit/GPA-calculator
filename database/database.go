package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/taofit/GPA-calculator/types"
)

type StudentGPA struct {
	SchoolID  int     `json:"school_id"`
	StudentID int     `json:"student_id"`
	GPA       float32 `json:"gpa"`
}

type Grade struct {
	SchoolID  int     `json:"school_id"`
	StudentID int     `json:"student_id"`
	CourseID  int     `json:"course_id"`
	Grade     string  `json:"grade"`
	Scale     float32 `json:"scale"`
	percent   int     `json:"percent"`
}

type DBInstance struct {
	db *sql.DB
}

type DBHandle interface {
	GetStudentsGPA() ([]StudentGPA, error)
	CreateGradeScale(types.GradeScale) error
	CreateGrade(types.Grade) error
	GetStudentsGrade(int64, int64) ([]Grade, error)
}

func NewDBInstance() (*DBInstance, error) {
	dbConnStr := fmt.Sprintf("host=database port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	// another way to connect to postgres
	// dbConnStr := fmt.Sprintf("postgres://%s:%s@database:5432/%s?sslmode=disable",
	// 	os.Getenv("POSTGRES_USER"),
	// 	os.Getenv("POSTGRES_PASSWORD"),
	// 	os.Getenv("POSTGRES_DB"),
	// )

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, err
	}
	return &DBInstance{
		db: db,
	}, nil
}

func (d *DBInstance) GetStudentsGPA() ([]StudentGPA, error) {
	rows, err := d.db.Query(`
		SELECT g.school_id, g.student_id, AVG(gs.scale)::NUMERIC(10,1) gpa
		FROM grade g
		INNER JOIN grade_scale gs ON
		g.school_id = gs.school_id AND
		g.grade = gs.grade
		GROUP BY g.school_id, g.student_id
		ORDER BY g.school_id, g.student_id;`)

	if err != nil {
		return nil, err
	}
	sdtsGPA := []StudentGPA{}
	for rows.Next() {
		sdtGPA, err := getStudentGPA(rows)
		if err != nil {
			return nil, err
		}
		sdtsGPA = append(sdtsGPA, sdtGPA)
	}

	return sdtsGPA, nil
}

func (d *DBInstance) CloseConn() {
	defer d.db.Close()
}

func getStudentGPA(row *sql.Rows) (StudentGPA, error) {
	stdntGPA := StudentGPA{}
	err := row.Scan(
		&stdntGPA.SchoolID,
		&stdntGPA.StudentID,
		&stdntGPA.GPA,
	)

	return stdntGPA, err
}

func (d *DBInstance) CreateGradeScale(gs types.GradeScale) error {
	query := `
		INSERT INTO grade_scale (school_id, scale, grade, percent)
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT ON CONSTRAINT grade_scale_school_id_grade_key DO UPDATE
		SET scale = EXCLUDED.scale, percent = EXCLUDED.percent;`
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
		ON CONFLICT ON CONSTRAINT grade_school_id_student_id_course_id_key DO UPDATE
		SET grade = EXCLUDED.grade;`
	_, err := d.db.Query(query, g.SchoolID, g.StudentID, g.CourseID, g.Grade)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBInstance) GetStudentsGrade(page, perPage int64) ([]Grade, error) {
	offSet := (page - 1) * perPage
	rows, err := d.db.Query(`
		SELECT g.school_id, g.student_id, g.course_id, g.grade, gs.scale, gs.percent
		FROM grade g
		INNER JOIN grade_scale gs ON
		g.school_id = gs.school_id AND
		g.grade = gs.grade
		OFFSET $1 LIMIT $2
	`, offSet, perPage)
	if err != nil {
		return nil, err
	}
	stdsGrade := []Grade{}
	for rows.Next() {
		stdGrade, err := getStudentGrade(rows)
		if err != nil {
			return nil, err
		}
		stdsGrade = append(stdsGrade, stdGrade)
	}
	return stdsGrade, nil
}

func getStudentGrade(row *sql.Rows) (Grade, error) {
	stdGrade := Grade{}
	err := row.Scan(
		&stdGrade.SchoolID,
		&stdGrade.StudentID,
		&stdGrade.CourseID,
		&stdGrade.Grade,
		&stdGrade.Scale,
		&stdGrade.percent,
	)

	return stdGrade, err
}

func (db *DBInstance) ClearTables() {
	db.db.Exec("DELETE FROM GRADE_SCALE")
	db.db.Exec("DELETE FROM grade")
}
