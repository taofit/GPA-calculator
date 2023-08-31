package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
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
}

func NewDBInstance() (*DBInstance, error) {
	dbConnStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, err
	}

	return &DBInstance{
		db: db,
	}, nil
}

func (d *DBInstance) GetStudentsGPA() ([]StudentGPA, error) {
	rows, err := d.db.Query("")
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
