package types

type Grade struct {
	SchoolID  int `json:"schoolID"`
	StudentID int `json:"studentID"`
	CourseID  int `json:"courseID"`
	Grade     int `json:"grade"`
}

type GradeScale struct {
	SchoolID   int     `json:"schoolID"`
	Grade      int     `json:"grade"`
	GPA        float32 `json:"gpa"`
	Percentage int     `json:"percentage"`
}
