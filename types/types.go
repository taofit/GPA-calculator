package types

type Grade struct {
	SchoolID  int    `json:"school_id"`
	StudentID int    `json:"student_id"`
	CourseID  int    `json:"course_id"`
	Grade     string `json:"grade"`
}

type GradeScale struct {
	SchoolID   int     `json:"school_id"`
	Scale      float32 `json:"scale"`
	Grade      string  `json:"grade"`
	Percentage int     `json:"percentage"`
}
