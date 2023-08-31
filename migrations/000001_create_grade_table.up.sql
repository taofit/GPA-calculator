CREATE TABLE IF NOT EXISTS grade (
    id INT PRIMARY KEY,
    school_id INT,
    student_id INT,
    course_id INT,
    grade CHAR(1),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);